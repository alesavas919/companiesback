package routes

import (
	"companies/analytic"
	"companies/controllers"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	routeBegin := gin.Default()

	// Configuración CORS
	config := cors.Config{
		AllowOrigins:     []string{"*"}, //NOT RECOMMENDED
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false, //NOT RECOMMENDED TRUE -> AllowOrigins != "*"
		//MaxAge:           12 * time.Hour,
		AllowOriginFunc: func(origin string) bool {
			return strings.HasPrefix(origin, "http://")
			//return true // ALLOW ANY ORIGIN
			// ONLY PROD: ADD VALIDTAION LOGIC HERE
			// allowedOrigins := []string{
			//       "https://dom1.com",
			//       "http://dom2:5173",
			//       "https://dom3.com",
			//   }
			//       for _, o := range allowedOrigins {
			//       if origin == o {
			//           return true
			//       }
			//   }
			//   return false
		},
	}

	routeBegin.Use(cors.New(config))

	{
		companiesApi := routeBegin.Group("/api/companies")
		companiesApi.GET("/", controllers.CompaniesReqGetAllDataFromDatabase)
		companiesApi.GET("", controllers.CompaniesReqGetAllDataFromDatabase)
		companiesApi.GET("/req", controllers.CompaniesReqGetAllDataFromRequest)
		companiesApi.GET("/load", controllers.CompaniesReqLoadAllDataFromRequest)
	}

	{
		apiusers := routeBegin.Group("/api/users")
		apiusers.GET("/", controllers.UserOpenTest)
		apiusers.GET("", controllers.UserOpenTest)
		//apiusers.GET("/sc", controllers.UserOpenSecretTest)
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
