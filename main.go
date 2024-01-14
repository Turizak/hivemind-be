package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type message struct {
	ID       string `json:"id"`
	Category string `json:"category"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Message  string `json:"message"`
	UUID     string `json:"uuid"`
	Created  string `json:"created"`
}

var messages = []message{
	{ID: "1", Category: "h/AITA", Title: "Am I the asshole for killing my cat?", Author: "rakazirut", Message: "I killed my cat because he was a bitch. Am I wrong?", UUID: uuid.NewString(), Created: time.Now().Format(time.RFC3339Nano)},
	{ID: "2", Category: "h/Fishing", Title: "I lost my sturgeon", Author: "tslanda", Message: "Jumped off the line. Damn!", UUID: uuid.NewString(), Created: time.Now().Format(time.RFC3339Nano)},
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

	messages = append(messages, newMsg)
	c.IndentedJSON(http.StatusCreated, newMsg)
}

func main() {
	router := gin.Default()
	router.GET("/messages", getMessages)
	router.GET("/message/:id", getMessageById)
	router.POST("/message", createMessage)
	router.Run("localhost:8080")
}
