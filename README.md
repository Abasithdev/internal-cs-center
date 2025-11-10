# Internal CS Center

A full-stack application for managing customer service payments with role-based access control.

## Tech Stack

**Backend:**
- Go 1.21+ (Gin framework)
- JWT authentication
- Swagger documentation
- In-memory storage (easily replaceable with DB)

**Frontend:**
- Vue 3 + TypeScript
- Vite
- Pinia (state management)
- Vue Router

---

## Project Structure

```
.
├── backend/          # Go backend API
│   ├── cmd/server/   # Main application entry
│   ├── internal/     # Internal packages
│   └── docs/         # Swagger documentation
├── frontend/         # Vue 3 frontend
│   └── src/          # Source code
└── package.json      # Root scripts for convenience
```

---

## Initial Setup

### Prerequisites

**Backend:**
- Go 1.21 or higher ([download](https://go.dev/dl/))
- Make (optional, can run go commands directly)

**Frontend:**
- Node.js 20.19.0+ or 22.12.0+ ([download](https://nodejs.org/))
- npm (comes with Node.js)

### Clone & Initialize

```bash
# Clone the repository
git clone https://github.com/Abasithdev/internal-cs-center.git
cd internal-cs-center

# Install backend dependencies
cd backend
go mod download
go mod tidy

# Install frontend dependencies
cd ../frontend
npm install
```

---

## Development (Running Locally)

### Backend

From the `backend` directory:

```bash
# Using Make (recommended)
make run

# Or directly with Go
go run -tags=skip_coverage ./cmd/server
```

The backend will start on **http://localhost:8080**

**Available Make commands:**
- `make run` - Start the server
- `make test` - Run tests with coverage
- `make coverage` - Generate HTML coverage report
- `make build` - Build binary to `bin/server`
- `make docs` - Regenerate Swagger docs

### Frontend

From the `frontend` directory:

```bash
npm run dev
```

The frontend will start on **http://localhost:5173**

**Available npm scripts:**
- `npm run dev` - Start dev server with hot reload
- `npm run build` - Build for production
- `npm run preview` - Preview production build
- `npm run test:unit` - Run unit tests
- `npm run lint` - Lint and auto-fix code
- `npm run type-check` - Check TypeScript types

### Running Both Together

From the **repository root**:

```bash
# Install frontend dependencies
npm run frontend:install

# Start frontend dev server
npm run dev

# In another terminal, start backend
npm run backend:run
```

---

## Environment Variables

### Backend (.env)

Create `backend/.env` (or use existing):

```env
# Server
SERVER_PORT=8080

# JWT
JWT_SECRET=sstttdonttellanyone

# CORS - allowed origins (comma-separated)
ALLOWED_ORIGINS=http://localhost:5173,http://127.0.0.1:5173

# Swagger documentation
SWAGGER_HOST=localhost:8080
SWAGGER_BASEPATH=/dashboard/v1
SWAGGER_TITLE=Internal CS Center API
SWAGGER_VERSION=1.0.0
```

### Frontend (.env)

Create `frontend/.env` (or use existing):

```env
VITE_API_BASE_URL=http://localhost:8080/dashboard/v1
```

**Note:** Vite requires the `VITE_` prefix for environment variables to be exposed to the browser.

---

## API Documentation

### Endpoints

**Authentication:**
- `POST /dashboard/v1/auth/login`
  - Body: `{ "email": "string", "password": "string" }`
  - Returns: `{ "token": "jwt_token", "role": "cs|operation" }`

**Payments (Protected):**
- `GET /dashboard/v1/payments`
  - Headers: `Authorization: Bearer <token>`
  - Query params: `page`, `size`, `status`, `search`
  - Returns: `{ meta: {...}, summary: {...} }`

- `PUT /dashboard/v1/payments/:id/review`
  - Headers: `Authorization: Bearer <token>`
  - Role required: `operation`
  - Marks payment as reviewed

**Health Check:**
- `GET /api` - Simple health check

**Swagger UI:**
- `GET /swagger/index.html` - Interactive API documentation

### Seeded Users

```
Email: cs@example.com
Password: password
Role: cs

Email: ops@example.com
Password: password
Role: operation
```

---

## Building for Production

### Backend

```bash
cd backend

# Build binary
make build

# The binary will be created at bin/server
# Run it:
./bin/server
```

**Or without Make:**
```bash
go build -tags=skip_coverage -o bin/server ./cmd/server
```

### Frontend

```bash
cd frontend

# Type-check and build
npm run build

# The production files will be in dist/
# Preview the build:
npm run preview
```

**Build output:** `frontend/dist/` contains static files ready for deployment.

---

## Deployment

### Backend Deployment

1. Build the binary: `make build`
2. Copy `bin/server` and `.env` to your server
3. Set production environment variables
4. Run: `./bin/server`

**Recommended:**
- Use a process manager (systemd, PM2, supervisor)
- Set `GIN_MODE=release` environment variable
- Use a reverse proxy (nginx, Caddy)
- Replace in-memory storage with a database

### Frontend Deployment

1. Build: `npm run build`
2. Deploy the `dist/` folder to:
   - Static hosting (Vercel, Netlify, Cloudflare Pages)
   - CDN + object storage (S3 + CloudFront)
   - Traditional web server (nginx, Apache)

**Example nginx config:**
```nginx
server {
    listen 80;
    server_name yourdomain.com;
    root /path/to/dist;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    # Proxy API requests to backend
    location /dashboard/v1 {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

---

## Testing

### Backend Tests

```bash
cd backend
make test

# With coverage report
make coverage
```

### Frontend Tests

```bash
cd frontend
npm run test:unit
```

---

## Troubleshooting

### Backend won't start

**Error: "build constraints exclude all Go files"**
- Solution: Use `go run -tags=skip_coverage ./cmd/server` or `make run`

**Error: "bad origin: origins must contain http://"**
- Solution: Check `ALLOWED_ORIGINS` in `.env` includes `http://` or `https://`

**Error: "cannot find package"**
- Solution: Run `go mod tidy` and `go mod download`

### Frontend won't start

**Error: "Cannot find module"**
- Solution: Run `npm install` in the frontend directory

**Error: "Port 5173 already in use"**
- Solution: Kill the process or change port in `vite.config.ts`

**API calls fail with CORS error**
- Solution: Ensure backend is running and `ALLOWED_ORIGINS` includes frontend URL

---

## Development Notes

- **In-memory storage:** Data resets on server restart. For persistence, implement a database adapter in `internal/storage/`.
- **JWT secret:** Change `JWT_SECRET` in production to a strong random value.
- **CORS:** The backend CORS middleware is configured via `ALLOWED_ORIGINS` environment variable.
- **Vite proxy:** The frontend dev server proxies `/dashboard/v1` requests to avoid CORS during development (configured in `vite.config.ts`).

---

## Contributors

Abdul Basith
