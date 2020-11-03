package main

type Quiz struct {
	// the ID of the quiz
	Id int `json:"id"`
	// the title of the quiz
	Title string `json:"title"`
	// the HTTP/HTTPS URL to the quiz
	HtmlUrl string `json:"html_url"`
	// A url that can be visited in the browser with a POST request to preview a
	// quiz as the teacher. Only present when the user may grade
	PreviewUrl string `json:"preview_url"`
	// the description of the quiz
	Description string `json:"description"`
	// type of quiz possible values: 'practice_quiz', 'assignment', 'graded_survey',
	// 'survey'
	QuizType string `json:"quiz_type"`
	// the ID of the quiz's assignment group:
	AssignmentGroupId int `json:"assignment_group_id"`
	// quiz time limit in minutes
	TimeLimit int `json:"time_limit"`
	// shuffle answers for students?
	ShuffleAnswers bool `json:"shuffle_answers"`
	// let students see their quiz responses? possible values: null, 'always',
	// 'until_after_last_attempt'
	HideResults string `json:"hide_results"`
	// show which answers were correct when results are shown? only valid if
	// hide_results=null
	ShowCorrectAnswers bool `json:"show_correct_answers"`
	// restrict the show_correct_answers option above to apply only to the last
	// submitted attempt of a quiz that allows multiple attempts. only valid if
	// show_correct_answers=true and allowed_attempts > 1
	ShowCorrectAnswersLastAttempt bool `json:"show_correct_answers_last_attempt"`
	// when should the correct answers be visible by students? only valid if
	// show_correct_answers=true
	ShowCorrectAnswersAt string `json:"show_correct_answers_at"`
	// prevent the students from seeing correct answers after the specified date has
	// passed. only valid if show_correct_answers=true
	HideCorrectAnswersAt string `json:"hide_correct_answers_at"`
	// prevent the students from seeing their results more than once (right after
	// they submit the quiz)
	OneTimeResults bool `json:"one_time_results"`
	// which quiz score to keep (only if allowed_attempts != 1) possible values:
	// 'keep_highest', 'keep_latest'
	ScoringPolicy string `json:"scoring_policy"`
	// how many times a student can take the quiz -1 = unlimited attempts
	AllowedAttempts int `json:"allowed_attempts"`
	// show one question at a time?
	OneQuestionAtATime bool `json:"one_question_at_a_time"`
	// the number of questions in the quiz
	QuestionCount int `json:"question_count"`
	// The total point value given to the quiz
	PointsPossible int `json:"points_possible"`
	// lock questions after answering? only valid if one_question_at_a_time=true
	CantGoBack bool `json:"cant_go_back"`
	// access code to restrict quiz access
	AccessCode string `json:"access_code"`
	// IP address or range that quiz access is limited to
	IpFilter string `json:"ip_filter"`
	// when the quiz is due
	DueAt string `json:"due_at"`
	// when to lock the quiz
	LockAt string `json:"lock_at"`
	// when to unlock the quiz
	UnlockAt string `json:"unlock_at"`
	// whether the quiz has a published or unpublished draft state.
	Published bool `json:"published"`
	// Whether the assignment's 'published' state can be changed to false. Will be
	// false if there are student submissions for the quiz.
	Unpublishable bool `json:"unpublishable"`
	// Whether or not this is locked for the user.
	LockedForUser bool `json:"locked_for_user"`
	// Whether survey submissions will be kept anonymous (only applicable to
	// 'graded_survey', 'survey' quiz types)
	AnonymousSubmissions bool `json:"anonymous_submissions"`
}

func pullQuizzes(id string) {
	if id == "all" {
		// TODO: get all
	} else {
		// TODO: get single
	}
}
