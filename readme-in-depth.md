# Blan Backend

Blan Backend is the execution and orchestration layer for the Blan ecosystem. It sits between the playground UI and the compiler runtime, handling asynchronous execution, worker scheduling, authentication, persistence, caching, and containerized deployment. The project is part of a systems-focused stack where each component solves a distinct infrastructure problem rather than a standalone portfolio app.

---

## Related Projects

### Blan Compiler

The core language runtime and compiler implementation written in C++.

<!-- - :contentReference[oaicite:0]{index=0} -->

The compiler executable is treated as a standalone runtime artifact and invoked by the Go backend during execution requests.

---

### StrataKV

Custom embedded LSM-inspired key-value engine written in Go.

<!-- - :contentReference[oaicite:1]{index=1} -->

Used inside the backend as a high-speed execution cache for repeated source code submissions.

---

## Why This Project Exists

Running arbitrary user-submitted code over HTTP is fundamentally a systems problem.

The backend needs to:

- control concurrency,
- avoid resource exhaustion,
- isolate execution,
- persist user state,
- reduce redundant work,
- and remain deployable in a constrained environment.

This repository focuses on solving those operational concerns while keeping the architecture intentionally small and understandable.

The goal was not to build a feature-heavy coding platform.

The goal was to build a clean execution system with deliberate infrastructure decisions.

---

# Current Status

### MVP Complete

The current backend supports:

- asynchronous code execution,
- bounded worker concurrency,
- JWT authentication,
- snippet persistence,
- execution result caching,
- Dockerized deployment,
- health/readiness probing.

The frontend playground is next.

---

# System Overview

```text
                    ┌─────────────────────┐
                    │   Playground UI     │
                    └──────────┬──────────┘
                               │
                               ▼
                    ┌─────────────────────┐
                    │      Gin API        │
                    │ Authentication      │
                    │ Validation          │
                    └──────────┬──────────┘
                               │
                    ┌──────────▼──────────┐
                    │    Worker Queue     │
                    │ Buffered Channel    │
                    └──────────┬──────────┘
                               │
                    ┌──────────▼──────────┐
                    │    Worker Pool      │
                    │ Fixed Goroutines    │
                    └──────────┬──────────┘
                               │
                     ┌─────────▼─────────┐
                     │  blan Binary      │
                     │ Compiler Runtime  │
                     └───────┬───────────┘
                             │
               ┌─────────────┴─────────────┐
               │                           │
               ▼                           ▼
      ┌─────────────────┐        ┌─────────────────┐
      │     MySQL       │        │    StrataKV     │
      │ Persistent Data │        │ Execution Cache │
      └─────────────────┘        └─────────────────┘
```

---

# Core Engineering Decisions

## 1. Bounded Concurrency

The backend intentionally avoids spawning unlimited execution paths.

A fixed worker pool backed by a buffered Go channel is used to control how many compilation jobs execute simultaneously.

### Why this matters

Without bounded concurrency:

- every incoming request could spawn a compiler process,
- CPU usage would spike uncontrollably,
- memory usage would become unstable,
- burst traffic could crash the server.

The worker pool adds controlled backpressure.

Current configuration:

- Worker Count: `3`
- Queue Capacity: `100`

This allows the server to absorb traffic bursts without immediately overwhelming the host machine.

---

## 2. Asynchronous Execution

Compilation is handled asynchronously.

`POST /compile` immediately returns a job identifier while the execution continues in the background.

Clients later poll:

```text
GET /api/v1/status/:id
```

for execution status.

### Why this matters

Keeping HTTP requests open during long-running execution creates:

- connection pressure,
- timeout issues,
- blocked goroutines,
- and poor user experience.

Asynchronous execution decouples request ingestion from execution time.

---

## 3. StrataKV Execution Cache

Repeated executions of identical source code are wasteful.

The backend hashes incoming source code using SHA-256 and checks StrataKV before invoking the compiler binary.

### Cache Hit

If a previous execution exists:

- the compiler is skipped entirely,
- the cached result is returned instantly,
- and execution latency becomes effectively O(1).

### Cache Miss

If no cached entry exists:

- the job enters the worker queue,
- execution proceeds normally,
- and the final result is persisted back into StrataKV.

This turns StrataKV into a lightweight execution acceleration layer.

---

## 4. Persistent Relational State

MySQL is used only for durable relational workloads.

Current persistent entities include:

- users,
- authentication,
- saved snippets.

The project intentionally avoids forcing StrataKV into relational responsibilities.

### Why this split exists

Different workloads require different storage models.

MySQL handles:

- structured querying,
- identity management,
- relational consistency,
- durable persistence.

StrataKV handles:

- fast execution caching,
- repeated lookup acceleration,
- lightweight runtime state.

This separation keeps the architecture clean and operationally sane.

---

## 5. Compiler Runtime Separation

The compiler itself is not embedded into the Go backend.

Instead:

- the backend invokes the standalone compiler binary,
- treats it as an external runtime component,
- and orchestrates execution around it.

This mirrors how real systems often treat runtimes and compilers internally.

The binary becomes an infrastructure artifact rather than application logic.

That separation is important.

---

# Tech Stack

| Layer              | Technology            |
| ------------------ | --------------------- |
| Backend API        | Go                    |
| HTTP Framework     | Gin                   |
| Authentication     | JWT                   |
| Persistent Storage | MySQL                 |
| Cache Layer        | StrataKV              |
| Compiler Runtime   | C++                   |
| Infrastructure     | Docker                |
| Concurrency Model  | Goroutines + Channels |

---

# API Surface

## Public Routes

### Compile Source

```http
POST /api/v1/compile
```

Submits source code for asynchronous execution.

Returns:

- job ID
- queue acknowledgement

---

### Poll Execution Status

```http
GET /api/v1/status/:id
```

Returns:

- queued
- running
- completed
- failed

along with execution output when available.

---

### StrataKV Health Probe

```http
GET /api/v1/health/strata
```

Performs a read/write verification against the cache engine.

Used for deployment readiness checks.

---

# Authentication Routes

### Signup

```http
POST /api/v1/signup
```

---

### Login

```http
POST /api/v1/login
```

Returns JWT token.

---

# Protected Routes

Require:

```text
Authorization: Bearer <token>
```

---

### Save Snippet

```http
POST /api/v1/snippets/
```

---

### Retrieve Snippets

```http
GET /api/v1/snippets/
```

---

# Local Development

## Prerequisites

- Go
- Docker
- MySQL
- Linux-compatible `blan` compiler binary

---

## Run Locally

```bash
git clone <repo-url>

cd blan-backend

docker compose up --build
```

---

## Verify Health

```bash
curl http://localhost:8080/api/v1/health/strata
```

---

# Documentation

Additional engineering documentation:

- `ARCHITECTURE.md`
- `BACKLOG.md`
- `DEPLOYMENT.md`

The architecture document explains:

- worker lifecycle,
- execution flow,
- cache integration,
- storage decisions,
- deployment model,
- and known failure domains.

---

# Scope Boundaries

This project intentionally avoids:

- distributed worker clusters,
- Kubernetes orchestration,
- autoscaling infrastructure,
- multi-node execution scheduling,
- distributed cache replication,
- language virtualization layers,
- AI-assisted compilation,
- and excessive platform features.

The current focus is:

- correctness,
- bounded execution,
- operational clarity,
- deployment simplicity,
- and infrastructure fundamentals.

---

# Known Limitations

Current constraints intentionally left unresolved:

- job state exists only in memory,
- execution isolation is process-level only,
- compiler execution still relies on raw `os/exec`,
- no strict memory quotas per execution,
- queue state does not survive server restart,
- no distributed worker coordination.

These are documented intentionally as future systems work rather than hidden limitations.

---

# Future Work

Planned next steps:

- frontend playground UI,
- stronger execution sandboxing,
- request rate limiting,
- persistent job state,
- metrics and tracing,
- testcase evaluation mode,
- deployment automation,
- container resource enforcement.

---

# Philosophy

This repository was built with a systems-first mindset.

The focus was never maximizing feature count.

The focus was understanding:

- execution pipelines,
- concurrency,
- orchestration,
- persistence,
- caching,
- and operational tradeoffs.

The architecture intentionally stays small enough to reason about while still exposing real infrastructure concerns encountered in production systems.
