package account

import (
	"example/hivemind-be/db"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Account struct {
	ID       int32       `json:"Id" gorm:"primaryKey:type:int32"`
	Username string      `json:"Username"`
	Email    string      `json:"Email"`
	Password string      `json:"Password"`
	UUID     string      `json:"Uuid"`
	Deleted  bool        `json:"Deleted"`
	Banned   bool        `json:"Banned"`
	Created  pq.NullTime `json:"Created"`
}

func CreateAccount(c *gin.Context) {
	var acc Account

	if err := c.BindJSON(&acc); err != nil {
		return
	}

	addr, err := mail.ParseAddress(strings.ToLower(acc.Email))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Error: Email address format is not valid. Please us a valid email address.",
		})
		return
	}

	hashedPassword, err := hashPassword(acc.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Error: A error occurred creating the account. Please try again.",
		})
		return
	}

	acc.Username = strings.ToLower(acc.Username)
	acc.Email = strings.ToLower(addr.Address)
	acc.Password = hashedPassword
	acc.UUID = uuid.NewString()
	acc.Deleted = false
	acc.Banned = false
	acc.Created = pq.NullTime{Time: time.Now(), Valid: true}

	if result := db.Db.Create(&acc); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Could not create account. Please try again.",
		})
		return
	}

	c.JSON(http.StatusCreated, acc)
}

// Hash password
func hashPassword(password string) (string, error) {
	// Convert password string to byte slice
	var passwordBytes = []byte(password)

	// Hash password with Bcrypt's min cost
	hashedPasswordBytes, err := bcrypt.
		GenerateFromPassword(passwordBytes, bcrypt.MinCost)

	return string(hashedPasswordBytes), err
}

// Check if two passwords match using Bcrypt's CompareHashAndPassword
// which return nil on success and an error on failure.
//   func doPasswordsMatch(hashedPassword, currPassword string) bool {
// 	err := bcrypt.CompareHashAndPassword(
// 	  []byte(hashedPassword), []byte(currPassword))
// 	return err == nil
//   }
