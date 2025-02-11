package routes

import (
	"github.com/gin-gonic/gin"
	"main.go/controller"
)

func MarkdownParserRoutes(router *gin.Engine) {
	router.MaxMultipartMemory = 8 << 20 // 8Mib
	router.POST("/api/v1/markdown", controller.SpellCheckMarkdown())
}
