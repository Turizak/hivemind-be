package messages

import (
	"example/hivemind-be/db"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

//GORM uses the name of your type as the DB table to query. Here the type is Message so gorm will use the messages table by default.
type Message struct {
	ID           int32  `json:"id" gorm:"primaryKey:type:int32"`
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

func GetMessages(c *gin.Context) {
	var messages []Message
	if result := db.Db.Order("id asc").Find(&messages); result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": result.Error.Error(),
		})
		return
	}
	c.IndentedJSON(http.StatusOK, &messages)
}

func GetMessageById(c *gin.Context) {
	var messages Message
	id := c.Param("id")
	if result := db.Db.First(&messages, id); result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": result.Error.Error(),
		})
		return
	}
	c.IndentedJSON(http.StatusOK, messages)
}

func CreateMessage(c *gin.Context) {
	var newMsg Message

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

	if result := db.Db.Create(&newMsg); result.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusCreated, newMsg)
}

func AddMessageUpvote(c *gin.Context) {
	var messages Message
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "missing id query parameter"})
	}

	if result := db.Db.First(&messages, id); result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	messages.Upvote += 1
	db.Db.Save(&messages)
	c.IndentedJSON(http.StatusOK, messages)
}

func RemoveMessageUpvote(c *gin.Context) {
	var messages Message
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "missing id query parameter"})
	}

	if result := db.Db.First(&messages, id); result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	if messages.Upvote <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "no upvotes to remove!"})
		return
	}

	messages.Upvote -= 1
	db.Db.Save(&messages)
	c.IndentedJSON(http.StatusOK, messages)
}

func AddMessageDownvote(c *gin.Context) {
	var messages Message
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "missing id query parameter"})
	}

	if result := db.Db.First(&messages, id); result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	messages.Downvote += 1
	db.Db.Save(&messages)
	c.IndentedJSON(http.StatusOK, messages)
}

func RemoveMessageDownvote(c *gin.Context) {
	var messages Message
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "missing id query parameter"})
	}

	if result := db.Db.First(&messages, id); result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	if messages.Downvote <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "no downvotes to remove!"})
		return
	}

	messages.Downvote -= 1
	db.Db.Save(&messages)
	c.IndentedJSON(http.StatusOK, messages)
}

func DeleteMessage(c *gin.Context) {
	var messages Message
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "missing id query parameter"})
	}

	if result := db.Db.First(&messages, id); result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	if messages.Deleted {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "message has already been deleted!"})
		return
	}

	messages.Deleted = true
	db.Db.Save(&messages)
	c.IndentedJSON(http.StatusOK, messages)
}

func UndeleteMessage(c *gin.Context) {
	var messages Message
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "missing id query parameter"})
	}

	if result := db.Db.First(&messages, id); result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	if !messages.Deleted {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "message has not been deleted!"})
		return
	}

	messages.Deleted = false
	db.Db.Save(&messages)
	c.IndentedJSON(http.StatusOK, messages)
}

func UpdateMessage(c *gin.Context) {
	var messages Message
	var updateMsg Message

	if err := c.BindJSON(&updateMsg); err != nil {
		return
	}

	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "missing id query parameter"})
	}

	if result := db.Db.First(&messages, id); result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	if val, ok := jsonDataHasKey(updateMsg, "category"); ok {
		messages.Category = val
	}
	if val, ok := jsonDataHasKey(updateMsg, "title"); ok {
		messages.Title = val
	}
	if val, ok := jsonDataHasKey(updateMsg, "message"); ok {
		messages.Message = val
	}

	db.Db.Save(&messages)
	c.IndentedJSON(http.StatusOK, messages)
}

func jsonDataHasKey(data Message, key string) (string, bool) {
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
