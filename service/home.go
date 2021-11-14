package service

import (
	"encoding/json"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/koron/go-dproxy"
	database "todolist.go/db"
)

// Home renders index.html
func Home(ctx *gin.Context) {
	session := sessions.Default(ctx)
	loginUserJson, err := dproxy.New(session.Get("user")).String()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Get User from session

	var loginInfo database.User
	err = json.Unmarshal([]byte(loginUserJson), &loginInfo)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Get tasks (and category if it isn't null) in DB
	var tasks []TaskCategory
	err = db.Select(&tasks, "SELECT tasks.*, categories.name FROM tasks LEFT JOIN categories ON tasks.category_id = categories.id WHERE tasks.user_id = ? ORDER BY tasks.created_at DESC", loginInfo.ID) // Use DB#Select for multiple entries
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	var categories []database.Category
	err = db.Select(&categories, "SELECT * FROM categories") // Use DB#Select for multiple entries
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Render tasks
	ctx.HTML(http.StatusOK, "index.html", gin.H{"Title": "HOME", "Tasks": tasks, "Categories": categories, "User": loginInfo})
}
