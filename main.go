package main

import (
	"example/hivemind-be/db"
	"example/hivemind-be/routes"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	db.ConnectDatabase()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"}, // Change this to the origin you want to allow
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
		MaxAge:       12 * time.Hour,
	}))
	routes.Routes(router)
	router.Run(os.Getenv("PORT"))
}
