package env

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"runtime"

	l4g "code.google.com/p/log4go"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v1"
)

var ConfigFile string
var ConfigEnvironment string

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	l4g.AddFilter("stdout", l4g.INFO, l4g.NewConsoleLogWriter())
	l4g.AddFilter("logfile", l4g.DEBUG, l4g.NewFileLogWriter("spored.log", false))
}

func ConfigFlags(f *flag.FlagSet) {
	f.StringVar(&ConfigFile, "config", "dbconfig.yml", "Configuration file to use.")
	f.StringVar(&ConfigEnvironment, "env", "development", "Environment to use.")
}

type Environment struct {
	Dialect      string `yaml:"dialect"`
	DataSource   string `yaml:"datasource"`
	Dir          string `yaml:"dir"`
	TableName    string `yaml:"table"`
	HTTP         string `yaml:"http"`
	SphinxServer string `yaml:"sphinx_server"`
	SphinxPort   int    `yaml:"sphinx_port"`
	RedisServer  string `yaml:"redis_server"`
	RedisPort    int    `yaml:"redis_port"`
}

func ReadConfig() (map[string]*Environment, error) {
	file, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		return nil, err
	}

	config := make(map[string]*Environment)
	err = yaml.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func GetEnvironment() (*Environment, error) {
	config, err := ReadConfig()
	if err != nil {
		return nil, err
	}

	env := config[ConfigEnvironment]
	if env == nil {
		return nil, errors.New("No environment: " + ConfigEnvironment)
	}

	if env.DataSource == "" {
		return nil, errors.New("No data source specified")
	}

	return env, nil
}

func GetConnection(env *Environment) (*sql.DB, error) {
	db, err := sql.Open("mysql", env.DataSource)
	if err != nil {
		return nil, fmt.Errorf("Cannot connect to database: %s", err)
	}

	return db, nil
}
