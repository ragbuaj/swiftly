# Swiftly

Swiftly is a modern, high-performance e-commerce platform built with a robust Go backend and a reactive Vue 3 frontend. Designed for speed and scalability, Swiftly leverages a monorepo architecture to streamline development and deployment.

## 🚀 Features (In Development)

- **Product Management:** Browse and search through a diverse catalog of items.
- **Shopping Cart:** Seamlessly add, remove, and manage items before purchase.
- **User Authentication:** Secure access for shoppers and administrators.
- **Checkout & Payments:** Integrated payment gateway for secure transactions.
- **Order Tracking:** Real-time updates on order status and delivery.

## 🛠 Tech Stack

### Backend

- **Language:** [Go 1.25.6](https://golang.org/)
- **Framework:** Standard Library (`net/http`)
- **Database:** [PostgreSQL 17](https://www.postgresql.org/)
- **Driver:** [pgx/v5](https://github.com/jackc/pgx) (Native PostgreSQL driver)
- **Migrations:** [golang-migrate](https://github.com/golang-migrate/migrate)
- **Hot Reload:** [Air](https://github.com/air-verse/air)
- **Architecture:** Feature-based modular architecture with Handler, Service, and Repository layers.

### Frontend

- **Framework:** [Vue 3](https://vuejs.org/) (Composition API)
- **Build Tool:** [Vite](https://vitejs.dev/)
- **Package Manager:** [pnpm](https://pnpm.io/)
- **Testing:** [Vitest](https://vitest.dev/) and [@vue/test-utils](https://test-utils.vuejs.org/)

## 📁 Project Structure

```text
.
├── backend/            # Go backend service
│   ├── cmd/
│   │   ├── api/        # API entry point
│   │   ├── migrate/    # Migration tool
│   │   └── seed/       # Database seeder
│   ├── internal/       # Modular features (e.g., user)
│   ├── migrations/     # SQL migration files (.up.sql, .down.sql)
│   ├── Makefile        # Backend management commands
│   ├── .air.toml       # Air configuration for hot reload
│   └── go.mod          # Go module definitions
├── frontend/           # Vue 3 frontend application
│   ├── src/            # Application source code
│   ├── public/         # Static assets
│   ├── vitest.config.js # Vitest configuration
│   └── package.json    # Frontend dependencies and scripts
├── docker-compose.yaml # Docker orchestration
└── README.md           # Project documentation
```

## 🏁 Getting Started

### Prerequisites

- [Go 1.25.6](https://golang.org/dl/) or later
- [Node.js](https://nodejs.org/) and [pnpm](https://pnpm.io/installation)
- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- [PostgreSQL](https://www.postgresql.org/download/) (if running locally without Docker)

### Development with Docker (Recommended)

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/yourusername/swiftly.git
    cd swiftly
    ```

2.  **Configure environment:**
    ```bash
    cp backend/.env.sample backend/.env
    # Update DATABASE_URL in backend/.env to point to your database
    ```

3.  **Run with Hot Reload:**
    ```bash
    docker compose up --build
    ```
    - **Backend:** `http://localhost:8080` (Auto-reloads on file change via Air)
    - **Frontend:** `http://localhost:5173` (HMR enabled)

4.  **Run Migrations:**
    ```bash
    docker compose run --rm migrate
    ```

### Local Development (Without Docker)

1.  **Backend:**
    ```bash
    cd backend
    go run cmd/migrate/main.go up
    go run cmd/api/main.go
    ```

2.  **Frontend:**
    ```bash
    cd frontend
    pnpm install
    pnpm run dev
    ```

## 🗄 Database Management

Database schema is managed via **golang-migrate**.

- **Run migrations:** `make migrate-up` (or `go run cmd/migrate/main.go up`)
- **Rollback migration:** `make migrate-down` (or `go run cmd/migrate/main.go down`)
- **Seed data:** `make seed` (or `go run cmd/seed/main.go`)

## 🧪 Testing

### Backend (Go)
```bash
cd backend
make test
```

### Frontend (Vue 3)
```bash
cd frontend
pnpm test
```

## 📝 API Standards

The project uses a standardized JSON response structure:
```json
{
  "success": true,
  "message": "Action completed successfully",
  "data": { ... },
  "errors": null
}
```

## 📄 License

This project is licensed under the MIT License.
