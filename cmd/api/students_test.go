package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sai29/one2n_sre_bootcamp/internal/data"
)

type mockStudentModel struct {
	insertFn func(s *data.Student) error
	getFn    func(id int64) (*data.Student, error)
	listFn   func() ([]*data.Student, error)
	updateFn func(s *data.Student) error
	deleteFn func(id int64) error
}

func (m *mockStudentModel) Insert(s *data.Student) error {
	return m.insertFn(s)
}

func (m *mockStudentModel) Get(id int64) (*data.Student, error) {
	return m.getFn(id)
}

func (m *mockStudentModel) ListAll() ([]*data.Student, error) {
	return m.listFn()
}

func (m *mockStudentModel) Update(s *data.Student) error {
	return m.updateFn(s)
}

func (m *mockStudentModel) Delete(id int64) error {
	return m.deleteFn(id)
}

func newTestApp(mock *mockStudentModel) *application {
	return &application{
		models: data.Models{Students: mock},
	}
}

func performRequest(r http.Handler, method, path string, body []byte) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	return w
}

func TestCreateStudentHandler(t *testing.T) {
	mock := &mockStudentModel{
		insertFn: func(s *data.Student) error {
			s.ID = 1
			return nil
		},
	}

	app := newTestApp(mock)
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.POST("/v1/students", app.createStudentHandler)

	body := []byte(`{"name":"John Doe","rollno": 10}`)
	w := performRequest(router, "POST", "/v1/students", body)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestShowStudentHandler(t *testing.T) {
	mock := &mockStudentModel{
		getFn: func(id int64) (*data.Student, error) {
			return &data.Student{ID: id, Name: "Bob", RollNo: 22}, nil
		},
	}

	app := newTestApp(mock)
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/v1/students/:id", app.showStudentHandler)

	w := performRequest(router, "GET", "/v1/students/1", nil)

	if w.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, w.Code)
	}
}

func TestListStudenthandler(t *testing.T) {
	mock := &mockStudentModel{
		listFn: func() ([]*data.Student, error) {
			return []*data.Student{
				{ID: 1, Name: "John"},
				{ID: 2, Name: "Jane"},
			}, nil
		},
	}

	app := newTestApp(mock)
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/v1/students", app.listStudentsHandler)

	w := performRequest(router, "GET", "/v1/students", nil)

	if w.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, w.Code)
	}
}

func TestUpdateStudentHandler(t *testing.T) {
	mock := &mockStudentModel{
		getFn: func(id int64) (*data.Student, error) {
			return &data.Student{
				ID: 1, Name: "Old Name", RollNo: 4,
			}, nil
		},
		updateFn: func(s *data.Student) error { return nil },
	}

	app := newTestApp(mock)
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.PUT("/v1/students/:id", app.updateStudentHandler)

	body := []byte(`{"name":"New","rollno":12}`)
	w := performRequest(router, "PUT", "/v1/students/1", body)

	if w.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, w.Code)
	}
}

func TestDeleteStudentHandler(t *testing.T) {
	mock := &mockStudentModel{
		deleteFn: func(id int64) error { return nil },
	}

	app := newTestApp(mock)
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.DELETE("/v1/students/:id", app.deleteStudentHandler)

	w := performRequest(router, "DELETE", "/v1/students/1", nil)

	if w.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, w.Code)
	}
}

func TestDeleteStudentHandler_NotFound(t *testing.T) {
	mock := &mockStudentModel{
		deleteFn: func(id int64) error { return data.ErrRecordNotFound },
	}

	app := newTestApp(mock)
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.DELETE("/v1/students/:id", app.deleteStudentHandler)

	w := performRequest(router, "DELETE", "/v1/students/99", nil)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected %d, got %d", http.StatusNotFound, w.Code)
	}
}
