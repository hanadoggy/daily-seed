# Daily Seed

Daily Seed is an integrated self-management web application designed for tracking Japanese language study, software development learning, and physical/mental habit cultivation (including dopamine detox, journaling, and habit tracking).

## Monorepo Structure

This repository is structured as a monorepo containing the following core services:

- `frontend/`: React 19 + TypeScript frontend powered by Vite, featuring a multi-theme engine (Light, Dark, British Green) using Shadcn UI, Tailwind CSS 4, and Zustand 5 for state management.
- `backend/`: Go backend providing a robust RESTful API following a pragmatic 2-Layer architecture (`handler` / `store`) per feature domain (`daily`, `task`, `habit`, `analytics`).

## Prerequisites

- **Go** (1.22+)
- **Node.js** (18+) with **pnpm**
- **Docker & Docker Compose** (for MongoDB, Mongo Express, and local dev)

## Getting Started

### 1. Full Stack (Docker Compose)
Start all services (MongoDB Replica Set, Backend, Frontend, Mongo Express) at once:
```bash
docker-compose up -d
```
- Frontend: `http://localhost:5173`
- Backend API: `http://localhost:8080`
- Mongo Express (DB Admin): `http://localhost:8081` (daily / seed)

### 2. Backend Only (Local Dev & Tests)
```bash
cd backend
go mod tidy
go run main.go
```
The server will start on `http://localhost:8080`.

To run backend tests (Slice & Integration tests powered by `testcontainers-go`):
```bash
go test ./...
```

### 3. Frontend Only (Local Dev & Tests)
```bash
cd frontend
pnpm install
pnpm dev
```
The frontend will be available at `http://localhost:5173`.

To run frontend tests (Vitest + React Testing Library):
```bash
pnpm test --run
```

## Architecture Highlights

- **Backend Architecture:** Idiomatic Go using Gin framework, `slog` for structured logging, and MongoDB Go Driver. Follows a pragmatic 2-Layer package-by-feature architecture (`handler` and `store` per domain inside `internal/<domain>`). Employs Go's *"Accept interfaces, return structs"* idiom with consumer-side interfaces. Tasks support customizable `unit` fields and MongoDB transactions for atomic task migration (with status guards for concurrency control). Input validation is explicitly enforced via struct `Validate() error` methods.
- **Backend Testing:** Fast, reliable Slice and Integration testing using `testcontainers-go` (real ephemeral MongoDB) and `httptest` without brittle mock frameworks.
- **Frontend Architecture:** React 19 with feature-sliced design principles, custom hooks, and global state management via Zustand 5. Implements optimistic updates with rollback logic on API failure, client-side request serialization queue for PATCH requests, in-flight task migration tracking, and pure TypeScript validation functions (`lib/validation.ts`).
- **Database:** MongoDB 7 (Replica Set mode) configured with indexes and collections for `tasks`, `habits`, and `dailyRecords`.

## API Endpoints Summary

| Group | Method | Endpoint | Description |
|---|---|---|---|
| System | `GET` | `/health` | Server health check |
| Tasks | `GET` | `/api/v1/tasks` | List all tasks |
| Tasks | `GET` | `/api/v1/tasks/progress` | Get cumulative task progress |
| Tasks | `GET` | `/api/v1/tasks/:id` | Get task by ID |
| Tasks | `POST` | `/api/v1/tasks` | Create task |
| Tasks | `PUT` | `/api/v1/tasks/:id` | Update task |
| Tasks | `DELETE` | `/api/v1/tasks/:id` | Soft-delete / archive task |
| Tasks | `POST` | `/api/v1/tasks/:id/migrate` | Atomic task migration |
| Habits | `GET` | `/api/v1/habits` | List all habits |
| Habits | `GET` | `/api/v1/habits/:id` | Get habit by ID |
| Habits | `POST` | `/api/v1/habits` | Create habit |
| Habits | `PUT` | `/api/v1/habits/:id` | Update habit |
| Habits | `DELETE` | `/api/v1/habits/:id` | Archive habit |
| Daily | `GET` | `/api/v1/daily/exists` | Get dates with existing daily records |
| Daily | `GET` | `/api/v1/daily/:date` | Get daily record for date |
| Daily | `PATCH` | `/api/v1/daily/:date` | Update daily record (context, tasks, habits, journal) |
| Analytics | `GET` | `/api/v1/analytics/heatmap` | Annual activity heatmap data |
| Analytics | `GET` | `/api/v1/analytics/summary` | Weekly/Monthly summary statistics |
| Analytics | `GET` | `/api/v1/analytics/streaks` | Habit streak statistics & milestone achievements |

## Current Status

Phases 0–3 of the project roadmap as well as Phase 4 Input Validation & Concurrency Control are fully completed:

- **Phase 0 (Foundation):** Monorepo structure, CRUD APIs, multi-theme UI, Docker Compose orchestration, and 2-Layer backend refactoring.
- **Phase 1 (Motivation & Reflection):** Micro-journaling (debounced auto-save), cumulative progress tracker, atomic task migration.
- **Phase 2 (Intelligence & Lifecycle):** Multi-select weather & context modes, conditional task filtering, automated migration prompts, task start/end date lifecycle, and archiving guards.
- **Phase 3 (Analytics & Insights):** GitHub-style Heatmap Dashboard, Weekly & Monthly Summary Views (task & habit completion rates, mode distributions, journal timeline), Habit Streak Tracking & Statistics (current/longest streaks, milestone celebration modals), calendar record indicators, and Admin Mode.
- **Phase 4 (Validation & Concurrency):** Form input validation via pure TypeScript functions (`lib/validation.ts`) & backend `Validate() error` methods, concurrency guards for atomic task migration (409 Conflict), and date-based PATCH request serialization queue.

