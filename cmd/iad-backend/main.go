package main

import (
	_ "iad-backend/docs"
	"iad-backend/internal/api"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// @title           IAD Backend API
// @version         1.0
// @description     This is the backend API for IAD (Information and Analytical Data) system.

// @contact.name   API Support
// @contact.url    https://github.com/your-org/iad-backend
// @contact.email  support@iad.example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

func main() {
	err := godotenv.Load("deploy/.env")
	if err != nil {
		panic(err)
	}

	logrus.SetLevel(logrus.ErrorLevel)
	api.StartServer()
}
