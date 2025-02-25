package routes

import (
	"go-markdown-parser/controller"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine) {
	router.POST("/auth/v1/signup", controller.SignUpController())
	router.POST("/auth/v1/login", controller.LogIn())

	// route to authenticate user
	router.GET("/auth/v1/authenticate", controller.AuthenticateUser())
}
