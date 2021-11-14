package service

import (
	"encoding/json"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/koron/go-dproxy"
	database "todolist.go/db"
)

// RGetTask renders task.html
func RGetTask(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "task.html", gin.H{})
}

// REditTask renders taskedit.html
func REditTask(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Get Categories
	var categories []database.Category
	err = db.Select(&categories, "SELECT * FROM categories")
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.HTML(http.StatusOK, "taskedit.html", gin.H{"Title": "Edit", "id": ctx.Param("id"), "Categories": categories})
}

// PostTask processes a POST request for a task
func PostTask(ctx *gin.Context) {
	session := sessions.Default(ctx)
	loginUserJson, err := dproxy.New(session.Get("user")).String()
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

	// Get User from session

	var loginInfo database.User
	err = json.Unmarshal([]byte(loginUserJson), &loginInfo)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Get task data
	var task database.Task
	task.Title = ctx.PostForm("title")
	categoryID := ctx.PostForm("category")

	if categoryID == "null" {
		// Insert task
		db.MustExec("INSERT INTO tasks (title, user_id) VALUES (?, ?)", task.Title, loginInfo.ID)
	} else {
		// Insert task
		db.MustExec("INSERT INTO tasks (title, category_id, user_id) VALUES (?, ?, ?)", task.Title, categoryID, loginInfo.ID)
	}

	// Render task
	ctx.Redirect(http.StatusSeeOther, "/")
}

// PutTask processes a PUT request for a task
func PutTask(ctx *gin.Context) {
	session := sessions.Default(ctx)
	loginUserJson, err := dproxy.New(session.Get("user")).String()
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

	// Get task data
	var task database.Task
	task.Title = ctx.PostForm("title")
	categoryID := ctx.PostForm("category")

	// Get User from session

	var loginInfo database.User
	err = json.Unmarshal([]byte(loginUserJson), &loginInfo)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Update task
	_, err = db.Exec("UPDATE tasks SET title = ?, category_id = ? WHERE id = ? AND user_id = ?", task.Title, categoryID, ctx.Param("id"), loginInfo.ID)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Render task
	ctx.Redirect(http.StatusSeeOther, "/")
}

// DeleteTask processes a DELETE request for a task
func DeleteTask(ctx *gin.Context) {
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

	// Delete task
	db.MustExec("DELETE FROM tasks WHERE id = ? AND user_id = ?", ctx.Param("id"), loginInfo.ID)

	// Render task
	ctx.Redirect(http.StatusSeeOther, "/")
}
