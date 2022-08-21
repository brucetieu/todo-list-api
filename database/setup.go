package database

import (
	"fmt"

	"github.com/brucetieu/todo-list-api/representations"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)


func ConnectToDB(pgVars representations.PGVars) (*gorm.DB, error) {	
	dbURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", pgVars.PgHost, pgVars.PgPort, pgVars.PgUser, pgVars.PgDbName, pgVars.PgPass)
	log.Info("dbUrl: ", dbURL)

	database, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return database, err
	}

	database.Logger.LogMode(logger.Info)

	_ = database.AutoMigrate(&representations.TodoList{})
	_ = database.AutoMigrate(&representations.Todo{})

	return database, nil
}
