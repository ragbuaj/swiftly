# Docker Setup Guide

This document explains how to containerize and run the Swiftly project using Docker and Docker Compose.

## 🐳 Prerequisites

- [Docker Desktop](https://www.docker.com/products/docker-desktop/) installed and running.
- [Docker Compose](https://docs.docker.com/compose/install/) (included with Docker Desktop).

## 🚀 Getting Started

To build and start all services (backend and frontend) in one command:

```bash
docker compose up --build
```

- **Backend:** Accessible at `http://localhost:8080`
- **Frontend:** Accessible at `http://localhost:5173`

## 🛠 Project Components

### 1. Backend (Go)
The backend uses a multi-stage Dockerfile (`backend/Dockerfile`) to create a lightweight production image:
- **Build Stage:** Compiles the Go binary using `golang:1.25.6-alpine`.
- **Run Stage:** Runs the binary on a minimal `alpine:latest` image.

### 2. Frontend (Vue 3)
The frontend uses a multi-stage Dockerfile (`frontend/Dockerfile`) to build and serve the application:
- **Build Stage:** Uses `node:20-alpine` and `pnpm` to build the static assets.
- **Production Stage:** Uses `nginx:stable-alpine` to serve the built files.

### 3. Orchestration (Docker Compose)
The `docker-compose.yaml` in the root directory manages both services and connects them via a shared network (`swiftly-network`).

## 🔧 Useful Commands

| Action | Command |
| :--- | :--- |
| Start services (detached mode) | `docker compose up -d` |
| Stop all services | `docker compose down` |
| View logs for a specific service | `docker compose logs -f [backend\|frontend]` |
| Rebuild a single service | `docker compose build [backend\|frontend]` |

## 📦 Environment Variables

Configuration can be customized in the `docker-compose.yaml` or by creating a `.env` file in the root directory:

```env
# Example .env file
PORT=8080
```
