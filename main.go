package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"todolist.go/db"
	"todolist.go/service"
)

const port = 8000

var SessionSecret string
var SessionName string = "sessions"

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d/%02d/%02d", year, month, day)
}

func main() {
	// initialize DB connection
	dsn := db.DefaultDSN(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))
	if err := db.Connect(dsn); err != nil {
		log.Fatal(err)
	}

	// initialize Gin engine
	engine := gin.Default()
	engine.SetFuncMap(template.FuncMap{
		"formatAsDate": formatAsDate,
	})
	engine.LoadHTMLGlob("views/*.html")

	// set session stores
	store := cookie.NewStore([]byte("secret"))
	engine.Use(sessions.Sessions(SessionName, store))

	// routing
	engine.Static("/assets", "./assets")
	engine.GET("/login", service.RLogin)
	engine.POST("/login", service.Login)
	engine.GET("/signup", service.RSignup)
	engine.POST("/signup", service.Signup)
	authGroup := engine.Group("/")
	authGroup.Use(service.LoginCheckMiddleware())
	engine.GET("/", service.Home)
	engine.GET("/task", service.RGetTask)
	engine.POST("/task", service.PostTask)
	engine.POST("/category", service.PostCategory)
	engine.GET("/list", service.TaskList)
	engine.POST("/list/search", service.Search)
	engine.GET("/task/:id", service.ShowTask)           // ":id" is a parameter
	engine.POST("/task/:id/finish", service.FinishTask) // ":id" is a parameter
	engine.POST("/task/:id/resume", service.ResumeTask) // ":id" is a parameter
	engine.GET("/task/:id/edit", service.REditTask)     // ":id" is a parameter
	engine.POST("/task/:id/edit", service.PutTask)      // ":id" is a parameter
	engine.GET("/task/:id/delete", service.DeleteTask)  // ":id" is a parameter
	engine.GET("/user/edit", service.REditUser)
	engine.POST("/user/edit", service.EditUser)
	engine.POST("/leave", service.Leave)
	engine.POST("/logout", service.Logout)

	// start server
	engine.Run(fmt.Sprintf(":%d", port))
}
