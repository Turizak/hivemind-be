package content

import (
	"example/hivemind-be/account"
	"example/hivemind-be/db"
	"example/hivemind-be/hive"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

// GORM uses the name of your type as the DB table to query. Here the type is Message so gorm will use the messages table by default.
type Content struct {
	ID           int32       `json:"Id" gorm:"primaryKey:type:int32"` //cannot be updated
	Hive         string      `json:"Hive"`                            //cannot be updated
	Title        string      `json:"Title"`                           //can be updated
	Author       string      `json:"Author"`                          //cannot be updated
	Message      string      `json:"Message"`                         //can be updated
	UUID         string      `json:"Uuid"`                            //cannot be update
	HiveUUID     string      `json:"HiveUuid"`                        //cannot be update
	AccountUUID  string      `json:"AccountUuid"`                     //cannot be update
	Link         string      `json:"Link" gorm:"default:null"`        //can be updated
	ImageLink    string      `json:"ImageLink" gorm:"default:null"`   //can be updated
	Upvote       int32       `json:"Upvote"`                          //cannot be updated
	Downvote     int32       `json:"Downvote"`                        //cannot be updated
	CommentCount int32       `json:"CommentCount"`                    //cannot be updated
	Deleted      bool        `json:"Deleted"`                         //can be updated
	Created      pq.NullTime `json:"Created"`                         //cannot be updated
	LastEdited   pq.NullTime `json:"LastEdited"`                      //updated when an update occurs
}

func GetContent(c *gin.Context) {
	var content []Content
	if result := db.Db.Order("id asc").Find(&content); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, &content)
}

func GetContentById(c *gin.Context) {
	var content Content
	id := c.Param("id")
	if result := db.Db.First(&content, id); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, content)
}

func GetContentByUuid(c *gin.Context) {
	var content Content
	uuid := c.Param("uuid")
	if result := db.Db.Where("uuid = ?", uuid).First(&content); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, content)
}

func GetContentByHiveUuid(c *gin.Context) {
	var content []Content
	uuid := c.Param("uuid")
	if result := db.Db.Where("hive_uuid = ?", uuid).Find(&content); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, content)
}

func CreateContent(c *gin.Context) {
	var content Content
	var hive hive.Hive
	var account account.Account

	if err := c.BindJSON(&content); err != nil {
		return
	}

	if result := db.Db.Where("name = ?", content.Hive).First(&hive); result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"Error": "Hive not found! Please use an existing hive or create a new hive first.",
		})
		return
	}

	if result := db.Db.Where("username = ?", content.Author).First(&account); result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"Error": "An error occurred. Please try again.",
		})
		return
	}

	content.UUID = uuid.NewString()
	content.HiveUUID = hive.UUID
	content.AccountUUID = account.UUID
	content.Upvote = 0
	content.Downvote = 0
	content.CommentCount = 0
	content.Deleted = false
	content.LastEdited = pq.NullTime{Valid: false}
	content.Created = pq.NullTime{Time: time.Now(), Valid: true}

	if result := db.Db.Create(&content); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": "There was an error creating this content. Please try again.",
		})
		return
	}

	hive.TotalContent += 1
	db.Db.Save(&hive)
	c.JSON(http.StatusCreated, content)
}

func AddContentUpvoteByUuid(c *gin.Context) {
	var content Content
	var hive hive.Hive
	uuid := c.Param("uuid")

	if result := db.Db.Where("uuid = ?", uuid).First(&content); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	if result := db.Db.Where("uuid = ?", content.HiveUUID).First(&hive); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	content.Upvote += 1
	hive.TotalUpvotes += 1
	db.Db.Save(&content)
	db.Db.Save(&hive)
	c.JSON(http.StatusOK, content)
}

func RemoveContentUpvoteByUuid(c *gin.Context) {
	var content Content
	var hive hive.Hive
	uuid := c.Param("uuid")

	if result := db.Db.Where("uuid = ?", uuid).First(&content); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	if content.Upvote <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "no upvotes to remove!"})
		return
	}

	if result := db.Db.Where("uuid = ?", content.HiveUUID).First(&hive); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	content.Upvote -= 1
	hive.TotalUpvotes -= 1
	db.Db.Save(&content)
	db.Db.Save(&hive)
	c.JSON(http.StatusOK, content)
}

func AddContentDownvoteByUuid(c *gin.Context) {
	var content Content
	var hive hive.Hive
	uuid := c.Param("uuid")

	if result := db.Db.Where("uuid = ?", uuid).First(&content); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	if result := db.Db.Where("uuid = ?", content.HiveUUID).First(&hive); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	content.Downvote += 1
	hive.TotalDownvotes += 1
	db.Db.Save(&content)
	db.Db.Save(&hive)
	c.JSON(http.StatusOK, content)
}

func RemoveContentDownvoteByUuid(c *gin.Context) {
	var content Content
	var hive hive.Hive
	uuid := c.Param("uuid")

	if result := db.Db.Where("uuid = ?", uuid).First(&content); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	if content.Downvote <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "no downvotes to remove!"})
		return
	}

	if result := db.Db.Where("uuid = ?", content.HiveUUID).First(&hive); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	content.Downvote -= 1
	hive.TotalDownvotes -= 1
	db.Db.Save(&content)
	db.Db.Save(&hive)
	c.JSON(http.StatusOK, content)
}

func DeleteContentByUuid(c *gin.Context) {
	var content Content
	var hive hive.Hive
	uuid := c.Param("uuid")

	if result := db.Db.Where("uuid = ?", uuid).First(&content); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	if content.Deleted {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "content has already been deleted!"})
		return
	}

	if result := db.Db.Where("uuid = ?", content.HiveUUID).First(&hive); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	hive.TotalContent -= 1
	content.Deleted = true
	content.LastEdited = pq.NullTime{Time: time.Now(), Valid: true}
	db.Db.Save(&content)
	db.Db.Save(&hive)
	c.JSON(http.StatusOK, content)
}

func UndeleteContentByUuid(c *gin.Context) {
	var content Content
	var hive hive.Hive
	uuid := c.Param("uuid")

	if result := db.Db.Where("uuid = ?", uuid).First(&content); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	if !content.Deleted {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "content has not been deleted!"})
		return
	}

	if result := db.Db.Where("uuid = ?", content.HiveUUID).First(&hive); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	hive.TotalContent += 1
	content.Deleted = false
	content.LastEdited = pq.NullTime{Time: time.Now(), Valid: true}
	db.Db.Save(&content)
	db.Db.Save(&hive)
	c.JSON(http.StatusOK, content)
}

func UpdateContentByUuid(c *gin.Context) {
	var content Content
	var updateContent Content

	if err := c.BindJSON(&updateContent); err != nil {
		return
	}

	uuid := c.Param("uuid")

	if result := db.Db.Where("uuid = ?", uuid).First(&content); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	if val, ok := jsonDataHasKey(updateContent, "title"); ok {
		content.Title = val
	}
	if val, ok := jsonDataHasKey(updateContent, "message"); ok {
		content.Message = val
	}
	if val, ok := jsonDataHasKey(updateContent, "link"); ok {
		content.Link = val
	}
	if val, ok := jsonDataHasKey(updateContent, "imageLink"); ok {
		content.ImageLink = val
	}

	content.LastEdited = pq.NullTime{Time: time.Now(), Valid: true}

	db.Db.Save(&content)
	c.JSON(http.StatusOK, content)
}

func jsonDataHasKey(data Content, key string) (string, bool) {
	switch key {
	case "title":
		return data.Title, true
	case "message":
		return data.Message, true
	case "link":
		return data.Link, true
	case "imageLink":
		return data.ImageLink, true
	default:
		return "null", false
	}
}
