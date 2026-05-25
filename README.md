# blan-backend

This is the Go backend for the [Blan Compiler](https://github.com/Adityarya11/Compiler-Blan). It sits between the playground UI and the compiler runtime, accepts source code, queues execution, and returns results asynchronously. Concurrency is bounded with a worker pool, repeat runs are accelerated via the [StrataKV](https://github.com/Adityarya11/StrataKV) cache, and user/snippet state is persisted in MySQL. The full execution flow, storage split, and deployment model are covered in [readme-in-depth](readme-in-depth.md) and [architecture](architecture.md).

The API is intentionally small and built around a single compile queue, polling for results, and a JWT-protected snippets surface. All routes below are under `/api/v1`.

| Route                       | Type      | Headers                                                           | Body | Description                                                                               |
| --------------------------- | --------- | ----------------------------------------------------------------- | ---- | ----------------------------------------------------------------------------------------- |
| `POST /api/v1/compile`      | Public    | `Content-Type: application/json`                                  | JSON | Enqueue a compile job. Body: `{"source_code":"..."}`. Returns a job id.                   |
| `GET /api/v1/status/:id`    | Public    | None                                                              | None | Fetch job status and output by id. Returns `queued`, `running`, `completed`, or `failed`. |
| `GET /api/v1/health/strata` | Public    | None                                                              | None | Read/write probe for StrataKV readiness and timestamped status.                           |
| `POST /api/v1/signup`       | Public    | `Content-Type: application/json`                                  | JSON | Create a user account. Body: `{"username":"...","email":"...","password":"..."}`.         |
| `POST /api/v1/login`        | Public    | `Content-Type: application/json`                                  | JSON | Issue a JWT token. Body: `{"email":"...","password":"..."}`.                              |
| `POST /api/v1/snippets/`    | Protected | `Authorization: Bearer <token>`, `Content-Type: application/json` | JSON | Save a snippet for the current user. Body: `{"source":"..."}`.                            |
| `GET /api/v1/snippets/`     | Protected | `Authorization: Bearer <token>`                                   | None | List saved snippets for the current user, with count and items.                           |

---

Local run assumes Docker and a reachable MySQL instance; set `DATABASE_URL` and start the compose stack.

```bash
docker compose up --build
```
