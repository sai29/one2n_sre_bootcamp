package data

import (
	"database/sql"
	"errors"
	"log"
	"os"
)

var (
	ErrRecordNotFound = errors.New("record not found")

	ErrEditConflict = errors.New("edit conflict")
)

type Models struct {
	Students StudentModel
}

func NewModels(db *sql.DB) Models {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)
	return Models{
		Students: StudentModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
	}
}
