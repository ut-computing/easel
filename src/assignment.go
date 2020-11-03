package main

type Assignment struct {
	Id int `json:"id"`
	// the name of the assignment
	Name string `json:"name"`
	// the assignment description, in an HTML fragment
	// TODO: Should we use markdown locally?
	Description string `json:"description"`
	// The time at which this assignment was originally created
	CreatedAt string `json:"created_at"`
	// The time at which this assignment was last modified in any way
	UpdatedAt string `json:"updated_at"`
	// the due date for the assignment. returns null if not present. NOTE: If this
	// assignment has assignment overrides, this field will be the due date as it
	// applies to the user requesting information from the API.
	DueAt string `json:"due_at"`
	// the lock date (assignment is locked after this date). returns null if not
	// present. NOTE: If this assignment has assignment overrides, this field will
	// be the lock date as it applies to the user requesting information from the
	// API.
	LockAt string `json:"lock_at"`
	// the unlock date (assignment is unlocked after this date) returns null if not
	// present NOTE: If this assignment has assignment overrides, this field will be
	// the unlock date as it applies to the user requesting information from the
	// API.
	UnlockAt string `json:"unlock_at"`
	// whether this assignment has overrides
	HasOverrides bool `json:"has_overrides"`
	// the ID of the course the assignment belongs to
	CourseId int `json:"course_id"`
	// the URL to the assignment's web page
	HtmlUrl string `json:"html_url"`
	// the URL to download all submissions as a zip
	SubmissionsDownloadUrl string `json:"submissions_download_url"`
	// the ID of the assignment's group
	AssignmentGroupId int `json:"assignment_group_id"`
	// Boolean flag indicating whether the assignment requires a due date based on
	// the account level setting
	DueDateRequired bool `json:"due_date_required"`
	// Allowed file extensions, which take effect if submission_types includes
	// 'online_upload'.
	AllowedExtensions []string `json:"allowed_extensions"`
	// An integer indicating the maximum length an assignment's name may be
	MaxNameLength int `json:"max_name_length"`
	// Boolean flag indicating whether or not Turnitin has been enabled for the
	// assignment. NOTE: This flag will not appear unless your account has the
	// Turnitin plugin available
	TurnitinEnabled bool `json:"turnitin_enabled"`
	// If this is a group assignment, boolean flag indicating whether or not
	// students will be graded individually.
	GradeGroupStudentsIndividually bool `json:"grade_group_students_individually"`
	// Boolean indicating if peer reviews are required for this assignment
	PeerReviews bool `json:"peer_reviews"`
	// Boolean indicating peer reviews are assigned automatically. If false, the
	// teacher is expected to manually assign peer reviews.
	AutomaticPeerReviews bool `json:"automatic_peer_reviews"`
	// Integer representing the amount of reviews each user is assigned. NOTE: This
	// key is NOT present unless you have automatic_peer_reviews set to true.
	PeerReviewCount int `json:"peer_review_count"`
	// String representing a date the reviews are due by. Must be a date that occurs
	// after the default due date. If blank, or date is not after the assignment's
	// due date, the assignment's due date will be used. NOTE: This key is NOT
	// present unless you have automatic_peer_reviews set to true.
	PeerReviewsAssignAt string `json:"peer_reviews_assign_at"`
	// Boolean representing whether or not members from within the same group on a
	// group assignment can be assigned to peer review their own group's work
	IntraGroupPeerReviews bool `json:"intra_group_peer_reviews"`
	// The ID of the assignment’s group set, if this is a group assignment. For
	// group discussions, set group_category_id on the discussion topic, not the
	// linked assignment.
	GroupCategoryId int `json:"group_category_id"`
	// if the requesting user has grading rights, the number of submissions that
	// need grading.
	NeedsGradingCount int `json:"needs_grading_count"`
	// if the requesting user has grading rights and the
	// 'needs_grading_count_by_section' flag is specified, the number of submissions
	// that need grading split out by section. NOTE: This key is NOT present unless
	// you pass the 'needs_grading_count_by_section' argument as true.  ANOTHER
	// NOTE: it's possible to be enrolled in multiple sections, and if a student is
	// setup that way they will show an assignment that needs grading in multiple
	// sections (effectively the count will be duplicated between sections)
	NeedsGradingCountBySection []map[string]string `json:"needs_grading_count_by_section"`
	// the sorting order of the assignment in the group
	Position int `json:"position"`
	// (optional, present if Sync Grades to SIS feature is enabled)
	Post_to_sis bool `json:"post_to_sis"`
	// (optional, Third Party unique identifier for Assignment)
	Integration_id string `json:"integration_id"`
	// (optional, Third Party integration data for assignment)
	Integration_data map[string]string `json:"integration_data"`
	// the maximum points possible for the assignment
	PointsPossible float64 `json:"points_possible"`
	// the types of submissions allowed for this assignment list containing one or
	// more of the following: 'discussion_topic', 'online_quiz', 'on_paper', 'none',
	// 'external_tool', 'online_text_entry', 'online_url', 'online_upload'
	// 'media_recording'
	SubmissionTypes []string `json:"submission_types"`
	// The type of grading the assignment receives; one of 'pass_fail', 'percent',
	// 'letter_grade', 'gpa_scale', 'points'
	GradingType string `json:"grading_type"`
	// Whether the assignment is published
	Published bool `json:"published"`
	// Whether the assignment's 'published' state can be changed to false. Will be
	// false if there are student submissions for the assignment.
	Unpublishable bool `json:"unpublishable"`
	// Whether the assignment is only visible to overrides.
	OnlyVisibleToOverrides bool `json:"only_visible_to_overrides"`
	// (Optional) id of the associated quiz (applies only when submission_types is
	// ['online_quiz'])
	QuizId int `json:"quiz_id"`
	// (Optional) whether anonymous submissions are accepted (applies only to quiz
	// assignments)
	AnonymousSubmissions bool `json:"anonymous_submissions"`
	// (Optional) the DiscussionTopic associated with the assignment, if
	// applicable TODO: is it the topic object or its id?
	DiscussionTopic int `json:"discussion_topic"`
	// (Optional) If true, the rubric is directly tied to grading the assignment.
	// Otherwise, it is only advisory. Included if there is an associated rubric.
	UseRubricForGrading bool `json:"use_rubric_for_grading"`
	// (Optional) An object describing the basic attributes of the rubric, including
	// the point total. Included if there is an associated rubric.
	RubricSettings map[string]interface{} `json:"rubric_settings"`
	// (Optional) A list of scoring criteria and ratings for each rubric criterion.
	// Included if there is an associated rubric.
	Rubric []RubricCriterion `json:"rubric"`
	// (Optional) If true, the assignment will be omitted from the student's final
	// grade
	OmitFromFinalGrade bool `json:"omit_from_final_grade"`
	// The number of submission attempts a student can make for this assignment. -1
	// is considered unlimited.
	AllowedAttempts int `json:"allowed_attempts"`
}

type RubricCriterion struct {
	Id              int    `json:"id"`
	Description     string `json:"description"`
	LongDescription string `json:"long_description"`
	Points          int    `json:"points"`
}

type AssignmentGroup struct {
	Id int
	// the name of the Assignment Group
	Name string `json:"name"`
	// the position of the Assignment Group
	Position int `json:"position"`
	// the weight of the Assignment Group
	GroupWeight int `json:"group_weight"`
}

func pullAssignments(id string) {
	if id == "all" {
		// TODO: get all assignments
	} else {
		// TODO: get single assignment
	}
}
