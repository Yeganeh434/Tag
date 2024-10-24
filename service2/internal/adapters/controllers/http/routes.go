package http

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func RunWebServer()error{
	router:=gin.New()
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	router.POST("/edit_tag",EditTag)

	err:=router.Run(":8081")
	if err!=nil {
		return err
	}
	return nil
}