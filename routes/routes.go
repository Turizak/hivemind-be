package routes

import (
	"example/hivemind-be/messages"
	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) {
	router.GET("/messages", messages.GetMessages)
	router.GET("/message/:id", messages.GetMessageById)
	router.POST("/message", messages.CreateMessage)
	router.PATCH("message/add-upvote", messages.AddMessageUpvote)
	router.PATCH("message/remove-upvote", messages.RemoveMessageUpvote)
	router.PATCH("message/add-downvote", messages.AddMessageDownvote)
	router.PATCH("message/remove-downvote", messages.RemoveMessageDownvote)
	router.PATCH("message/delete", messages.DeleteMessage)
	router.PATCH("message/undelete", messages.UndeleteMessage)
	router.PATCH("message/update", messages.UpdateMessage)
}
