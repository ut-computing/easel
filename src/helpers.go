package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const (
	easelDb = ".easeldb"
)

func writeFile(filename, metadata, html string) error {
	text := html
	if metadata != "" {
		text = fmt.Sprintf("```\n%s```\n", metadata)
	}
	return ioutil.WriteFile(filename, []byte(text), 0644)
}

func mustCreateDb() {
	directory, err := filepath.Abs(".")
	if err != nil {
		log.Fatalf("error finding directory: %v", err)
	}

	_, err = os.Stat(filepath.Join(directory, easelDb))
	if err == nil {
		log.Fatal("Database already exists")
	}

	db, err := sql.Open("sqlite3", easelDb)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// TODO: create all tables (execute .sql file)
	mustCreateCoursesTable(db)
}

func findDb() *sql.DB {
	directory, err := filepath.Abs(".")
	if err != nil {
		log.Fatalf("error finding directory: %v", err)
	}

	dbName := easelDb
	stepDir := directory
	for {
		path := filepath.Join(directory, dbName)
		_, err := os.Stat(path)
		if err == nil {
			break
		}
		if !os.IsNotExist(err) {
			log.Fatalf("error searching for %s in %s: %v", dbName, directory, err)
		}

		// try moving up a directory
		stepDir = directory
		directory = filepath.Dir(directory)
		if directory == stepDir {
			log.Fatal("No database found.")
		}
	}

	dbName = directory + "/" + dbName
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
