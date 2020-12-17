package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	pagesTable = "pages"
	pagesDir   = "pages" // TODO: make configable
)

type Page struct {
	Id int `json:"-" yaml:"-" meddler:"id,pk"`
	// the unique locator for the page, e.g., "my-page-title"
	Url string `json:"url" yaml:"url" meddler:"url"`
	// the title of the page
	Title string `json:"title" yaml:"title" meddler:"title" `
	// the creation date for the page
	CreatedAt string `json:"created_at" yaml:"created_at" meddler:"created_at" `
	// the date the page was last updated
	UpdatedAt string `json:"updated_at" yaml:"updated_at" meddler:"updated_at" `
	// the page content, in HTML (present when requesting a single page; omitted
	// when listing pages)
	Body string `json:"body" yaml:"-" meddler:"body" `
	// whether the page is published (true) or draft state (false).
	Published bool `json:"published" yaml:"published" meddler:"published" `
	// whether this page is the front page for the wiki
	FrontPage bool `json:"front_page" yaml:"front_page" meddler:"front_page" `
}

func mustCreatePagesTable(db *sql.DB) {
	command := `CREATE TABLE pages (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"url" text NOT NULL,
		"title" integer NOT NULL,
		"created_at" TEXT NOT NULL,
		"updated_at" TEXT NOT NULL,
		"body" TEXT NOT NULL
		"published" integer NOT NULL,
		"front_page" integer NOT NULL
	  );`

	statement, err := db.Prepare(command)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
}

func getAllPages(db *sql.DB) []*Page {
	pages := make([]*Page, 0)
	courses, _ := findCourses(db)
	values := url.Values{}
	values.Add("per_page", "100")
	// TODO: do it for all courses
	courseId := courses[0].CanvasId
	reqUrl := fmt.Sprintf(pagesPath, courseId)
	mustGetObject(reqUrl, values, &pages)
	return pages
}

func getPageUrlFromFilepath(pagefilepath string) string {
	ext := filepath.Ext(pagefilepath)
	basename := filepath.Base(pagefilepath)
	return basename[:len(basename)-len(ext)] // remove extension
}

func pullPage(db *sql.DB, pageUrl string) {
	page := new(Page)
	courses, _ := findCourses(db)
	// TODO: do it for all courses
	courseId := courses[0].CanvasId
	pageFullPath := fmt.Sprintf(pagePath, courseId, pageUrl)
	fmt.Println("Pulling page", pageFullPath)
	mustGetObject(pageFullPath, url.Values{}, page)
	page.Dump()
}

func pullAllPages(db *sql.DB) {
	pagesMeta := getAllPages(db)
	// TODO: prompt for overwrite, etc.

	for _, page := range pagesMeta {
		pullPage(db, page.Url)
	}
}

func (page *Page) Dump() error {
	metadata, err := yaml.Marshal(page)
	if err != nil {
		return err
	}
	pageFilePath := fmt.Sprintf("%s/%s.md", pagesDir, page.Url)
	return writeFile(pageFilePath, string(metadata), page.Body)
}
