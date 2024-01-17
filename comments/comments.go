package comments

import (
	"example/hivemind-be/content"
	"example/hivemind-be/db"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Comment struct {
	ID          int32       `json:"Id" gorm:"primaryKey:type:int32"`
	Author      string      `json:"Author"`
	Message     string      `json:"Message"`
	UUID        string      `json:"Uuid"`
	ContentUUID string      `json:"ContentUuid" gorm:"foreignKey:ContentUuid"` //foreign key gorm associations to content type table Uuid
	Upvote      int32       `json:"Upvote"`
	Downvote    int32       `json:"Downvote"`
	Deleted     bool        `json:"Deleted"`
	Created     pq.NullTime `json:"Created"`
	LastEdited  pq.NullTime `json:"LastEdited"`
}

func CreateComment(c *gin.Context) {
	var newComment Comment
	var content content.Content
	uid := c.Param("uuid")

	if err := c.BindJSON(&newComment); err != nil {
		return
	}

	if result := db.Db.Where("uuid = ?", uid).First(&content); result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	newComment.UUID = uuid.NewString()
	newComment.ContentUUID = content.UUID
	newComment.Upvote = 0
	newComment.Downvote = 0
	newComment.Deleted = false
	newComment.Created = pq.NullTime{Time: time.Now(), Valid: true}
	newComment.LastEdited = pq.NullTime{Valid: false}

	if result := db.Db.Create(&newComment); result.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	content.CommentCount += 1
	db.Db.Save(&content)
	c.IndentedJSON(http.StatusCreated, newComment)
}

func GetCommentsByContentUuid(c *gin.Context) {
	var comment []Comment
	uuid := c.Param("uuid")

	if result := db.Db.Where("content_uuid = ?", uuid).Find(&comment); result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, comment)
}

func GetCommentByUuid(c *gin.Context) {
	var comment Comment
	uuid := c.Param("uuid")

	if result := db.Db.Where("uuid = ?", uuid).First(&comment); result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, comment)
}

func DeleteCommentByUuid(c *gin.Context) {
	var comment Comment
	uuid := c.Param("uuid")

	if result := db.Db.Where("uuid = ?", uuid).First(&comment); result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	if comment.Deleted {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"Error": "comment has already been deleted!"})
		return
	}

	comment.Deleted = true
	comment.LastEdited = pq.NullTime{Time: time.Now(), Valid: true}
	db.Db.Save(&comment)
	c.IndentedJSON(http.StatusOK, comment)
}

func UndeleteCommentByUuid(c *gin.Context) {
	var comment Comment
	uuid := c.Param("uuid")

	if result := db.Db.Where("uuid = ?", uuid).First(&comment); result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	if !comment.Deleted {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"Error": "comment has not been deleted!"})
		return
	}

	comment.Deleted = false
	comment.LastEdited = pq.NullTime{Time: time.Now(), Valid: true}
	db.Db.Save(&comment)
	c.IndentedJSON(http.StatusOK, comment)
}

func UpdateCommentByUuid(c *gin.Context) {
	var comment Comment
	var updateComment Comment

	if err := c.BindJSON(&updateComment); err != nil {
		return
	}

	uuid := c.Param("uuid")

	if result := db.Db.Where("uuid = ?", uuid).First(&comment); result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	if val, ok := jsonDataHasKey(updateComment, "message"); ok {
		comment.Message = val
	}

	comment.LastEdited = pq.NullTime{Time: time.Now(), Valid: true}

	db.Db.Save(&comment)
	c.IndentedJSON(http.StatusOK, comment)
}

func jsonDataHasKey(data Comment, key string) (string, bool) {
	switch key {
	case "message":
		return data.Message, true
	default:
		return "null", false
	}
}
