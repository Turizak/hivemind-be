package routes

import (
	"example/hivemind-be/content"
	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) {
	router.GET("/content", content.GetContent)
	router.GET("/content/id/:id", content.GetContentById)
	router.GET("/content/uuid/:uuid", content.GetContentByUuid)
	router.POST("/content", content.CreateContent)
	router.PATCH("content/add-upvote", content.AddContentUpvote)
	router.PATCH("content/remove-upvote", content.RemoveContentUpvote)
	router.PATCH("content/add-downvote", content.AddContentDownvote)
	router.PATCH("content/remove-downvote", content.RemoveContentDownvote)
	router.PATCH("content/delete", content.DeleteContent)
	router.PATCH("content/undelete", content.UndeleteContent)
	router.PATCH("content/update", content.UpdateContent)
}
