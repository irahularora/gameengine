package models

const (
	AnswerYes = "yes"
	AnswerNo  = "no"
)

type AnswerType string

type UserResponse struct {
	UserID string     `json:"user_id"`
	Answer AnswerType `json:"answer"`
}
