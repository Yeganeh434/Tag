package http

import (
	"errors"
	"log"
	"strconv"
	"tag_project/internal/application/usecases"
	"tag_project/internal/domain/entity"
	"tag_project/internal/domain/service"

	"github.com/gin-gonic/gin"
)

type TaxonomyHandler struct {
	usecase *usecases.TaxonomyManagementUseCase
}

func NewTaxonomyHandler(usecase *usecases.TaxonomyManagementUseCase) *TaxonomyHandler {
	return &TaxonomyHandler{usecase: usecase}
}

type Taxonomy struct {
	FromTag          uint64 `json:"fromTag"`
	ToTag            uint64 `json:"toTag"`
	RelationshipType string `json:"relationshipType"`
}

type TagRelationship struct {
	ID               uint64 `json:"id"`
	RelationshipType string `json:"relationshipType"`
}

func (h *TaxonomyHandler) RegisterTagRelationship(c *gin.Context) {
	var requestBody Taxonomy
	err := c.BindJSON(&requestBody)
	if err != nil {
		log.Printf("error binding JSON:%v", err)
		c.Status(400)
		return
	}
	if requestBody.FromTag==requestBody.ToTag {
		c.JSON(400,gin.H{
			"error":"a tag cannot have a relationship with itself",
		})
		return
	}
	taxonomyID, err := usecases.GenerateID()
	if err != nil {
		log.Printf("error generating ID:%v", err)
		c.Status(400)
		return
	}
	taxonomyInfo := entity.TaxonomyEntity{
		ID:               taxonomyID,
		FromTag:          requestBody.FromTag,
		ToTag:            requestBody.ToTag,
		RelationshipType: requestBody.RelationshipType,
		Status:           "active",
	}
	err = h.usecase.RegisterTagRelationship(taxonomyInfo)
	if err != nil {
		if errors.Is(err, service.ErrInvalidTagID) {
			c.JSON(400, gin.H{
				"error": "invalid tag ID",
			})
			return 
		}
		if errors.Is(err,service.ErrInvalidRelationshipType){
			c.JSON(400,gin.H{
				"error":"invalid relationship type",
			})
			return
		}
		if errors.Is(err,service.ErrNoTagExistsWithThisID){
			c.JSON(400,gin.H{
				"error":"no tag exists with this from tag ID or to tag ID",
			})
			return
		}
		if errors.Is(err,service.ErrThisRelationshipAlreadyExists){
			c.JSON(400,gin.H{
				"error":"this relationship already exists",
			})
			return
		}
		log.Printf("error registering tag relationship:%v", err)
		c.Status(400)
		return
	}
	c.JSON(200, gin.H{
		"message": "tags relationship registered seccussfully",
	})
}

func (h *TaxonomyHandler) SetTagRelationship(c *gin.Context) {
	var requestBody TagRelationship
	err := c.BindJSON(&requestBody)
	if err != nil {
		log.Printf("error binding JSON:%v", err)
		c.Status(400)
		return
	}
	err = h.usecase.SaveTagRelationship(requestBody.ID, requestBody.RelationshipType)
	if err != nil {
		if errors.Is(err, service.ErrInvalidTagID) {
			c.JSON(400, gin.H{
				"error": "invalid tag ID",
			})
			return
		}
		if errors.Is(err,service.ErrRelationshipTypeCannotBeEmpty){
			c.JSON(400,gin.H{
				"error":"relationship type cannot be empty",
			})
			return
		}
		if errors.Is(err,service.ErrInvalidRelationshipType){
			c.JSON(400,gin.H{
				"error":"invalid relationship type",
			})
			return
		}
		if errors.Is(err,service.ErrNoRelationExistsWithThisID){
			c.JSON(400,gin.H{
				"error":"no relation exists with this ID",
			})
			return
		}
		log.Printf("error saving tag relationship:%v", err)
		c.Status(400)
		return
	}
	c.JSON(200, gin.H{
		"message": "tag relationship setup was successful",
	})
}

func (h *TaxonomyHandler) GetRelatedTagsByKey(c *gin.Context) {
	key := c.Param("key")
	ID, err := h.usecase.GetIDByKey(key)
	if err != nil {
		if errors.Is(err,service.ErrKeyCannotBeEmpty){
			c.JSON(400,gin.H{
				"error":"key cannot be empty",
			})
			return
		}
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
	IDs, err := h.usecase.GetRelatedTagsByID(ID)
	if err != nil {
		if errors.Is(err,service.ErrInvalidTagID){
			c.JSON(400,gin.H{
				"error":"invalid tag ID",
			})
			return
		}
		if errors.Is(err,service.ErrNoTagExistsWithThisID){
			c.JSON(400,gin.H{
				"error":"no tag exists with this ID",
			})
			return
		}
		log.Printf("error getting related tags:%v", err)
		c.Status(400)
		return
	}
	if IDs == nil {
		c.JSON(400, gin.H{
			"message": "no related tags exist for this key",
		})
		return
	}
	c.JSON(200, IDs)
}

func (h *TaxonomyHandler) GetRelatedTagsByID(c *gin.Context) {
	id := c.Param("ID")
	ID, _ := strconv.ParseUint(id, 10, 64)
	IDs, err := h.usecase.GetRelatedTagsByID(ID)
	if err != nil {
		if errors.Is(err,service.ErrInvalidTagID){
			c.JSON(400,gin.H{
				"error":"invalid tag ID",
			})
			return
		}
		if errors.Is(err,service.ErrNoTagExistsWithThisID){
			c.JSON(400,gin.H{
				"error":"no tag exists with this ID",
			})
			return
		}
		log.Printf("error getting related tags:%v", err)
		c.Status(400)
		return
	}
	if IDs == nil {
		c.JSON(400, gin.H{
			"message": "no related tags exist for this ID",
		})
		return
	}
	c.JSON(200, IDs)
}

func (h *TaxonomyHandler) SearchTagByTitle(c *gin.Context) {
	title := c.Param("title")
	IDs, err := h.usecase.GetIDsByTitle(title)
	if err != nil {
		if errors.Is(err,service.ErrTitleCannotBeEmpty){
			c.JSON(400,gin.H{
				"error":"title cannot be empty",
			})
			return
		}
		log.Printf("error retrieving IDs:%v", err)
		c.Status(400)
		return
	}
	if IDs == nil {
		c.JSON(400, gin.H{
			"message": "no record exists with this title",
		})
		return
	}
	var relatedTagsID []uint64
	for _, ID := range IDs {
		tempIDs, err := h.usecase.GetRelatedTagsByID(ID)
		if err != nil {
			if errors.Is(err,service.ErrInvalidTagID){
				c.JSON(400,gin.H{
					"error":"invalid tag ID",
				})
				return
			}
			if errors.Is(err,service.ErrNoTagExistsWithThisID){
				c.JSON(400,gin.H{
					"error":"no tag exists with this ID",
				})
				return
			}
			log.Printf("error getting related tags:%v", err)
			c.Status(400)
			return
		}
		relatedTagsID = append(relatedTagsID, tempIDs...)
	}
	if relatedTagsID == nil {
		c.JSON(400, gin.H{
			"message": "no related tags exist for this title",
		})
		return
	}
	c.JSON(200, relatedTagsID)
}
