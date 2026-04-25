package test

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateTest(ctx context.Context, t *Test) error {
	sql := `INSERT INTO tests (title, questions, generated_code, created_by) VALUES ($1, $2, $3, $4) RETURNING id, created_at`

	questionsJSON, err := json.Marshal(t.Questions)
	if err != nil {
		return err
	}

	err = r.db.QueryRow(ctx, sql, t.Title, questionsJSON, t.GeneratedCode, t.CreatedBy).Scan(&t.ID, &t.CreatedAt)
	return err
}

func (r *Repository) GetTestByID(ctx context.Context, id uuid.UUID) (*Test, error) {
	sql := `SELECT id, title, questions, generated_code, created_at, created_by FROM tests WHERE id = $1`

	var t Test
	var questionsJSON []byte
	err := r.db.QueryRow(ctx, sql, id).Scan(&t.ID, &t.Title, &questionsJSON, &t.GeneratedCode, &t.CreatedAt, &t.CreatedBy)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("Test not found")
		}
		return nil, err
	}

	err = json.Unmarshal(questionsJSON, &t.Questions)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *Repository) UpdateTest(ctx context.Context, t *Test) error {
	sql := `UPDATE tests SET title = $1, questions = $2, generated_code = $3 WHERE id = $4`

	questionsJSON, err := json.Marshal(t.Questions)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, sql, t.Title, questionsJSON, t.GeneratedCode, t.ID)
	return err
}

func (r *Repository) DeleteTest(ctx context.Context, id uuid.UUID) error {
	sql := `DELETE FROM tests WHERE id = $1`
	_, err := r.db.Exec(ctx, sql, id)
	return err
}

func (r *Repository) GetTestByCode(ctx context.Context, code string) (*Test, error) {
	sql := `SELECT id, title, questions, generated_code, created_at, created_by FROM tests WHERE generated_code = $1`

	var t Test
	var questionsJSON []byte
	err := r.db.QueryRow(ctx, sql, code).Scan(&t.ID, &t.Title, &questionsJSON, &t.GeneratedCode, &t.CreatedAt, &t.CreatedBy)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("Test not found")
		}
		return nil, err
	}

	err = json.Unmarshal(questionsJSON, &t.Questions)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *Repository) CreateSubmission(ctx context.Context, s *Submission) error {
	sql := `INSERT INTO submissions (test_id, user_id, answers, score, total) VALUES ($1, $2, $3, $4, $5) RETURNING id, submitted_at`

	answersJSON, err := json.Marshal(s.Answers)
	if err != nil {
		return err
	}

	err = r.db.QueryRow(ctx, sql, s.TestID, s.UserID, answersJSON, s.Score, s.Total).Scan(&s.ID, &s.SubmittedAt)
	return err
}

func (r *Repository) GetSubmissionsByTestID(ctx context.Context, testID uuid.UUID) ([]*Submission, error) {
	sql := `SELECT id, test_id, user_id, answers, score, total, submitted_at FROM submissions WHERE test_id = $1 ORDER BY submitted_at DESC`

	rows, err := r.db.Query(ctx, sql, testID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var submissions []*Submission
	for rows.Next() {
		var s Submission
		var answersJSON []byte
		err := rows.Scan(&s.ID, &s.TestID, &s.UserID, &answersJSON, &s.Score, &s.Total, &s.SubmittedAt)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(answersJSON, &s.Answers)
		if err != nil {
			return nil, err
		}

		submissions = append(submissions, &s)
	}

	return submissions, nil
}

func (r *Repository) GetSubmissionByUserAndTest(ctx context.Context, userID, testID uuid.UUID) (*Submission, error) {
	sql := `SELECT id, test_id, user_id, answers, score, total, submitted_at FROM submissions WHERE user_id = $1 AND test_id = $2`

	var s Submission
	var answersJSON []byte
	err := r.db.QueryRow(ctx, sql, userID, testID).Scan(&s.ID, &s.TestID, &s.UserID, &answersJSON, &s.Score, &s.Total, &s.SubmittedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("Submission not found")
		}
		return nil, err
	}

	err = json.Unmarshal(answersJSON, &s.Answers)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *Repository) ListTests(ctx context.Context) ([]*Test, error) {
	sql := `SELECT id, title, questions, generated_code, created_at, created_by FROM tests ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tests []*Test
	for rows.Next() {
		var t Test
		var questionsJSON []byte
		err := rows.Scan(&t.ID, &t.Title, &questionsJSON, &t.GeneratedCode, &t.CreatedAt, &t.CreatedBy)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(questionsJSON, &t.Questions)
		if err != nil {
			return nil, err
		}

		tests = append(tests, &t)
	}

	return tests, nil
}
