# Swiftly Project Overview

Swiftly is a full-stack e-commerce web application featuring a Go-based backend, a Vue 3 frontend, and a containerized infrastructure. This project is structured as a monorepo with separate directories for backend and frontend services.

## Full Project Architecture

### 1. Backend (/backend)
- **Language:** Go 1.25.6
- **Framework:** Standard library (
et/http) with custom routing.
- **Authentication:** JWT (JSON Web Tokens) with Access & Refresh tokens, OTP via Email/Phone, Google Social Auth.
- **Structure:**
  - cmd/api/: Application entry point and dependency injection.
  - cmd/migrate/: Database migration runner.
  - internal/database/: Database (PostgreSQL) and Redis connection managers.
  - internal/middleware/: HTTP middlewares (Auth, Rate Limiting, CORS, Logging).
  - internal/user/: User domain module (Handlers, Services, Repositories, Models).
  - internal/pkg/: Shared utilities (auth, captcha, response format, sanitizer, socialauth, storage).
  - migrations/: SQL up/down migration files.

### 2. Frontend (/frontend)
- **Framework:** Vue 3 (Composition API) with TypeScript.
- **UI Library:** Shadcn-Vue (based on Radix Vue/Reka UI) and Tailwind CSS v4.
- **Build Tool:** Vite.
- **Package Manager:** pnpm.
- **State Management:** Pinia (e.g., stores/auth.ts).
- **Routing:** Vue Router with Protected Routes.
- **Architecture:** Component-based UI with API integration layers (src/api).
- **Development Mode:** Run locally on host OS (Windows) for optimal performance.

### 3. Infrastructure & Services (Dockerized)
- **Database:** PostgreSQL for persistent relational data (Users, Products, Orders, etc.).
- **Cache & Session:** Redis for Token Blacklisting (Logout), OTP caching, and Password Reset tokens.
- **Object Storage:** MinIO (S3-compatible) for handling file uploads (Avatars, Product Images).
- **Orchestration:** docker-compose.yaml manages ackend (with Air for hot-reload), edis, minio, and migrate services.

## Getting Started

### Prerequisites
- Docker and Docker Compose
- Go 1.25.6 or later (for local dev)
- Node.js and pnpm (for frontend dev)

### Building and Running
1. Start the infrastructure and backend using Docker:
`ash
docker-compose up -d
`
2. Start the frontend development server locally on Windows:
`ash
cd frontend
pnpm install
pnpm run dev
`

- **Backend API:** http://localhost:8080
- **Frontend App:** http://localhost:5173
- **MinIO Console:** http://localhost:9001 (User/Pass: minioadmin)

## Development Conventions
- **Backend:** Follow standard Go idioms. Logic should be organized within the internal/ directory. Use interface-based design for external services (e.g., storage.Uploader).
- **Frontend:** Use Vue 3 <script setup> syntax for Single File Components (SFCs).

## Development Workflow
- **Testing Mandate**: Always create, update, or perfect test cases whenever adding a new function, modifying an existing feature, or fixing a bug. Never consider a task complete without comprehensive test coverage (positive and negative cases) for the newly introduced logic, and always implement a cleanup mechanism for test data.
- **PowerShell Compatibility**: When executing shell commands on this project (which uses a Windows environment), **never use the && operator** to chain commands. Always use the semicolon ; as the statement separator instead.
- **Reusable Components**: Always strive to create modular, reusable UI components (in the frontend) or helper functions (in the backend) whenever a piece of UI or logic is likely to be used in more than one place. Avoid code duplication by abstracting common patterns.
- **Environment Templates**: Always update `.env.sample` (backend) or `.env.example` (frontend) whenever a new environment variable is introduced. This ensures that the configuration templates remain in sync with the application's requirements.
