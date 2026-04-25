# APNT API

Backend for a quiz/test platform. Teachers create tests, students join them by code and submit answers. Built with Go, chi, and PostgreSQL.

---

## Stack

- **Go** — `github.com/go-chi/chi/v5` router
- **PostgreSQL** — via `pgx/v5`
- **Docker / Docker Compose** for local dev
- **godotenv** for config

---

## Getting Started

### Prerequisites

- Go 1.25+
- Docker

### Run with Docker

```bash
make app-up
```

This builds the server image and starts both the `db` and `server` containers. The server will be available at `http://localhost:8080`.

```bash
make app-down      # stop and remove containers + volumes
make app-reload    # app-down then app-up
```

### Run locally (without Docker)

1. Copy `.env` and set your `DATABASE_URL` and `PORT`.
2. Run migrations from `db_init/` against your Postgres instance.
3. `go run ./cmd/api`

---

## Configuration

Environment variables (loaded from `.env`):

| Variable       | Example                          | Description              |
|----------------|----------------------------------|--------------------------|
| `PORT`         | `:8080`                          | Listen address           |
| `DATABASE_URL` | `postgresql://user:pass@host/db` | Postgres connection string |

---

## API Routes

### Auth

| Method | Path           | Auth required | Description                                      |
|--------|----------------|---------------|--------------------------------------------------|
| POST   | `/auth/register` | —           | Register a new user. Returns session cookie + token. |
| POST   | `/auth/login`    | —           | Login. Returns session cookie + token.           |
| POST   | `/auth/logout`   | —           | Clears the session cookie.                       |
| GET    | `/auth/me`       | ✓           | Returns the authenticated user's info.           |

Auth uses `HttpOnly` session cookies (`session_token`). The token is also returned in the JSON response body.

**Register / Login body:**
```json
{
  "username": "alice",
  "password": "secret",
  "role": "student"
}
```
`role` is either `"student"` or `"teacher"`.

**`/auth/me` response:**
```json
{
  "id": "uuid",
  "username": "alice",
  "role": "student"
}
```

---

### Tests

All `/tests` routes require authentication (`session_token` cookie).

| Method | Path             | Role        | Description                                         |
|--------|------------------|-------------|-----------------------------------------------------|
| GET    | `/tests/`        | teacher     | List all tests.                                     |
| POST   | `/tests/`        | teacher     | Create a new test.                                  |
| GET    | `/tests/join`    | any         | Join a test by code. Returns test with questions.   |
| POST   | `/tests/submit`  | any         | Submit answers for a test.                          |
| GET    | `/tests/results` | any         | Get results. Teachers see all submissions; students see their own. |

**Create test body:**
```json
{
  "title": "Math Quiz",
  "questions": [
    {
      "id": "q1",
      "text": "What is 2 + 2?",
      "options": ["3", "4", "5"],
      "answer": "4"
    }
  ]
}
```
A random 5-character `generated_code` is assigned automatically.

**Join test:** `GET /tests/join?code=AbC12`

**Submit answers body:**
```json
{
  "test_id": "uuid",
  "answers": {
    "q1": "4",
    "q2": "Paris"
  }
}
```

**Get results:** `GET /tests/results?test_id=uuid`

Submission response:
```json
{
  "id": "uuid",
  "test_id": "uuid",
  "user_id": "uuid",
  "answers": { "q1": "4" },
  "score": 1,
  "total": 1,
  "submitted_at": "2026-01-01T00:00:00Z"
}
```

---

## Database Schema

```
users        — id, username, password, role (student|teacher), created_at
sessions     — token, user_id, expires_at (7 days)
tests        — id, title, questions (JSONB), generated_code, created_by, created_at
submissions  — id, test_id, user_id, answers (JSONB), score, total, submitted_at
```

Migrations are in `db_init/` and run automatically when using Docker Compose.

---

## Notes

- Passwords are stored in plain text. In future we will replace this with hash.
- CORS is configured for `http://localhost:3000` only.
- Session tokens expire after 7 days.
