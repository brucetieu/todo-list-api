package tests

import (
	"fmt"

	"github.com/brucetieu/todo-list-api/representations"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)


func ConnectToDB() (*gorm.DB, error) {	
	dbURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", "localhost", "5432", "postgres", "testdb", "admin")
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
