package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func UserOpenTest(c *gin.Context) {
	godotenv.Load()
	c.JSON(http.StatusOK, os.Getenv("VERSION"))
}
