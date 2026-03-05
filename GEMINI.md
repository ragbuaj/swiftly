# Swiftly Project Overview

Swiftly is a full-stack e-commerce web application featuring a Go-based backend, a Vue 3 frontend, and a containerized infrastructure. This project is structured as a monorepo with separate directories for backend and frontend services.

## Full Project Architecture

### 1. Backend (/backend)

- **Language:** Go 1.25.6
- **Framework:** Standard library (net/http) with custom routing.
- **Architecture (Modular Clean Architecture):**
  - **`internal/app`**: Dependency Container (App struct) for managing global resources (DB, S3, Redis).
  - **`internal/api/routes`**: Centralized route registration for all modules.
  - **`internal/config`**: Structured and validated configuration management (Fail-Fast).
  - **`internal/pkg/apperror`**: Domain-Driven Error Handling (Mapping domain errors to HTTP).
  - **`internal/pkg/validator`**: Centralized struct validation using `go-playground/validator`.
  - **`internal/database`**: Atomic transaction support via `WithTransaction` helper.
- **Standard Library Hooks:** Uses `log/slog` for structured logging.

### 2. Frontend (/frontend)

- **Framework:** Vue 3 (Composition API) with TypeScript.
- **UI Library:** Shadcn-Vue (based on Radix Vue/Reka UI) and Tailwind CSS v4.
- **Build Tool:** Vite.
- **Package Manager:** pnpm.
- **State Management:** Pinia.

### 3. Infrastructure & Services (Dockerized)

- **Database:** PostgreSQL for persistent relational data.
- **Cache & Session:** Redis for Token Blacklisting, OTP, and session management.
- **Object Storage:** MinIO (S3-compatible) for handling file uploads.

## Development Workflow

- **Collaboration Protocol**:
  1. Sebelum menulis kode, jelaskan pendekatan yang akan diambil dan tunggu persetujuan.
  2. Jika persyaratan ambigu, ajukan pertanyaan klarifikasi sebelum menulis kode.
  3. Setelah selesai menulis kode, daftar kemungkinan edge cases dan usulkan test case untuknya.
  4. Jika tugas memerlukan perubahan pada lebih dari 3 file, berhenti dan bagi menjadi tugas-tugas yang lebih kecil terlebih dahulu.
  5. Jika bug dilaporkan, mulai dengan menulis reproduction test case, lalu perbaiki hingga tes tersebut lulus.
- **Testing Mandate**: Selalu buat, perbarui, atau sempurnakan test case setiap kali ada perubahan kode. Tugas belum selesai tanpa cakupan tes yang komprehensif.
- **Syntax Integrity**: **Selalu lakukan pengecekan syntax** (misal `go build` atau `tsc`) setiap kali ada penambahan atau perubahan kode.
- **Test-Driven Correction**: Jika test case gagal, lakukan analisis mendalam mengapa gagal, lalu usulkan rencana perbaikan kepada user sebelum melanjutkan.
- **PowerShell Compatibility**: Gunakan titik koma `;` untuk pemisah perintah (jangan gunakan `&&`).
- **Environment Templates**: Perbarui `.env.sample` (backend) atau `.env.sample` (frontend) saat menambahkan variabel lingkungan baru.
