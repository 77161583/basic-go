package web

import "github.com/gin-gonic/gin"

type handler interface {
	RegisterRoutes(serve *gin.Engine)
}
