package command

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/chekun/spore/spored/env"
	"github.com/chekun/spore/spored/schema"
	"github.com/mitchellh/cli"
	"github.com/polaris1119/times"
)

// StatCommand Stat Command
type StatCommand struct {
	UI *cli.BasicUi
}

// Help Stat Command Help
func (c *StatCommand) Help() string {
	helpText := `
Usage: spored serve [options] ...

  Spored Stat.

Options:

  -config=config.yml   Configuration file to use.
  -env=development     Environment.
  -start=2012-02-02    Start Date.
  -end=2015-05-05      End Date.
`

	return strings.TrimSpace(helpText)
}

// Synopsis Serve Command Synopsis
func (c *StatCommand) Synopsis() string {
	return "Spored Stat"
}

// Run Stat Command Implementation
func (c *StatCommand) Run(args []string) int {

	var startDate string
	var endDate string

	cmdFlags := flag.NewFlagSet("stat", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.UI.Output(c.Help()) }
	env.ConfigFlags(cmdFlags)

	cmdFlags.StringVar(&startDate, "start", "", "Start Date.")
	cmdFlags.StringVar(&endDate, "end", "", "End Date.")

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

	startTime := time.Now().Add(time.Hour * -24)
	endTime := time.Now().Add(time.Hour * -24)

	if startDate != "" {
		startTime = parseDate(startDate)
	}

	if endDate != "" {
		endTime = parseDate(endDate).Add(time.Hour * 24)
	}

	fixPostsWithNoGroup()

	for startTime.Before(endTime) {

		currentDate := times.Format("Y-m-d", startTime)

		fmt.Println("Running date: ", currentDate, " ...")

		//loop users
		statUsers(currentDate)
		//loop groups
		statGroups(currentDate)

		startTime = startTime.Add(time.Hour * 24)
	}

	return 0
}

func fixPostsWithNoGroup() {
	sql := "SELECT id, thread_id FROM posts WHERE group_id=0"
	stmt, _ := mysql.Prepare(sql)
	defer stmt.Close()

	rows, _ := stmt.Query()
	defer rows.Close()

	for rows.Next() {
		var postID, threadID, groupID int
		rows.Scan(&postID, &threadID)
		sql = "SELECT group_id FROM threads WHERE id = ?"
		stmt, _ = mysql.Prepare(sql)
		stmt.QueryRow(threadID).Scan(&groupID)
		sql = "UPDATE posts SET group_id = ? WHERE id = ?"
		stmt, _ = mysql.Prepare(sql)
		stmt.Exec(groupID, postID)
		stmt.Close()
	}
}

func statGroups(date string) {
	dateStart := date + " 00:00:00"
	dateEnd := date + " 23:59:59"

	sql := "SELECT id FROM groups ORDER BY id ASC"
	stmt, _ := mysql.Prepare(sql)
	defer stmt.Close()
	rows, _ := stmt.Query()
	defer rows.Close()

	for rows.Next() {
		var groupID int
		rows.Scan(&groupID)
		//count threads
		sql = "SELECT count(*) FROM threads WHERE group_id=? AND created_at >= ? AND created_at <= ?"
		stmt, _ = mysql.Prepare(sql)
		var threadsCount int
		stmt.QueryRow(groupID, dateStart, dateEnd).Scan(&threadsCount)
		//count posts
		sql = "SELECT count(*) FROM posts WHERE group_id=? AND created_at >= ? AND created_at <= ?"
		stmt, _ = mysql.Prepare(sql)
		var postsCount int
		stmt.QueryRow(groupID, dateStart, dateEnd).Scan(&postsCount)
		//count groups
		sql = `
		SELECT count(*) FROM (
			SELECT distinct user_id FROM (
				SELECT user_id FROM threads t WHERE group_id=? AND created_at >= ? AND created_at <= ?
				UNION
				SELECT user_id FROM posts p WHERE group_id=? AND created_at >= ? AND created_at <= ?
			) tp
		) tptp
		`
		var groupsCount int
		stmt, _ = mysql.Prepare(sql)
		stmt.QueryRow(groupID, dateStart, dateEnd, groupID, dateStart, dateEnd).Scan(&groupsCount)
		//@todo group attachments

		saveStat(date, groupID, schema.OwnerGroup, schema.StatThreadsPerday, threadsCount)
		saveStat(date, groupID, schema.OwnerGroup, schema.StatPostsPerday, postsCount)
		saveStat(date, groupID, schema.OwnerGroup, schema.StatLivesPerday, groupsCount)

		//Update Total
		sql = "SELECT stat_type, SUM(stat_value) AS stat_value FROM stats WHERE owner_type=? AND owner_id=?  AND stat_type IN (?, ?, ?) GROUP BY stat_type"
		stmt, _ = mysql.Prepare(sql)
		totalRows, _ := stmt.Query(schema.OwnerGroup, groupID, schema.StatThreadsPerday, schema.StatPostsPerday, schema.StatLivesPerday)
		defer totalRows.Close()

		for totalRows.Next() {
			var statType, statValue int
			totalRows.Scan(&statType, &statValue)
			saveStat(fmt.Sprintf("1970-01-0%d", schema.Stat(statType)+1), groupID, schema.OwnerGroup, schema.Stat(statType)+1, statValue)
		}
	}
}

func statUsers(date string) {

	dateStart := date + " 00:00:00"
	dateEnd := date + " 23:59:59"

	sql := "SELECT id FROM users ORDER BY id ASC"
	stmt, _ := mysql.Prepare(sql)
	defer stmt.Close()
	rows, _ := stmt.Query()
	defer rows.Close()

	for rows.Next() {
		var userID int
		rows.Scan(&userID)
		//count threads
		sql = "SELECT count(*) FROM threads WHERE user_id=? AND created_at >= ? AND created_at <= ?"
		stmt, _ = mysql.Prepare(sql)
		var threadsCount int
		stmt.QueryRow(userID, dateStart, dateEnd).Scan(&threadsCount)
		//count posts
		sql = "SELECT count(*) FROM posts WHERE user_id=? AND created_at >= ? AND created_at <= ?"
		stmt, _ = mysql.Prepare(sql)
		var postsCount int
		stmt.QueryRow(userID, dateStart, dateEnd).Scan(&postsCount)
		//count groups
		sql = `
		SELECT count(*) FROM (
			SELECT distinct group_id FROM (
				SELECT group_id FROM threads t WHERE user_id=? AND created_at >= ? AND created_at <= ?
				UNION
				SELECT group_id FROM posts p WHERE user_id=? AND created_at >= ? AND created_at <= ?
			) tp
		) tptp
		`
		var groupsCount int
		stmt, _ = mysql.Prepare(sql)
		stmt.QueryRow(userID, dateStart, dateEnd, userID, dateStart, dateEnd).Scan(&groupsCount)
		//count attachments
		sql = "SELECT count(*) FROM attachments WHERE user_id = ?"
		stmt, _ = mysql.Prepare(sql)
		var attachmentsCount int
		stmt.QueryRow(userID).Scan(&attachmentsCount)

		saveStat(date, userID, schema.OwnerUser, schema.StatThreadsPerday, threadsCount)
		saveStat(date, userID, schema.OwnerUser, schema.StatPostsPerday, postsCount)
		saveStat(date, userID, schema.OwnerUser, schema.StatLivesPerday, groupsCount)
		saveStat(date, userID, schema.OwnerUser, schema.StatAttachmentsPerday, attachmentsCount)

		//Update Total
		sql = "SELECT stat_type, SUM(stat_value) AS stat_value FROM stats WHERE owner_type=? AND owner_id=?  AND stat_type IN (?, ?, ?, ?) GROUP BY stat_type"
		stmt, _ = mysql.Prepare(sql)
		totalRows, _ := stmt.Query(schema.OwnerUser, userID, schema.StatThreadsPerday, schema.StatPostsPerday, schema.StatLivesPerday, schema.StatAttachmentsPerday)
		defer totalRows.Close()

		for totalRows.Next() {
			var statType, statValue int
			totalRows.Scan(&statType, &statValue)
			saveStat(fmt.Sprintf("1970-01-0%d", schema.Stat(statType)+1), userID, schema.OwnerUser, schema.Stat(statType)+1, statValue)
		}
	}
}

func saveStat(date string, ownerID int, ownerType schema.Stat, valueType schema.Stat, value int) {
	sql := "REPLACE INTO stats VALUES(?, ?, ?, ?, ?)"
	stmt, _ := mysql.Prepare(sql)
	stmt.Exec(ownerID, ownerType, date, valueType, value)
	stmt.Close()
}

func parseDate(dateString string) time.Time {
	segments := strings.Split(dateString, "-")
	year, _ := strconv.Atoi(segments[0])
	month, _ := strconv.Atoi(segments[1])
	day, _ := strconv.Atoi(segments[2])
	timeZone, _ := time.LoadLocation("PRC")
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, timeZone)

	return date
}
