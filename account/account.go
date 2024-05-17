package account

import (
	"example/hivemind-be/db"
	"example/hivemind-be/token"
	"example/hivemind-be/utils"
	"net/http"
	"net/mail"
	"strings"
	"time"

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

type ResponseAccount struct {
	Username string      `json:"Username"`
	Email    string      `json:"Email"`
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

	hashedPassword, err := utils.HashPassword(acc.Password)
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

func AccountLogin(c *gin.Context) {
	var acc Account
	var account Account

	if err := c.BindJSON(&acc); err != nil {
		return
	}

	if result := db.Db.Where("email = ?", strings.ToLower(acc.Email)).First(&account); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Account not found. Please try again.",
		})
		return
	}

	if !utils.DoPasswordsMatch(account.Password, acc.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Password is incorrect. Please try again.",
		})
		return
	}

	token, err := token.CreateToken(account.Username, account.UUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Error: A error occurred creating the token. Please try again.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Token": token,
	})
}

func ValidateAccountToken(c *gin.Context) {
	authToken := c.GetHeader("Authorization")
	_, validToken := utils.ValidateAuthentication(c, authToken)
	if !validToken {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Message": "Token is valid.",
	})
}

func GetAccount(c *gin.Context) {
	authToken := c.GetHeader("Authorization")
	claims, validToken := utils.ValidateAuthentication(c, authToken)
	if !validToken {
		return
	}

	var account Account
	if result := db.Db.Where("uuid = ?", claims.AccountUUID).First(&account); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Account not found. Please try again.",
		})
		return
	}

	var responseAccount = ResponseAccount{
		Username: account.Username,
		Email:    account.Email,
		UUID:     account.UUID,
		Deleted:  account.Deleted,
		Banned:   account.Banned,
		Created:  account.Created,
	}

	c.JSON(http.StatusOK, responseAccount)
}
