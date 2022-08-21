package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/brucetieu/todo-list-api/representations"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
)

type todoHandler struct {
	db *gorm.DB
}

func NewTodoHandler(db *gorm.DB) *todoHandler {
	return &todoHandler{
		db: db,
	}
}

// Wrap HTTPError
func NewError(ctx *gin.Context, status int, err error) {
	er := HTTPError{
		Code:    status,
		Message: err.Error(),
	}
	ctx.JSON(status, er)
}

// Wrap ctx.JSON
func JSON(ctx *gin.Context, code int, obj interface{}) {
    ctx.Header("Content-Type", "application/json")
    ctx.JSON(code, obj)
}

func newFalse() *bool {
    b := false
    return &b
}

type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}

// Create a single todo list
func (t *todoHandler) CreateTodoList(ctx *gin.Context) {
	// Validate input
	var todoList representations.TodoList
	if err := ctx.ShouldBindJSON(&todoList); err != nil {
		log.WithField("error", err.Error()).Error("Error validating input")
		NewError(ctx, http.StatusBadRequest, err)
		return
	}

	if todoList.Title == "" {
		NewError(ctx, http.StatusBadRequest, errors.New("todo list must have a 'title' property"))
		return
	}

	id := uuid.Must(uuid.NewRandom()).String()

	if len(todoList.Todos) != 0 {
		for i := 0; i < len(todoList.Todos); i++ {
			if todoList.Todos[i].Complete == nil {
				todoList.Todos[i].Complete = newFalse()
			}
			todoList.Todos[i].ItemID = uuid.Must(uuid.NewRandom()).String()
			todoList.Todos[i].TodoListID = id
		}
	}

	todoList.ID = id

	dt := time.Now()

	todoList.CreatedAt = dt
	todoList.UpdatedAt = dt

	if err := t.db.Create(&todoList).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i := 0; i < len(todoList.Todos); i++ {
		todoList.Todos[i].TodoListID = ""
	}
	
	JSON(ctx, http.StatusCreated, todoList)
}

// Get all todo lists
func (t *todoHandler) GetTodoLists(ctx *gin.Context) {
	var todoList []representations.TodoList

	if err := t.db.Preload("Todos").Find(&todoList).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i := 0; i < len(todoList); i++ {
		for j := 0; j < len(todoList[i].Todos); j++ {
			todoList[i].Todos[j].TodoListID = ""
		}
	}

	JSON(ctx, http.StatusOK, todoList)

}

// Get a single todo list by its id
func (t *todoHandler) GetTodoList(ctx *gin.Context) {
	id := ctx.Param("id")

	var todoList representations.TodoList

	res := t.db.
		// Set("gorm:auto_preload", true).
		Preload("Todos").
		Where("id = ?", id).
		First(&todoList)
	if res.Error != nil {
		NewError(ctx, http.StatusNotFound, errors.New(res.Error.Error()))
		return
	}


	for i := 0; i < len(todoList.Todos); i++ {
		todoList.Todos[i].TodoListID = ""
	}

	JSON(ctx, http.StatusOK, todoList)

}

// Update a single todo list by its id
func (t *todoHandler) UpdateTodoList(ctx *gin.Context) {
	id := ctx.Param("id")

	var existingList representations.TodoList

	// Grab existing entry
	res := t.db.
		Preload("Todos").
		Where("id = ?", id).
		First(&existingList)
	if res.Error != nil {
		NewError(ctx, http.StatusNotFound, errors.New(res.Error.Error()))
		return
	}

	// Validate input
	var updatedTodoList representations.TodoList
	if err := ctx.ShouldBindJSON(&updatedTodoList); err != nil {
		log.WithField("error", err.Error()).Error("Error validating input")
		NewError(ctx, http.StatusBadRequest, err)
		return
	}

	tx := t.db.Begin()

	// Update fields
	var todoList representations.TodoList
	todoList.Title = updatedTodoList.Title
	todoList.CreatedAt = existingList.CreatedAt
	todoList.UpdatedAt = time.Now()
	todoList.ID = existingList.ID

	// Perform update
	if err := tx.Save(&todoList).Error; err != nil {
		tx.Rollback()
		NewError(ctx, http.StatusInternalServerError, errors.New(err.Error()))
		return
	}

	if errCommit := tx.Commit().Error; errCommit != nil {
		NewError(ctx, http.StatusInternalServerError, errors.New(errCommit.Error()))
		return
	}

	// Delete previous todos (to replace with new ones)
	if errDel := t.db.Where("todo_list_id = ?", existingList.ID).Delete(representations.Todo{}).Error; errDel != nil {
		tx.Rollback()
		NewError(ctx, http.StatusInternalServerError, errors.New(errDel.Error()))
		return
	}

	if len(updatedTodoList.Todos) != 0 {
		todoList.Todos = make([]representations.Todo, len(updatedTodoList.Todos))
		for i, updatedTodo := range updatedTodoList.Todos {
			todoList.Todos[i] = representations.Todo{
				ItemID: uuid.Must(uuid.NewRandom()).String(),
				TodoListID: existingList.ID,
				Title: updatedTodo.Title,
				Complete: updatedTodo.Complete,
			}

			if todoList.Todos[i].Complete == nil {
				todoList.Todos[i].Complete = newFalse()
			}
		}

		for i := 0; i < len(todoList.Todos); i++ {
			if err := t.db.Save(&todoList.Todos[i]).Error; err != nil {
				NewError(ctx, http.StatusInternalServerError, errors.New(err.Error()))
				return
			}
		}
	} else {
		todoList.Todos = []representations.Todo{}
	}


	for i := 0; i < len(todoList.Todos); i++ {
		todoList.Todos[i].TodoListID = ""
	}

	JSON(ctx, http.StatusOK, todoList)
}

// Delete a single todo list by its id
func (t *todoHandler) DeleteTodoList(ctx *gin.Context) {
	id := ctx.Param("id")

	var existingList representations.TodoList

	// Grab existing entry
	res := t.db.
		Preload("Todos").
		Where("id = ?", id).
		First(&existingList)
	if res.Error != nil {
		NewError(ctx, http.StatusNotFound, errors.New(res.Error.Error()))
		return
	}

	// If there are any Todos, delete those first
	// TODO: figure out if there's cascading delete option for GORM
	for i := 0; i < len(existingList.Todos); i++ {
		_ = t.db.Delete(&existingList.Todos[i])
	}

	// Delete the todo list
	res = t.db.Delete(&existingList)
	if res.Error != nil {
		NewError(ctx, http.StatusInternalServerError, errors.New(res.Error.Error()))
		return
	}

	ctx.Status(http.StatusNoContent)
}

// Create a todo for a single todo list
func (t *todoHandler) CreateTodo(ctx *gin.Context) {
	id := ctx.Param("id")

	var existingList representations.TodoList

	// Check there is a list associated with this id
	res := t.db.
		Where("id = ?", id).
		First(&existingList)
	if res.Error != nil {
		NewError(ctx, http.StatusNotFound, errors.New(res.Error.Error()))
		return
	}

	// Validate input
	var todo representations.Todo
	if err := ctx.ShouldBindJSON(&todo); err != nil {
		log.WithField("error", err.Error()).Error("Error validating input")
		NewError(ctx, http.StatusBadRequest, err)
		return
	}
	if todo.Title == "" {
		NewError(ctx, http.StatusBadRequest, errors.New("todo must have a 'title' property"))
		return
	}

	todo.TodoListID = id
	todo.ItemID = uuid.Must(uuid.NewRandom()).String()

	// Persist
	if err := t.db.Create(&todo).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	JSON(ctx, http.StatusCreated, todo)
}

// Get all todos in a todo list
func (t *todoHandler) GetTodos(ctx *gin.Context) {
	id := ctx.Param("id")

	var todoList representations.TodoList

	// Query todolist
	res := t.db.
		Preload("Todos").
		Where("id = ?", id).
		First(&todoList)
	if res.Error != nil {
		NewError(ctx, http.StatusNotFound, errors.New(res.Error.Error()))
		return
	}

	// Clear TodoListID for cleaner output
	for i := 0; i < len(todoList.Todos); i++ {
		todoList.Todos[i].TodoListID = ""
	}

	JSON(ctx, http.StatusOK, todoList.Todos)
}

// Get a single todo in a todo list
func (t *todoHandler) GetTodo(ctx *gin.Context) {
	id := ctx.Param("id")
	todoId := ctx.Param("todoId")

	var todo representations.Todo

	// Query for todo
	res := t.db.
		Where("item_id = ? AND todo_list_id = ?", todoId, id).
		First(&todo)
	if res.Error != nil {
		NewError(ctx, http.StatusNotFound, errors.New(res.Error.Error()))
		return
	}

	todo.TodoListID = ""

	JSON(ctx, http.StatusOK, todo)
}

// Delete a single todo in a todo list
func (t *todoHandler) DeleteTodo(ctx *gin.Context) {
	id := ctx.Param("id")
	todoId := ctx.Param("todoId")

	var todo representations.Todo
	
	// Check todo exists
	res := t.db.
	Where("item_id = ? AND todo_list_id = ?", todoId, id).
	First(&todo)
	if res.Error != nil {
		NewError(ctx, http.StatusNotFound, errors.New(res.Error.Error()))
		return
	}

	// Delete todo
	res = t.db.
		Where("item_id = ? AND todo_list_id = ?", todoId, id).
		Delete(&todo)
	if res.Error != nil {
		NewError(ctx, http.StatusNotFound, errors.New(res.Error.Error()))
		return
	}

	ctx.Status(http.StatusNoContent)
}

// Update a single todo in a todolist
func (t *todoHandler) UpdateTodo(ctx *gin.Context) {
	id := ctx.Param("id")
	todoId := ctx.Param("todoId")

	var existingList representations.TodoList
	var todo representations.Todo

	// Grab existing list where todo belongs
	res := t.db.
		Where("id = ?", id).
		First(&existingList)
	if res.Error != nil {
		NewError(ctx, http.StatusNotFound, errors.New(res.Error.Error()))
		return
	}

	// Get existing todo within list
	res = t.db.
		Where("item_id = ? AND todo_list_id = ?", todoId, id).
		First(&todo)
	if res.Error != nil {
		NewError(ctx, http.StatusNotFound, errors.New(res.Error.Error()))
		return
	}

	// Validate input
	if err := ctx.ShouldBindJSON(&todo); err != nil {
		log.WithField("error", err.Error()).Error("Error validating input")
		NewError(ctx, http.StatusBadRequest, err)
		return
	}

	var updatedTodo representations.Todo
	updatedTodo.ItemID = uuid.Must(uuid.NewRandom()).String()
	updatedTodo.TodoListID = existingList.ID
	updatedTodo.Title = todo.Title

	if todo.Complete == nil {
		updatedTodo.Complete = newFalse()
	} else {
		updatedTodo.Complete = todo.Complete
	}

	// Delete old todo
	if err := t.db.Delete(&todo).Error; err != nil {
		NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	// Save new todo
	if err := t.db.Save(&updatedTodo).Error; err != nil {
		NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	
	updatedTodo.TodoListID = ""

	// Update timestamp of todo list
	if err := t.db.Model(&representations.TodoList{}).
			Where("id = ?", id).
			Update("updated_at", time.Now()).Error; err != nil {
				NewError(ctx, http.StatusInternalServerError, err)
	}

	JSON(ctx, http.StatusOK, updatedTodo)
}