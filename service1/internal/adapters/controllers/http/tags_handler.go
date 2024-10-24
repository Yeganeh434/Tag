package http

import (
	"context"
	"errors"
	"log"
	"service1/internal/application/usecases"
	"service1/internal/config"
	"service1/internal/domain/entity"
	"service1/internal/domain/service"

	"github.com/gin-gonic/gin"
	// "go.opentelemetry.io/otel/trace"
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
	config.RequestsCounter.Add(config.Ctx, 1)

	ctx := context.Background()
	ctx, span := config.Tracer.Start(ctx, "RegisterApprovedTag_handler")
	defer span.End()

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
	err = h.usecase.RegisterTag(tagInfo, ctx)
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
	config.RequestsCounter.Add(config.Ctx, 1)

	ctx := context.Background()
	ctx, span := config.Tracer.Start(ctx, "RegisterTagAsDraft_handler")
	defer span.End()

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
	err = h.usecase.RegisterTag(tagInfo, ctx)
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
	config.RequestsCounter.Add(config.Ctx, 1)

	ctx := context.Background()
	ctx, span := config.Tracer.Start(ctx, "ApproveOrRejectTag_handler")
	defer span.End()

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
	err = h.usecase.UpdateTagStatus(requestBody.ID, isApproved, ctx)
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
	config.RequestsCounter.Add(config.Ctx, 1)

	ctx := context.Background()
	ctx, span := config.Tracer.Start(ctx, "MergeTags_handler")
	defer span.End()

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
	err = h.usecase.RegisterTag(tagInfo, ctx)
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
	err = h.usecase.MergeTags(requestBody.OriginalTagID, tagInfo.ID, ctx)
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
