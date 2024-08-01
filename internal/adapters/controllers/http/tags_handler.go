package http

import (
	"errors"
	"log"
	"tag_project/internal/application/usecases"
	"tag_project/internal/domain/entity"
	"tag_project/internal/domain/service"

	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	usecase *usecases.TagManagementUseCase
}

func NewTagHandler(usecase *usecases.TagManagementUseCase) *TagHandler {
	return &TagHandler{usecase: usecase}
}

type Tag struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Picture     string `json:"picture"`
	Key         string `json:"key"`
}

type TagStatus struct {
	ID         uint64 `json:"id"`
	IsApproved bool   `json:"isApproved"`
}

type TagMerge struct {
	OriginalTagID uint64 `json:"originalTagID"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Picture       string `json:"picture"`
	Key           string `json:"key"`
}

func (h *TagHandler) RegisterApprovedTag(c *gin.Context) {
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
	tagInfo := entity.TagEntity{
		ID:          tagID,
		Title:       requestBody.Title,
		Description: requestBody.Description,
		Picture:     requestBody.Picture,
		Key:         requestBody.Key,
		Status:      "approved",
	}
	err = h.usecase.RegisterTag(tagInfo)
	if err != nil {
		if errors.Is(err, service.ErrTitleCannotBeEmpty) {
			c.JSON(400, gin.H{
				"error": "tag title cannot be empty",
			})
			return
		}
		if errors.Is(err, service.ErrTagKeyAlreadyExists) {
			c.JSON(400, gin.H{
				"error": "tag key already exists",
			})
			return
		}
		log.Printf("error registering tag:%v", err)
		c.Status(400)
		return
	}
	c.JSON(200, gin.H{
		"message": "tag registered successfully",
	})
}

func (h *TagHandler) RegisterTagAsDraft(c *gin.Context) {
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
	tagInfo := entity.TagEntity{
		ID:          tagID,
		Title:       requestBody.Title,
		Description: requestBody.Description,
		Picture:     requestBody.Picture,
		Key:         requestBody.Key,
		Status:      "under_review",
	}
	err = h.usecase.RegisterTag(tagInfo)
	if err != nil {
		if errors.Is(err, service.ErrTitleCannotBeEmpty) {
			c.JSON(400, gin.H{
				"error": "tag title cannot be empty",
			})
			return
		}
		if errors.Is(err, service.ErrTagKeyAlreadyExists) {
			c.JSON(400, gin.H{
				"error": "tag key already exists",
			})
			return
		}
		log.Printf("error registering tag:%v", err)
		c.Status(400)
		return
	}
	c.JSON(200, gin.H{
		"message": "tag registered successfully",
	})
}

func (h *TagHandler) ApproveOrRejectTag(c *gin.Context) {
	var requestBody TagStatus
	err := c.BindJSON(&requestBody)
	if err != nil {
		log.Printf("error binding JSON:%v", err)
		c.Status(400)
		return
	}
	isApproved := "rejected"
	if requestBody.IsApproved {
		isApproved = "approved"
	}
	err = h.usecase.UpdateTagStatus(requestBody.ID, isApproved)
	if err != nil {
		if errors.Is(err, service.ErrInvalidTagID) {
			c.JSON(400, gin.H{
				"error": "invalid tag ID",
			})
			return
		}
		if errors.Is(err, service.ErrNoTagExistsWithThisID) {
			c.JSON(400, gin.H{
				"error": "no tag exists with this ID",
			})
			return
		}
		log.Printf("error updating tag status:%v", err)
		c.Status(400)
		return
	}
	c.JSON(200, gin.H{
		"message": "tag status updated seccessfully",
	})
}

func (h *TagHandler) MergeTags(c *gin.Context) {
	var requestBody TagMerge
	err := c.BindJSON(&requestBody)
	if err != nil {
		log.Printf("error binding json:%v", err)
		c.Status(400)
		return
	}
	mergeTagID, err := usecases.GenerateID()
	if err != nil {
		log.Printf("error generating ID:%v", err)
		c.Status(400)
		return
	}
	tagInfo := entity.TagEntity{
		ID:          mergeTagID,
		Title:       requestBody.Title,
		Description: requestBody.Description,
		Picture:     requestBody.Picture,
		Key:         requestBody.Key,
		Status:      "approved",
	}
	err = h.usecase.RegisterTag(tagInfo)
	if err != nil {
		if errors.Is(err, service.ErrTitleCannotBeEmpty) {
			c.JSON(400, gin.H{
				"error": "tag title cannot be empty",
			})
			return
		}
		if errors.Is(err, service.ErrTagKeyAlreadyExists) {
			c.JSON(400, gin.H{
				"error": "tag key already exists",
			})
			return
		}
		log.Printf("error registering tag:%v", err)
		c.Status(400)
		return
	}
	err = h.usecase.MergeTags(requestBody.OriginalTagID, tagInfo.ID)
	if err != nil {
		if errors.Is(err, service.ErrInvalidTagID) {
			c.JSON(400, gin.H{
				"error": "invalid tag ID",
			})
			return
		}
		if errors.Is(err, service.ErrNoTagExistsWithThisID) {
			c.JSON(400, gin.H{
				"error": "no tag exists with this original ID",
			})
			return
		}
		log.Printf("error merging tags:%v", err)
		c.Status(400)
		return
	}
	c.JSON(200, gin.H{
		"message": "tags merged successfully",
	})
}
