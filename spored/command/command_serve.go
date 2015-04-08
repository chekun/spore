package command

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/astaxie/goredis"
	"github.com/chekun/spore/spored/env"
	"github.com/chekun/spore/spored/schema"
	"github.com/mitchellh/cli"
	"github.com/slvmnd/gosphinx"
)

// ServeCommand Serve Command
type ServeCommand struct {
	UI *cli.BasicUi
}

// Help Serve Command Help
func (c *ServeCommand) Help() string {
	helpText := `
Usage: spored serve [options] ...

  Spored Server.

Options:

  -config=dbconfig.yml   Configuration file to use.
  -env=development       Environment.
`

	return strings.TrimSpace(helpText)
}

// Synopsis Serve Command Synopsis
func (c *ServeCommand) Synopsis() string {
	return "Spored Server"
}

var mysql *sql.DB

// Run Serve Command Implementation
func (c *ServeCommand) Run(args []string) int {

	cmdFlags := flag.NewFlagSet("serve", flag.ContinueOnError)
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

	mysql, err = env.GetConnection(environment)
	if err != nil {
		fmt.Println("Could not Get DB Connection: ", err)
		return 1
	}

	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))
	http.HandleFunc("/", HomeEntry)
	http.HandleFunc("/search", SearchEntry)
	http.HandleFunc("/search/do", DoSearchHandler)
	err = http.ListenAndServe(environment.HTTP, nil)
	if err != nil {
		fmt.Println("Server Failed to Start, ", err.Error())
		return 1
	}
	return 0
}

// HomeEntry Home Handler
func HomeEntry(w http.ResponseWriter, r *http.Request) {
	t := template.New("home.html")
	t, _ = t.ParseFiles("../resources/views/home.html")
	t.Execute(w, r.Cookies())
	//fmt.Fprintf(w, "OK")
}

// SearchEntry Search Handler
func SearchEntry(w http.ResponseWriter, r *http.Request) {
	t := template.New("search.html")
	t, _ = t.ParseFiles("../resources/views/search.html")
	t.Execute(w, nil)
}

// DoSearchHandler Do Search Handler
func DoSearchHandler(w http.ResponseWriter, r *http.Request) {

	type SearchResult struct {
		Total   int             `json:"total"`
		Groups  map[string]uint `json:"groups"`
		Results []interface{}   `json:"results"`
		Page    int             `json:"page"`
		Error   string          `json:"error"`
	}

	environment, _ := env.GetEnvironment()
	sphinx := gosphinx.NewSphinxClient()
	sphinx.SetServer(environment.SphinxServer, environment.SphinxPort)
	sphinx.Open()
	defer sphinx.Close()

	searchType := r.URL.Query().Get("t")
	searchWord := r.URL.Query().Get("q")
	page := 1
	if r.URL.Query().Get("page") != "" {
		page, _ = strconv.Atoi(r.URL.Query().Get("page"))
		if page == 0 {
			page = 1
		}
	}

	if searchWord == "" {
		result := SearchResult{Error: "请输入要搜索的关键词"}
		resultJSON, _ := json.Marshal(result)
		fmt.Fprintf(w, string(resultJSON))
		return
	}

	result := SearchResult{Error: ""}

	if searchType == "" {
		sphinx.SetGroupBy("id_type", gosphinx.SPH_GROUPBY_ATTR, "@group asc")
		res, err := sphinx.Query(searchWord, "spored", "Group Query for "+searchWord)
		if err != nil {
			result.Error = "Sphinx Error" + err.Error()
			resultJSON, _ := json.Marshal(result)
			fmt.Fprintf(w, string(resultJSON))
			return
		}

		groups := map[string]uint{}
		for _, r := range res.Matches {
			typeRef := reflect.ValueOf(r.AttrValues[0])
			switch typeRef.Uint() {
			case 1:
				ref := reflect.ValueOf(r.AttrValues[2])
				groups["users"] = uint(ref.Uint())
			case 2:
				ref := reflect.ValueOf(r.AttrValues[2])
				groups["groups"] = uint(ref.Uint())
			case 3:
				ref := reflect.ValueOf(r.AttrValues[2])
				groups["threads"] = uint(ref.Uint())
			}
		}
		result.Groups = groups
	}

	sphinx.ResetFilters()
	sphinx.ResetGroupBy()

	if searchType != "" {
		searchTypeInt, _ := strconv.Atoi(searchType)
		sphinx.SetFilter("id_type", []uint64{uint64(searchTypeInt)}, false)
	}
	sphinx.SetLimits((page-1)*10, 10, 5000, 5000)
	res, err := sphinx.Query(searchWord, "spored", "Query for "+searchWord)
	if err != nil {
		result.Error = "Sphinx Error" + err.Error()
		resultJSON, _ := json.Marshal(result)
		fmt.Fprintf(w, string(resultJSON))
		return
	}

	results := []interface{}{}
	for _, r := range res.Matches {
		redis := goredis.Client{}
		redis.Addr = environment.RedisServer + ":" + strconv.Itoa(environment.RedisPort)
		idTypeRef := reflect.ValueOf(r.AttrValues[0])
		idType := idTypeRef.Uint()
		redisKey := fmt.Sprintf("spored:%d_%d", idType, r.DocId)
		redisValue, _ := redis.Get(redisKey)
		if redisValue == nil || string(redisValue) == "" {
			//update cache
			var jsonBytes []byte
			switch idType {
			case 1:
				object := schema.User{}
				object.New(mysql, int(r.DocId))
				jsonBytes, _ = json.Marshal(object)
				results = append(results, object)
			case 2:
				object := schema.Group{}
				object.New(mysql, int(r.DocId))
				jsonBytes, _ = json.Marshal(object)
				results = append(results, object)
			case 3:
				object := schema.Message{}
				object.New(mysql, int(r.DocId))
				jsonBytes, _ = json.Marshal(object)
				results = append(results, object)
			}
			if jsonBytes != nil {
				redis.Set(redisKey, jsonBytes)
				redis.Expire(redisKey, 24*3600)
			}
		} else {
			switch idType {
			case 1:
				object := schema.User{}
				json.Unmarshal(redisValue, &object)
				results = append(results, object)
			case 2:
				object := schema.Group{}
				json.Unmarshal(redisValue, &object)
				results = append(results, object)
			case 3:
				object := schema.Message{}
				json.Unmarshal(redisValue, &object)
				results = append(results, object)
			}
		}
	}
	result.Results = results
	result.Page = page
	result.Total = res.Total
	resultJSON, _ := json.Marshal(result)
	fmt.Fprintf(w, string(resultJSON))

}
