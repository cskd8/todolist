package service

import (
	"encoding/json"
	"fmt"
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

	// Get task data
	var task database.Task
	err = db.Get(&task, "SELECT * FROM tasks WHERE id = ? and user_id = ?", ctx.Param("id"), loginInfo.ID)
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

	ctx.HTML(http.StatusOK, "taskedit.html", gin.H{"Title": "Edit", "id": ctx.Param("id"), "Categories": categories, "Name": task.Title})
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
	expires := ctx.PostForm("expires")

	if categoryID == "" {
		// Insert task
		db.MustExec("INSERT INTO tasks (title, user_id, expires) VALUES (?, ?, ?)", task.Title, loginInfo.ID, expires)
	} else {
		// Insert task
		db.MustExec("INSERT INTO tasks (title, category_id, user_id, expires) VALUES (?, ?, ?, ?)", task.Title, categoryID, loginInfo.ID, expires)
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
	expires := ctx.PostForm("expires")
	fmt.Println(expires)
	if categoryID == "" {
		// Update task
		_, err = db.Exec("UPDATE tasks SET title = ?, expires = ? WHERE id = ? AND user_id = ?", task.Title, expires, ctx.Param("id"), loginInfo.ID)
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}
		// Render task
		ctx.Redirect(http.StatusSeeOther, "/")
		return
	}

	// Update task
	_, err = db.Exec("UPDATE tasks SET title = ?, category_id = ?, expires = ? WHERE id = ? AND user_id = ?", task.Title, categoryID, expires, ctx.Param("id"), loginInfo.ID)
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
