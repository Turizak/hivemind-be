package hive

import (
	"example/hivemind-be/db"
	"example/hivemind-be/token"
	"example/hivemind-be/utils"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Hive struct {
	ID             int32       `json:"Id" gorm:"primaryKey:type:int32"`
	Name           string      `json:"Name"`
	Creator        string      `json:"Creator"`
	Description    string      `json:"Description"`
	UUID           string      `json:"Uuid"`
	AccountUUID    string      `json:"AccountUUID"`
	MemberCount    int32       `json:"MemberCount"`
	TotalUpvotes   int32       `json:"TotalUpvotes"`
	TotalDownvotes int32       `json:"TotalDownvotes"`
	TotalComments  int32       `json:"TotalComments"`
	TotalContent   int32       `json:"TotalContent"`
	Archived       bool        `json:"Archived"`
	Banned         bool        `json:"Banned"`
	Created        pq.NullTime `json:"Created"`
	LastEdited     pq.NullTime `json:"LastEdited"`
}

// CreateHive creates a new hive based on the provided JSON data in the request body.
// It requires a valid authorization token in the "Authorization" header.
// The function parses the token and retrieves the claims, which include the username and account UUID.
// The function then binds the JSON data to the `hive` variable.
// It sets the creator, UUID, account UUID, member count, upvotes, downvotes, comments, archived, banned,
// created time, and last edited time for the hive.
// If there is an error during the creation process, it returns an error response with the corresponding status code.
// Otherwise, it returns a success response with the created hive.
func CreateHive(c *gin.Context) {
	var hive Hive

	authToken := c.GetHeader("Authorization")
	validToken := token.CheckToken(c, authToken)

	if !validToken {
		return
	}
	claims, err := token.ParseToken(authToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized.",
		})
		return
	}

	if err := c.BindJSON(&hive); err != nil {
		return
	}

	pattern := "^[a-zA-Z]{1,30}$"
	regex, err := regexp.Compile(pattern)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Unknown error occurred. Please try again.",
		})
		return
	}
	// Check if the test string matches the pattern
	if !regex.MatchString(hive.Name) {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Hive name should be between 1 and 30 characters long and contain only letters.",
		})
		return
	}

	hive.Creator = claims.Username
	hive.UUID = uuid.NewString()
	hive.AccountUUID = claims.AccountUUID
	hive.MemberCount = 0
	hive.TotalUpvotes = 0
	hive.TotalDownvotes = 0
	hive.TotalComments = 0
	hive.Archived = false
	hive.Banned = false
	hive.Created = pq.NullTime{Time: time.Now(), Valid: true}
	hive.LastEdited = pq.NullTime{Valid: false}

	if result := db.Db.Create(&hive); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, hive)
}

func GetHive(c *gin.Context) {
	var hive []Hive
	authToken := c.GetHeader("Authorization")
	validToken := token.CheckToken(c, authToken)

	if !validToken {
		return
	}

	if result := db.Db.Order("id asc").Find(&hive); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, &hive)
}

func BanHiveByUuid(c *gin.Context) {
	var hive Hive

	authToken := c.GetHeader("Authorization")
	validToken := token.CheckToken(c, authToken)

	if !validToken {
		return
	}

	uuid := c.Param("uuid")
	if result := db.Db.Where("uuid = ?", uuid).First(&hive); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	if hive.Banned {
		mes := fmt.Sprintf("%s is already banned!", hive.Name)
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": mes,
		})
		return
	}
	hive.Banned = true
	db.Db.Save(&hive)
	mes := fmt.Sprintf("%s has been banned!", hive.Name)
	c.JSON(http.StatusOK, gin.H{
		"Message": mes,
	})
}

func UnBanHiveByUuid(c *gin.Context) {
	var hive Hive

	authToken := c.GetHeader("Authorization")
	validToken := token.CheckToken(c, authToken)

	if !validToken {
		return
	}

	uuid := c.Param("uuid")
	if result := db.Db.Where("uuid = ?", uuid).First(&hive); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	if !hive.Banned {
		mes := fmt.Sprintf("%s has not been banned!", hive.Name)
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": mes,
		})
		return
	}
	hive.Banned = false
	db.Db.Save(&hive)
	mes := fmt.Sprintf("%s has been unbanned!", hive.Name)
	c.JSON(http.StatusOK, gin.H{
		"Message": mes,
	})
}

func ArchiveHiveByUuid(c *gin.Context) {
	var hive Hive

	authToken := c.GetHeader("Authorization")
	validToken := token.CheckToken(c, authToken)

	if !validToken {
		return
	}

	uuid := c.Param("uuid")
	if result := db.Db.Where("uuid = ?", uuid).First(&hive); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	if hive.Archived {
		mes := fmt.Sprintf("%s is already archived!", hive.Name)
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": mes,
		})
		return
	}
	hive.Archived = true
	db.Db.Save(&hive)
	mes := fmt.Sprintf("%s has been archived!", hive.Name)
	c.JSON(http.StatusOK, gin.H{
		"Message": mes,
	})
}

func UnArchiveHiveByUuid(c *gin.Context) {
	var hive Hive

	authToken := c.GetHeader("Authorization")
	validToken := token.CheckToken(c, authToken)

	if !validToken {
		return
	}

	uuid := c.Param("uuid")
	if result := db.Db.Where("uuid = ?", uuid).First(&hive); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	if !hive.Archived {
		mes := fmt.Sprintf("%s has not been archived!", hive.Name)
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": mes,
		})
		return
	}
	hive.Archived = false
	db.Db.Save(&hive)
	mes := fmt.Sprintf("%s has been unarchived!", hive.Name)
	c.JSON(http.StatusOK, gin.H{
		"Message": mes,
	})
}

func UpdateHiveByUuid(c *gin.Context) {
	var hive Hive
	var updateHive Hive

	authToken := c.GetHeader("Authorization")
	validToken := token.CheckToken(c, authToken)

	if !validToken {
		return
	}

	if err := c.BindJSON(&updateHive); err != nil {
		return
	}

	uuid := c.Param("uuid")

	if result := db.Db.Where("uuid = ?", uuid).First(&hive); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	if val, ok := utils.JsonDataHasKey(updateHive, "Description"); ok {
		hive.Description, _ = val.(string)
	}

	hive.LastEdited = pq.NullTime{Time: time.Now(), Valid: true}

	db.Db.Save(&hive)
	c.JSON(http.StatusOK, hive)
}
