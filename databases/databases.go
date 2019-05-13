package databases

import (
	"demo-go-server/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
)

func Init() (*gorm.DB, error){
	db, err := gorm.Open("mysql", os.Getenv("DB_CONFIG"))
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Todo{})

	return db, nil
}
