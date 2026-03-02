# Swiftly

Swiftly is a modern, high-performance e-commerce platform built with a robust Go backend and a type-safe Vue 3 frontend. Designed for speed and scalability, Swiftly leverages a monorepo architecture to streamline development and deployment.

## 🚀 Features (In Development)

- **User Authentication:** Secure Sign In, Register, and Google Login integration.
- **Account Recovery:** Advanced "Forgot Password" flow supporting Email, Username, or Phone Number.
- **Security:** Token Blacklisting with Redis and full-stack input sanitization.
- **Type Safety:** Full-stack type safety from Go backend to TypeScript frontend.

## 🛠 Tech Stack

### Backend
- **Language:** Go 1.25.6
- **Database:** PostgreSQL 17
- **Caching/Security:** Redis (Alpine)
- **Hot Reload:** [Air](https://github.com/air-verse/air)

### Frontend
- **Framework:** Vue 3 (TypeScript, Composition API)
- **State:** Pinia
- **Styling:** Tailwind CSS v4
- **Build Tool:** Vite

---

## 🐳 Docker Workflow (Detailed Guide)

This project uses Docker Compose to orchestrate all services. Using Docker ensures that every developer has the exact same environment.

### 1. Core Commands

| Action | Command |
| :--- | :--- |
| **Start Everything** | `docker compose up -d` |
| **Start & Rebuild** | `docker compose up -d --build` |
| **Stop Everything** | `docker compose down` |
| **Stop & Wipe Data** | `docker compose down -v` |
| **Check Status** | `docker compose ps` |
| **View Logs** | `docker compose logs -f` |

### 2. Managing Individual Services
You don't always need to restart everything. You can target specific services: `backend`, `frontend`, `redis`, or `migrate`.

- **Rebuild only Backend:** `docker compose up -d --build backend`
- **Rebuild only Frontend:** `docker compose up -d --build frontend`
- **Restart Redis:** `docker compose restart redis`

### 3. Database Migrations
Migrations are handled by a dedicated `migrate` service.

- **Run Up Migrations:** `docker compose run --rm migrate up`
- **Rollback (Down):** `docker compose run --rm migrate down`
- **Fix "Dirty" State:** `docker compose run --rm migrate force <version>`

### 4. Interactive CLI Access
Sometimes you need to run commands inside the containers:

- **Redis CLI:** `docker compose exec redis redis-cli`
- **Database (psql):** `docker compose exec db psql -U user -d swiftly`
- **Backend Shell:** `docker compose exec backend sh`

---

## 🏁 Getting Started

### Prerequisites
- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- [Node.js](https://nodejs.org/) & [pnpm](https://pnpm.io/) (for local dev)

### Setup Steps

1.  **Clone & Enter:**
    ```bash
    git clone https://github.com/yourusername/swiftly.git
    cd swiftly
    ```

2.  **Environment Setup:**
    *   `cp backend/.env.sample backend/.env`
    *   `cp frontend/.env.example frontend/.env`
    *   *Make sure to update `GOOGLE_CLIENT_ID` and `JWT_SECRET`.*

3.  **Launch:**
    ```bash
    docker compose up -d --build
    ```

4.  **Migrate:**
    ```bash
    docker compose run --rm migrate
    ```

---

## 📁 Project Structure

```text
.
├── backend/            # Go backend service
│   ├── cmd/            # Entry points (API, Migrations)
│   ├── internal/       # Modular logic (Handler, Service, Repository)
│   └── migrations/     # SQL migration files
├── frontend/           # Vue 3 TypeScript application
│   ├── src/            # Components, stores, types, views
│   └── types/          # Modular TypeScript interfaces
├── docker-compose.yaml # Docker orchestration
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
