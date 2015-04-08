package schema

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestParseGroup(t *testing.T) {
	groupJsonString, _ := ioutil.ReadFile("../../resources/samples/group.json")
	var group GroupBase
	json.Unmarshal(groupJsonString, &group)
	t.Log("Group ID is: ", group.Id)
	t.Log("Group Name is: ", group.Group.Name)
}

func TestParseThreadMessage(t *testing.T) {
	threadJsonString, _ := ioutil.ReadFile("../../resources/samples/thread.json")
	var thread MessageBase
	json.Unmarshal(threadJsonString, &thread)
	t.Log("Thread Id is: ", thread.Id)
	t.Log("Thread Title is: ", thread.Message.Title)
	t.Log("Thread Author is: ", thread.Message.Author.User.ScreenName)
}

func TestParsePostMessage(t *testing.T) {
	postJsonString, _ := ioutil.ReadFile("../../resources/samples/post.json")
	var post MessageBase
	json.Unmarshal(postJsonString, &post)
	t.Log("Post Id is: ", post.Id)
	t.Log("Post to Thread: ", post.Message.ThreadId)
	t.Log("Post Content is: ", post.Message.Content)
}
