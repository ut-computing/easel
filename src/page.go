package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"path/filepath"

	"github.com/russross/blackfriday"
	"gopkg.in/yaml.v2"
)

const (
	pagesTable = "pages"
	pagesDir   = "pages" // TODO: make configable
)

type Page struct {
	Id             int    `json:"-" yaml:"-" meddler:"id,pk"`
	Url            string `json:"url" yaml:"url" meddler:"url"` // the unique locator for the page, e.g., "my-page-title"
	Title          string `json:"title" yaml:"title" meddler:"title" `
	CreatedAt      string `json:"created_at" yaml:"created_at" meddler:"created_at" `
	UpdatedAt      string `json:"updated_at" yaml:"updated_at" meddler:"updated_at" `
	Body           string `json:"body" yaml:"-" meddler:"body" `                      // the page content, in HTML
	Published      bool   `json:"published" yaml:"published" meddler:"published" `    // whether the page is published (true) or draft state (false).
	FrontPage      bool   `json:"front_page" yaml:"front_page" meddler:"front_page" ` // whether this page is the front page for the wiki
	TodoDate       string `json:"todo_date" yaml:"todo_date" meddler:"todo_date"`
	EditingRoles   string `json:"editing_roles" yaml:"editing_roles" meddler:"editing_roles"` // command separated string: "teachers,students,members,public"
	NotifyOfUpdate bool   `json:"notify_of_update" yaml:"-" meddler:"-"`
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
	mdBody := string(blackfriday.MarkdownCommon([]byte(page.Body)))
	// Using a map here because Canvas doesn't like it when we PUT with fields
	// such as CreatedAt and I can't figure out how to remove them only for
	// marshaling.
	wikiPage := map[string]interface{}{
		"wiki_page": map[string]interface{}{
			"title":            page.Title,
			"body":             mdBody,
			"editing_roles":    page.EditingRoles,
			"notify_of_update": false, // TODO: make configurable (cobra flag)
			"published":        page.Published,
			"front_page":       page.FrontPage,
			"todo_date":        page.TodoDate, // TODO: canvas not accepting this for some reason
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
