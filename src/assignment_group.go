package main

import (
	"database/sql"
	"fmt"
	"net/url"

	"gopkg.in/yaml.v2"
)

const (
	assignmentGroupsTable = "assignment_groups"
	assignmentGroupsDir   = "assignment_groups" // TODO: make configable
)

type AssignmentGroup struct {
	Id          int     `json:"-" yaml:"-" meddler:"id,pk"`
	CanvasId    int     `json:"id" yaml:"id" meddler:"canvas_id"`
	Name        string  `json:"name" yaml:"name" meddler:"name"`
	Position    int     `json:"position" yaml:"position" meddler:"position"`
	GroupWeight float64 `json:"group_weight" yaml:"group_weight" meddler:"group_weight"`
}

func getAssignmentGroups(db *sql.DB) []*AssignmentGroup {
	ags := make([]*AssignmentGroup, 0)
	courses, _ := findCourses(db)
	values := url.Values{}
	values.Add("per_page", "100")
	// TODO: do it for all courses
	courseId := courses[0].CanvasId
	reqUrl := fmt.Sprintf(assignmentGroupsPath, courseId)
	mustGetObject(reqUrl, values, &ags)
	return ags
}

func pullAssignmentGroups(db *sql.DB) {
	ags := getAssignmentGroups(db)
	for _, ag := range ags {
		ag.Dump()
	}
}

func (ag *AssignmentGroup) Dump() error {
	metadata, err := yaml.Marshal(ag)
	if err != nil {
		return err
	}
	assignmentGroupFilePath := fmt.Sprintf("%s/%s.md", assignmentGroupsDir, ag.Slug())
	return writeYamlFile(assignmentGroupFilePath, string(metadata))
}

func (ag *AssignmentGroup) Pull(db *sql.DB) error {
	return pullComponent(db, assignmentGroupPath, ag.CanvasId, ag)
}

func (ag *AssignmentGroup) Slug() string {
	return slug(ag.Name)
}
