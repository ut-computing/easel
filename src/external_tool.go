package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"path/filepath"
)

const (
	externalToolsDir = "external_tools"
)

type ExternalTool struct {
	Name         string `json:"name" yaml:"name"`
	ConsumerKey  string `json:"consumer_key" yaml:"consumer_key"`
	SharedSecret string `json:"shared_secret" yaml:"shared_secret"`
	ConfigType   string `json:"config_type" yaml:"config_type"`
	ConfigUrl    string `json:"config_url" yaml:"config_url"`
}

func pushExternalTool(db *sql.DB, filename string) {
	et := new(ExternalTool)
	// read yaml
	err := readYamlFile(filename, et)
	if err != nil {
		log.Fatalf("Failed to read yaml file %s: %v\n", filename, err)
	}
	et.Push(db)
}

func pushExternalTools(db *sql.DB) {
	files, err := ioutil.ReadDir(externalToolsDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fullPath := fmt.Sprintf("%s/%s", externalToolsDir, f.Name())
		if filepath.Ext(fullPath) == ".yaml" {
			pushExternalTool(db, fullPath)
		}
	}
}

func (et *ExternalTool) Push(db *sql.DB) {
	courses, _ := findCourses(db)
	for _, course := range courses {
		etFullPath := fmt.Sprintf(externalToolsPath, course.CanvasId)
		fmt.Printf("Pushing %T %s\n", et, et.Name)
		mustPostObject(etFullPath, url.Values{}, et, nil)
	}
}
