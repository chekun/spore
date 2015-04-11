package schema

import (
	"database/sql"
	"encoding/json"

	"github.com/go-sql-driver/mysql"
	"github.com/polaris1119/times"
)

//Base Object Define
type Base struct {
	Origin string `json:"root"`
	//ignore the tmpl attribute
}

//BaseDefine Object Define
type BaseDefine struct {
	Base Base `json:"base"`
}

//Avatar Object Define
type Avatar struct {
	X           int    `json:"x"`
	Y           int    `json:"y"`
	Size        int    `json:"size"`
	SourceFile  string `json:"source"`
	CroppedFile string `json:"crop"`
}

//User Object Define
type User struct {
	ID         int    `json:"id"`
	ScreenName string `json:"name"`
	Avatar     Avatar `json:"icon"`
}

//UserBase Object Define
type UserBase struct {
	ID   int  `json:"id"`
	Base Base `json:"base"`
	User User `json:"user"`
}

//Group Object Define
type Group struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Founder     UserBase   `json:"user"`
	Team        []UserBase `json:"admins"`
	Description string     `json:"intro"`
	Avatar      Avatar     `json:"icon"`
	Clubs       []string   `json:"clubs"`
}

//GroupBase Object Define
type GroupBase struct {
	ID    int   `json:"id"`
	Base  Base  `json:"base"`
	Group Group `json:"group"`
	//just ignore the bbs attributes, will generate our own hot threads later
}

//Attachment Object Define
type Attachment struct {
	FileName string `json:"source"`
	Height   int    `json:"height"`
	Width    int    `json:"width"`
	Type     string `json:"type"`
}

//Message Object Define
type Message struct {
	ID int `json:"id"`
	//go is typed language while javascript is not.
	//we make it interface{} in case we fall into pitfalls like this:
	//{"msgid": "3"}
	//It may hurt like hell.
	ThreadID    interface{}  `json:"msgid"`
	Title       string       `json:"title"`
	Content     string       `json:"content"`
	Attachments []Attachment `json:"attachment2"`
	Author      UserBase     `json:"user"`
	Group       GroupBase    `json:"object"`
	CreatedAt   string       `json:"created"`
}

//MessageBase Object Define
type MessageBase struct {
	ID      int     `json:"id"`
	Base    Base    `json:"base"`
	Message Message `json:"message"`
}

//Save save user to database
func (user *User) Save(db *sql.DB) {
	stmt, _ := db.Prepare("REPLACE INTO users(id, screen_name, avatar) VALUES(?, ?, ?)")
	defer stmt.Close()
	avatar, _ := json.Marshal(user.Avatar)
	stmt.Exec(user.ID, user.ScreenName, avatar)
}

//New make new user object from database
func (user *User) New(db *sql.DB, id int) {
	stmt, _ := db.Prepare("SELECT * FROM users WHERE id = ?")
	defer stmt.Close()
	var avatar string
	stmt.QueryRow(id).Scan(&user.ID, &user.ScreenName, &avatar)
	json.Unmarshal([]byte(avatar), &user.Avatar)
}

//NewUserBase make new UserBase object from database
func NewUserBase(db *sql.DB, id int) UserBase {
	userBase := UserBase{}
	userBase.User.New(db, id)
	userBase.ID = userBase.User.ID
	userBase.Base = Base{"user"}
	return userBase
}

//SaveTeam save management teams to database
func (group *Group) SaveTeam(db *sql.DB) {
	stmt, _ := db.Prepare("REPLACE INTO users_groups(user_id, group_id) VALUES(?, ?)")
	defer stmt.Close()
	for _, user := range group.Team {
		stmt.Exec(user.ID, group.ID)
	}
}

//SaveClub save club to database
func (group *Group) SaveClub(db *sql.DB) {
	//first get club ID using club name
	queryClubStmt, _ := db.Prepare("SELECT id FROM clubs WHERE name=?")
	defer queryClubStmt.Close()
	insertStmt, _ := db.Prepare("REPLACE INTO clubs(name) VALUES(?)")
	defer insertStmt.Close()
	insertRelationStmt, _ := db.Prepare("REPLACE INTO clubs_groups(club_id, group_id) VALUES(?, ?)")
	defer insertRelationStmt.Close()
	for _, club := range group.Clubs {
		var clubID int64
		queryClubStmt.QueryRow(club).Scan(&clubID)
		if clubID == 0 {
			//insert club first
			result, _ := insertStmt.Exec(club)
			var err error
			clubID, err = result.LastInsertId()
			if err != nil {
				clubID = 0
			}
		}
		if clubID > 0 {
			insertRelationStmt.Exec(clubID, group.ID)
		}
	}
}

//NewClub make new club object from database
func NewClub(db *sql.DB, id int) string {
	stmt, _ := db.Prepare("SELECT name FROM clubs WHERE id = ?")
	defer stmt.Close()
	var club string
	stmt.QueryRow(id).Scan(&club)
	return club
}

// Save save group to database
func (group *Group) Save(db *sql.DB) {
	stmt, _ := db.Prepare("REPLACE INTO groups(id, name, avatar, description, user_id) VALUES(?, ?, ?, ?, ?)")
	defer stmt.Close()
	avatar, _ := json.Marshal(group.Avatar)
	stmt.Exec(group.ID, group.Name, avatar, group.Description, group.Founder.ID)
	//save admins
	group.SaveTeam(db)
	//save club
	group.SaveClub(db)
}

//New Make new group object from database
func (group *Group) New(db *sql.DB, id int) {
	stmt, _ := db.Prepare("SELECT * FROM groups WHERE id = ?")
	defer stmt.Close()
	var avatar string
	var ownerID int
	stmt.QueryRow(id).Scan(&group.ID, &group.Name, &avatar, &group.Description, &ownerID)
	json.Unmarshal([]byte(avatar), &group.Avatar)
	group.Founder = NewUserBase(db, ownerID)
	//fetch admins
	stmt, _ = db.Prepare("SELECT user_id FROM users_groups WHERE group_id = ?")
	rows, _ := stmt.Query(id)
	defer rows.Close()
	teams := []UserBase{}
	for rows.Next() {
		var adminID int
		rows.Scan(&adminID)
		teams = append(teams, NewUserBase(db, adminID))
	}
	group.Team = teams
	//fetch clubs
	stmt, _ = db.Prepare("SELECT club_id FROM clubs_groups WHERE group_id = ?")
	rows, _ = stmt.Query(id)
	clubs := []string{}
	for rows.Next() {
		var clubID int
		rows.Scan(&clubID)
		clubs = append(clubs, NewClub(db, clubID))
	}
	group.Clubs = clubs
}

//NewGroupBase make new NewGroupBase object from database
func NewGroupBase(db *sql.DB, id int) GroupBase {
	groupBase := GroupBase{}
	groupBase.Group.New(db, id)
	groupBase.ID = groupBase.Group.ID
	groupBase.Base = Base{"group"}
	return groupBase
}

// SaveAttachments Save Attachments to database
func (message *Message) SaveAttachments(db *sql.DB) {
	var ownerType = 0
	if message.ThreadID != nil {
		ownerType = 2
	} else {
		ownerType = 1
	}
	stmt, _ := db.Prepare("REPLACE INTO attachments(owner_id, owner_type, file_name, width, height, file_type, user_id) VALUES(?, ?, ?, ?, ?, ?, ?)")
	defer stmt.Close()
	for _, attachment := range message.Attachments {
		stmt.Exec(message.ID, ownerType, attachment.FileName, attachment.Width, attachment.Height, attachment.Type, message.Author.ID)
	}
}

// Save save message(thread, post) to database
func (message *Message) Save(db *sql.DB) {
	if message.ThreadID != nil {
		//this is a post message
		stmt, _ := db.Prepare("REPLACE INTO posts(id, thread_id, content, user_id, group_id, created_at) VALUES(?, ?, ?, ?, ?, ?)")
		defer stmt.Close()
		//Warning!
		//Some posts have no Group Object, eg. 10255 and 11903
		//I will fix the data before stats
		stmt.Exec(message.ID, message.ThreadID, message.Content, message.Author.ID, message.Group.ID, message.CreatedAt)
	} else {
		stmt, _ := db.Prepare("REPLACE INTO threads(id, title, content, user_id, group_id, created_at) VALUES(?, ?, ?, ?, ?, ?)")
		defer stmt.Close()
		stmt.Exec(message.ID, message.Title, message.Content, message.Author.ID, message.Group.ID, message.CreatedAt)
	}
	//save User
	message.Author.User.Save(db)
	//save attachments
	message.SaveAttachments(db)
}

//New make new Thread Object from database
func (message *Message) New(db *sql.DB, id int) {
	stmt, _ := db.Prepare("SELECT * FROM threads WHERE id = ?")
	defer stmt.Close()
	var userID int
	var groupID int
	var nt mysql.NullTime
	stmt.QueryRow(id).Scan(&message.ID, &message.Title, &message.Content, &userID, &groupID, &nt)
	message.CreatedAt = times.Format("Y-m-d H:i:s", nt.Time)
	message.Author = NewUserBase(db, userID)
	message.Group = NewGroupBase(db, groupID)
}
