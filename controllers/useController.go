package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserOpenTest(c *gin.Context) {
	c.JSON(http.StatusOK, "V1.6.1")
}
