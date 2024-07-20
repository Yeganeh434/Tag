package http

import (
	"log"
	"tag_project/internal/adapters/databases/mysql"
	"tag_project/internal/application/usecases"
	"tag_project/internal/domain/entity"

	"github.com/gin-gonic/gin"
)

type Taxonomy struct {
	FromTag          uint64 `json:"fromTag"`
	ToTag            uint64 `json:"toTag"`
	RelationshipType string `json:"relationshipType"`
}

type TagRelationship struct {
	ID               uint64 `json:"id"`
	RelationshipType string `json:"relationshipType"`
}

func RegisterTagRelationship(c *gin.Context) {
	var requestBody Taxonomy
	err := c.BindJSON(&requestBody)
	if err != nil {
		log.Printf("error binding JSON:%v", err)
		c.Status(400)
		return
	}
	taxonomyID, err := usecases.GenerateID()
	if err != nil {
		log.Printf("error generating ID:%v", err)
		c.Status(400)
		return
	}
	taxonomyInfo := entity.Taxonomy{
		ID:               taxonomyID,
		FromTag:          requestBody.FromTag,
		ToTag:            requestBody.ToTag,
		RelationshipType: requestBody.RelationshipType,
		Status:           "active",
	}
	err = mysql.TagDB.RegisterTagRelationship(taxonomyInfo)
	if err != nil {
		log.Printf("error registering tag relationship:%v", err)
		c.Status(400)
		return
	}
	c.JSON(200, gin.H{
		"message": "tags relationship registered seccussfully",
	})
}

func SetTagRelationship(c *gin.Context) {
	var requestBody TagRelationship
	err:=c.BindJSON(&requestBody)
	if err!=nil {
		log.Printf("error binding JSON:%v",err)
		c.Status(400)
		return
	}
	err=mysql.TagDB.SaveTagRelationship(requestBody.ID,requestBody.RelationshipType)
	if err!=nil {
		log.Printf("error saving tag relationship:%v",err)
		c.Status(400)
		return
	}
	c.JSON(200, gin.H{
		"message": "tag relationship setup was successful",
	})
}

func GetRelatedTagsByKey(c *gin.Context) {

}

func GetRelatedTagsByID(c *gin.Context) {

}

func SearchTagByTitle(c *gin.Context) {

}
