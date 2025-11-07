# CS Dashboard - Backend (Go)

Requirements:
- Go 1.21+

Run:
  make run

Tests:
  make test

API:
- POST /dashboard/v1/auth/login  { email, password }
  -> returns { token, role }
  seeded users:
    - cs@example.com / password  (role: cs)
    - ops@example.com / password (role: operation)

- GET /dashboard/v1/payments
  - protected (Authorization: Bearer <token>)
  - query params: page, per_page, status, search, sort_by (date|amount), order (asc|desc)
  - response: { meta: { ...paged }, summary: { total, completed, processing, failed } }

- PUT /dashboard/v1/payment/:id/review
  - protected, only role operation can call

Notes:
- In-memory store used for simplicity. Replace `internal/storage/memory.go` with SQLite or Redis impl for persistence.
- JWT secret is hardcoded for demo (`verysecretkey`). Use env var in production.
