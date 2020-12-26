package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
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

func getPages(db *sql.DB) []*Page {
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

func loadPage(db *sql.DB, pageUrl string) *Page {
	page := new(Page)
	body, err := readFile(fmt.Sprintf("%s/%s.md", pagesDir, pageUrl), page)
	if err != nil {
		log.Fatalf("Failed to load page %s\n", pageUrl)
	}
	page.Body = body
	return page
}

func pullPages(db *sql.DB) {
	pages := getPages(db)
	for _, page := range pages {
		page.Pull(db)
	}
}

func pushPage(db *sql.DB, pageUrl string) {
	page := loadPage(db, pageUrl)
	wikiPage := map[string]interface{}{
		"wiki_page": map[string]interface{}{
			"title":            page.Title,
			"body":             page.Body,
			"editing_roles":    "teachers", // TODO: make configurable
			"notify_of_update": false,      // TODO: make configurable (cobra flag)
			"published":        page.Published,
			"front_page":       page.FrontPage,
		},
	}
	courses, _ := findCourses(db)
	for _, course := range courses {
		courseId := courses[0].CanvasId
		pageFullPath := fmt.Sprintf(pagePath, courseId, pageUrl)
		fmt.Printf("Pushing page %s to %s\n", pageUrl, course.Name)
		mustPutObject(pageFullPath, url.Values{}, wikiPage, nil)
	}
}

func pushPages(db *sql.DB) {
	files, err := ioutil.ReadDir(pagesDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		pushPage(db, getPageUrlFromFilepath(f.Name()))
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

func (page *Page) Pull(db *sql.DB) error {
	return pullComponent(db, pagePath, page.Url, page)
}
