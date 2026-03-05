# Swiftly

Swiftly is a modern, high-performance e-commerce platform built with a robust Go backend and a type-safe Vue 3 frontend. Designed for speed and scalability, Swiftly leverages a **Modular Clean Architecture** to streamline development and ensure long-term maintainability.

## 🚀 Features (In Development)

- **User Authentication:** Secure Sign In, Register, Google Login, and Social Auth expansion.
- **Modular Backend:** Dependency injection container, centralized route management, and domain-driven error handling.
- **Product & Store Management:** Multi-vendor support with dynamic notes, store customization, and variant-based product tracking (Planned).
- **Security:** Fail-fast configuration, structured logging (`slog`), input validation, and bot protection.
- **Object Storage:** Integrated MinIO/S3 for high-performance asset management.

## 🛠 Tech Stack

### Backend
- **Language:** Go 1.25.6
- **Database:** PostgreSQL 17 (via `pgx/v5`)
- **Caching:** Redis (with session rotation and replay protection)
- **Object Storage:** MinIO (S3-compatible)
- **Validation:** `go-playground/validator/v10`
- **Logging:** Structured `log/slog`

### Frontend
- **Framework:** Vue 3 (TypeScript, Composition API)
- **UI Library:** Shadcn-Vue & Tailwind CSS v4
- **State:** Pinia
- **Build Tool:** Vite

---

## 🏗 Backend Architecture

Swiftly's backend follows a strict modular structure to separate concerns:
- **`internal/app`**: Global dependency container (Database, Storage, Redis).
- **`internal/api/routes`**: Centralized routing registry.
- **`internal/config`**: Validated fail-fast configuration.
- **`internal/pkg/apperror`**: Domain errors mapped to HTTP.
- **`internal/database`**: Atomic transactions support.

---

## 💻 Development Workflow

### 1. Infrastructure & Backend (Docker)
```bash
docker-compose up -d
```

### 2. Frontend (Local)
```bash
cd frontend
pnpm install
pnpm run dev
```

---

## 🏁 Getting Started

1.  **Clone:** `git clone https://github.com/ragbuaj/swiftly.git`
2.  **Env:** Copy `.env.sample` files and update keys.
3.  **Migrate:** `docker compose up migrate`
4.  **Run:** `docker compose up -d`

## 📝 API Standards
All responses follow a unified JSON format:
```json
{
  "success": true,
  "message": "Action completed successfully",
  "data": { ... },
  "errors": { "code": "STATUS_CODE", "details": "..." }
}
```

## 📄 License
This project is licensed under the MIT License.
