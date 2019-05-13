package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"demo-go-server/controllers"
	"demo-go-server/databases"
	"demo-go-server/lib/middlewares"
)

//CORSMiddleware ...
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Origin, Authorization, Accept, Client-Security-Token, Accept-Encoding, x-access-token")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			fmt.Println("OPTIONS")
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

func main() {
	// load .env environment variables
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	r.Use(CORSMiddleware())
	db, _ := databases.Init()

	r.Use(databases.Inject(db))
	r.Use(middlewares.JWTMiddleware())

	v1 := r.Group("api/v1")
	{
		/*** START USER ***/
		user := new(controllers.UserController)

		v1.POST("/user/register", user.Register)
		v1.POST("/user/login", user.Login)

		todo := new(controllers.TodoController)
		v1.GET("/todo", todo.FetchAll)
		v1.POST("/todo", middlewares.Authorized, todo.Create)
		v1.PATCH("/todo/:id", middlewares.Authorized, todo.Update)
		v1.DELETE("/todo/:id", middlewares.Authorized, todo.Delete)
	}

	r.Run()
}