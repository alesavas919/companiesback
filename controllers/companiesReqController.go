package controllers

import (
	"net/http"

	"companies/service"

	"github.com/gin-gonic/gin"
)

func CompaniesReqGetAllDataFromDatabase(c *gin.Context) {
	companies := service.CompaniesReqGetAllDataFromDatabaseService()
	c.JSON(http.StatusAccepted, companies)
}

// ///////////////////////////////////////// GET FROM PETITION ///////////////////////////////////////////
func CompaniesReqGetAllDataFromRequest(c *gin.Context) {
	resData := service.CompaniesReqGetAllDataFromRequestService()
	c.String(http.StatusCreated, string([]byte(resData)))
}

func CompaniesReqLoadAllDataFromRequest(c *gin.Context) {
	resData, err := service.CompaniesReqLoadAllDataFromRequestService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.String(http.StatusCreated, string([]byte(resData))) //{created:'created'}
}
