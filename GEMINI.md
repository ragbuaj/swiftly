# Swiftly Project Overview

Swiftly is a full-stack web application featuring a Go-based backend and a Vue 3 frontend. This project is structured as a monorepo with separate directories for backend and frontend services.

## Project Architecture

### Backend (`/backend`)
- **Language:** Go 1.25.6
- **Framework:** Standard library (`net/http`)
- **Structure:**
  - `cmd/api/`: Application entry point.
  - `internal/`: Core logic, models, handlers, and repositories.
- **API Endpoint:** Currently features a health check at `http://localhost:8080/api/health`.

### Frontend (`/frontend`)
- **Framework:** Vue 3 (Composition API)
- **Build Tool:** Vite
- **Package Manager:** pnpm
- **Main Dependencies:** `vue`, `@vitejs/plugin-vue`.

## Getting Started

### Prerequisites
- Go 1.25.6 or later
- Node.js and `pnpm`

### Building and Running

#### Backend
To start the backend server:
```bash
cd backend
go run cmd/api/main.go
```
The server will be available at `http://localhost:8080`.

#### Frontend
To start the frontend development server:
```bash
cd frontend
pnpm install
pnpm run dev
```
The application will typically be available at `http://localhost:5173`.

## Development Conventions

- **Backend:** Follow standard Go idioms. Logic should be organized within the `internal/` directory to enforce encapsulation.
- **Frontend:** Use Vue 3 `<script setup>` syntax for Single File Components (SFCs).
- **Tooling:** Use `vite` for frontend development and `pnpm` for dependency management.
