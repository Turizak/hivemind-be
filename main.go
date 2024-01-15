package main

import (
	"example/hivemind-be/routes"
	"example/hivemind-be/db"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	db.ConnectDatabase()
	routes.Routes(router)
	router.Run(":8080")
}
