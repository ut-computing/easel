package main

type Module struct {
	Id int `json:"id"`
	// the position of this module in the course (1-based)
	Position int `json:"position"`
	// the name/text of this module
	Name string `json:"name"`
	// (Optional) the date this module will unlock
	UnlockAt string `json:"unlock_at"`
	// Whether module items must be unlocked in order
	RequireSequentialProgress bool `json:"require_sequential_progress"`
	// IDs of Modules that must be completed before this one is unlocked
	PrerequisiteModuleIds []int `json:"prerequisite_module_ids"`
	// The number of items in the module TODO: view only?
	ItemsCount int `json:"items_count"`
	// The API URL to retrive this module's items TODO: view only?
	ItemsUrl string `json:"items_url"`
	// The contents of this module, as an array of Module Items. (Present only if
	// requested via include[]=items AND the module is not deemed too large by
	// Canvas.)
	Items []ModuleItem `json:"items"`
	// (Optional) Whether this module is published. This field is present only if
	// the caller has permission to view unpublished modules.
	Published bool `json:"published"`
}

type ModuleItem struct {
	Id int `json:"id"`
	// the id of the Module this item appears in
	ModuleId int `json:"module_id"`
	// the position of this item in the module (1-based)
	Position int `json:"position"`
	// the title of this item
	Title string `json:"title"`
	// 0-based indent level; module items may be indented to show a hierarchy
	Indent int `json:"indent"`
	// the type of object referred to one of 'File', 'Page', 'Discussion',
	// 'Assignment', 'Quiz', 'SubHeader', 'ExternalUrl', 'ExternalTool'
	Type string `json:"type"`
	// the id of the object referred to applies to 'File', 'Discussion',
	// 'Assignment', 'Quiz', 'ExternalTool' types
	ContentId int `json:"content_id"`
	// link to the item in Canvas
	HtmlUrl string `json:"html_url"`
	// (Optional) link to the Canvas API object, if applicable
	Url string `json:"url"`
	// (only for 'Page' type) unique locator for the linked wiki page, e.g.,
	// "my-page-title"
	PageUrl string `json:"page_url"`
	// (only for 'ExternalUrl' and 'ExternalTool' types) external url that the item
	// points to
	ExternalUrl string `json:"external_url"`
	// (only for 'ExternalTool' type) whether the external tool opens in a new tab
	NewTab bool `json:"new_tab"`
	// Completion requirement for this module item
	CompletionRequirement CompletionRequirement `json:"completion_requirement"`
	// (Optional) Whether this module item is published. This field is present only
	// if the caller has permission to view unpublished items.
	Published bool `json:"published"`
}

type CompletionRequirement struct {
	Type      string `json:type"`
	MinScore  int    `json:"min_score"`
	Completed bool   `json:"completed"`
}

func pullModules(id string) {
	if id == "all" {
		// TODO: get all
	} else {
		// TODO: get single
	}
}
