package http

import (
	"strconv"
	"tag_project/internal/adapters/databases/mysql"
	"tag_project/internal/application/usecases"
	"tag_project/internal/domain/service"

	"github.com/gin-gonic/gin"
)

func RunWebServer(port int) error {
	addr := ":" + strconv.Itoa(port)
	router := gin.New()

	tagRepo := mysql.NewMySQLTagRepository(mysql.TagDB.DB)
	tagService := service.NewTagService(tagRepo)
	tagManagementUseCase := usecases.NewTagManagementUseCase(tagService)
	tagHandler := NewTagHandler(tagManagementUseCase)

	//tags management
	router.POST("/register_approved_tag", tagHandler.RegisterApprovedTag)
	router.POST("/register_tag_as_draft", tagHandler.RegisterTagAsDraft)
	router.POST("/approve_or_reject_tag", tagHandler.ApproveOrRejectTag)
	router.POST("/merge_tags", tagHandler.MergeTags)

	//tags relationship management
	router.POST("/register_tag_relationship", RegisterTagRelationship)
	router.POST("/set_tag_relationship", SetTagRelationship)
	router.POST("/get_related_tags_by_key/:key", GetRelatedTagsByKey)
	router.POST("/get_related_tags_by_id/:id", GetRelatedTagsByID)
	router.POST("/search_tag_by_title/:title", SearchTagByTitle)

	err := router.Run(addr)
	return err
}
