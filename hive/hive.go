package hive

import (
	"example/hivemind-be/account"
	"example/hivemind-be/db"
	"fmt"
	"net/http"
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

func CreateHive(c *gin.Context) {
	var hive Hive
	var account account.Account

	if err := c.BindJSON(&hive); err != nil {
		return
	}

	if result := db.Db.Where("username = ?", hive.Creator).First(&account); result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"Error": "An error occurred. Please try again.",
		})
		return
	}

	hive.UUID = uuid.NewString()
	hive.AccountUUID = account.UUID
	hive.MemberCount = 0
	hive.TotalUpvotes = 0
	hive.TotalDownvotes = 0
	hive.TotalComments = 0
	hive.Archived = false
	hive.Banned = false
	hive.Created = pq.NullTime{Time: time.Now(), Valid: true}
	hive.LastEdited = pq.NullTime{Valid: false}

	if result := db.Db.Create(&hive); result.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusCreated, hive)
}

func GetHive(c *gin.Context) {
	var hive []Hive
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

	if val, ok := jsonDataHasKey(updateHive, "description"); ok {
		hive.Description = val
	}

	hive.LastEdited = pq.NullTime{Time: time.Now(), Valid: true}

	db.Db.Save(&hive)
	c.JSON(http.StatusOK, hive)
}

func jsonDataHasKey(data Hive, key string) (string, bool) {
	switch key {
	case "description":
		return data.Description, true
	default:
		return "null", false
	}
}
