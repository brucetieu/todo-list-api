package main

import (
	"os"

	"github.com/brucetieu/todo-list-api/database"
	"github.com/brucetieu/todo-list-api/representations"
	"github.com/brucetieu/todo-list-api/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	LoadConfig()

	pgVars := representations.PGVars{
		PgHost: os.Getenv("POSTGRES_HOST"),
		PgUser: os.Getenv("POSTGRES_USER"),
		PgPort: os.Getenv("POSTGRES_PORT"),
		PgPass: os.Getenv("POSTGRES_PASS"),
		PgDbName: os.Getenv("POSTGRES_DB"),
	}

	db, err := database.ConnectToDB(pgVars)
	if err != nil {
		log.Fatal("Failed connecting to database: " + err.Error())
	}

	router := gin.Default()
	routes.InitializeRoutes(router, db)

	_ = router.Run(":" + os.Getenv("PORT"))

}