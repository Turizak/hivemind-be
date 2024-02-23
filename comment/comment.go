package comment

import (
	"example/hivemind-be/account"
	"example/hivemind-be/content"
	"example/hivemind-be/db"
	"example/hivemind-be/hive"
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
	AccountUUID string      `json:"AccountUUID"`
	ContentUUID string      `json:"ContentUuid" gorm:"foreignKey:ContentUuid"` //foreign key gorm associations to content type table Uuid
	ParentUUID  string      `json:"ParentUuid" gorm:"default:null"`            //if comment is a reply, the ParentUUID will be the UUID of the parent comment
	Upvote      int32       `json:"Upvote"`
	Downvote    int32       `json:"Downvote"`
	Deleted     bool        `json:"Deleted"`
	Created     pq.NullTime `json:"Created"`
	LastEdited  pq.NullTime `json:"LastEdited"`
}

type CommentWithReplies struct {
	Parent  Comment   `json:"Comment"`
	Replies []Comment `json:"Replies"`
}

func CreateComment(c *gin.Context) {
	var newComment Comment
	var content content.Content
	var hive hive.Hive
	var account account.Account
	uid := c.Param("uuid")

	if err := c.BindJSON(&newComment); err != nil {
		return
	}

	if result := db.Db.Where("uuid = ?", uid).First(&content); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	if result := db.Db.Where("username = ?", newComment.Author).First(&account); result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"Error": "An error occurred. Please try again.",
		})
		return
	}

	newComment.UUID = uuid.NewString()
	newComment.ContentUUID = content.UUID
	newComment.AccountUUID = account.UUID
	newComment.Upvote = 0
	newComment.Downvote = 0
	newComment.Deleted = false
	newComment.Created = pq.NullTime{Time: time.Now(), Valid: true}
	newComment.LastEdited = pq.NullTime{Valid: false}

	if result := db.Db.Create(&newComment); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	if result := db.Db.Where("uuid = ?", content.HiveUUID).First(&hive); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	hive.TotalComments += 1
	content.CommentCount += 1
	db.Db.Save(&content)
	db.Db.Save(&hive)
	c.JSON(http.StatusCreated, newComment)
}

func CreateCommentReply(c *gin.Context) {
	var newComment Comment
	var parentComment Comment
	var content content.Content
	var hive hive.Hive
	var account account.Account
	uid := c.Param("uuid")
	pid := c.Param("parentuuid")

	if err := c.BindJSON(&newComment); err != nil {
		return
	}

	if result := db.Db.Where("uuid = ?", uid).First(&content); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	if result := db.Db.Where("uuid = ?", pid).First(&parentComment); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	if parentComment.ParentUUID != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Cannot reply to a reply. Please reply to the parent comment.",
		})
		return
	}

	if result := db.Db.Where("username = ?", newComment.Author).First(&account); result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"Error": "An error occurred. Please try again.",
		})
		return
	}

	newComment.UUID = uuid.NewString()
	newComment.ParentUUID = parentComment.UUID
	newComment.ContentUUID = content.UUID
	newComment.AccountUUID = account.UUID
	newComment.Upvote = 0
	newComment.Downvote = 0
	newComment.Deleted = false
	newComment.Created = pq.NullTime{Time: time.Now(), Valid: true}
	newComment.LastEdited = pq.NullTime{Valid: false}

	if result := db.Db.Create(&newComment); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	if result := db.Db.Where("uuid = ?", content.HiveUUID).First(&hive); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	hive.TotalComments += 1
	content.CommentCount += 1
	db.Db.Save(&content)
	db.Db.Save(&hive)
	c.JSON(http.StatusCreated, newComment)
}

func GetCommentsByContentUuid(c *gin.Context) {
	var comment []Comment
	uuid := c.Param("uuid")

	if result := db.Db.Where("content_uuid = ?", uuid).Find(&comment); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, comment)
}

func GetCommentByUuid(c *gin.Context) {
	var comment Comment
	uuid := c.Param("uuid")

	if result := db.Db.Where("uuid = ?", uuid).First(&comment); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, comment)
}

func GetCommentByUuidWithReplies(c *gin.Context) {
	var comment Comment
	var replies []Comment
	uuid := c.Param("uuid")

	if result := db.Db.Where("uuid = ?", uuid).First(&comment); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	if result := db.Db.Where("parent_uuid = ?", uuid).Find(&replies); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	commentWithReplies := CommentWithReplies{
		Parent:  comment,
		Replies: replies,
	}

	c.JSON(http.StatusOK, commentWithReplies)
}

func DeleteCommentByUuid(c *gin.Context) {
	var comment Comment
	var content content.Content
	var hive hive.Hive
	uuid := c.Param("uuid")

	if result := db.Db.Where("uuid = ?", uuid).First(&comment); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	if result := db.Db.Where("uuid = ?", comment.ContentUUID).First(&content); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	if comment.Deleted {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "comment has already been deleted!"})
		return
	}

	if result := db.Db.Where("uuid = ?", content.HiveUUID).First(&hive); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	hive.TotalComments -= 1
	comment.Deleted = true
	comment.LastEdited = pq.NullTime{Time: time.Now(), Valid: true}
	content.CommentCount -= 1
	db.Db.Save(&comment)
	db.Db.Save(&content)
	db.Db.Save(&hive)
	c.JSON(http.StatusOK, comment)
}

func UndeleteCommentByUuid(c *gin.Context) {
	var comment Comment
	var content content.Content
	var hive hive.Hive
	uuid := c.Param("uuid")

	if result := db.Db.Where("uuid = ?", uuid).First(&comment); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	if result := db.Db.Where("uuid = ?", comment.ContentUUID).First(&content); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	if !comment.Deleted {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "comment has not been deleted!"})
		return
	}

	if result := db.Db.Where("uuid = ?", content.HiveUUID).First(&hive); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	hive.TotalComments += 1
	comment.Deleted = false
	comment.LastEdited = pq.NullTime{Time: time.Now(), Valid: true}
	content.CommentCount += 1
	db.Db.Save(&comment)
	db.Db.Save(&content)
	db.Db.Save(&hive)
	c.JSON(http.StatusOK, comment)
}

func UpdateCommentByUuid(c *gin.Context) {
	var comment Comment
	var updateComment Comment

	if err := c.BindJSON(&updateComment); err != nil {
		return
	}

	uuid := c.Param("uuid")

	if result := db.Db.Where("uuid = ?", uuid).First(&comment); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	if val, ok := jsonDataHasKey(updateComment, "message"); ok {
		comment.Message = val
	}

	comment.LastEdited = pq.NullTime{Time: time.Now(), Valid: true}

	db.Db.Save(&comment)
	c.JSON(http.StatusOK, comment)
}

func AddCommentUpvoteByUuid(c *gin.Context) {
	var comment Comment
	uuid := c.Param("uuid")

	if result := db.Db.Where("uuid = ?", uuid).First(&comment); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	comment.Upvote += 1
	db.Db.Save(&comment)
	c.JSON(http.StatusOK, comment)
}

func RemoveCommentUpvoteByUuid(c *gin.Context) {
	var comment Comment
	uuid := c.Param("uuid")

	if result := db.Db.Where("uuid = ?", uuid).First(&comment); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	if comment.Upvote <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "no upvotes to remove!"})
		return
	}

	comment.Upvote -= 1
	db.Db.Save(&comment)
	c.JSON(http.StatusOK, comment)
}

func AddCommentDownvoteByUuid(c *gin.Context) {
	var comment Comment
	uuid := c.Param("uuid")

	if result := db.Db.Where("uuid = ?", uuid).First(&comment); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	comment.Downvote += 1
	db.Db.Save(&comment)
	c.JSON(http.StatusOK, comment)
}

func RemoveCommentDownvoteByUuid(c *gin.Context) {
	var comment Comment
	uuid := c.Param("uuid")

	if result := db.Db.Where("uuid = ?", uuid).First(&comment); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	if comment.Downvote <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "no downvotes to remove!"})
		return
	}

	comment.Downvote -= 1
	db.Db.Save(&comment)
	c.JSON(http.StatusOK, comment)
}

func jsonDataHasKey(data Comment, key string) (string, bool) {
	switch key {
	case "message":
		return data.Message, true
	default:
		return "null", false
	}
}
