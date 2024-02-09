package main

import (
	"example/hivemind-be/db"
	"example/hivemind-be/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	db.ConnectDatabase()
	router.Use(cors.Default())
	routes.Routes(router)
	router.Run(":8080")
}
