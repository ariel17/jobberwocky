package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

	"github.com/ariel17/jobberwocky/api"
	"github.com/ariel17/jobberwocky/configs"
)

type Handler interface {
	ConfigureRoutes(router *gin.Engine)
}

type swaggerHandler struct{}

func NewSwaggerHandler() Handler {
	return &swaggerHandler{}
}

func (s *swaggerHandler) ConfigureRoutes(router *gin.Engine) {
	api.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", configs.GetHTTPPort())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}