package controller

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func SpellCheckMarkdown() gin.HandlerFunc {
	return func(c *gin.Context) {

		file, err := c.FormFile("markdownfile")
		if err != nil {
			c.JSON(http.StatusBadRequest,
				gin.H{
					"message": "Bad Request",
					"error":   err.Error(),
				})
			return
		}

		filename := file.Filename
		filetype := strings.Split(file.Header.Get("Content-Type"), "/")[1]

		if filetype != "markdown" {
			c.JSON(http.StatusBadRequest,
				gin.H{
					"message": "invalid file type. API supports only markdown files `.md`",
				})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":  "SpellCheck complete from controller",
			"filename": filename,
			"filetype": filetype,
		})
	}
}
