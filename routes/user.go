package routes

import (
	"github.com/SanjaySinghRajpoot/remote-crawler/controller"
	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.Engine) {
	router.GET("/", controller.UserController)
}
