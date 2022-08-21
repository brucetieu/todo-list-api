package representations

import (
	"time"
)

type TodoList struct {
	ID string  `gorm:"primary_key;type:char(36)" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Title string `json:"title"`
	Todos []Todo `json:"todos" gorm:"foreignKey:TodoListID;ssociation_foreignkey:ID"`
}

type Todo struct {
	ItemID string `gorm:"primary_key;type:char(36)" json:"itemID"`
	TodoListID string `json:"todoListID,omitempty"`
	Title string `json:"title"`
	Complete *bool `json:"complete"`
}