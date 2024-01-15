package db

import (
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // don't forget to add it. It doesn't be added automatically
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"strconv"
)

var Db *gorm.DB //created outside to make it global.

// make sure your function start with uppercase to call outside of the directory.
func ConnectDatabase() {
	err := godotenv.Load() //by default, it is .env so we don't have to write
	if err != nil {
		fmt.Println("Error is occurred  on .env file please check")
	}
	//we read our .env file
	host := os.Getenv("HOST")
	port, _ := strconv.Atoi(os.Getenv("PORT")) // don't forget to convert int since port is int type.
	user := os.Getenv("ACCOUNT")
	dbname := os.Getenv("DB_NAME")
	pass := os.Getenv("PASSWORD")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=require", host, user, pass, dbname, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	
	if err != nil {
		fmt.Println("There is an error while connecting to the database ", err)
		panic(err)
	} else {
		Db = db
		fmt.Println("Successfully connected to database!")
	}
}
