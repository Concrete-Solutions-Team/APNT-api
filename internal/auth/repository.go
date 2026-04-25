package auth

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(ctx context.Context, u *User) error {
	sql := `INSERT INTO users (username, password, role) VALUES ($1, $2, $3);`

	_, err := r.db.Exec(ctx, sql, u.Username, u.Password, u.Role)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return errors.New("Username is taken")
		}
		return err
	}
	return nil
}

func (r *Repository) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	sql := `SELECT id, username, password, role FROM users WHERE username = $1`

	var u User
	err := r.db.QueryRow(ctx, sql, username).Scan(&u.ID, &u.Username, &u.Password, &u.Role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("User not found")
		}
		return nil, err
	}
	return &u, nil
}

func (r *Repository) CreateSession(ctx context.Context, s *Session) error {
	sql := `INSERT INTO sessions (token, user_id, expires_at) VALUES ($1, $2, $3);`

	_, err := r.db.Exec(ctx, sql, &s.Token, &s.UserId, &s.ExpiresAt)
	return err
}

func (r *Repository) DeleteSession(ctx context.Context, token string) error {
	sql := `DELETE FROM sessions WHERE token = $1`
	_, err := r.db.Exec(ctx, sql, token)
	return err
}

func (r *Repository) GetUserBySessionToken(ctx context.Context, token string) (*User, error) {
	sql := `
		SELECT u.id, u.username, u.password, u.role 
		FROM users u 
		JOIN sessions s ON u.id = s.user_id
		WHERE s.token = $1 AND s.expires_at > $2
	`

	var u User
	err := r.db.QueryRow(ctx, sql, token, time.Now()).Scan(&u.ID, &u.Username, &u.Password, &u.Role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("User not found")
		}
		return nil, err
	}
	return &u, nil
}
