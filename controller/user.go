package controller

import (
	"github.com/SanjaySinghRajpoot/remote-crawler/config"
	"github.com/SanjaySinghRajpoot/remote-crawler/models"
	"github.com/gin-gonic/gin"
)

func UserController(c *gin.Context) {

	// user := make([]models.Users, 0)

	jobs := []models.Job{}
	config.DB.Find(&jobs)

	c.JSON(200, &jobs)
}

func TestUserController(c *gin.Context) {

	c.String(200, "status working")
}
