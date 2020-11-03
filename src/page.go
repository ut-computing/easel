package main

type Page struct {
	// the unique locator for the page, e.g., "my-page-title"
	Url string `json:"url"`
	// the title of the page
	Title string `json:"title"`
	// the creation date for the page
	CreatedAt string `json:"created_at"`
	// the date the page was last updated
	UpdatedAt string `json:"updated_at"`
	// the page content, in HTML (present when requesting a single page; omitted
	// when listing pages)
	Body string `json:"body"`
	// whether the page is published (true) or draft state (false).
	Published bool `json:"published"`
	// whether this page is the front page for the wiki
	FrontPage bool `json:"front_page"`
}

func pullPages(id string) {
	if id == "all" {
		// TODO: get all
	} else {
		// TODO: get single
	}
}
