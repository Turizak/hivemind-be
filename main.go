package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type message struct {
	ID           string `json:"id"`
	Category     string `json:"category"`
	Title        string `json:"title"`
	Author       string `json:"author"`
	Message      string `json:"message"`
	UUID         string `json:"uuid"`
	Created      string `json:"created"`
	Type         string `json:"type"`
	Upvote       int32  `json:"upvote"`
	Downvote     int32  `json:"downvote"`
	CommentCount int32  `json:"commentCount"`
	Deleted      bool   `json:"deleted"`
}

var messages = []message{
	{ID: "1", Category: "h/AITA", Title: "Am I the asshole for killing my cat?", Author: "rakazirut", Message: "I killed my cat because he was a bitch. Am I wrong?", UUID: uuid.NewString(), Created: time.Now().Format(time.RFC3339Nano), Type: "message", Upvote: 30, Downvote: 1, CommentCount: 50, Deleted: true},
	{ID: "2", Category: "h/Fishing", Title: "I lost my sturgeon", Author: "tslanda", Message: "Jumped off the line. Damn!", UUID: uuid.NewString(), Created: time.Now().Format(time.RFC3339Nano), Type: "message", Upvote: 300, Downvote: 14, CommentCount: 450, Deleted: false},
}

func getMessages(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, messages)
}

func helperGetMessageById(id string) (*message, error) {
	for i, m := range messages {
		if m.ID == id {
			return &messages[i], nil
		}
	}

	return nil, errors.New("message not found")
}

func getMessageById(c *gin.Context) {
	id := c.Param("id")
	message, err := helperGetMessageById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Message not found!"})
		return
	}

	c.IndentedJSON(http.StatusOK, message)
}

func createMessage(c *gin.Context) {
	var newMsg message

	if err := c.BindJSON(&newMsg); err != nil {
		return
	}

	newMsg.UUID = uuid.NewString()
	newMsg.Created = time.Now().Format(time.RFC3339Nano)
	newMsg.Type = "message"
	newMsg.Upvote = 0
	newMsg.Downvote = 0
	newMsg.CommentCount = 0
	newMsg.Deleted = false

	messages = append(messages, newMsg)
	c.IndentedJSON(http.StatusCreated, newMsg)
}

func addMessageUpvote(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter"})
	}

	message, err := helperGetMessageById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Message not found!"})
		return
	}

	message.Upvote += 1
	c.IndentedJSON(http.StatusOK, message)
}

func removeMessageUpvote(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter"})
	}

	message, err := helperGetMessageById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Message not found!"})
		return
	}

	if message.Downvote <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "No upvotes to remove!"})
		return
	}

	message.Upvote -= 1
	c.IndentedJSON(http.StatusOK, message)
}

func addMessageDownvote(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter"})
	}

	message, err := helperGetMessageById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Message not found!"})
		return
	}

	message.Downvote += 1
	c.IndentedJSON(http.StatusOK, message)
}

func removeMessageDownvote(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter"})
	}

	message, err := helperGetMessageById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Message not found!"})
		return
	}

	if message.Downvote <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "No downvotes to remove!"})
		return
	}

	message.Downvote -= 1
	c.IndentedJSON(http.StatusOK, message)
}

func deleteMessage(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter"})
	}

	message, err := helperGetMessageById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Message not found!"})
		return
	}

	if message.Deleted {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Message has already beem deleted!"})
		return
	}

	message.Deleted = true
	c.IndentedJSON(http.StatusOK, message)
}

func updateMessage(c *gin.Context) {
	var updateMsg message

	if err := c.BindJSON(&updateMsg); err != nil {
		return
	}

	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter"})
	}

	message, err := helperGetMessageById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Message not found!"})
		return
	}

	if val, ok := jsonDataHasKey(updateMsg, "category"); ok {
		message.Category = val
	}
	if val, ok := jsonDataHasKey(updateMsg, "title"); ok {
		message.Title = val
	}
	if val, ok := jsonDataHasKey(updateMsg, "message"); ok {
		message.Message = message.Message + "\n\nEdit: " + val
	}

	c.IndentedJSON(http.StatusOK, message)
}

func jsonDataHasKey(data message, key string) (string, bool) {
	switch key {
	case "category":
		return data.Category, true
	case "title":
		return data.Title, true
	case "message":
		return data.Message, true
	default:
		return "null", false
	}
}

func main() {
	router := gin.Default()
	router.GET("/messages", getMessages)
	router.GET("/message/:id", getMessageById)
	router.POST("/message", createMessage)
	router.PATCH("message/add-upvote", addMessageUpvote)
	router.PATCH("message/remove-upvote", removeMessageUpvote)
	router.PATCH("message/add-downvote", addMessageDownvote)
	router.PATCH("message/remove-downvote", removeMessageDownvote)
	router.PATCH("message/delete", deleteMessage)
	router.PATCH("message/update", updateMessage)
	router.Run("localhost:8080")
}
