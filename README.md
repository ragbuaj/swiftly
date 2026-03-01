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
- **Architecture:** Clean Architecture with handlers, services, and repositories.

### Frontend
- **Framework:** [Vue 3](https://vuejs.org/) (Composition API)
- **Build Tool:** [Vite](https://vitejs.dev/)
- **Package Manager:** [pnpm](https://pnpm.io/)
- **Styling:** Vanilla CSS for maximum flexibility and performance.

## 📁 Project Structure

```text
.
├── backend/            # Go backend service
│   ├── cmd/api/        # Application entry point
│   ├── internal/       # Business logic, models, and repositories
│   └── go.mod          # Go module definitions
├── frontend/           # Vue 3 frontend application
│   ├── src/            # Application source code
│   ├── public/         # Static assets
│   └── package.json    # Frontend dependencies and scripts
└── GEMINI.md           # Internal development guidelines
```

## 🏁 Getting Started

### Prerequisites

- [Go 1.25.6](https://golang.org/dl/) or later
- [Node.js](https://nodejs.org/) and [pnpm](https://pnpm.io/installation)

### Installation

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/yourusername/swiftly.git
    cd swiftly
    ```

2.  **Backend Setup:**
    ```bash
    cd backend
    # No dependencies to install (standard library only)
    go run cmd/api/main.go
    ```
    The API will be available at `http://localhost:8080/api/health`.

3.  **Frontend Setup:**
    ```bash
    cd ../frontend
    pnpm install
    pnpm run dev
    ```
    The application will typically be available at `http://localhost:5173`.

## 📖 Development Conventions

- **Backend:** Follows standard Go idioms. Logic is organized within the `internal/` directory to enforce encapsulation.
- **Frontend:** Uses Vue 3 `<script setup>` syntax for Single File Components (SFCs).

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
