package main

import (
	"example/hivemind-be/db"
	"example/hivemind-be/routes"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	db.ConnectDatabase()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Change this to the origin you want to allow
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	routes.Routes(router)
	router.Run(":8080")
}
