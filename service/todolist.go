package service

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/koron/go-dproxy"
	database "todolist.go/db"
)

// TaskList renders list of tasks in DB
func TaskList(ctx *gin.Context) {
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

	// Get tasks in DB
	var tasks []database.Task
	err = db.Select(&tasks, "SELECT * FROM tasks WHERE user_id = ?", loginInfo.ID) // Use DB#Select for multiple entries
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Render tasks
	ctx.HTML(http.StatusOK, "task_list.html", gin.H{"Title": "Task list", "Tasks": tasks})
}

// ShowTask renders a task with given ID
func ShowTask(ctx *gin.Context) {
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

	// parse ID given as a parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	// Get a task and category with given ID
	var task TaskCategory
	err = db.Get(&task, "SELECT tasks.*, categories.name FROM tasks LEFT JOIN categories ON tasks.category_id = categories.id WHERE tasks.id=? AND tasks.user_id=?", id, loginInfo.ID) // Use DB#Get for one entry
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	task.Remaining = task.Expires.Unix() - time.Now().Unix()

	// Render task
	ctx.HTML(http.StatusOK, "task.html", gin.H{"Title": task.Title, "Task": task})
}

// FinishTask processes a POST request for a task
func FinishTask(ctx *gin.Context) {
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

	// parse ID given as a parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	// Update task with given ID
	_, err = db.Exec("UPDATE tasks SET is_done=1 WHERE id=? AND user_id = ?", id, loginInfo.ID) // Use DB#Exec for one entry
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Redirect to task list
	ctx.Redirect(http.StatusFound, "/")
}

// ResumeTask processes a POST request for a task
func ResumeTask(ctx *gin.Context) {
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

	// parse ID given as a parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	// Update task with given ID
	_, err = db.Exec("UPDATE tasks SET is_done=0 WHERE id=? AND user_id=?", id, loginInfo.ID) // Use DB#Exec for one entry
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Redirect to task list
	ctx.Redirect(http.StatusFound, "/")
}

// tasks with category name struct
type TaskCategory struct {
	database.Task
	CategoryName *string `db:"name"`
	Remaining    int64
}

// Searching Tasks
func Search(ctx *gin.Context) {
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

	// Get search query
	query := ctx.PostForm("search")

	// Get filter
	filter := ctx.PostForm("filter")

	// Get category
	category_name := ctx.PostForm("category")

	// Get Start date
	start := ctx.PostForm("start")

	// Get End date
	end := ctx.PostForm("end")

	if query == "" && (filter == "" || filter == "all") && category_name == "" && start == "" && end == "" {
		ctx.Redirect(http.StatusFound, "/")
		return
	}

	// Get tasks with its category using join in DB
	var tasks []TaskCategory
	periodquery := ""
	if start != "" && end != "" {
		periodquery = " AND expires BETWEEN ? AND ?"
	}
	orderquery := " ORDER BY tasks.created_at DESC"
	// category_id is in tasks table and name is in categories table
	// use TaskCategory struct to get category name
	if filter == "all" || filter == "" {
		if category_name == "" {
			err = db.Select(&tasks, "SELECT tasks.*, categories.name FROM tasks LEFT JOIN categories ON tasks.category_id = categories.id WHERE tasks.user_id = ? AND tasks.title LIKE ?"+periodquery+orderquery, loginInfo.ID, "%"+query+"%", start, end)
		} else {
			err = db.Select(&tasks, "SELECT tasks.*, categories.name FROM tasks LEFT JOIN categories ON tasks.category_id = categories.id WHERE tasks.user_id = ? AND tasks.title LIKE ? AND categories.id = ?"+periodquery+orderquery, loginInfo.ID, "%"+query+"%", category_name, start, end)
		}
	} else if filter == "todo" {
		if category_name == "" {
			err = db.Select(&tasks, "SELECT tasks.*, categories.name FROM tasks LEFT JOIN categories ON tasks.category_id = categories.id WHERE tasks.user_id = ? AND tasks.title LIKE ? AND tasks.is_done=0"+periodquery+orderquery, loginInfo.ID, "%"+query+"%", start, end)
		} else {
			err = db.Select(&tasks, "SELECT tasks.*, categories.name FROM tasks LEFT JOIN categories ON tasks.category_id = categories.id WHERE tasks.user_id = ? AND tasks.title LIKE ? AND tasks.is_done=0 AND categories.id = ?"+periodquery+orderquery, loginInfo.ID, "%"+query+"%", category_name, start, end)
		}
	} else {
		if category_name == "" {
			err = db.Select(&tasks, "SELECT tasks.*, categories.name FROM tasks LEFT JOIN categories ON tasks.category_id = categories.id WHERE tasks.user_id = ? AND tasks.title LIKE ? AND tasks.is_done=1"+periodquery+orderquery, loginInfo.ID, "%"+query+"%", start, end)
		} else {
			err = db.Select(&tasks, "SELECT tasks.*, categories.name FROM tasks LEFT JOIN categories ON tasks.category_id = categories.id WHERE tasks.user_id = ? AND tasks.title LIKE ? AND tasks.is_done=1 AND categories.id = ?"+periodquery+orderquery, loginInfo.ID, "%"+query+"%", category_name, start, end)
		}
	}
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	var categories []database.Category
	err = db.Select(&categories, "SELECT * FROM categories")
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// tasks[i].expires - now
	for i := 0; i < len(tasks); i++ {
		tasks[i].Remaining = tasks[i].Expires.Unix() - time.Now().Unix()
	}

	// Render tasks
	ctx.HTML(http.StatusOK, "index.html", gin.H{"Title": "HOME", "Tasks": tasks, "Query": query, "Categories": categories, "User": loginInfo})
}
