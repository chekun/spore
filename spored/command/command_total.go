package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/chekun/spore/spored/env"
	"github.com/chekun/spore/spored/schema"
	"github.com/mitchellh/cli"
)

// TotalCommand Stat Command
type TotalCommand struct {
	UI *cli.BasicUi
}

// Help Stat Command Help
func (c *TotalCommand) Help() string {
	helpText := `
Usage: spored serve [options] ...

  Spored Stat.

Options:

  -config=config.yml   Configuration file to use.
  -env=development     Environment.
`

	return strings.TrimSpace(helpText)
}

// Synopsis Serve Command Synopsis
func (c *TotalCommand) Synopsis() string {
	return "Spored Stat"
}

// Run Stat Command Implementation
func (c *TotalCommand) Run(args []string) int {

	cmdFlags := flag.NewFlagSet("total", flag.ContinueOnError)
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

	fixPostsWithNoGroup()

	//loop users
	totalStatUsers()
	//loop groups
	totalStatGroups()

	return 0
}

func totalStatGroups() {

	statSql := "REPLACE INTO stats VALUES(?, ?, ?, ?, ?)"
	statStmt, _ := mysql.Prepare(statSql)
	defer statStmt.Close()

	sql := "SELECT count(id), group_id FROM threads GROUP BY group_id"
	stmt, _ := mysql.Prepare(sql)
	rows, _ := stmt.Query()
	for rows.Next() {
		var value, groupID int
		rows.Scan(&value, &groupID)
		statStmt.Exec(groupID, schema.OwnerGroup, fmt.Sprintf("1970-01-0%d", schema.StatThreads), schema.StatThreads, value)
	}
	rows.Close()
	stmt.Close()

	sql = "SELECT count(id), group_id FROM posts GROUP BY group_id"
	stmt, _ = mysql.Prepare(sql)
	rows, _ = stmt.Query()
	for rows.Next() {
		var value, groupID int
		rows.Scan(&value, &groupID)
		statStmt.Exec(groupID, schema.OwnerGroup, fmt.Sprintf("1970-01-0%d", schema.StatPosts), schema.StatPosts, value)
	}
	rows.Close()
	stmt.Close()

	sql = `
    SELECT count(user_id),
       group_id
    FROM
    ( SELECT DISTINCT user_id,
                    group_id
    FROM
     ( SELECT user_id,
              group_id
      FROM threads t
      UNION SELECT user_id,
                   group_id
      FROM posts p) tp ) tptp
    GROUP BY group_id
    `
	stmt, _ = mysql.Prepare(sql)
	rows, _ = stmt.Query()
	for rows.Next() {
		var value, groupID int
		rows.Scan(&value, &groupID)
		statStmt.Exec(groupID, schema.OwnerGroup, fmt.Sprintf("1970-01-0%d", schema.StatLives), schema.StatLives, value)
	}
	rows.Close()
	stmt.Close()
}

func totalStatUsers() {
	statSql := "REPLACE INTO stats VALUES(?, ?, ?, ?, ?)"
	statStmt, _ := mysql.Prepare(statSql)
	defer statStmt.Close()

	sql := "SELECT count(id), user_id FROM threads GROUP BY user_id"
	stmt, _ := mysql.Prepare(sql)
	rows, _ := stmt.Query()
	for rows.Next() {
		var value, userID int
		rows.Scan(&value, &userID)
		statStmt.Exec(userID, schema.OwnerUser, fmt.Sprintf("1970-01-0%d", schema.StatThreads), schema.StatThreads, value)
	}
	rows.Close()
	stmt.Close()

	sql = "SELECT count(id), user_id FROM posts GROUP BY user_id"
	stmt, _ = mysql.Prepare(sql)
	rows, _ = stmt.Query()
	for rows.Next() {
		var value, userID int
		rows.Scan(&value, &userID)
		statStmt.Exec(userID, schema.OwnerUser, fmt.Sprintf("1970-01-0%d", schema.StatPosts), schema.StatPosts, value)
	}
	rows.Close()
	stmt.Close()

	sql = `
    SELECT count(group_id),
           user_id
    FROM
      (SELECT DISTINCT user_id,
                       group_id
       FROM
         (SELECT user_id,
                 group_id
          FROM threads t
          UNION SELECT user_id,
                       group_id
          FROM posts p) tp) tptp
    GROUP BY user_id
    `
	stmt, _ = mysql.Prepare(sql)
	rows, _ = stmt.Query()
	for rows.Next() {
		var value, userID int
		rows.Scan(&value, &userID)
		statStmt.Exec(userID, schema.OwnerUser, fmt.Sprintf("1970-01-0%d", schema.StatLives), schema.StatLives, value)
	}
	rows.Close()
	stmt.Close()

	sql = "SELECT count(*), user_id FROM attachments GROUP BY user_id"
	stmt, _ = mysql.Prepare(sql)
	rows, _ = stmt.Query()
	for rows.Next() {
		var value, userID int
		rows.Scan(&value, &userID)
		statStmt.Exec(userID, schema.OwnerUser, fmt.Sprintf("1970-01-0%d", schema.StatAttachments), schema.StatAttachments, value)
	}
	rows.Close()
	stmt.Close()
}
