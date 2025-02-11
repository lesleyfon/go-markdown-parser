package routes

import (
	"github.com/gin-gonic/gin"
	"main.go/controller"
)

func MarkdownParserRoutes(router *gin.Engine) {
	router.POST("/api/v1/markdown", controller.SpellCheckMarkdown())
}
