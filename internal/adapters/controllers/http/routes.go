package http

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func RunWebServer(port int) error {
	addr := ":" + strconv.Itoa(port)
	router := gin.New()

	//tags management
	router.POST("/register_approved_tag", RegisterApprovedTag)
	router.POST("/register_tag_as_draft", RegisterTagAsDraft)
	router.POST("/approve_or_reject_tag", ApproveOrRejectTag)
	router.POST("/merge_tags", MergeTags)

	//tags relationship management
	router.POST("/register_tag_relationship", RegisterTagRelationship)
	router.POST("/set_tag_relationship", SetTagRelationship)
	router.POST("/get_related_tags_by_key", GetRelatedTagsByKey)
	router.POST("/get_related_tags_by_id", GetRelatedTagsByID)
	router.POST("/search_tag_by_title", SearchTagByTitle)

	err := router.Run(addr)
	return err
}
