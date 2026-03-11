# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

TFT Oracle is an ultra-lightweight desktop app (~50MB RAM) that acts as an AI-powered coach for Teamfight Tactics. It uses AI to translate complex game data into actionable tactical tips in real-time.

**Solo developer project** — @MeninoNias is the sole contributor. No team reviews, no approval gates.

## Architecture

```
┌──────────────────┐     gRPC/Connect     ┌──────────────────┐
│   Tauri v2 App   │◄───────────────────►│   Go Backend     │
│   React + Vite   │                      │   (Goroutines)   │
│   TailwindCSS    │                      │                  │
└──────────────────┘                      ├──────────────────┤
                                          │  PostgreSQL      │
                                          │  Redis (Phase 2) │
                                          │  OpenAI API      │
                                          │  Riot API        │
                                          └──────────────────┘
```

### Stack

| Layer | Tech | Notes |
|-------|------|-------|
| Desktop shell | Tauri v2 (Rust) | ~50MB RAM target |
| UI | React + Vite + TailwindCSS + shadcn/ui | Dark mode gamer aesthetic |
| State | TanStack Query + Zustand | Cache + client state |
| DnD | @dnd-kit | Phase 3 board builder |
| Backend | Go + Connect RPC (Buf) | Protobuf contracts → auto-generated React hooks |
| DB access | sqlc + pgx | Compiled SQL, no ORM |
| Database | PostgreSQL v16+ | Relational storage |
| Cache | Redis | Riot API rate limits (Phase 2+) |
| AI | OpenAI GPT-4o-mini | Structured Outputs (JSON Schema) |
| Data sync | CommunityDragon | Game data (champions, items, traits) |
| Crawler | Go (Goroutines) | MetaTFT, TFTactics, Mobalytics (Phase 4) |

## Project Structure (planned)

```
tft-oracle/
├── proto/                  # Protobuf contracts (.proto files)
├── backend/                # Go backend (Connect RPC server)
│   ├── cmd/server/         # Entry point
│   ├── internal/           # Business logic
│   ├── gen/                # Generated protobuf code
│   └── sqlc/               # Generated SQL queries
├── frontend/               # React + Vite app
│   ├── src/
│   │   ├── components/     # UI components
│   │   ├── hooks/          # Custom hooks
│   │   ├── stores/         # Zustand stores
│   │   └── gen/            # Generated Connect-Query hooks
│   └── index.html
├── src-tauri/              # Tauri v2 Rust shell
├── migrations/             # PostgreSQL migrations
├── docs/                   # Specification and docs
└── .github/                # Templates, CI, dependabot
```

## Development Phases

Tracked via GitHub milestones and issues (#1-#24):

- **Phase 1** (#1-#6, #19-#21): Foundation — monorepo setup, Protobuf contracts, CommunityDragon sync, Go backend, Tauri+React frontend, PostgreSQL schema, champion/item display
- **Phase 2** (#7-#11): Analytics — Riot API integration, match history, player dashboard, analytics module, Redis caching
- **Phase 3** (#12-#15): AI Simulator — DnD board builder, prompt engineering, OpenAI integration, simulation results UI
- **Phase 4** (#16-#18, #22-#24): Crawler & Launch — Go crawler (MetaTFT/TFTactics/Mobalytics), consolidation engine, Super Tier List, UX polish, Tauri build (.exe)

## Commands

Commands will be added as the codebase is built. Planned:

```bash
# Backend (Go)
cd backend && go run ./cmd/server          # Run backend
cd backend && go test ./...                 # Run all tests
buf generate                                # Generate protobuf code
buf lint                                    # Lint proto files

# Frontend (React)
cd frontend && npm run dev                  # Dev server
cd frontend && npm run build                # Production build
cd frontend && npm run lint                 # ESLint
cd frontend && npm run format               # Prettier

# Tauri
cd src-tauri && cargo tauri dev             # Dev mode (frontend + backend + shell)
cd src-tauri && cargo tauri build           # Production .exe

# Database
sqlc generate                               # Generate Go from SQL queries
```

## Code Conventions

### Git

- **Branch naming**: `feature/`, `fix/`, `infra/`, `docs/` prefixes
- **Commits**: Conventional style — `feat:`, `fix:`, `docs:`, `chore:`, `refactor:`, `test:`, `ci:`
- **PRs**: Always link issues with `Closes #N`
- **Main branch**: Protected — changes go through PRs

### Go (backend)

- Format with `gofmt`, lint with `golangci-lint`
- Use `internal/` for unexported packages
- No ORMs — only sqlc-generated code for DB access
- Error handling: wrap with `fmt.Errorf("context: %w", err)`

### TypeScript/React (frontend)

- ESLint + Prettier enforced
- Functional components only, hooks for logic
- TanStack Query for server state, Zustand for client state
- No `any` types — use generated protobuf types

### Protobuf

- Lint with `buf lint`
- Single source of truth for API contracts
- Generate both Go server stubs and TypeScript client hooks

### SQL

- Write raw SQL in `backend/sqlc/queries/`
- Run `sqlc generate` to produce type-safe Go code
- Never write raw SQL strings in Go code

## Key Constraints

- **Performance**: Must stay under ~50MB RAM with zero FPS impact during gameplay
- **Polyglot**: Rust (desktop shell), Go (backend), TypeScript/React (UI), Protobuf (contracts), SQL (database)
- **Out of MVP scope**: No memory/screen reading (Vanguard anti-cheat risk), no auto lobby detection (manual Riot ID), no deterministic combat sim (AI heuristics instead)
- **Security**: Never commit API keys (.env files). Riot API key, OpenAI key, and DB credentials go in environment variables only.

## Key Files

- `docs/SPEC.md` — Full technical specification (Portuguese)
- `.github/ISSUE_TEMPLATE/` — Bug report and feature request templates
- `.github/pull_request_template.md` — PR template with checklist
- `.github/CONTRIBUTING.md` — Contribution guide with conventions
- `.github/SECURITY.md` — Vulnerability reporting policy
