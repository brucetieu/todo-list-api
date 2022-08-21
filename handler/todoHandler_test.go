package handler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brucetieu/todo-list-api/representations"
	"github.com/brucetieu/todo-list-api/tests"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"
)


func Test_CreateTodoList(t *testing.T) {
	db, err := tests.ConnectToDB()
	if err != nil {
		panic(err)
	}

	handler := NewTodoHandler(db)
	router := gin.Default()

	todos := []representations.Todo{
		{
			Title: "test item 1",
		},
		{
			Title: "test item 2",
		},
	}

	values := map[string]interface{}{"title": "test todo list title", "todos": todos}
	data, err := json.Marshal(values)
	if err != nil {
		log.Fatal(err)
	}

	var todoList representations.TodoList
	req, _ := http.NewRequest("POST", "/api/todolists", bytes.NewBuffer(data))
	router.POST("/api/todolists", handler.CreateTodoList)

	defer clear(db)

	w := httptest.NewRecorder()
	res := w.Result()
    defer res.Body.Close()

    body, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(body, &todoList)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.NotNil(t, w.Body.String())
}

func clear(db *gorm.DB) {
	db.Delete(&representations.TodoList{})
	db.Delete(&representations.Todo{})
}