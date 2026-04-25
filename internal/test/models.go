package test

import (
	"time"

	"github.com/google/uuid"
)

type Question struct {
	ID      string   `json:"id"`
	Text    string   `json:"text"`
	Options []string `json:"options,omitempty"`
	Answer  string   `json:"answer,omitempty"`
}

type Test struct {
	ID            uuid.UUID  `json:"id"`
	Title         string     `json:"title"`
	Questions     []Question `json:"questions"`
	GeneratedCode string     `json:"generated_code,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	CreatedBy     *uuid.UUID `json:"created_by,omitempty"`
}

type Submission struct {
	ID          uuid.UUID         `json:"id"`
	TestID      uuid.UUID         `json:"test_id"`
	UserID      uuid.UUID         `json:"user_id"`
	Answers     map[string]string `json:"answers"`
	Score       int               `json:"score"`
	Total       int               `json:"total"`
	SubmittedAt time.Time         `json:"submitted_at"`
}

type createTestRequest struct {
	Title     string     `json:"title"`
	Questions []Question `json:"questions"`
}

type updateTestRequest struct {
	Title     *string     `json:"title,omitempty"`
	Questions *[]Question `json:"questions,omitempty"`
}

type joinTestRequest struct {
	Code string `json:"code"`
}

type submitTestRequest struct {
	TestID  string            `json:"test_id"`
	Answers map[string]string `json:"answers"`
}
