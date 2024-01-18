package routes

import (
	"example/hivemind-be/comments"
	"example/hivemind-be/content"
	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) {
	router.GET("/content", content.GetContent)
	router.GET("/content/id/:id", content.GetContentById)
	router.GET("/content/uuid/:uuid", content.GetContentByUuid)

	router.POST("/content", content.CreateContent)

	router.PATCH("content/uuid/:uuid/add-upvote", content.AddContentUpvoteByUuid)
	router.PATCH("content/uuid/:uuid/remove-upvote", content.RemoveContentUpvoteByUuid)
	router.PATCH("content/uuid/:uuid/add-downvote", content.AddContentDownvoteByUuid)
	router.PATCH("content/uuid/:uuid/remove-downvote", content.RemoveContentDownvoteByUuid)
	router.PATCH("content/uuid/:uuid/delete", content.DeleteContentByUuid)
	router.PATCH("content/uuid/:uuid/undelete", content.UndeleteContentByUuid)
	router.PATCH("content/uuid/:uuid/update", content.UpdateContentByUuid)

	router.GET("/content/uuid/:uuid/comments", comments.GetCommentsByContentUuid)
	router.POST("/content/uuid/:uuid/comment", comments.CreateComment)

	router.GET("/comment/uuid/:uuid", comments.GetCommentByUuid)
	router.PATCH("comment/uuid/:uuid/delete", comments.DeleteCommentByUuid)
	router.PATCH("comment/uuid/:uuid/undelete", comments.UndeleteCommentByUuid)
	router.PATCH("comment/uuid/:uuid/update", comments.UpdateCommentByUuid)
	router.PATCH("comment/uuid/:uuid/add-upvote", comments.AddCommentUpvoteByUuid)
	router.PATCH("comment/uuid/:uuid/remove-upvote", comments.RemoveCommentUpvoteByUuid)
	router.PATCH("comment/uuid/:uuid/add-downvote", comments.AddCommentDownvoteByUuid)
	router.PATCH("comment/uuid/:uuid/remove-downvote", comments.RemoveCommentDownvoteByUuid)
}
