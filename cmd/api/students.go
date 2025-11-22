package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sai29/one2n_sre_bootcamp/internal/data"
)

func (app *application) createStudentHandler(c *gin.Context) {
	var student data.Student

	if err := c.ShouldBindJSON(&student); err != nil {
		app.badRequestResponse(c, err)
		return
	}

	if err := app.models.Students.Insert(&student); err != nil {
		app.serverErrorResponse(c, err)
		return
	}

	location := fmt.Sprintf("/v1/students/%d", student.ID)

	c.Header("Location", location)
	c.JSON(http.StatusCreated, gin.H{
		"student": student,
	})
}

func (app *application) showStudentHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id < 1 {
		app.notFoundResponse(c)
		return
	}

	student, err := app.models.Students.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			app.notFoundResponse(c)
		default:
			app.serverErrorResponse(c, err)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"student": student,
	})
}

func (app *application) listStudentsHandler(c *gin.Context) {
	students, err := app.models.Students.ListAll()
	if err != nil {
		app.serverErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"students": students,
	})
}

func (app *application) updateStudentHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id < 1 {
		app.notFoundResponse(c)
		return
	}

	studentRecord, err := app.models.Students.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(c)
		default:
			app.serverErrorResponse(c, err)
		}
		return
	}

	var input struct {
		Name   *string `json:"name"`
		RollNo *int32  `json:"rollno"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		app.badRequestResponse(c, err)
		return
	}

	if input.Name != nil {
		studentRecord.Name = *input.Name
	}

	if input.RollNo != nil {
		studentRecord.RollNo = *input.RollNo
	}

	if err := app.models.Students.Update(studentRecord); err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(c)
		default:
			app.serverErrorResponse(c, err)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"student": studentRecord})
}

func (app *application) deleteStudentHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id < 1 {
		app.notFoundResponse(c)
		return
	}

	err = app.models.Students.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(c)
		default:
			app.serverErrorResponse(c, err)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "student deleted",
	})
}
