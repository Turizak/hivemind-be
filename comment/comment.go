package comment

import (
	"example/hivemind-be/content"
	"example/hivemind-be/db"
	"example/hivemind-be/hive"
	"example/hivemind-be/token"
	"example/hivemind-be/utils"
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

type CommentVote struct {
	ID          int32
	AccountUUID string
	CommentUUID string
	Upvote      bool
	Downvote    bool
	LastEdited  pq.NullTime
}

func CreateComment(c *gin.Context) {
	var newComment Comment
	var content content.Content
	var hive hive.Hive

	authToken := c.GetHeader("Authorization")
	validToken := token.CheckToken(c, authToken)

	if !validToken {
		return
	}
	claims, err := token.ParseToken(authToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized.",
		})
		return
	}

	uid := c.Param("uuid")

	if err := c.BindJSON(&newComment); err != nil {
		return
	}

	if !validateCommentMessage(newComment.Message) {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Message must be between 1 and 2048 characters.",
		})
		return
	}

	if result := db.Db.Where("uuid = ?", uid).First(&content); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	newComment.Author = claims.Username
	newComment.UUID = uuid.NewString()
	newComment.ContentUUID = content.UUID
	newComment.AccountUUID = claims.AccountUUID
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

	authToken := c.GetHeader("Authorization")
	validToken := token.CheckToken(c, authToken)

	if !validToken {
		return
	}
	claims, err := token.ParseToken(authToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized.",
		})
		return
	}

	uid := c.Param("uuid")
	pid := c.Param("parentuuid")

	if err := c.BindJSON(&newComment); err != nil {
		return
	}

	if !validateCommentMessage(newComment.Message) {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Message must be between 1 and 2048 characters.",
		})
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

	newComment.Author = claims.Username
	newComment.UUID = uuid.NewString()
	newComment.ParentUUID = parentComment.UUID
	newComment.ContentUUID = content.UUID
	newComment.AccountUUID = claims.AccountUUID
	newComment.Upvote = 0
	newComment.Downvote = 0
	newComment.Deleted = false
	newComment.Created = pq.NullTime{Time: time.Now(), Valid: true}
	newComment.LastEdited = pq.NullTime{Valid: false}

	if result := db.Db.Create(&newComment); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": result.Error.Error(),
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

	authToken := c.GetHeader("Authorization")

	validToken := token.CheckToken(c, authToken)

	if !validToken {
		return
	}

	uuid := c.Param("uuid")

	if result := db.Db.Where("content_uuid = ?", uuid).Order("created DESC").Find(&comment); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	for i := 0; i < len(comment); i++ {
		if comment[i].Deleted {
			comment[i].Message = "This comment has been deleted."
		}

	}

	c.JSON(http.StatusOK, comment)
}

func GetCommentByUuid(c *gin.Context) {
	var comment Comment

	authToken := c.GetHeader("Authorization")
	validToken := token.CheckToken(c, authToken)

	if !validToken {
		return
	}

	uuid := c.Param("uuid")

	if result := db.Db.Where("uuid = ?", uuid).First(&comment); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	if comment.Deleted {
		comment.Message = "This comment has been deleted."
	}

	c.JSON(http.StatusOK, comment)
}

func GetCommentByUuidWithReplies(c *gin.Context) {
	var comment Comment
	var replies []Comment

	authToken := c.GetHeader("Authorization")
	validToken := token.CheckToken(c, authToken)

	if !validToken {
		return
	}

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

	if comment.Deleted {
		comment.Message = "This comment has been deleted."
	}

	for i := 0; i < len(replies); i++ {
		if replies[i].Deleted {
			replies[i].Message = "This comment has been deleted."
		}

	}

	c.JSON(http.StatusOK, commentWithReplies)
}

func DeleteCommentByUuid(c *gin.Context) {
	var comment Comment
	var content content.Content
	var hive hive.Hive

	authToken := c.GetHeader("Authorization")
	validToken := token.CheckToken(c, authToken)

	if !validToken {
		return
	}

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

	authToken := c.GetHeader("Authorization")
	validToken := token.CheckToken(c, authToken)

	if !validToken {
		return
	}

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

	authToken := c.GetHeader("Authorization")
	validToken := token.CheckToken(c, authToken)

	if !validToken {
		return
	}

	if err := c.BindJSON(&updateComment); err != nil {
		return
	}

	if !validateCommentMessage(updateComment.Message) {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Message must be between 1 and 2048 characters.",
		})
		return
	}

	uuid := c.Param("uuid")

	if result := db.Db.Where("uuid = ?", uuid).First(&comment); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	if val, ok := utils.JsonDataHasKey(updateComment, "Message"); ok {
		comment.Message, _ = val.(string)
	}

	comment.LastEdited = pq.NullTime{Time: time.Now(), Valid: true}

	db.Db.Save(&comment)
	c.JSON(http.StatusOK, comment)
}

func AddCommentUpvoteByUuid(c *gin.Context) {
	var comment Comment
	var commentVote CommentVote

	authToken := c.GetHeader("Authorization")
	validToken := token.CheckToken(c, authToken)

	if !validToken {
		return
	}
	claims, err := token.ParseToken(authToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized.",
		})
		return
	}

	uuid := c.Param("uuid")

	//check comment exist
	if result := db.Db.Where("uuid = ?", uuid).First(&comment); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	voteQuery := map[string]interface{}{
		"account_uuid": claims.AccountUUID,
		"comment_uuid": uuid,
	}

	//check account vote
	if result := db.Db.Where(voteQuery).First(&commentVote); result.Error != nil {
		//user has no record
		comment.Upvote += 1
		commentVote.AccountUUID = claims.AccountUUID
		commentVote.CommentUUID = uuid
		commentVote.Upvote = true
		commentVote.Downvote = false
		commentVote.LastEdited = pq.NullTime{Time: time.Now(), Valid: true}
		db.Db.Save(&comment)
		db.Db.Save(&commentVote)
		c.JSON(http.StatusOK, gin.H{
			"Message": "User successfully upvoted!",
		})
		return
	}

	//error if user has already voted
	if commentVote.Upvote || commentVote.Downvote {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "User has already voted on this comment!",
		})
		return
	}

	//user has false for both upvote and downvote
	if !commentVote.Upvote && !commentVote.Downvote {
		comment.Upvote += 1
		commentVote.Upvote = true
		commentVote.LastEdited = pq.NullTime{Time: time.Now(), Valid: true}
		db.Db.Save(&comment)
		db.Db.Save(&commentVote)
		c.JSON(http.StatusOK, gin.H{
			"Message": "User successfully upvoted!",
		})
		return
	}
}

func RemoveCommentUpvoteByUuid(c *gin.Context) {
	var comment Comment
	var commentVote CommentVote

	authToken := c.GetHeader("Authorization")
	validToken := token.CheckToken(c, authToken)

	if !validToken {
		return
	}
	claims, err := token.ParseToken(authToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized.",
		})
		return
	}

	uuid := c.Param("uuid")

	//check comment exist
	if result := db.Db.Where("uuid = ?", uuid).First(&comment); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	voteQuery := map[string]interface{}{
		"account_uuid": claims.AccountUUID,
		"comment_uuid": uuid,
	}

	//check account vote
	if result := db.Db.Where(voteQuery).First(&commentVote); result.Error != nil {
		//user has not voted at all
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "User has not voted on this content!",
		})
		return
	}

	//error if user has not already upvoted or has downvoted
	if !commentVote.Upvote || commentVote.Downvote {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "User has not upvoted on this comment!",
		})
		return
	}

	comment.Upvote -= 1
	commentVote.Upvote = false
	commentVote.LastEdited = pq.NullTime{Time: time.Now(), Valid: true}
	db.Db.Save(&comment)
	db.Db.Save(&commentVote)
	c.JSON(http.StatusOK, gin.H{
		"Message": "User upvote removed sucessfully!",
	})
}

func AddCommentDownvoteByUuid(c *gin.Context) {
	var comment Comment
	var commentVote CommentVote

	authToken := c.GetHeader("Authorization")
	validToken := token.CheckToken(c, authToken)

	if !validToken {
		return
	}
	claims, err := token.ParseToken(authToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized.",
		})
		return
	}

	uuid := c.Param("uuid")

	//check comment exist
	if result := db.Db.Where("uuid = ?", uuid).First(&comment); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	voteQuery := map[string]interface{}{
		"account_uuid": claims.AccountUUID,
		"comment_uuid": uuid,
	}

	//check account vote
	if result := db.Db.Where(voteQuery).First(&commentVote); result.Error != nil {
		//user has no record
		comment.Downvote += 1
		commentVote.AccountUUID = claims.AccountUUID
		commentVote.CommentUUID = uuid
		commentVote.Upvote = false
		commentVote.Downvote = true
		commentVote.LastEdited = pq.NullTime{Time: time.Now(), Valid: true}
		db.Db.Save(&comment)
		db.Db.Save(&commentVote)
		c.JSON(http.StatusOK, gin.H{
			"Message": "User successfully downvoted!",
		})
		return
	}

	//error if user has already voted
	if commentVote.Upvote || commentVote.Downvote {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "User has already voted on this comment!",
		})
		return
	}

	//user has false for both upvote and downvote
	if !commentVote.Upvote && !commentVote.Downvote {
		comment.Downvote += 1
		commentVote.Downvote = true
		commentVote.LastEdited = pq.NullTime{Time: time.Now(), Valid: true}
		db.Db.Save(&comment)
		db.Db.Save(&commentVote)
		c.JSON(http.StatusOK, gin.H{
			"Message": "User successfully downvoted!",
		})
		return
	}
}

func RemoveCommentDownvoteByUuid(c *gin.Context) {
	var comment Comment
	var commentVote CommentVote

	authToken := c.GetHeader("Authorization")
	validToken := token.CheckToken(c, authToken)

	if !validToken {
		return
	}
	claims, err := token.ParseToken(authToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized.",
		})
		return
	}

	uuid := c.Param("uuid")

	//check comment exist
	if result := db.Db.Where("uuid = ?", uuid).First(&comment); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	voteQuery := map[string]interface{}{
		"account_uuid": claims.AccountUUID,
		"comment_uuid": uuid,
	}

	//check account vote
	if result := db.Db.Where(voteQuery).First(&commentVote); result.Error != nil {
		//user has not voted at all
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "User has not voted on this comment!",
		})
		return
	}

	//error if user has not already downvoted or has upvoted
	if commentVote.Upvote || !commentVote.Downvote {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "User has not downvoted on this comment!",
		})
		return
	}

	comment.Downvote -= 1
	commentVote.Downvote = false
	commentVote.LastEdited = pq.NullTime{Time: time.Now(), Valid: true}
	db.Db.Save(&comment)
	db.Db.Save(&commentVote)
	c.JSON(http.StatusOK, gin.H{
		"Message": "User downvote removed sucessfully!",
	})
}

func validateCommentMessage(message string) bool {
	if len(message) >= 1 && len(message) <= 2048 {
		return true
	} else {
		return false
	}
}
