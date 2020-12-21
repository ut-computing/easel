package main

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	assignmentsTable = "assignments"
	assignmentsDir   = "assignments" // TODO: make configable
)

type Assignment struct {
	Id             int     `json:"-" yaml:"-" meddler:"id,pk"`
	CanvasId       int     `json:"id" yaml:"id" meddler:"canvas_id"`
	CourseId       int     `json:"course_id" yaml:"course_id" meddler:"course_id"`                   // the ID of the course the assignment belongs to
	Name           string  `json:"name" yaml:"name" meddler:"name"`                                  // the name of the assignment
	Published      bool    `json:"published" yaml:"published" meddler:"published"`                   // Whether the assignment is published
	Description    string  `json:"description" yaml:"-" meddler:"description"`                       // the assignment description, in an HTML fragment
	PointsPossible float64 `json:"points_possible" yaml:"points_possible" meddler:"points_possible"` // the maximum points possible for the assignment
	Position       int     `json:"position" yaml:"position" meddler:"position"`                      // the sorting order of the assignment in the group

	// the types of submissions allowed for this assignment list containing one or
	// more of the following: 'discussion_topic', 'online_quiz', 'on_paper', 'none',
	// 'external_tool', 'online_text_entry', 'online_url', 'online_upload'
	// 'media_recording'
	SubmissionTypes []string `json:"submission_types" yaml:"submission_types" meddler:"submission_types"`

	// Allowed file extensions, which take effect if submission_types includes
	// 'online_upload'.
	AllowedExtensions []string `json:"allowed_extensions" yaml:"allowed_extensions" meddler:"allowed_extensions"`

	// (Optional) assignment's settings for external tools if submission_types
	// include 'external_tool'. Only url and new_tab are included (new_tab defaults
	// to false).  Use the 'External Tools' API if you need more information about
	// an external tool.
	ExternalToolTagAttributes ExternalToolTagAttributes `json:"external_tool_tag_attributes" yaml:"external_tool_tag_attributes" meddler:"external_tool_tag_attributes"`

	// The number of submission attempts a student can make for this assignment. -1
	// is considered unlimited.
	AllowedAttempts int `json:"allowed_attempts" yaml:"allowed_attempts" meddler:"allowed_attempts"`

	// The type of grading the assignment receives; one of 'pass_fail', 'percent',
	// 'letter_grade', 'gpa_scale', 'points'
	GradingType       string `json:"grading_type" yaml:"grading_type" meddler:"grading_type"`
	AssignmentGroupId int    `json:"assignment_group_id" yaml:"assignment_group_id" meddler:"assignment_group_id"` // the ID of the assignment's group
	QuizId            int    `json:"quiz_id" yaml:"quiz_id" meddler:"quiz_id"`                                     // (Optional) id of the associated quiz (applies only when submission_types is ['online_quiz'])

	// the due date for the assignment. returns null if not present. NOTE: If this
	// assignment has assignment overrides, this field will be the due date as it
	// applies to the user requesting information from the API.
	DueAt string `json:"due_at" yaml:"due_at" meddler:"due_at"`

	// the unlock date (assignment is unlocked after this date) returns null if not
	// present NOTE: If this assignment has assignment overrides, this field will be
	// the unlock date as it applies to the user requesting information from the
	// API.
	UnlockAt string `json:"unlock_at" yaml:"unlock_at" meddler:"unlock_at"`

	// the lock date (assignment is locked after this date). returns null if not
	// present. NOTE: If this assignment has assignment overrides, this field will
	// be the lock date as it applies to the user requesting information from the
	// API.
	LockAt    string `json:"lock_at" yaml:"lock_at" meddler:"lock_at"`
	CreatedAt string `json:"created_at" yaml:"created_at" meddler:"created_at"` // The time at which this assignment was originally created
	UpdatedAt string `json:"updated_at" yaml:"updated_at" meddler:"updated_at"` // The time at which this assignment was last modified in any way

	// If this is a group assignment, boolean flag indicating whether or not
	// students will be graded individually.
	GradeGroupStudentsIndividually bool `json:"grade_group_students_individually" yaml:"grade_group_students_individually" meddler:"grade_group_students_individually"`

	// The ID of the assignment’s group set, if this is a group assignment. For
	// group discussions, set group_category_id on the discussion topic, not the
	// linked assignment.
	GroupCategoryId int  `json:"group_category_id" yaml:"group_category_id" meddler:"group_category_id"`
	PeerReviews     bool `json:"peer_reviews" yaml:"peer_reviews" meddler:"peer_reviews"` // Boolean indicating if peer reviews are required for this assignment

	// Boolean indicating peer reviews are assigned automatically. If false, the
	// teacher is expected to manually assign peer reviews.
	AutomaticPeerReviews bool `json:"automatic_peer_reviews" yaml:"automatic_peer_reviews" meddler:"automatic_peer_reviews"`

	// Integer representing the amount of reviews each user is assigned. NOTE: This
	// key is NOT present unless you have automatic_peer_reviews set to true.
	PeerReviewCount int `json:"peer_review_count" yaml:"peer_review_count" meddler:"peer_review_count"`

	// String representing a date the reviews are due by. Must be a date that occurs
	// after the default due date. If blank, or date is not after the assignment's
	// due date, the assignment's due date will be used. NOTE: This key is NOT
	// present unless you have automatic_peer_reviews set to true.
	PeerReviewsAssignAt string `json:"peer_reviews_assign_at" yaml:"peer_reviews_assign_at" meddler:"peer_reviews_assign_at"`

	// Boolean representing whether or not members from within the same group on a
	// group assignment can be assigned to peer review their own group's work
	IntraGroupPeerReviews bool              `json:"intra_group_peer_reviews" yaml:"intra_group_peer_reviews" meddler:"intra_group_peer_reviews"`
	HtmlUrl               string            `json:"html_url" yaml:"html_url" meddler:"html_url"`                                        // the URL to the assignment's web page
	IntegrationId         string            `json:"integration_id" yaml:"integration_id" meddler:"integration_id"`                      // (optional, Third Party unique identifier for Assignment)
	IntegrationData       map[string]string `json:"integration_data" yaml:"integration_data" meddler:"integration_data"`                // (optional, Third Party integration data for assignment)
	AnonymousSubmissions  bool              `json:"anonymous_submissions" yaml:"anonymous_submissions" meddler:"anonymous_submissions"` // (Optional) whether anonymous submissions are accepted (applies only to quiz assignments)

	// (Optional) If true, the assignment will be omitted from the student's final
	// grade
	OmitFromFinalGrade bool `json:"omit_from_final_grade" yaml:"omit_from_final_grade" meddler:"omit_from_final_grade"`

	// (Optional) the DiscussionTopic associated with the assignment, if
	// applicable TODO: is it the topic object or its id?
	DiscussionTopic int `json:"discussion_topic" yaml:"discussion_topic" meddler:"discussion_topic"`

	// (Optional) If true, the rubric is directly tied to grading the assignment.
	// Otherwise, it is only advisory. Included if there is an associated rubric.
	UseRubricForGrading bool `json:"use_rubric_for_grading" yaml:"use_rubric_for_grading" meddler:"use_rubric_for_grading"`

	// (Optional) An object describing the basic attributes of the rubric, including
	// the point total. Included if there is an associated rubric.
	RubricSettings map[string]interface{} `json:"rubric_settings" yaml:"rubric_settings" meddler:"rubric_settings"`

	// (Optional) A list of scoring criteria and ratings for each rubric criterion.
	// Included if there is an associated rubric.
	Rubric []RubricCriterion `json:"rubric" yaml:"rubric" meddler:"rubric"`

	// if the requesting user has grading rights, the number of submissions that
	// need grading.
	NeedsGradingCount      int    `json:"needs_grading_count" yaml:"needs_grading_count" meddler:"needs_grading_count"`
	SubmissionsDownloadUrl string `json:"submissions_download_url" yaml:"submissions_download_url" meddler:"submissions_download_url"` // the URL to download all submissions as a zip
}

type ExternalToolTagAttributes struct {
	Id     int    `json:"-" yaml:"-" meddler:"id,pk"`
	Url    string `json:"url" yaml:"url" meddler:"url"`             // URL to the external tool
	NewTab bool   `json:"new_tab" yaml:"new_tab" meddler:"new_tab"` // Whether or not there is a new tab for the external tool
}

type RubricCriterion struct {
	Id              int    `json:"-" yaml:"-" meddler:"id,pk"`
	CanvasId        int    `json:"id" yaml:"id" meddler:"canvas_id"`
	Description     string `json:"description" yaml:"description" meddler:"description"`
	LongDescription string `json:"long_description" yaml:"long_description" meddler:"long_description"`
	Points          int    `json:"points" yaml:"points" meddler:"points"`
}

func getAssignments(db *sql.DB) []*Assignment {
	assignments := make([]*Assignment, 0)
	courses, _ := findCourses(db)
	values := url.Values{}
	values.Add("per_page", "100")
	// TODO: do it for all courses
	courseId := courses[0].CanvasId
	reqUrl := fmt.Sprintf(assignmentsPath, courseId)
	mustGetObject(reqUrl, values, &assignments)
	return assignments
}

func loadAssignmentFromFile(assignmentFilepath string) (*Assignment, error) {
	assignment := new(Assignment)
	err := readYamlFile(assignmentFilepath, assignment)
	return assignment, err
}

func pullAssignments(db *sql.DB) {
	assignmentsMeta := getAssignments(db)
	// TODO: prompt for overwrite, etc.

	for _, assignment := range assignmentsMeta {
		assignment.Pull(db)
	}
}

func (assignment *Assignment) Dump() error {
	metadata, err := yaml.Marshal(assignment)
	if err != nil {
		return err
	}
	assignmentFilePath := fmt.Sprintf("%s/%s.md", assignmentsDir, assignment.Slug())
	return writeFile(assignmentFilePath, string(metadata), assignment.Description)
}

func (assignment *Assignment) Pull(db *sql.DB) {
	courses, _ := findCourses(db)
	// TODO: do it for all courses
	courseId := courses[0].CanvasId
	assignmentFullPath := fmt.Sprintf(assignmentPath, courseId, assignment.CanvasId)

	fmt.Println("Pulling assignment", assignmentFullPath)
	mustGetObject(assignmentFullPath, url.Values{}, assignment)
	assignment.Dump()
}

func (assignment *Assignment) Slug() string {
	return strings.ToLower(strings.ReplaceAll(assignment.Name, " ", "-"))
}
