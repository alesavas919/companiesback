package main

import (
	"log"
	"os"

	"companies/database"
	"companies/routes"
	"companies/security"

	"github.com/joho/godotenv"
)

func main() {
	// URL CONNECT DATABASE
	godotenv.Load()
	var connString string = os.Getenv("DB_CONNECTION_STRING")
	if connString == "" {
		//connString = security.ResourceSecurityData("DB_CONNECTION_STRING")
		connString, _ = security.SecretString(0)
	}

	// START DATABASE
	database.InitDB(connString)
	defer database.CloseDB()

	// ROUTES
	r := routes.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Server running on port %s", port)
	r.SetTrustedProxies(nil)
	r.Run(":" + port)
}
