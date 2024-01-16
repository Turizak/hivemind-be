package messages

import (
	"example/hivemind-be/db"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

//GORM uses the name of your type as the DB table to query. Here the type is Message so gorm will use the messages table by default.
type Message struct {
	ID           int32       `json:"Id" gorm:"primaryKey:type:int32"`
	Category     string      `json:"Category"`
	Title        string      `json:"Title"`
	Author       string      `json:"Author"`
	Message      string      `json:"Message"`
	UUID         string      `json:"Uuid"`
	Created      pq.NullTime `json:"Created"`
	LastEdited   pq.NullTime `json:"LastEdited"`
	Type         string      `json:"Type"`
	Upvote       int32       `json:"Upvote"`
	Downvote     int32       `json:"Downvote"`
	CommentCount int32       `json:"CommentCount"`
	Deleted      bool        `json:"Deleted"`
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
	newMsg.Created = pq.NullTime{Time: time.Now(), Valid: true}
	newMsg.Type = "message"
	newMsg.Upvote = 0
	newMsg.Downvote = 0
	newMsg.CommentCount = 0
	newMsg.Deleted = false
	newMsg.LastEdited = pq.NullTime{Valid: false}

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
	messages.LastEdited = pq.NullTime{Time: time.Now(), Valid: true}
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
	messages.LastEdited = pq.NullTime{Time: time.Now(), Valid: true}
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

	messages.LastEdited = pq.NullTime{Time: time.Now(), Valid: true}

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
