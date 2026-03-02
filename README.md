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
- **Framework:** Standard Library (`net/http`) for lightweight, efficient API development.
- **Database:** [Supabase](https://supabase.com/) (PostgreSQL)
- **Migrations:** [Goose](https://github.com/pressly/goose)
- **Architecture:** Clean Architecture with handlers, services, and repositories.

### Frontend

- **Framework:** [Vue 3](https://vuejs.org/) (Composition API)
- **Build Tool:** [Vite](https://vitejs.dev/)
- **Package Manager:** [pnpm](https://pnpm.io/)
- **Styling:** Vanilla CSS for maximum flexibility and performance.
- **Testing:** [Vitest](https://vitest.dev/) and [@vue/test-utils](https://test-utils.vuejs.org/)

## 📁 Project Structure

```text
.
├── backend/            # Go backend service
│   ├── cmd/api/        # Application entry point & tests
│   ├── cmd/seed/       # Database seeder
│   ├── migrations/     # SQL migration files
│   ├── internal/       # Business logic, models, and repositories
│   ├── Makefile        # Backend management commands
│   └── go.mod          # Go module definitions
├── frontend/           # Vue 3 frontend application
│   ├── src/            # Application source code
│   │   └── components/__tests__/ # Component unit tests
│   ├── public/         # Static assets
│   └── package.json    # Frontend dependencies and scripts
└── GEMINI.md           # Internal development guidelines
```

## 🏁 Getting Started

### Prerequisites

- [Go 1.25.6](https://golang.org/dl/) or later
- [Node.js](https://nodejs.org/) and [pnpm](https://pnpm.io/installation)
- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- [Goose CLI](https://github.com/pressly/goose) (optional, for manual migrations)

### Installation & Development

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/yourusername/swiftly.git
    cd swiftly
    ```

2.  **Backend Setup:**

    ```bash
    cd backend
    cp .env.sample .env # Update with your Supabase credentials and DATABASE_URL
    make migrate-up     # Apply database migrations
    make seed           # (Optional) Seed the database with initial data
    make run            # Start the API server
    ```

    The API will be available at `http://localhost:8080/api/health`.

3.  **Frontend Setup:**

    ```bash
    cd ../frontend
    pnpm install
    pnpm run dev
    ```

    The application will typically be available at `http://localhost:5173`.

## 🗄 Database Management

Backend database schema is managed via **Goose** migrations. All migration files are located in `backend/migrations/`.

- **Run migrations:** `make migrate-up`
- **Rollback migration:** `make migrate-down`
- **Check status:** `make migrate-status`
- **Seed data:** `make seed`
- **Create new migration:** `make migrate-create NAME=your_migration_name`

> **Note:** Ensure `DATABASE_URL` in your `.env` is set to your Supabase direct connection string.

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

## 📝 Logging

### Backend

The backend includes a custom `LoggingMiddleware` that automatically logs all incoming HTTP requests:
`[METHOD] PATH REMOTE_ADDR STATUS_CODE DURATION`

## 📄 License

This project is licensed under the MIT License.
