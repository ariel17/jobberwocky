package http

import "github.com/gin-gonic/gin"

type Handler interface {
	ConfigureRoutes(router *gin.Engine)
}