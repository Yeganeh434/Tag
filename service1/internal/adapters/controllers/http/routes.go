package http

import (
	"strconv"
	"service1/internal/adapters/databases/mysql"
	"service1/internal/application/usecases"
	"service1/internal/domain/service"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func RunWebServer(port int) error {
	addr := ":" + strconv.Itoa(port)
	router := gin.New()

	tagRepo := mysql.NewMySQLTagRepository(mysql.TagDB.DB)
	tagService := service.NewTagService(tagRepo)
	tagManagementUseCase := usecases.NewTagManagementUseCase(tagService)
	tagHandler := NewTagHandler(tagManagementUseCase)

	taxonomyRepo := mysql.NewMySQLTaxonomyRepository(mysql.TagDB.DB)
	taxonomyService := service.NewTaxonomyService(taxonomyRepo)
	taxonomyManagementUseCase := usecases.NewTaxonomyManagementUseCase(taxonomyService)
	taxonomyHandler := NewTaxonomyHandler(taxonomyManagementUseCase)

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	//tags management
	router.POST("/register_approved_tag", tagHandler.RegisterApprovedTag)
	router.POST("/register_tag_as_draft", tagHandler.RegisterTagAsDraft)
	router.POST("/approve_or_reject_tag", tagHandler.ApproveOrRejectTag)
	router.POST("/merge_tags", tagHandler.MergeTags)

	//tags relationship management
	router.POST("/register_tag_relationship", taxonomyHandler.RegisterTagRelationship)
	router.POST("/set_tag_relationship", taxonomyHandler.SetTagRelationship)
	router.GET("/get_related_tags_by_key/:key", taxonomyHandler.GetRelatedTagsByKey)
	router.GET("/get_related_tags_by_id/:ID", taxonomyHandler.GetRelatedTagsByID)
	router.GET("/search_tag_by_title/:title", taxonomyHandler.SearchTagByTitle)
	router.GET("/get_tags_with_same_title/:title", taxonomyHandler.GetTagsWithSameTitle)

	err := router.Run(addr)
	// err := router.Run(":8080")
	return err
}
