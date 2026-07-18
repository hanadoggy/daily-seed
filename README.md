# Daily Seed

Daily Seed is an integrated self-management web application designed for tracking Japanese language study, software development learning, and physical/mental habit cultivation (including dopamine detox, journaling, and habit tracking).

## Monorepo Structure

This repository is structured as a monorepo containing the following core services:

- `frontend/`: React + TypeScript frontend powered by Vite, featuring a multi-theme engine (Light, Dark, British Green) using Shadcn UI and Tailwind CSS.
- `backend/`: Go backend providing a robust RESTful API following layered architecture (`handler`, `service`, `repository`).

## Prerequisites

- **Go** (1.20+)
- **Node.js** (18+) with **pnpm**
- **Docker & Docker Compose** (for MongoDB, Mongo Express, and local dev)

## Getting Started

### 1. Full Stack (Docker Compose)
Start all services (MongoDB, Backend, Frontend, Mongo Express) at once:
```bash
docker-compose up -d
```
- Frontend: `http://localhost:5173`
- Backend API: `http://localhost:8080`
- Mongo Express (DB Admin): `http://localhost:8081` (daily / seed)

### 2. Backend Only (Local Dev)
```bash
cd backend
go mod tidy
go run main.go
```
The server will start on `http://localhost:8080`.

### 3. Frontend Only (Local Dev)
```bash
cd frontend
pnpm install
pnpm dev
```
The frontend will be available at `http://localhost:5173`.

## Architecture Highlights
- **Backend:** Idiomatic Go using Gin framework, `slog` for structured logging, and MongoDB Go Driver. Follows layered architecture with separated handlers, services, and repositories. Task migration uses MongoDB transactions for atomicity.
- **Frontend:** React with feature-sliced design principles and global state management via Zustand. Optimistic updates with rollback on API failure.
- **Database:** MongoDB 7 (Replica Set mode) configured with collections for `tasks`, `habits`, and `dailyRecords`.

## Current Status

Phases 0–2 of the project roadmap are completed:

- **Phase 0 (Foundation):** Monorepo structure, CRUD APIs, multi-theme UI, Docker Compose orchestration.
- **Phase 1 (Motivation & Reflection):** Micro-journaling (auto-save), cumulative progress tracker, atomic task migration.
- **Phase 2 (Intelligence):** Multi-select weather & context mode, array-based conditional task filtering, automated migration prompts, task start/end date lifecycle, archiving guards.

The frontend (React/TypeScript) and backend (Go/MongoDB) are fully integrated with consistent data models and stable API contracts.
