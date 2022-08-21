package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	log "github.com/sirupsen/logrus"
	"github.com/brucetieu/todo-list-api/handler"
)

func GetHealth(ctx *gin.Context) {
	log.Info("Checking health of API")
	ctx.JSON(http.StatusOK, "Healthy!")	
} 

func InitializeRoutes(router *gin.Engine, db *gorm.DB) {
	groupRoute := router.Group("/api")

	todoHandler := handler.NewTodoHandler(db)

	// Health check
	groupRoute.GET("/health", GetHealth)
	
	// Todo list methods
	groupRoute.POST("/todolists", todoHandler.CreateTodoList)
	groupRoute.GET("/todolists", todoHandler.GetTodoLists)
	groupRoute.GET("/todolists/:id", todoHandler.GetTodoList)
	groupRoute.PUT("/todolists/:id", todoHandler.UpdateTodoList)
	groupRoute.DELETE("/todolists/:id", todoHandler.DeleteTodoList)

	// Todos inside of todo list methods
	groupRoute.POST("/todolists/:id/todos", todoHandler.CreateTodo)
	groupRoute.GET("/todolists/:id/todos", todoHandler.GetTodos)
	groupRoute.GET("/todolists/:id/todos/:todoId", todoHandler.GetTodo)
	groupRoute.DELETE("/todolists/:id/todos/:todoId", todoHandler.DeleteTodo)
	groupRoute.PUT("/todolists/:id/todos/:todoId", todoHandler.UpdateTodo)
}