package schema

import (
	"database/sql"
	"encoding/json"
)

type Base struct {
	Origin string `json:"root"`
	//ignore the tmpl attribute
}

type BaseDefine struct {
	Base Base `json:"base"`
}

type Avatar struct {
	X           int    `json:"x"`
	Y           int    `json:"y"`
	Size        int    `json:"size"`
	SourceFile  string `json:"source"`
	CroppedFile string `json:"crop"`
}

type User struct {
	Id         int    `json:"id"`
	ScreenName string `json:"name"`
	Avatar     Avatar `json:"icon"`
}

type UserBase struct {
	Id   int  `json:"id"`
	Base Base `json:"base"`
	User User `json:"user"`
}

type Group struct {
	Id          int        `json:"id"`
	Name        string     `json:"name"`
	Founder     UserBase   `json:"user"`
	Team        []UserBase `json:"admins"`
	Description string     `json:"intro"`
	Avatar      Avatar     `json:"icon"`
	Clubs       []string   `json:"clubs"`
}

type GroupBase struct {
	Id    int   `json:"id"`
	Base  Base  `json:"base"`
	Group Group `json:"group"`
	//just ignore the bbs attributes, will generate our own hot threads later
}

type Attachment struct {
	FileName string `json:"source"`
	Height   int    `json:"height"`
	Width    int    `json:"width"`
	Type     string `json:"type"`
}

type Message struct {
	Id int `json:"id"`
	//go is typed language while javascript is not.
	//we make it interface{} in case we fall into pitfalls like this:
	//{"msgid": "3"}
	//It may hurt like hell.
	ThreadId    interface{}  `json:"msgid"`
	Title       string       `json:"title"`
	Content     string       `json:"content"`
	Attachments []Attachment `json:"attachment2"`
	Author      UserBase     `json:"user"`
	Group       GroupBase    `json:"object"`
	CreatedAt   string       `json:"created"`
}

type MessageBase struct {
	Id      int     `json:"id"`
	Base    Base    `json:"base"`
	Message Message `json:"message"`
}

func (this *User) Save(db *sql.DB) {
	stmt, _ := db.Prepare("REPLACE INTO users(id, screen_name, avatar) VALUES(?, ?, ?)")
	defer stmt.Close()
	avatar, _ := json.Marshal(this.Avatar)
	stmt.Exec(this.Id, this.ScreenName, avatar)
}

func (this *Group) SaveTeam(db *sql.DB) {
	stmt, _ := db.Prepare("REPLACE INTO users_groups(user_id, group_id) VALUES(?, ?)")
	defer stmt.Close()
	for _, user := range this.Team {
		stmt.Exec(user.Id, this.Id)
	}
}

func (this *Group) SaveClub(db *sql.DB) {
	//first get club ID using club name
	queryClubStmt, _ := db.Prepare("SELECT id FROM clubs WHERE name=?")
	defer queryClubStmt.Close()
	insertStmt, _ := db.Prepare("REPLACE INTO clubs(name) VALUES(?)")
	defer insertStmt.Close()
	insertRelationStmt, _ := db.Prepare("REPLACE INTO clubs_groups(club_id, group_id) VALUES(?, ?)")
	defer insertRelationStmt.Close()
	for _, club := range this.Clubs {
		var clubId int64
		queryClubStmt.QueryRow(club).Scan(&clubId)
		if clubId == 0 {
			//insert club first
			result, _ := insertStmt.Exec(club)
			var err error
			clubId, err = result.LastInsertId()
			if err != nil {
				clubId = 0
			}
		}
		if clubId > 0 {
			insertRelationStmt.Exec(clubId, this.Id)
		}
	}
}

func (this *Group) Save(db *sql.DB) {
	stmt, _ := db.Prepare("REPLACE INTO groups(id, name, avatar, description, user_id) VALUES(?, ?, ?, ?, ?)")
	defer stmt.Close()
	avatar, _ := json.Marshal(this.Avatar)
	stmt.Exec(this.Id, this.Name, avatar, this.Description, this.Founder.Id)
	//save admins
	this.SaveTeam(db)
	//save club
	this.SaveClub(db)
}

func (this *Message) SaveAttachments(db *sql.DB) {
	var ownerType = 0
	if this.ThreadId != nil {
		ownerType = 2
	} else {
		ownerType = 1
	}
	stmt, _ := db.Prepare("REPLACE INTO attachments(owner_id, owner_type, file_name, width, height, file_type, user_id) VALUES(?, ?, ?, ?, ?, ?, ?)")
	defer stmt.Close()
	for _, attachment := range this.Attachments {
		stmt.Exec(this.Id, ownerType, attachment.FileName, attachment.Width, attachment.Height, attachment.Type, this.Author.Id)
	}
}

func (this *Message) Save(db *sql.DB) {
	if this.ThreadId != nil {
		//this is a post message
		stmt, _ := db.Prepare("REPLACE INTO posts(id, thread_id, content, user_id, group_id, created_at) VALUES(?, ?, ?, ?, ?, ?)")
		defer stmt.Close()
		stmt.Exec(this.Id, this.ThreadId, this.Content, this.Author.Id, this.Group.Id, this.CreatedAt)
	} else {
		stmt, _ := db.Prepare("REPLACE INTO threads(id, title, content, user_id, group_id, created_at) VALUES(?, ?, ?, ?, ?, ?)")
		defer stmt.Close()
		stmt.Exec(this.Id, this.Title, this.Content, this.Author.Id, this.Group.Id, this.CreatedAt)
	}
	//save User
	this.Author.User.Save(db)
	//save attachments
	this.SaveAttachments(db)
}
