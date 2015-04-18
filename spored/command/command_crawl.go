package command

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	l4g "code.google.com/p/log4go"
	"github.com/chekun/spore/spored/env"
	"github.com/chekun/spore/spored/schema"
	"github.com/mitchellh/cli"
)

//CrawlCommand Command Object
type CrawlCommand struct {
	UI *cli.BasicUi
}

//Help Crawl Command Help
func (c *CrawlCommand) Help() string {
	helpText := `
Usage: spored crawl [options] ...

  Crawl Baoz.cn for data.

Options:

  -config=config.yml  Configuration file to use.
  -env=development    Environment.
`

	return strings.TrimSpace(helpText)
}

//Synopsis Crawl Command Synopsis
func (c *CrawlCommand) Synopsis() string {
	return "Crawl Baoz.cn for data"
}

var messageChan chan *schema.MessageBase
var groupChan chan *schema.GroupBase
var quitChan chan int
var resumeChan chan int
var isRunning bool
var idChan chan int
var sharedHTTPClient *http.Client
var continuousFailedIds []int

//Run Crawl Command Run
func (c *CrawlCommand) Run(args []string) int {

	cmdFlags := flag.NewFlagSet("crawl", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.UI.Output(c.Help()) }
	env.ConfigFlags(cmdFlags)

	if err := cmdFlags.Parse(args); err != nil {
		fmt.Println("Could not parse config: ", err)
		return 1
	}

	environment, err := env.GetEnvironment()
	if err != nil {
		fmt.Println("Could not parse config ", err)
		return 1
	}

	db, err := env.GetConnection(environment)
	if err != nil {
		l4g.Error("Could not Get DB Connection: ", err)
		return 1
	}

	sharedHTTPClient = &http.Client{}
	messageChan = make(chan *schema.MessageBase)
	groupChan = make(chan *schema.GroupBase)
	quitChan = make(chan int)
	resumeChan = make(chan int)
	idChan = make(chan int)
	isRunning = false
	continuousFailedIds = []int{}

	go stopper()
	go scheduler()

	go func() {
		resumeChan <- 1
	}()

	for {
		select {
		case id := <-idChan:
			go crawler(id + 1)
		case m := <-messageChan:
			m.Message.Save(db)
			go next(m.ID)
		case g := <-groupChan:
			g.Group.Save(db)
			go next(g.ID)
		case <-quitChan:
			time.Sleep(1 * time.Second)
			l4g.Info("[crawl.scheduler] Goodbye!")
			time.Sleep(100 * time.Millisecond)
			return 0
		case status := <-resumeChan:
			// status = 1 will trigger crawling start.
			l4g.Info("[crawl.scheduler] Resume Signal %d", status)
			if status == 1 && !isRunning {
				//check db for max(id).
				var maxID int64
				statment, err := db.Prepare(`
					SELECT MAX(id) FROM
						(SELECT MAX(id) AS id FROM users u
						UNION
						SELECT MAX(id) AS id FROM threads t
						UNION
						SELECT MAX(id) AS id FROM posts p)
					utp
					`)
				if err != nil {
					l4g.Error("[mysql.getMax] %s", err)
				}
				err = statment.QueryRow().Scan(&maxID)
				if err != nil {
					l4g.Error("[mysql.getMax.scan] %s", err)
				}
				statment.Close()
				l4g.Info("[mysql.getMax] Get MaxId %d", maxID)
				l4g.Info("[crawl.scheduler] Resumed from %d", maxID)
				isRunning = true
				go next(int(maxID))
			}
			if status == 0 && isRunning {
				l4g.Info("[crawl.scheduler] Paused")
				isRunning = false
			}
		}

	}
}

func stopper() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			l4g.Info("[keyboard.interrupt] %s", sig)
			quitChan <- 1
		}
	}()
}

func scheduler() {
	for {
		select {
		case <-time.NewTicker(5 * time.Minute).C:
			resumeChan <- 1
		}
	}
}

func newRequest(method, urlString string) (*http.Request, error) {
	request, err := http.NewRequest(method, urlString, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/35.0.1916.114 Safari/537.36")
	return request, nil
}

func next(id int) {
	//crawl slowly ^_^
	time.Sleep(100 * time.Millisecond)
	idChan <- id
}

func crawler(id int) {
	url := fmt.Sprintf("http://i.baoz.cn/rest/object/get.%d", id)
	request, err := newRequest("GET", url)
	if err != nil {
		l4g.Error("[crawler.request] %s", err)
	}
	response, err := sharedHTTPClient.Do(request)
	if err != nil {
		l4g.Error("[crawler.request] %d %s", id, err)
		go next(id)
		return
	}
	l4g.Debug("[crawler.request] %s %s", url, response.Status)
	if response.StatusCode != 200 {
		l4g.Error("[crawler.request] %d %s", id, response.Status)
		//when to determin pause signal
		continuousFailedIds = append(continuousFailedIds, id)
		length := len(continuousFailedIds)
		if length > 50 {
			if length == (continuousFailedIds[length-1] - continuousFailedIds[0] + 1) {
				l4g.Info("[crawler.scheduler] Pause Signal Detected")
				//let's pause it
				resumeChan <- 0
			}
			continuousFailedIds = nil
			l4g.Info("[crawler.scheduler] flushed continuousFailedIds")
		} else {
			go next(id)
		}
		return
	} else {
		if len(continuousFailedIds) > 0 {
			l4g.Info("[crawler.scheduler] soft flushed continuousFailedIds")
			continuousFailedIds = nil
		}
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		l4g.Error("[crawler.read] %d %s", id, err)
	}
	defer response.Body.Close()
	//get rid of define() staff
	body = body[7 : len(body)-2]
	if string(body) == `{"error":"noPower"}` {
		//pass
		go next(id)
		return
	}
	var define schema.BaseDefine
	err = json.Unmarshal(body, &define)
	if err != nil {
		l4g.Error("[json.unmarshal.define] %d %s", id, err)
		go next(id)
		return
	}

	switch define.Base.Origin {
	case "message":
		var message schema.MessageBase
		err = json.Unmarshal(body, &message)
		if err != nil {
			l4g.Error("[json.unmarshal.message] %d %s", id, err)
			go next(id)
			return
		}
		messageChan <- &message
	case "group":
		var group schema.GroupBase
		err = json.Unmarshal(body, &group)
		if err != nil {
			l4g.Error("[json.unmarshal.group] %d %s", id, err)
			go next(id)
			return
		}
		groupChan <- &group
	default:
		//pass the object we don't know yet!
		l4g.Error("[unknown.base.origin] %s", define.Base.Origin)
		go next(id)
	}
}
