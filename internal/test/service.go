package test

import (
	"context"
	"math/rand/v2"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateTest(ctx context.Context, title string, questions []Question, createdBy *uuid.UUID) (*Test, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 5)
	for i := range b {
		b[i] = charset[rand.IntN(len(charset))]
	}
	test := &Test{
		Title:         title,
		Questions:     questions,
		GeneratedCode: string(b),
		CreatedBy:     createdBy,
	}

	err := s.repo.CreateTest(ctx, test)
	if err != nil {
		return nil, err
	}

	return test, nil
}

func (s *Service) GetTest(ctx context.Context, id uuid.UUID) (*Test, error) {
	return s.repo.GetTestByID(ctx, id)
}

func (s *Service) UpdateTest(ctx context.Context, id uuid.UUID, title *string, questions *[]Question) (*Test, error) {
	test, err := s.repo.GetTestByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if title != nil {
		test.Title = *title
		test.GeneratedCode = "" + *title
	}
	if questions != nil {
		test.Questions = *questions
	}

	err = s.repo.UpdateTest(ctx, test)
	if err != nil {
		return nil, err
	}

	return test, nil
}

func (s *Service) DeleteTest(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteTest(ctx, id)
}

func (s *Service) ListTests(ctx context.Context) ([]*Test, error) {
	return s.repo.ListTests(ctx)
}

func (s *Service) JoinTest(ctx context.Context, code string) (*Test, error) {
	return s.repo.GetTestByCode(ctx, code)
}

func (s *Service) SubmitTest(ctx context.Context, testID uuid.UUID, userID uuid.UUID, answers map[string]string) (*Submission, error) {
	test, err := s.repo.GetTestByID(ctx, testID)
	if err != nil {
		return nil, err
	}

	score := 0
	total := len(test.Questions)
	for _, q := range test.Questions {
		if ans, ok := answers[q.ID]; ok && ans == q.Answer {
			score++
		}
	}

	submission := &Submission{
		TestID:  testID,
		UserID:  userID,
		Answers: answers,
		Score:   score,
		Total:   total,
	}

	err = s.repo.CreateSubmission(ctx, submission)
	if err != nil {
		return nil, err
	}

	return submission, nil
}

func (s *Service) GetResults(ctx context.Context, testID uuid.UUID) ([]*Submission, error) {
	return s.repo.GetSubmissionsByTestID(ctx, testID)
}

func (s *Service) GetMyResult(ctx context.Context, testID, userID uuid.UUID) (*Submission, error) {
	return s.repo.GetSubmissionByUserAndTest(ctx, userID, testID)
}
