# System Architecture and Data Flow

This document outlines how requests are parsed, queued, executed, and cached inside the Blan Cloud Engine.

## 1. The Concurrency Model (Worker Pool)

To avoid C++ compiler runs blocking the Gin HTTP server, the system uses a strict worker pool pattern.

1. **Ingestion:** A request hits `POST /compile`.
2. **Buffering:** The request is assigned a UUID and pushed into a buffered Go channel (`Capacity: 100`). The server immediately returns a `202 Accepted` to the client.
3. **Execution:** A fixed pool of goroutines (`Workers: 3`) pulls jobs, writes `.bl` code to the `tmpfs` workspace, and executes the `blan` binary via `os/exec` with a strict 3-second context timeout to prevent infinite loops.
4. **State Management:** Job state (`queued` -> `running` -> `completed`/`failed`) is stored in an in-memory Go map protected by a `sync.RWMutex`.

## 2. The Storage Strategy (Workload Specialization)

The platform uses two storage engines, separating relational identity data from high-speed blob storage.

### MySQL (Relational Layer)

Handles ACID-compliant transactions for users and metadata.

- `users` table: Identity and hashed credentials.
- `snippets` table: Saved source code linked via foreign keys.

### [StrataKV](https://github.com/Adityarya11/StrataKV) (Caching Layer)

A custom embedded Log-Structured Merge (LSM) tree database written in Go.

- **Mechanism:** Upon a compilation request, the Go API hashes the source code using SHA-256.
- **Read Path:** It queries StrataKV for the hash. If found, it returns the output instantly, bypassing the worker pool.
- **Write Path:** If the hash is not found, the worker pool executes the code. Upon completion, the result is asynchronously written to StrataKV's MemTable and appended to the Write-Ahead Log (WAL).

## 3. Container Security & Isolation

Production deployment relies on Docker to enforce resource constraints and security.

- **Multi-Stage Builds:** The final container image is stripped of all Go tooling, containing only the static Go server and the stripped `blan` C++ ELF binary in a lightweight Debian environment.
- **Non-Root Execution:** The application runs under a restricted `appuser`.
- **Ephemeral Memory (`tmpfs`):** The `/workspace` directory, where user code is temporarily written to disk for the C++ compiler to read, is mounted directly to RAM. This eliminates disk I/O bottlenecks and ensures no artifacts survive container restarts.

## Architectural Backlog (Future Scope)

These features are designed but deferred until operational pressure requires their implementation:

1. **Job State Persistence:** Migrate the in-memory Worker Pool map (`sync.RWMutex`) to the MySQL `jobs` table or Redis to ensure job state survives server reboots/deployments.
2. **Advanced Sandboxing:** Replace raw `os/exec` with `cgroups` or an unprivileged sub-container (e.g., `nsjail`) to strictly limit CPU cores and RAM allocation per execution, preventing malicious OOM attacks.
3. **Memory Cleanup Daemon:** Implement a background Go `time.Ticker` to prune stale/completed jobs from the in-memory map to prevent memory leaks over long uptimes.
4. **Testcase Evaluation (The "LeetCode" Feature):** Expand the `/compile` route to accept an array of STDIN inputs and expected outputs, running the binary multiple times to evaluate algorithmic accuracy.
5. **Token Rotation:** Implement a dual-token system (short-lived access token + secure HttpOnly refresh token) for seamless client authentication.
