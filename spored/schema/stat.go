package schema

//Stat Stat Type
type Stat int

const (
	//StatThreadsPerday Threads Perday
	StatThreadsPerday Stat = iota + 1
	//StatThreads Total Threads
	StatThreads
	//StatPostsPerday Posts Perday
	StatPostsPerday
	//StatPosts Total Posts
	StatPosts
	//StatLivesPerday Stat Lives Perday
	StatLivesPerday
	//StatLives Total Lives
	StatLives
	//StatAttachmentsPerday Attachments Perday
	StatAttachmentsPerday
	//StatAttachments Total Attachments
	StatAttachments
)

const (
	//OwnerUser Owner User Type
	OwnerUser Stat = iota + 1
	//OwnerThread Owner Thread Type
	OwnerThread
)
