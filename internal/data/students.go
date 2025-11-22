package data

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
)

type Student struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Name      string    `json:"name"`
	RollNo    int32     `json:"rollno"`
	Version   int32     `json:"version"`
}

type StudentModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m StudentModel) Insert(student *Student) error {
	query := `
	INSERT INTO students (name, rollno)
	VALUES ($1, $2)
	RETURNING id, created_at, version`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{student.Name, student.RollNo}

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&student.ID, &student.CreatedAt, &student.Version)
}

func (m StudentModel) Get(id int64) (*Student, error) {
	query := `
	SELECT id, created_at, name, rollno, version 
	FROM students 
	WHERE id = $1`

	var student Student

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&student.ID,
		&student.CreatedAt,
		&student.Name,
		&student.RollNo,
		&student.Version,
	)

	if err != nil {
		return nil, err
	}

	return &student, nil
}

func (m StudentModel) ListAll() ([]*Student, error) {
	query := `
		SELECT * FROM students ORDER BY id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			m.ErrorLog.Println(err)
		}
	}()

	students := []*Student{}

	for rows.Next() {
		var s Student

		err := rows.Scan(
			&s.ID,
			&s.CreatedAt,
			&s.Name,
			&s.RollNo,
			&s.Version,
		)
		if err != nil {
			return nil, err
		}
		students = append(students, &s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return students, rows.Err()

}

func (m StudentModel) Update(student *Student) error {
	query := `
		UPDATE students
		SET name = $1, rollno = $2, version = version + 1
		WHERE id = $3 AND version = $4
		RETURNING version, created_at`

	args := []interface{}{
		student.Name,
		student.RollNo,
		student.ID,
		student.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newVersion int32
	var createdAt time.Time

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&newVersion, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrEditConflict
		}
		return err
	}

	student.Version = newVersion
	student.CreatedAt = createdAt
	return nil
}

func (m StudentModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM students
		WHERE id = $1
		`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
