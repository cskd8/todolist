package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/koron/go-dproxy"
	database "todolist.go/db"
)

// Max returns the larger of x or y.
func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

// Min returns the smaller of x or y.
func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

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

	limit := 10
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		page = 1
	}
	fmt.Println(page)

	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Get tasks (and category if it isn't null) in DB
	var tasks []TaskCategory
	err = db.Select(&tasks, "SELECT tasks.*, categories.name FROM tasks LEFT JOIN categories ON tasks.category_id = categories.id WHERE tasks.user_id = ? ORDER BY tasks.created_at DESC LIMIT ? OFFSET ?", loginInfo.ID, limit, (page-1)*limit) // Use DB#Select for multiple entries
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

	// Get count of tasks
	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM tasks WHERE user_id = ?", loginInfo.ID)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Calc page count
	pageCount := count / limit
	if count%limit != 0 {
		pageCount++
	}

	pages := make([]struct{ Number int }, pageCount)
	for i := 0; i < pageCount; i++ {
		pages[i] = struct{ Number int }{i + 1}
	}

	// Render tasks
	ctx.HTML(http.StatusOK, "index.html", gin.H{"Title": "HOME", "Tasks": tasks, "Categories": categories, "User": loginInfo, "Pages": pages, "Page": page, "Prev": page - 1, "Next": page + 1, "PageCount": pageCount, "Url": "/"})
}
