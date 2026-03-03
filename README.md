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

4.  **Environment Refinement:**
    *   **Backend:** Ensure `JWT_SECRET` is unique and strong. Set `FRONTEND_URL` to match your client (default: `http://localhost:5173`).
    *   **Frontend:** Verify `VITE_API_URL` points to your Backend API.
    *   *See [Environment Setup](#environment-setup) for details.*

5.  **Launch Stack:**
    ```bash
    docker compose up -d
    ```

---

## 🔐 Environment Setup

Properly configuring your `.env` files is crucial for security and session management.

### Backend (`backend/.env`)
- `JWT_SECRET`: Used for signing session tokens. Generate one using: `openssl rand -base64 32`.
- `REDIS_URL`: Connection string for session storage and blacklisting.
- `REDIS_DB`: Redis database index (default: `0` for dev, `1` is used during tests).
- `FRONTEND_URL`: Used for CORS and OAuth redirects.

### Frontend (`frontend/.env`)
- `VITE_API_URL`: The base URL for all API requests.
- `VITE_GOOGLE_CLIENT_ID`: Required for Google Social Login.

---

## 🔑 Obtaining API Keys & Secrets

Follow these steps to properly configure your security credentials:

### 1. JWT Secret
Generate a unique, strong secret key for token signing:
- **Command:** Run `openssl rand -base64 32` in your terminal.
- **Action:** Copy the resulting string into `JWT_SECRET` in `backend/.env`.

### 2. Google OAuth Credentials (Social Login)
1. Go to [Google Cloud Console](https://console.cloud.google.com/).
2. Create a new project or select an existing one.
3. Navigate to **APIs & Services > Credentials**.
4. Click **Create Credentials > OAuth client ID**.
5. Select **Web application** as the Application type.
6. **Origins:** Add `http://localhost:5173` to *Authorized JavaScript origins*.
7. **Redirects:** Add `http://localhost:8080/api/auth/google/callback` to *Authorized redirect URIs*.
8. Copy the **Client ID** (for both envs) and **Client Secret** (for backend).

### 3. Cloudflare Turnstile (Bot Protection)
1. Log in to the [Cloudflare Dashboard](https://dash.cloudflare.com/).
2. Navigate to **Turnstile** in the sidebar.
3. Click **Add site**.
4. **Domain:** Enter `localhost` for development.
5. Copy the **Site Key** (into `VITE_TURNSTILE_SITE_KEY` in frontend) and **Secret Key** (into `TURNSTILE_SECRET_KEY` in backend).

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
