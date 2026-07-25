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

- **Backend Architecture:** Idiomatic Go using Gin framework, `slog` for structured logging, and MongoDB Go Driver. Follows a pragmatic 2-Layer package-by-feature architecture (`handler` and `store` per domain inside `internal/<domain>`). Employs Go's *"Accept interfaces, return structs"* idiom with consumer-side interfaces. Task migration uses MongoDB transactions for atomicity.
- **Backend Testing:** Fast, reliable Slice and Integration testing using `testcontainers-go` (real ephemeral MongoDB) and `httptest` without brittle mock frameworks.
- **Frontend Architecture:** React 19 with feature-sliced design principles, custom hooks, and global state management via Zustand 5. Implements optimistic updates with rollback logic on API failure.
- **Database:** MongoDB 7 (Replica Set mode) configured with indexes and collections for `tasks`, `habits`, and `dailyRecords`.

## Current Status

Phases 0–3 of the project roadmap are fully completed:

- **Phase 0 (Foundation):** Monorepo structure, CRUD APIs, multi-theme UI, Docker Compose orchestration, and 2-Layer backend refactoring.
- **Phase 1 (Motivation & Reflection):** Micro-journaling (debounced auto-save), cumulative progress tracker, atomic task migration.
- **Phase 2 (Intelligence & Lifecycle):** Multi-select weather & context modes, conditional task filtering, automated migration prompts, task start/end date lifecycle, and archiving guards.
- **Phase 3 (Analytics & Insights):** GitHub-style Heatmap Dashboard, calendar record indicators, and Admin Mode.
