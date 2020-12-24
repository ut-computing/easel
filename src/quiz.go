package main

import (
	"database/sql"
	"fmt"
	"net/url"

	"gopkg.in/yaml.v2"
)

const (
	quizzesTable        = "quizzes"
	quizzesDir          = "quizzes" // TODO: make configable
	quizQuestionsSuffix = "_questions"
)

type Quiz struct {
	Id                 int     `json:"-" yaml:"-" meddler:"id,pk"`
	CanvasId           int     `json:"id" yaml:"id" meddler:"canvas_id"`
	Title              string  `json:"title" yaml:"title" meddler:"title"`
	Description        string  `json:"description" yaml:"-" meddler:"description"` // the description of the quiz
	QuestionCount      int     `json:"question_count" yaml:"question_count" meddler:"question_count"`
	PointsPossible     float64 `json:"points_possible" yaml:"points_possible" meddler:"points_possible"`
	HtmlUrl            string  `json:"html_url" yaml:"html_url" meddler:"html_url"` // the HTTP/HTTPS URL to the quiz
	AccessCode         string  `json:"access_code" yaml:"access_code" meddler:"access_code"`
	OneQuestionAtATime bool    `json:"one_question_at_a_time" yaml:"one_question_at_a_time" meddler:"one_question_at_a_time"`
	CantGoBack         bool    `json:"cant_go_back" yaml:"cant_go_back" meddler:"cant_go_back"` // lock questions after answering? only valid if one_question_at_a_time=true
	TimeLimit          int     `json:"time_limit" yaml:"time_limit" meddler:"time_limit"`       // quiz time limit in minutes
	ShuffleAnswers     bool    `json:"shuffle_answers" yaml:"shuffle_answers" meddler:"shuffle_answers"`
	IpFilter           string  `json:"ip_filter" yaml:"ip_filter" meddler:"ip_filter"`
	DueAt              string  `json:"due_at" yaml:"due_at" meddler:"due_at"`
	LockAt             string  `json:"lock_at" yaml:"lock_at" meddler:"lock_at"`
	UnlockAt           string  `json:"unlock_at" yaml:"unlock_at" meddler:"unlock_at"`
	Published          bool    `json:"published" yaml:"published" meddler:"published"`
	AssignmentGroupId  int     `json:"assignment_group_id" yaml:"assignment_group_id" meddler:"assignment_group_id"` // the ID of the quiz's assignment group:

	// how many times a student can take the quiz (-1 = unlimited attempts)
	AllowedAttempts int `json:"allowed_attempts" yaml:"allowed_attempts" meddler:"allowed_attempts"`

	// let students see their quiz responses? possible values: null, 'always',
	// 'until_after_last_attempt'
	HideResults string `json:"hide_results" yaml:"hide_results" meddler:"hide_results"`

	// show which answers were correct when results are shown? only valid if
	// hide_results=null
	ShowCorrectAnswers bool `json:"show_correct_answers" yaml:"show_correct_answers" meddler:"show_correct_answers"`

	// restrict the show_correct_answers option above to apply only to the last
	// submitted attempt of a quiz that allows multiple attempts. only valid if
	// show_correct_answers=true and allowed_attempts > 1
	ShowCorrectAnswersLastAttempt bool `json:"show_correct_answers_last_attempt" yaml:"show_correct_answers_last_attempt" meddler:"show_correct_answers_last_attempt"`

	// when should the correct answers be visible by students? only valid if
	// show_correct_answers=true
	ShowCorrectAnswersAt string `json:"show_correct_answers_at" yaml:"show_correct_answers_at" meddler:"show_correct_answers_at"`

	// prevent the students from seeing correct answers after the specified date has
	// passed. only valid if show_correct_answers=true
	HideCorrectAnswersAt string `json:"hide_correct_answers_at" yaml:"hide_correct_answers_at" meddler:"hide_correct_answers_at"`

	// prevent the students from seeing their results more than once (right after
	// they submit the quiz)
	OneTimeResults bool `json:"one_time_results" yaml:"one_time_results" meddler:"one_time_results"`

	// which quiz score to keep (only if allowed_attempts != 1) possible values:
	// 'keep_highest', 'keep_latest'
	ScoringPolicy string `json:"scoring_policy" yaml:"scoring_policy" meddler:"scoring_policy"`

	// type of quiz possible values: 'practice_quiz', 'assignment', 'graded_survey',
	// 'survey'
	QuizType string `json:"quiz_type" yaml:"quiz_type" meddler:"quiz_type"`

	// A url that can be visited in the browser with a POST request to preview a
	// quiz as the teacher. Only present when the user may grade
	PreviewUrl string `json:"preview_url" yaml:"preview_url" meddler:"preview_url"`

	// Whether survey submissions will be kept anonymous (only applicable to
	// 'graded_survey', 'survey' quiz types)
	AnonymousSubmissions bool            `json:"anonymous_submissions" yaml:"anonymous_submissions" meddler:"anonymous_submissions"`
	QuizQuestions        []*QuizQuestion `json:"-" yaml:"-" meddler:"-"`
}

func getQuizzes(db *sql.DB) []*Quiz {
	quizzes := make([]*Quiz, 0)
	courses, _ := findCourses(db)
	values := url.Values{}
	values.Add("per_page", "100")
	// TODO: do it for all courses
	courseId := courses[0].CanvasId
	reqUrl := fmt.Sprintf(quizzesPath, courseId)
	mustGetObject(reqUrl, values, &quizzes)
	// get the quiz's questions while we're here
	for _, quiz := range quizzes {
		quiz.QuizQuestions = getQuizQuestions(courseId, quiz.CanvasId)
	}
	return quizzes
}

func pullQuizzes(db *sql.DB) {
	quizzes := getQuizzes(db)
	for _, quiz := range quizzes {
		quiz.Dump()
	}
}

func (quiz *Quiz) Dump() error {
	metadata, err := yaml.Marshal(quiz)
	if err != nil {
		return err
	}

	quizFilePath := fmt.Sprintf("%s/%s.md", quizzesDir, quiz.Slug())
	err = writeFile(quizFilePath, string(metadata), quiz.Description)
	if err != nil {
		return err
	}

	// Dump the questions too, for now put all in one file
	qqs, err := yaml.Marshal(quiz.QuizQuestions)
	if err != nil {
		return err
	}
	quizQuestionsFilePath := fmt.Sprintf("%s/%s%s.md", quizzesDir, quiz.Slug(), quizQuestionsSuffix)
	return writeYamlFile(quizQuestionsFilePath, string(qqs))
}

func (quiz *Quiz) Pull(db *sql.DB) error {
	// get the quiz questions
	courses, _ := findCourses(db)
	courseId := courses[0].CanvasId
	quiz.QuizQuestions = getQuizQuestions(courseId, quiz.CanvasId)
	// pull the quiz and then dump it and its questions
	return pullComponent(db, quizPath, quiz.CanvasId, quiz)
}

func (quiz *Quiz) Slug() string {
	return slug(quiz.Title)
}
