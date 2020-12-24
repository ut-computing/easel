package main

import (
	// "encoding/json"
	"errors"
	"fmt"
	"net/url"
)

const (
	quizQuestionsTable = "quiz_questions"
	quizQuestionsDir   = "quizzes/questions" // TODO: make configable
)

type QuizQuestion struct {
	Id                int                      `json:"-" yaml:"-" meddler:"id,pk"`
	CanvasId          int                      `json:"id" yaml:"id" meddler:"canvas_id"`
	QuizId            int                      `json:"quiz_id" yaml:"quiz_id" meddler:"quiz_id"`
	Position          int                      `json:"position" yaml:"position" meddler:"position"`                // The order in which the question will be retrieved and displayed.
	QuestionName      string                   `json:"question_name" yaml:"question_name" meddler:"question_name"` // The name of the question.
	QuestionType      string                   `json:"question_type" yaml:"question_type" meddler:"question_type"` // Allowed values: calculated_question, essay_question, file_upload_question, fill_in_multiple_blanks_question, matching_question, multiple_answers_question, multiple_choice_question, multiple_dropdowns_question, numerical_question, short_answer_question, text_only_question, true_false_question
	QuestionText      string                   `json:"question_text" yaml:"question_text" meddler:"question_text"`
	PointsPossible    float64                  `json:"points_possible" yaml:"points_possible" meddler:"points_possible"`
	CorrectComments   string                   `json:"correct_comments" yaml:"correct_comments" meddler:"correct_comments"`       // The comments to display if the student answers the question correctly.
	IncorrectComments string                   `json:"incorrect_comments" yaml:"incorrect_comments" meddler:"incorrect_comments"` // The comments to display if the student answers incorrectly.
	NeutralComments   string                   `json:"neutral_comments" yaml:"neutral_comments" meddler:"neutral_comments"`       // The comments to display regardless of how the student answered.
	Answers           []map[string]interface{} `json:"answers" yaml:"answers" meddler:"answers"`                                  // An array of available answers to display to the student.
	Matches           []map[string]interface{} `json:"matches" yaml:"matches" meddler:"matches"`                                  // The possible matches for matching_question types
	// Answers           []*QuizAnswer      `json:"answers" yaml:"answers" meddler:"answers"`                                  // An array of available answers to display to the student.
	// Matches           []*QuizAnswerMatch `json:"matches" yaml:"matches" meddler:"matches"`                                  // The possible matches for matching_question types
}

type QuizAnswer struct {
	Id       int    `json:"-" yaml:"-" meddler:"id,pk"`
	CanvasId int    `json:"id" yaml:"id" meddler:"canvas_id"`
	Weight   int    `json:"weight" yaml:"weight" meddler:"weight"`
	Text     string `json:"text" yaml:"text" meddler:"text"`
	Comments string `json:"comments" yaml:"comments" meddler:"comments"`
}

// for matching_question question type
type QuizAnswerMatch struct {
	Id       int    `json:"-" yaml:"-" meddler:"id,pk"`
	CanvasId int    `json:"id" yaml:"id" meddler:"canvas_id"`
	MatchId  int    `json:"match_id" yaml:"match_id" meddler:"match_id"`
	Text     string `json:"text" yaml:"text" meddler:"text"`
}

func getQuizQuestions(courseId, quizId int) []*QuizQuestion {
	qqs := make([]*QuizQuestion, 0)
	reqUrl := fmt.Sprintf(quizQuestionsPath, courseId, quizId)
	values := url.Values{}
	values.Add("per_page", "100")
	mustGetObject(reqUrl, values, &qqs)
	return qqs
}

func (qq *QuizQuestion) Dump() error {
	return errors.New("for now, put all quiz questions in one file")
}

func (qq *QuizQuestion) Slug() string {
	// TODO: how to uniquely identify a question?
	return slug(qq.QuestionName)
}

/*
func (qq *QuizQuestion) UnmarshalJSON(raw []byte) error {
	var tmp map[string]interface{}
	err := json.Unmarshal(raw, &tmp)
	if err != nil {
		return err
	}
	qq.CanvasId = int(tmp["id"].(float64))
	qq.QuizId = int(tmp["quiz_id"].(float64))
	qq.Position = int(tmp["position"].(float64))
	qq.QuestionName = tmp["question_name"].(string)
	qq.QuestionType = tmp["question_type"].(string)
	qq.QuestionText = tmp["question_text"].(string)
	qq.PointsPossible = tmp["points_possible"].(float64)
	qq.CorrectComments = tmp["correct_comments"].(string)
	qq.IncorrectComments = tmp["incorrect_comments"].(string)
	qq.NeutralComments = tmp["neutral_comments"].(string)
	qq.Answers = make([]*QuizAnswer, 0)
	for _, answer := range tmp["answers"].([]interface{}) {
		quizAnswer := answer.(QuizAnswer)
		qq.Answers = append(qq.Answers, &quizAnswer)
	}
	qq.Matches = make([]*QuizAnswerMatch, 0)
	for _, answer := range tmp["answers"].([]interface{}) {
		quizAnswerMatch := answer.(QuizAnswerMatch)
		qq.Matches = append(qq.Matches, &quizAnswerMatch)
	}
	return nil
}
*/
