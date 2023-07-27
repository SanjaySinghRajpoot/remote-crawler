package controller

import (
	"github.com/SanjaySinghRajpoot/remote-crawler/config"
	"github.com/SanjaySinghRajpoot/remote-crawler/models"
	"github.com/gin-gonic/gin"
)

func UserController(c *gin.Context) {

	// user := make([]models.Users, 0)

	users := []models.User{}
	config.DB.Find(&users)

	c.JSON(200, &users)
}

func PostUserController(c *gin.Context) {
	c.String(200, "hello world")
}
