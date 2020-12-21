package main

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"
)

type Component interface {
	Dump() error
}

func pullComponent(db *sql.DB, path string, id interface{}, component Component) error {
	courses, _ := findCourses(db)
	// TODO: do it for all courses
	courseId := courses[0].CanvasId
	fullPath := fmt.Sprintf(path, courseId, id)

	fmt.Printf("Pulling %T %s\n", component, fullPath)
	mustGetObject(fullPath, url.Values{}, component)
	return component.Dump()
}

func slug(text string) string {
	return strings.ToLower(strings.ReplaceAll(text, " ", "-"))
}
