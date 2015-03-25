package schema

type Base struct {
	Origin string `json:"root"`
	//ignore the tmpl attribute
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
	//we make it interface{} in case we fall into pits like this:
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
