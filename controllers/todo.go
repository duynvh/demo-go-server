package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"go-server-demo/lib/common"
	"go-server-demo/models"
	"net/http"
	"os"
	"strconv"
)

type Todo = models.Todo

// JSON type alias
type JSON = common.JSON
type TodoController struct {}

func (ctrl TodoController) Create(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	if c.PostForm("title") == "" || c.PostForm("completed") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Please input all field"})
		return
	}

	completed, _ := strconv.Atoi(c.PostForm("completed"))
	user := c.MustGet("user").(models.User)
	todo := Todo{Title: c.PostForm("title"), Completed: completed, UserID: user.ID}
	db.NewRecord(todo)
	db.Create(&todo)
	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Todo item created successfully!", "resourceId": todo.ID})
}

func (ctrl TodoController) FetchAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	offset := (page - 1) * 2;
	limit := os.Getenv("DB_LIMIT")
	db := c.MustGet("db").(*gorm.DB)
	var todos []Todo

	db.Offset(offset).Limit(limit).Find(&todos)

	if len(todos) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}

	length := len(todos)
	serialized := make([]JSON, length, length)

	for i := 0; i < length; i++ {
		serialized[i] = todos[i].Serialize()
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": serialized, "page": page})
}

func (ctrl TodoController) Update(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")
	user := c.MustGet("user").(models.User)
	if c.PostForm("title") == "" || c.PostForm("completed") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Please input all field"})
		return
	}
	var todo Todo
	if err := db.Preload("User").Where("id = ?", id).First(&todo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Not Found"})
		return
	}

	if todo.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"status": http.StatusForbidden, "message": "Not Forbidden!"})
		return
	}

	completed, _ := strconv.Atoi(c.PostForm("completed"))
	todo.Title = c.PostForm("title")
	todo.Completed = completed
	db.Save(&todo)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo updated successfully!"})
}

func (ctrl TodoController) Delete(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")
	user := c.MustGet("user").(models.User)

	var todo Todo
	if err := db.Where("id = ?", id).First(&todo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Not Found"})
		return
	}

	if todo.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"status": http.StatusForbidden, "message": "Not Forbidden!"})
		return
	}

	db.Delete(&todo)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo deleted successfully!"})
}
