package main

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"

	"github.com/russross/meddler"
	"gopkg.in/yaml.v2"
)

const (
	modulesTable = "modules"
	modulesDir   = "modules" // TODO: make configable
)

type Module struct {
	Id       int `json:"-" yaml:"-" meddler:"id,pk"`
	CanvasId int `json:"id" yaml:"id" meddler:"canvas_id"`
	// the position of this module in the course (1-based)
	Position int `json:"position" yaml:"position" meddler:"position"`
	// the name/text of this module
	Name string `json:"name" yaml:"name" meddler:"name"`
	// (Optional) the date this module will unlock
	UnlockAt string `json:"unlock_at" yaml:"unlock_at" meddler:"unlock_at"`
	// Whether module items must be unlocked in order
	RequireSequentialProgress bool `json:"require_sequential_progress" yaml:"require_sequential_progress" meddler:"require_sequential_progress"`
	// IDs of Modules that must be completed before this one is unlocked
	PrerequisiteModuleIds []int `json:"prerequisite_module_ids" yaml:"prerequisite_module_ids" meddler:"prerequisite_module_ids"`
	// The number of items in the module TODO: view only?
	ItemsCount int `json:"items_count" yaml:"items_count" meddler:"items_count"`
	// The API URL to retrive this module's items TODO: view only?
	ItemsUrl string `json:"items_url" yaml:"items_url" meddler:"items_url"`
	// The contents of this module, as an array of Module Items. (Present only if
	// requested via include[]=items AND the module is not deemed too large by
	// Canvas.)
	Items []ModuleItem `json:"items" yaml:"items" meddler:"items"`
	// (Optional) Whether this module is published. This field is present only if
	// the caller has permission to view unpublished modules.
	Published bool `json:"published" yaml:"published" meddler:"published"`
}

type ModuleItem struct {
	Id       int `json:"-" yaml:"-" meddler:"id,pk"`
	CanvasId int `json:"id" yaml:"id" meddler:"canvas_id"`
	// the id of the Module this item appears in
	ModuleId int `json:"module_id" yaml:"module_id" meddler:"module_id"`
	// the position of this item in the module (1-based)
	Position int `json:"position" yaml:"position" meddler:"position"`
	// the title of this item
	Title string `json:"title" yaml:"title" meddler:"title"`
	// 0-based indent level; module items may be indented to show a hierarchy
	Indent int `json:"indent" yaml:"indent" meddler:"indent"`
	// the type of object referred to one of 'File', 'Page', 'Discussion',
	// 'Assignment', 'Quiz', 'SubHeader', 'ExternalUrl', 'ExternalTool'
	Type string `json:"type" yaml:"type" meddler:"type"`
	// the id of the object referred to applies to 'File', 'Discussion',
	// 'Assignment', 'Quiz', 'ExternalTool' types
	ContentId int `json:"content_id" yaml:"content_id" meddler:"content_id"`
	// link to the item in Canvas
	HtmlUrl string `json:"html_url" yaml:"html_url" meddler:"html_url"`
	// (Optional) link to the Canvas API object, if applicable
	Url string `json:"url" yaml:"url" meddler:"url"`
	// (only for 'Page' type) unique locator for the linked wiki page, e.g.,
	// "my-page-title"
	PageUrl string `json:"page_url" yaml:"page_url" meddler:"page_url"`
	// (only for 'ExternalUrl' and 'ExternalTool' types) external url that the item
	// points to
	ExternalUrl string `json:"external_url" yaml:"external_url" meddler:"external_url"`
	// (only for 'ExternalTool' type) whether the external tool opens in a new tab
	NewTab bool `json:"new_tab" yaml:"new_tab" meddler:"new_tab"`
	// Completion requirement for this module item
	CompletionRequirement CompletionRequirement `json:"completion_requirement" yaml:"completion_requirement" meddler:"completion_requirement"`
	// (Optional) Whether this module item is published. This field is present only
	// if the caller has permission to view unpublished items.
	Published bool `json:"published" yaml:"published" meddler:"published"`
}

type CompletionRequirement struct {
	Type      string `json:"type" yaml:"type" meddler:"type"`
	MinScore  int    `json:"min_score" yaml:"min_score" meddler:"min_score"`
	Completed bool   `json:"completed" yaml:"completed" meddler:"completed"`
}

func loadModuleFromFile(moduleFilepath string) (*Module, error) {
	module := new(Module)
	err := readYamlFile(moduleFilepath, module)
	return module, err
}

func getModules(db *sql.DB) []*Module {
	modules := make([]*Module, 0)
	courses, _ := findCourses(db)
	values := url.Values{}
	values.Add("per_page", "100")
	// TODO: do it for all courses
	courseId := courses[0].CanvasId
	reqUrl := fmt.Sprintf(modulesPath, courseId)
	mustGetObject(reqUrl, values, &modules)
	return modules
}

func pullModule(db *sql.DB, moduleId int) error {
	module := new(Module)
	module.Id = moduleId
	err := module.Pull(db)
	if err != nil {
		return err
	}
	return module.Save(db)
}

func pullModules(db *sql.DB) {
	modulesMeta := getModules(db)
	// TODO: prompt for overwrite, etc.

	for _, module := range modulesMeta {
		pullModule(db, module.CanvasId)
	}
}

func (module *Module) Slug() string {
	return strings.ToLower(strings.ReplaceAll(module.Name, " ", "-"))
}

func (module *Module) Dump() error {
	data, err := yaml.Marshal(module)
	if err != nil {
		return err
	}
	moduleFilePath := fmt.Sprintf("%s/%s.yaml", modulesDir, module.Slug())
	return writeYamlFile(moduleFilePath, string(data))
}

func (module *Module) Pull(db *sql.DB) error {
	courses, _ := findCourses(db)
	// TODO: do it for all courses
	courseId := courses[0].CanvasId
	moduleFullPath := fmt.Sprintf(modulePath, courseId, module.Id)
	values := url.Values{}
	values.Add("include[]", "items, content_details")
	fmt.Println("Pulling module", moduleFullPath, values)
	mustGetObject(moduleFullPath, values, module)
	return module.Dump()
}

func (module *Module) Save(db *sql.DB) error {
	return meddler.Insert(db, modulesTable, module)
}
