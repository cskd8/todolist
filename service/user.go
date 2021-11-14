package service

import (
	"encoding/json"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/koron/go-dproxy"
	"golang.org/x/crypto/bcrypt"
	database "todolist.go/db"
)

func RLogin(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login.html", gin.H{"Title": "Login"})
}

func RSignup(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "signup.html", gin.H{"Title": "Signup"})
}

func Login(ctx *gin.Context) {
	session := sessions.Default(ctx)
	// Login with username and password using SHA2

	username := ctx.PostForm("username")
	password := ctx.PostForm("password")

	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Get user from DB
	var user database.User
	err = db.Get(&user, "SELECT * FROM users WHERE name = ?", username)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}
		ctx.String(http.StatusBadRequest, "Invalid username")
		return
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Invalid password")
		return
	}

	// json user
	userJson, err := json.Marshal(user)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Set session
	session.Set("user", string(userJson))

	// Save session
	session.Save()

	// Redirect to index
	ctx.Redirect(http.StatusFound, "/")
}

func Signup(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	username := ctx.PostForm("username")
	password := ctx.PostForm("password")

	// Get user from DB
	var user database.User
	err = db.Get(&user, "SELECT * FROM users WHERE name = ?", username)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err == nil {
		ctx.String(http.StatusInternalServerError, "Username already exists")
		return
	}

	// Create user
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Insert user
	_, err = db.Exec("INSERT INTO users (name, password) VALUES (?, ?)", username, string(hashedPassword))
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Redirect to login
	ctx.Redirect(http.StatusFound, "/login")
}

func Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Save()

	ctx.Redirect(http.StatusFound, "/login")
}

func LoginCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		loginUserJson, err := dproxy.New(session.Get("user")).String()

		if err != nil {
			c.Status(http.StatusUnauthorized)
			c.Abort()
		} else {
			var loginInfo database.User
			err := json.Unmarshal([]byte(loginUserJson), &loginInfo)
			if err != nil {
				c.Status(http.StatusUnauthorized)
				c.Abort()
			} else {
				c.Next()
			}
		}
	}
}
