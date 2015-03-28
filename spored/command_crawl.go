package main

import (
	l4g "code.google.com/p/log4go"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"spore"
	"spore/schema"
	"strings"
	"time"
)

type CrawlCommand struct {
}

func (c *CrawlCommand) Help() string {
	helpText := `
Usage: spored crawl [options] ...

  Crawl Baoz.cn for data.

Options:

  -config=dbconfig.yml   Configuration file to use.
  -env=development    Environment.
`

	return strings.TrimSpace(helpText)
}

func (c *CrawlCommand) Synopsis() string {
	return "Crawl Baoz.cn for data"
}

var messageChan chan *schema.MessageBase
var groupChan chan *schema.GroupBase
var quitChan chan int
var resumeChan chan int
var isRunning bool
var idChan chan int
var sharedHttpClient *http.Client
var continuousFailedIds []int

func (c *CrawlCommand) Run(args []string) int {

	cmdFlags := flag.NewFlagSet("crawl", flag.ContinueOnError)
	cmdFlags.Usage = func() { ui.Output(c.Help()) }
	spore.ConfigFlags(cmdFlags)

	if err := cmdFlags.Parse(args); err != nil {
		fmt.Println("Could not parse config: %s", err)
		return 1
	}

	env, err := spore.GetEnvironment()
	if err != nil {
		fmt.Println("Could not parse config ", err)
		return 1
	}

	db, err := spore.GetConnection(env)
	if err != nil {
		l4g.Error("Could not Get DB Connection: %s", err)
		return 1
	}

	sharedHttpClient = &http.Client{}
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
			go next(m.Id)
		case g := <-groupChan:
			g.Group.Save(db)
			go next(g.Id)
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
				var maxId int64
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
				err = statment.QueryRow().Scan(&maxId)
				if err != nil {
					l4g.Error("[mysql.getMax.scan] %s", err)
				}
				statment.Close()
				l4g.Info("[mysql.getMax] Get MaxId %d", maxId)
				l4g.Info("[crawl.scheduler] Resumed from %d", maxId)
				go next(int(maxId))
			} else {
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
	response, err := sharedHttpClient.Do(request)
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
		if len(continuousFailedIds) > 20 {
			isResumeDetected := true
			matchFailedId := 0
			for _, failedId := range continuousFailedIds {
				if matchFailedId == 0 {
					matchFailedId = failedId
					continue
				}
				if failedId-matchFailedId != 1 {
					isResumeDetected = false
					break
				}
				matchFailedId = failedId
			}
			continuousFailedIds = nil
			l4g.Info("[crawler.scheduler] flush continuousFailedIds")
			if isResumeDetected {
				l4g.Info("[crawler.scheduler] Pause Signal Detected")
				//let's pause it
				resumeChan <- 0
			}
		} else {
			go next(id)
		}
		return
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
	}
}
