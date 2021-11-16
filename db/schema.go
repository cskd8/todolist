package db

// schema.go provides data models in DB
import (
	"time"
)

// Task corresponds to a row in `tasks` table
type Task struct {
	ID         uint64    `db:"id"`
	Title      string    `db:"title"`
	CategoryID *uint64   `db:"category_id"`
	UserID     uint64    `db:"user_id"`
	Expires    time.Time `db:"expires"`
	CreatedAt  time.Time `db:"created_at"`
	IsDone     bool      `db:"is_done"`
}

type Category struct {
	ID        uint64    `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}

type User struct {
	ID        uint64    `db:"id" json:"id" bson:"id"`
	Name      string    `db:"name" json:"name" bson:"name"`
	Password  string    `db:"password" json:"password" bson:"password"`
	CreatedAt time.Time `db:"created_at" json:"created_at" bson:"created_at"`
}
