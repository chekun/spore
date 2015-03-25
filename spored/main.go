package main

import (
	"github.com/mitchellh/cli"
	"log"
	"os"
)

var ui cli.Ui

func main() {

	ui = &cli.BasicUi{Writer: os.Stdout}

	app := &cli.CLI{
		HelpFunc: cli.BasicHelpFunc("pored"),
		Args:     os.Args[1:],
		Version:  "1.0.0",
		Commands: map[string]cli.CommandFactory{
			"crawl": func() (cli.Command, error) {
				return &CrawlCommand{}, nil
			},
		},
	}

	exitCode, err := app.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitCode)
}
