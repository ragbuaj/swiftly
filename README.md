# Swiftly

Swiftly is a modern, high-performance e-commerce platform built with a robust Go backend and a type-safe Vue 3 frontend. Designed for speed and scalability, Swiftly leverages a monorepo architecture to streamline development and deployment.

## 🚀 Features (In Development)

- **User Authentication:** Secure Sign In, Register, Google Login, and Social Auth expansion (Facebook, X).
- **User Profile:** Dedicated Dashboard for managing personal info, bio, and preferences.
- **Avatar Management:** Deferred upload system with local preview and MinIO/S3 object storage integration.
- **Account Recovery:** Advanced "Forgot Password" flow supporting Email, Username, or Phone Number.
- **Verification:** OTP-based Phone/Email verification loop.
- **Security:** Token Blacklisting with Redis, full-stack input sanitization, and Bot Protection (Cloudflare Turnstile).
- **Type Safety:** Full-stack type safety from Go backend to TypeScript frontend.

## 🛠 Tech Stack

### Backend
- **Language:** Go 1.25.6
- **Database:** PostgreSQL 17
- **Caching/Security:** Redis (Alpine)
- **Object Storage:** MinIO (S3-compatible)
- **Hot Reload:** [Air](https://github.com/air-verse/air)

### Frontend
- **Framework:** Vue 3 (TypeScript, Composition API)
- **UI Library:** Shadcn-Vue (Radix UI) & Tailwind CSS v4
- **State:** Pinia
- **Build Tool:** Vite

---

## 💻 Development Workflow (Hybrid Mode)

For the best developer experience on Windows/macOS, we recommend running the **Backend & Infrastructure in Docker** and the **Frontend locally** on your host OS.

### 1. Infrastructure & Backend (Docker)
Ensure Docker is running, then start the core services:
```bash
docker compose up -d
```
This starts: `PostgreSQL`, `Redis`, `MinIO`, and the `Go API`.

### 2. Frontend (Local)
Run the Vite development server on your host machine for instant HMR and better performance:
```bash
cd frontend
pnpm install
pnpm run dev
```
*App will be available at: `http://localhost:5173`*

### 3. Object Storage (MinIO)
Access the MinIO web console to manage buckets and uploaded files:
*   **URL:** `http://localhost:9001`
*   **Credentials:** `minioadmin` / `minioadmin`

---

## 🏁 Getting Started

### Prerequisites
- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- [Go 1.25.6+](https://go.dev/) (optional, for local testing)
- [Node.js](https://nodejs.org/) & [pnpm](https://pnpm.io/) (required for frontend)

### Setup Steps

1.  **Clone & Enter:**
    ```bash
    git clone https://github.com/ragbuaj/swiftly.git
    cd swiftly
    ```

2.  **Environment Setup:**
    *   `cp backend/.env.sample backend/.env`
    *   `cp frontend/.env.example frontend/.env`
    *   *Update `GOOGLE_CLIENT_ID`, `JWT_SECRET`, and `S3_PUBLIC_URL`.*

3.  **Database Migration:**
    ```bash
    docker compose up migrate
    ```

4.  **Launch Stack:**
    ```bash
    docker compose up -d
    ```

---

## 📁 Project Structure

```text
.
├── backend/            # Go backend service
│   ├── cmd/            # Entry points (API, Migrations)
│   ├── internal/       # Modular logic (Handler, Service, Repository)
│   ├── pkg/            # Shared internal packages (auth, storage, sanitizer)
│   └── migrations/     # SQL migration files
├── frontend/           # Vue 3 TypeScript application
│   ├── src/            # Components, stores, types, views
│   └── public/         # Static assets
├── docker-compose.yaml # Infrastructure orchestration
└── README.md           # Project documentation
```

## 📝 API Standards
Standard JSON response:
```json
{
  "success": true,
  "message": "Action completed successfully",
  "data": { ... }
}
```

## 📄 License
This project is licensed under the MIT License.
