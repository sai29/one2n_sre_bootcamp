package data

import (
	"database/sql"
	"log"
	"time"
)

type Student struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Name      string    `json:"title"`
	RollNo    int32     `json:"rollno"`
}

type StudentModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}
