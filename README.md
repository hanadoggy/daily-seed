# Daily Seed

Daily Seed is an integrated self-management web application designed for tracking Japanese language study, software development learning, and physical/mental habit cultivation (including dopamine detox, journaling, and habit tracking).

## Monorepo Structure

This repository is structured as a monorepo containing the following core services:

- `frontend/`: React + TypeScript frontend powered by Vite, featuring a multi-theme engine (Light, Dark, British Green) using Shadcn UI and Tailwind CSS.
- `backend/`: Go backend providing a robust RESTful API following layered architecture (`handler`, `service`, `repository`).
- `mongo-init/`: Initialization scripts and configuration for the MongoDB database.

## Prerequisites

- **Go** (1.20+)
- **Node.js** (18+)
- **Docker & Docker Compose** (for MongoDB and local dev)

## Getting Started

### 1. Database Setup
Start the local MongoDB instance using Docker Compose:
```bash
docker-compose up -d
```

### 2. Backend Setup
Navigate to the backend directory and run the server:
```bash
cd backend
go mod tidy
go run main.go
```
The server will start on `http://localhost:8080`.

### 3. Frontend Setup
Navigate to the frontend directory, install dependencies, and start the development server:
```bash
cd frontend
npm install
npm run dev
```
The frontend will be available at `http://localhost:5173`.

## Architecture Highlights
- **Backend:** Idiomatic Go using Gin framework, `slog` for structured logging, and MongoDB Go Driver. Follows layered architecture with separated handlers, services, and repositories.
- **Frontend:** React with feature-sliced design principles and global state management via Zustand.
- **Database:** MongoDB configured with collections for `Tasks`, `Habits`, and `DailyRecords`.

## Current Status
Phase 0 of the project roadmap is completed. The codebase provides the foundational MVP for habit tracking and daily record management.
