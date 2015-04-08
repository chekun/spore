package main

import (
	"log"
	"os"

	"github.com/chekun/spore/spored/command"
	"github.com/mitchellh/cli"
)

var ui *cli.BasicUi

func main() {

	ui = &cli.BasicUi{Writer: os.Stdout}

	app := &cli.CLI{
		HelpFunc: cli.BasicHelpFunc("spored"),
		Args:     os.Args[1:],
		Version:  "1.0.0",
		Commands: map[string]cli.CommandFactory{
			"crawl": func() (cli.Command, error) {
				return &command.CrawlCommand{ui}, nil
			},
			"serve": func() (cli.Command, error) {
				return &command.ServeCommand{ui}, nil
			},
		},
	}

	exitCode, err := app.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitCode)
}
