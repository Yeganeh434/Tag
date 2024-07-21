package http

import (
	"log"
	"strconv"
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
	err := c.BindJSON(&requestBody)
	if err != nil {
		log.Printf("error binding JSON:%v", err)
		c.Status(400)
		return
	}
	err = mysql.TagDB.SaveTagRelationship(requestBody.ID, requestBody.RelationshipType)
	if err != nil {
		log.Printf("error saving tag relationship:%v", err)
		c.Status(400)
		return
	}
	c.JSON(200, gin.H{
		"message": "tag relationship setup was successful",
	})
}

func GetRelatedTagsByKey(c *gin.Context) {
	key := c.Param("key")
	ID, err := mysql.TagDB.GetIDByKey(key)
	if err != nil {
		log.Printf("error retrieving ID:%v", err)
		c.Status(400)
		return
	}
	if ID == 0 {
		c.JSON(400, gin.H{
			"message": "no record exists with this key",
		})
		return
	}
	IDs, err := mysql.TagDB.GetRelatedTagsByID(ID)
	if err != nil {
		log.Printf("error getting related tags:%v", err)
		c.Status(400)
		return
	}
	if IDs == nil {
		c.JSON(400, gin.H{
			"message": "no related tags exist for this tag",
		})
		return
	}
	c.JSON(200, IDs)
}

func GetRelatedTagsByID(c *gin.Context) {
	id := c.Param("ID")
	ID, _ := strconv.ParseUint(id,10,64)
	IDs, err := mysql.TagDB.GetRelatedTagsByID(ID)
	if err != nil {
		log.Printf("error getting related tags:%v", err)
		c.Status(400)
		return
	}
	if IDs == nil {
		c.JSON(400, gin.H{
			"message": "no related tags exist for this tag",
		})
		return
	}
	c.JSON(200, IDs)
}

func SearchTagByTitle(c *gin.Context) {

}
