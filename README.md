# Social Engine

> A modern, full-stack social platform — a mini Twitter/blog — built as a portfolio project to demonstrate clean architecture across a **SvelteKit** frontend and a **Go** API.

<p>
  <img alt="SvelteKit" src="https://img.shields.io/badge/SvelteKit-2-FF3E00?logo=svelte&logoColor=white">
  <img alt="Svelte 5" src="https://img.shields.io/badge/Svelte-5-FF3E00?logo=svelte&logoColor=white">
  <img alt="TypeScript" src="https://img.shields.io/badge/TypeScript-3178C6?logo=typescript&logoColor=white">
  <img alt="TailwindCSS" src="https://img.shields.io/badge/Tailwind_CSS-v4-06B6D4?logo=tailwindcss&logoColor=white">
  <img alt="Go" src="https://img.shields.io/badge/Go-1.25-00ADD8?logo=go&logoColor=white">
  <img alt="Fiber" src="https://img.shields.io/badge/Fiber-v2-00ADD8">
  <img alt="PostgreSQL" src="https://img.shields.io/badge/PostgreSQL-16-4169E1?logo=postgresql&logoColor=white">
  <img alt="License: MIT" src="https://img.shields.io/badge/License-MIT-green">
</p>

Social Engine lets people register, publish short posts with hashtags, like and comment, follow feeds, explore trending topics, and personalize their profile. It exists to showcase a **production-minded full-stack setup**: a typed, animated SvelteKit UI backed by a layered, tested Go API.

A distinctive feature is that the frontend runs in two modes: **offline** (rich mocked data, no backend needed) and **online** (talking to the real Go API + PostgreSQL), toggled by a single environment flag.

---

## ✨ Features

- **Authentication** — register & login with JWT access tokens, bcrypt password hashing, and httpOnly server-side session cookies.
- **Posts & feed** — publish posts, automatic `#hashtag` extraction, public feed, and single-post view with comments.
- **Engagement** — like/unlike (with per-user "liked by me" state) and threaded comments.
- **Profiles** — public profiles with post/follower/following stats and customizable avatars & banners.
- **Discovery** — Explore (suggested people + posts, filterable by tag) and Trending (hottest hashtags + top posts).
- **Polished UX** — micro-animations, page transitions, and an ambient animated background, all respecting `prefers-reduced-motion`.
- **Offline / Online modes** — develop the UI with zero backend, then flip one flag to run against the live API.

---

## 🧱 Technologies

**Frontend**
- SvelteKit 2 · Svelte 5 (runes) · TypeScript
- TailwindCSS v4 · Vite

**Backend**
- Go 1.25 · Fiber v2 (HTTP) · GORM (ORM)
- PostgreSQL 16 · golang-migrate (migrations)
- JWT (`golang-jwt/v5`) · bcrypt · `go-playground/validator`
- Zap (structured logging) · Swagger / swaggo (API docs)

**Tooling & Infra**
- Docker Compose · Make · GitHub Actions (tests + 80% coverage gate)

---

## 🏗️ Architecture

High-level flow — the SvelteKit server layer talks to the Go API, which owns the database. Auth tokens never reach the browser (they live in an httpOnly cookie, used server-side).

```
┌─────────────┐     httpOnly session      ┌──────────────────┐      REST /json      ┌──────────────┐
│   Browser   │◄────── cookie ───────────►│  SvelteKit server │◄──── + Bearer ──────►│    Go API    │
│  (Svelte 5) │                           │ (loads & actions) │                      │  (Fiber v2)  │
└─────────────┘                           └──────────────────┘                      └──────┬───────┘
                                                                                            │ GORM
                                                                                     ┌──────▼───────┐
                                                                                     │  PostgreSQL  │
                                                                                     └──────────────┘
```

The **Go API** is layered: `handlers` (HTTP + auth middleware) → `models` (domain logic) → `repositories` (a generic GORM repository over typed entities), with a shared catalog of API errors, request validation, and structured logging. The **web app** keeps a typed API client, server-only auth/session handling, and a mock layer that mirrors the API contract.

---

## 🚀 Getting Started

### Prerequisites

- [Node.js](https://nodejs.org/) 20.19+ or 22+
- [Go](https://go.dev/) 1.25+
- [Docker](https://www.docker.com/) & Docker Compose
- [golang-migrate](https://github.com/golang-migrate/migrate) CLI (only for **online** mode)
- `make`

### Installation

```bash
# 1. Clone
git clone <repo-url> social-engine
cd social-engine

# 2. Configure environment (each app has its own .env)
cp api/.env.example api/.env
cp web/.env.example web/.env

# 3. Install web dependencies
cd web && npm install && cd ..
```

### Environment Variables

The API and the web app each read their own `.env` file.

**`api/.env`**

| Variable               | Description                                          | Example                                                                 |
| ---------------------- | ---------------------------------------------------- | ----------------------------------------------------------------------- |
| `PORT`                 | Port the Go API listens on                           | `8080`                                                                   |
| `DATABASE_URL`         | PostgreSQL connection string                         | `postgres://social:social@localhost:5432/social_engine?sslmode=disable` |
| `JWT_SECRET`           | Secret used to sign JWT tokens — **min. 32 bytes** (the API refuses to start otherwise); generate with `openssl rand -base64 32` | *(long, random value)* |
| `CORS_ALLOWED_ORIGINS` | Comma-separated allowed origins                      | `http://localhost:5173`                                                  |

**`web/.env`**

| Variable       | Description                                          | Example                 |
| -------------- | --------------------------------------------------- | ----------------------- |
| `USE_API`      | `false` = mocked data · `true` = call the real API  | `false`                 |
| `API_BASE_URL` | Base URL the web uses to reach the API              | `http://localhost:8080` |

### Running the Project

**Offline** — only the web app, with mocked data. No database or API required. Great for exploring the UI. Keep `USE_API=false` in `web/.env`, then:

```bash
cd web && npm run dev
```

**Online** — run the database, API, and web app together (set `USE_API=true` in `web/.env` first):

```bash
# 1. Start PostgreSQL and apply migrations
cd api
docker compose up -d db
make migrate_up

# 2. Run the API (loads api/.env automatically)
go run .

# 3. In another terminal, run the web app
cd web && npm run dev
```

The web app runs at **http://localhost:5173** and the API at **http://localhost:8080** (Swagger docs at `/docs`).

---

## 🛠️ Available Commands

**API** (from `api/`)

| Command                   | Description                              |
| ------------------------- | ---------------------------------------- |
| `docker compose up -d db` | Start the PostgreSQL container           |
| `go run .`                | Run the API server                       |
| `make migrate_up`         | Apply database migrations                |
| `make tests`              | Run all Go tests with coverage           |
| `bash scripts/coverage.sh`| Run tests with the 80% coverage gate     |
| `make generate_docs`      | Regenerate the Swagger documentation     |

**Web** (from `web/`)

| Command         | Description                        |
| --------------- | ---------------------------------- |
| `npm run dev`   | Start the development server       |
| `npm run build` | Production build                   |
| `npm run check` | Type-check with `svelte-check`     |

---

## 📁 Project Structure

```
social-engine/
├── api/                    # Go API (Fiber + GORM + PostgreSQL)
│   ├── common/             # apiErrors, models, repositories, validation, logger, migrations
│   ├── handlers/           # HTTP handlers + JWT auth middleware
│   ├── docs/               # generated Swagger docs
│   ├── scripts/            # coverage gate script
│   ├── docker-compose.yml  # PostgreSQL
│   ├── .env.example        # API configuration
│   └── Makefile            # migrations, tests, docs
├── web/                    # SvelteKit app (Svelte 5 + TypeScript + Tailwind v4)
│   ├── src/
│   │   ├── lib/
│   │   │   ├── components/ # UI components (feed, nav, auth, background)
│   │   │   ├── server/     # typed API client, auth/session, mock data
│   │   │   └── ...         # appearance, formatting, shared types
│   │   └── routes/         # (auth) and (main) route groups
│   └── .env.example        # web configuration
└── .github/workflows/      # CI: tests + coverage gate
```

---

## 📄 License

Released under the [MIT License](./LICENSE).
