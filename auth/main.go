package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/saddagada1/oz-auth/controllers"
	"github.com/saddagada1/oz-auth/middleware"
	"github.com/saddagada1/oz-auth/utils"
)

func init() {
	utils.LoadSecrets()
	utils.ConnectToDB()

	if os.Getenv("ENV") == "PROD" {
		utils.SyncDB()
	}
}

func main() {
	r := gin.Default()

	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.GET("/validate", middleware.Auth, controllers.Validate)

	r.Run()
}
