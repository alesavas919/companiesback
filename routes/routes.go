package routes

import (
	"companies/analytic"
	"companies/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	routeBegin := gin.Default()

	{
		companiesApi := routeBegin.Group("/api/companies")
		companiesApi.GET("/", controllers.CompaniesReqGetAllDataFromDatabase)
		companiesApi.GET("/req", controllers.CompaniesReqGetAllDataFromRequest)
		companiesApi.GET("/load", controllers.CompaniesReqLoadAllDataFromRequest)
	}

	{
		apiusers := routeBegin.Group("/api/users")
		apiusers.POST("/", controllers.UserCreateUser)
		apiusers.GET("/", controllers.UserOpenTest)
		//apiusers.GET("/sc", controllers.UserOpenSecretTest)
		apiusers.GET("/users", controllers.UserGetUsers)
		apiusers.GET("/:id", controllers.UserGetUserByID)
	}

	{
		apiCalculateDate := routeBegin.Group("/api/calculate")
		apiCalculateDate.GET("/", analytic.AnalyticCaculatedResponse)
	}
	/*
		https://es.investing.com/indices/nyse-composite-historical-data
	*/

	return routeBegin
}
