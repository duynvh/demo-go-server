package models

import (
	"github.com/jinzhu/gorm"
	"go-server-demo/lib/common"
)

type Todo struct {
	gorm.Model
	Title string `json:"title"`
	Completed int `json:"completed"`
	User   User   `gorm:"foreignkey:UserID"`
	UserID uint
}

// Serialize serializes post data
func (todo Todo) Serialize() common.JSON {
	return common.JSON{
		"id":         todo.ID,
		"title":      todo.Title,
		"completed":      todo.Completed,
		"user_id":       todo.UserID,
		"created_at": todo.CreatedAt,
	}
}