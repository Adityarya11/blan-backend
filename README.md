# blan-backend

The execution and orchestration layer for the Blan ecosystem. Sits between the playground UI and the C++ compiler runtime, handling asynchronous job execution, bounded concurrency, JWT authentication, snippet persistence, and execution result caching.

Part of a three-repository stack:

| Repository                                                    | Role                                                       |
| ------------------------------------------------------------- | ---------------------------------------------------------- |
| [Compiler-Blan](https://github.com/Adityarya11/Compiler-Blan) | C++ tree-walking interpreter for the Blan language         |
| **blan-backend**                                              | Go API server — execution orchestration, auth, persistence |
| blan-playground                                               | Next.js frontend — Monaco editor, terminal output, auth UI |

---

## API Reference

All routes are prefixed with `/api/v1`.

### Public Routes

| Method | Route            | Body                                                     | Description                                                                              |
| ------ | ---------------- | -------------------------------------------------------- | ---------------------------------------------------------------------------------------- |
| `POST` | `/compile`       | `{"source_code": "..."}`                                 | Enqueue a compile job. Returns a job ID on cache miss, or immediate output on cache hit. |
| `GET`  | `/status/:id`    | —                                                        | Poll job status by ID. Returns `queued`, `running`, `completed`, or `failed`.            |
| `GET`  | `/health`        | —                                                        | Server liveness probe.                                                                   |
| `GET`  | `/health/strata` | —                                                        | StrataKV read/write readiness probe.                                                     |
| `POST` | `/signup`        | `{"username": "...", "email": "...", "password": "..."}` | Create a user account.                                                                   |
| `POST` | `/login`         | `{"email": "...", "password": "..."}`                    | Issue a JWT. Returns `{"token": "..."}`.                                                 |

### Protected Routes

All protected routes require `Authorization: Bearer <token>` header.

| Method | Route        | Body                | Description                                         |
| ------ | ------------ | ------------------- | --------------------------------------------------- |
| `POST` | `/snippets/` | `{"source": "..."}` | Save a code snippet for the authenticated user.     |
| `GET`  | `/snippets/` | —                   | List all saved snippets for the authenticated user. |

### Response — `POST /compile`

**Cache hit** (source code was previously executed):

```json
{
  "output": "10\nhello\n",
  "cached": true
}
```

**Cache miss** (job enqueued for execution):

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "queued"
}
```

### Response — `GET /status/:id`

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "completed",
  "output": "10\nhello\n"
}
```

---

## Local Development

### Prerequisites

- Go 1.26+
- Docker and Docker Compose
- The `blan` compiler binary (Linux) or `blan.exe` (Windows) placed in the project root

### Environment

Copy `.env.example` to `.env` and fill in the values:

```bash
cp .env.example .env
```

| Variable              | Description                                                          |
| --------------------- | -------------------------------------------------------------------- |
| `DATABASE_URL`        | MySQL DSN                                                            |
| `JWT_SECRET`          | Secret key for signing JWT tokens. Use a long random string.         |
| `ALLOWED_ORIGIN`      | Frontend origin for CORS. Use `http://localhost:3000` for local dev. |
| `APP_PORT`            | Port the server listens on. Defaults to `8080`.                      |
| `MYSQL_ROOT_PASSWORD` | MySQL root password (Docker Compose only).                           |
| `MYSQL_DATABASE`      | MySQL database name (Docker Compose only).                           |
| `MYSQL_USER`          | MySQL user (Docker Compose only).                                    |
| `MYSQL_PASSWORD`      | MySQL user password (Docker Compose only).                           |

### Run with Docker Compose

```bash
docker compose up --build
```

This starts MySQL and the Go server together. The server waits for MySQL to pass its health check before connecting.

### Verify

```bash
curl http://localhost:8080/api/v1/health
curl http://localhost:8080/api/v1/health/strata
```

---

## Tech Stack

| Layer              | Technology                        |
| ------------------ | --------------------------------- |
| Language           | Go 1.26                           |
| HTTP Framework     | Gin                               |
| Authentication     | JWT (golang-jwt/jwt v5)           |
| Persistent Storage | MySQL via GORM                    |
| Cache Layer        | StrataKV (custom LSM-tree engine) |
| Compiler Runtime   | Blan C++ binary via `os/exec`     |
| Infrastructure     | Docker, Docker Compose            |
| Concurrency        | Goroutines + buffered channels    |

---

## Project Structure

```
blan-backend/
├── api/            # HTTP handlers
│   ├── auth.go
│   ├── compile.go
│   ├── health.go
│   └── snippet.go
├── cache/          # StrataKV integration
│   └── strata.go
├── database/       # GORM connection and migration
│   └── db.go
├── middleware/     # JWT auth middleware
│   └── auth.go
├── models/         # GORM schema and request/response DTOs
│   ├── dto.go
│   └── schema.go
├── runner/         # Worker pool and compiler execution
│   ├── executor.go
│   └── pool.go
├── utils/          # JWT lifecycle, cache key generation
│   ├── hash.go
│   └── jwt.go
├── workspace/      # Ephemeral .bl files written before execution
├── main.go
├── Dockerfile
└── docker-compose.yml
```

---

## Further Reading

- [Architecture and system design](architecture.md)
- [Blan language specification and examples](https://github.com/Adityarya11/Compiler-Blan)
- [StrataKV — custom LSM-tree cache engine](https://github.com/Adityarya11/StrataKV)
