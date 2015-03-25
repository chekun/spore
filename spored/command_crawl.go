package main

import (
	"strings"
)

type CrawlCommand struct {
}

func (c *CrawlCommand) Help() string {
	helpText := `
Usage: spored crawl [options] ...

  Crawl Baoz.cn for data.

Options:

  -env=development    Environment.
`

	return strings.TrimSpace(helpText)
}

func (c *CrawlCommand) Synopsis() string {
	return "Crawl Baoz.cn for data"
}

func (c *CrawlCommand) Run(arsg []string) int {
	return 0
}
