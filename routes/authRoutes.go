package routes

import (
	"github.com/gin-gonic/gin"
	"main.go/controller"
)

func AuthRoutes(router *gin.Engine) {
	router.POST("/auth/v1/signup", controller.SignUpController())
	router.POST("/auth/v1/login", controller.LogIn())

	// route to authenticate user
	router.GET("/auth/v1/authenticate", controller.AuthenticateUser())
}
