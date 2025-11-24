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
	Students StudentStore
}

type StudentStore interface {
	Insert(*Student) error
	Get(int64) (*Student, error)
	ListAll() ([]*Student, error)
	Update(*Student) error
	Delete(int64) error
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
