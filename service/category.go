package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	database "todolist.go/db"
)

func PostCategory(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	category := ctx.PostForm("name")

	// Insert category into database
	db.MustExec("INSERT INTO categories (name) VALUES (?)", category)

	ctx.Redirect(http.StatusSeeOther, "/")
}
