package controllers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"go-server-demo/lib/common"
	"go-server-demo/models"
	"io/ioutil"
	"net/http"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

type User = models.User 
type UserController struct {}

func hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func checkHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateToken(data common.JSON) (string, error) {

	//  token is valid for 7days
	date := time.Now().Add(time.Hour * 24)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": data,
		"exp":  date.Unix(),
	})

	// get path from root dir
	pwd, _ := os.Getwd()
	keyPath := pwd + "/jwtsecret.key"

	key, readErr := ioutil.ReadFile(keyPath)
	if readErr != nil {
		return "", readErr
	}
	tokenString, err := token.SignedString(key)
	return tokenString, err
}

func (ctrl UserController) Register(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	if c.PostForm("email") == "" || c.PostForm("password") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Please input all field"})
		return
	}

	// check existancy
	var exists User
	if err := db.Where("email = ?", c.PostForm("email")).First(&exists).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"status": http.StatusConflict, "message": "User is exist"})
		return
	}

	hash, hashErr := hash(c.PostForm("password"))
	if hashErr != nil {
		c.AbortWithStatus(500)
		return
	}
	user := User{Email: c.PostForm("email"), Password: hash}
	db.Save(&user)

	// Generate token
	serialized := user.Serialize()
	token, _ := generateToken(serialized)
	c.SetCookie("token", token, 60*60*24, "/", "", false, true)

	c.JSON(http.StatusCreated, common.JSON{
		"user":  user.Serialize(),
		"token": token,
	})
	//c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "User created successfully!"})
}

func (ctrl UserController) Login(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	if c.PostForm("email") == "" || c.PostForm("password") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Please input all field"})
		return
	}

	// check existancy
	var user User
	if err := db.Where("email = ?", c.PostForm("email")).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "User not found"})
		return
	}

	if !checkHash(c.PostForm("password"), user.Password ) {
		c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Password is not correct"})
		return
	}

	// Generate token
	serialized := user.Serialize()
	token, _ := generateToken(serialized)
	c.SetCookie("token", token, 60*60*24, "/", "", false, true)

	c.JSON(http.StatusCreated, common.JSON{
		"user":  user.Serialize(),
		"token": token,
	})
}
