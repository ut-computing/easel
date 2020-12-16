package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/russross/meddler"
)

const (
	coursesTable = "courses"
)

type Course struct {
	Id       int    `json:"-" yaml:"-" meddler:"id,pk"`
	CanvasId int    `json:"id" yaml:"id" meddler:"canvas_id"`
	Name     string `json:"name" yaml:"name" meddler:"name"`
	Code     string `json:"course_code" yaml:"course_code" meddler:"code"`
	// the current state of the course one of 'unpublished', 'available',
	// 'completed', or 'deleted'
	WorkflowState string `json:"workflow_state" yaml:"workflow_state" meddler:"workflow_state"`
	Syllabus      string `json:"syllabus_body" yaml:"-" meddler:"-"`
}

func mustCreateCoursesTable(db *sql.DB) {
	command := `CREATE TABLE courses (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"canvas_id" integer NOT NULL,
		"name" TEXT NOT NULL,
		"code" TEXT NOT NULL,
		"workflow_state" TEXT NOT NULL
	  );`

	statement, err := db.Prepare(command)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
}

func getCourseIdFromUrl(courseUrl string) (int, error) {
	parsed, err := url.Parse(courseUrl)
	if err != nil {
		return -1, err
	}
	path := strings.Split(parsed.Path, "/")
	courseId, err := strconv.Atoi(path[len(path)-1])
	if err != nil {
		return -1, err
	}
	return courseId, nil
}

func findCourse(db *sql.DB, courseId int) (*Course, error) {
	course := new(Course)
	err := meddler.Load(db, coursesTable, course, int64(courseId))
	return course, err
}

func pullCourse(courseId int) *Course {
	course := new(Course)
	values := url.Values{}
	values.Add("include[]", "syllabus_body")
	mustGetObject(fmt.Sprintf(coursePath, courseId), values, course)
	return course
}

func createCourse(db *sql.DB, courseId int) (*Course, error) {
	course := pullCourse(courseId)
	err := meddler.Insert(db, coursesTable, course)
	return course, err
}

func (course *Course) dump() error {
	return writeFile("syllabus.html", "", course.Syllabus)
}
