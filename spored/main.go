package main

import (
	"log"
	"os"

	"github.com/chekun/spore/spored/command"
	"github.com/chekun/spore/spored/env"
	"github.com/mitchellh/cli"
)

var ui *cli.BasicUi

func main() {

	ui = &cli.BasicUi{Writer: os.Stdout}

	app := &cli.CLI{
		HelpFunc: cli.BasicHelpFunc("spored"),
		Args:     os.Args[1:],
		Version:  env.VERSION,
		Commands: map[string]cli.CommandFactory{
			"crawl": func() (cli.Command, error) {
				return &command.CrawlCommand{ui}, nil
			},
			"serve": func() (cli.Command, error) {
				return &command.ServeCommand{ui}, nil
			},
			"stat": func() (cli.Command, error) {
				return &command.StatCommand{ui}, nil
			},
			"total": func() (cli.Command, error) {
				return &command.TotalCommand{ui}, nil
			},
		},
	}

	exitCode, err := app.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitCode)
}
