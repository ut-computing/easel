package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
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

func findCourse(db *sql.DB, courseId int) (*Course, error) {
	course := new(Course)
	err := meddler.QueryRow(db, course, "select * from "+coursesTable+" where canvas_id = ?", courseId)
	return course, err
}

func findCourses(db *sql.DB) ([]*Course, error) {
	var courses []*Course
	err := meddler.QueryAll(db, &courses, "select * from "+coursesTable)
	return courses, err
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

func matchCourse(db *sql.DB, search string) ([]*Course, error) {
	courses := []*Course{}

	// assume searching by canvas id first
	err := meddler.QueryAll(db, &courses, `SELECT * FROM `+coursesTable+` WHERE canvas_id = ?`, search)

	// if we didn't find the course, user might have used section number so
	// search for it by name
	if err != nil || len(courses) == 0 {
		// build search terms
		where := ""
		args := []interface{}{}
		for _, term := range strings.Split(search, " ") {
			where, args = addWhereLike(where, args, "name", term)
		}

		err = meddler.QueryAll(db, &courses, `SELECT * FROM `+coursesTable+where, args...)
	}

	return courses, err
}

func pullCourse(db *sql.DB, courseId int) (*Course, error) {
	course := new(Course)
	values := url.Values{}
	values.Add("include[]", "syllabus_body")
	if !getObject(fmt.Sprintf(coursePath, courseId), values, course) {
		return course, errors.New("Failed to pull course from Canvas")
	}

	// TODO: prompt for overwrite, manually merge, abort
	err := meddler.Insert(db, coursesTable, course)
	if err != nil {
		return course, err
	}

	if _, err := os.Stat("syllabus.md"); os.IsNotExist(err) {
		err = course.Dump()
		if err != nil {
			return course, err
		}
	}
	return course, nil
}

func pullCourses(db *sql.DB) ([]*Course, error) {
	courses, err := findCourses(db)
	if err != nil {
		return courses, err
	}

	for i := range courses {
		courses[i], err = pullCourse(db, courses[i].Id)
		if err != nil {
			return courses, err
		}
	}

	return courses, nil
}

func pushCourses(db *sql.DB) {
	courses, err := findCourses(db)
	if err != nil {
		log.Fatalf("Error finding courses: %v", err)
	}
	for _, course := range courses {
		syllabus, err := readFile("syllabus.md", nil)
		if err != nil {
			log.Fatalf("Failed to load syllabus\n")
		}
		course.Syllabus = syllabus
		course.Push()
	}
}

func (course *Course) Dump() error {
	return writeFile("syllabus.html", "", course.Syllabus)
}

func (course *Course) GetCourseNumber() string {
	values := strings.Split(course.Name, " ")
	return values[0]
}

func (course *Course) Push() error {
	return nil
}

func (course *Course) Remove(db *sql.DB) error {
	_, err := db.Exec("DELETE from "+coursesTable+" WHERE id = ?", course.Id)
	return err
}

func (course *Course) String() string {
	// TODO: grind student tabular output?
	return fmt.Sprintf("%d\t%s\t%s\t%s", course.CanvasId, course.Name, course.Code, course.WorkflowState)
}
