package http

import (
	"log"
	"tag_project/internal/adapters/databases/mysql"
	"tag_project/internal/application/usecases"
	"tag_project/internal/domain/entity"

	"github.com/gin-gonic/gin"
)

type Tag struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Picture     string `json:"picture"`
}

func RegisterApprovedTag(c *gin.Context) {
	var requestBody Tag
	err := c.BindJSON(&requestBody)
	if err != nil {
		log.Printf("error binding json:%v", err)
		c.Status(400)
		return
	}
	tagID, err := usecases.GenerateID()
	if err != nil {
		log.Printf("error generating ID:%v", err)
		c.Status(400)
		return
	}
	tagInfo := entity.Tag{
		ID:          tagID,
		Title:       requestBody.Title,
		Description: requestBody.Description,
		Picture:     requestBody.Picture,
		Key:         "........", //////////////////////////////////////////////
		Status:      "approved",
	}
	err = mysql.TagDB.RegisterTag(tagInfo)
	if err != nil {
		log.Printf("error registering tag:%v", err)
		c.Status(400)
		return
	}
	c.JSON(200, gin.H{
		"message": "tag registered successfully",
	})
}

func RegisterTagAsDraft(c *gin.Context) {
	var requestBody Tag
	err := c.BindJSON(&requestBody)
	if err != nil {
		log.Printf("error binding json:%v", err)
		c.Status(400)
		return
	}
	tagID, err := usecases.GenerateID()
	if err != nil {
		log.Printf("error generating ID:%v", err)
		c.Status(400)
		return
	}
	tagInfo := entity.Tag{
		ID:          tagID,
		Title:       requestBody.Title,
		Description: requestBody.Description,
		Picture:     requestBody.Picture,
		Key:         "........", //////////////////////////////////////////////
		Status:      "under_review",
	}
	err = mysql.TagDB.RegisterTag(tagInfo)
	if err != nil {
		log.Printf("error registering tag:%v", err)
		c.Status(400)
		return
	}
	c.JSON(200, gin.H{
		"message": "tag registered successfully",
	})
}

func ApproveOrRejectTag(c *gin.Context) {

}

func MergeTags(c *gin.Context) {

}
