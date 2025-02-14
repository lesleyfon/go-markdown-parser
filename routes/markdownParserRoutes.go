package routes

import (
	"github.com/gin-gonic/gin"
	"main.go/controller"
)

func MarkdownParserRoutes(router *gin.Engine) {
	router.MaxMultipartMemory = 8 << 20 // 8Mib
	router.POST("/api/v1/markdown", controller.SpellCheckMarkdown())
	// Get all files names for a user
	router.GET("/api/v1/markdown/files", controller.GetAllFiles())
	// Get a file by id
	router.GET("/api/v1/markdown/files/:file_id", controller.GetFileById())
}
