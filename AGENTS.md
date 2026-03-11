# AGENTS.md

Instructions for AI agents (Claude Code, Copilot, Cursor, etc.) working on this repository.

## General Rules

- This is a **solo developer project**. Do not suggest team workflows, code review processes, or multi-contributor patterns.
- Read existing code before modifying. Understand context before making changes.
- Keep changes minimal and focused. Do not refactor, add comments, or "improve" code beyond what was asked.
- Do not create files unless strictly necessary. Prefer editing existing files.
- Never commit `.env`, API keys, credentials, or secrets.
- Follow the conventions defined in CLAUDE.md and `.github/CONTRIBUTING.md`.

## Language-Specific Instructions

### Go (backend/)

- Use `gofmt` formatting. Run `golangci-lint` before committing.
- All DB access goes through sqlc — never write raw SQL in Go.
- Use `internal/` packages for business logic. Keep `cmd/` thin.
- Wrap errors with context: `fmt.Errorf("fetchChampions: %w", err)`.
- Use Connect RPC handlers, not raw HTTP handlers.
- Protobuf types are the API contract — do not create parallel DTOs.

### TypeScript/React (frontend/)

- Functional components only. No class components.
- Use generated Connect-Query hooks for API calls — do not use `fetch` or `axios` directly.
- Server state via TanStack Query. Client state via Zustand. Do not mix them.
- Style with TailwindCSS utility classes. Use shadcn/ui components where applicable.
- No `any` types. Use the generated protobuf types from `frontend/src/gen/`.

### Protobuf (proto/)

- Proto files are the single source of truth for the API.
- Run `buf lint` before committing any `.proto` changes.
- Run `buf generate` after modifying protos to regenerate Go + TypeScript code.
- Never hand-edit files in `backend/gen/` or `frontend/src/gen/`.

### SQL (migrations/, backend/sqlc/)

- Write migrations in `migrations/` as sequential numbered files.
- Write queries in `backend/sqlc/queries/`. Run `sqlc generate` to produce Go code.
- Never write inline SQL strings in Go — always use sqlc-generated functions.

### Rust (src-tauri/)

- Minimal Tauri shell — only IPC commands and window management.
- Format with `cargo fmt`, lint with `cargo clippy`.
- Do not add heavy Rust logic. The backend is Go.

## Commit Conventions

- Use conventional commits: `feat:`, `fix:`, `docs:`, `chore:`, `refactor:`, `test:`, `ci:`
- Branch names: `feature/`, `fix/`, `infra/`, `docs/` prefixes
- Always reference issues: `Closes #N` or `Relates to #N`
- Keep commits atomic — one concern per commit

## What NOT to Do

- Do not add ORMs (GORM, Prisma, etc.) — this project uses sqlc.
- Do not add Electron or Overwolf — this project uses Tauri.
- Do not add REST endpoints — this project uses Connect RPC (gRPC).
- Do not add Redux, MobX, or Recoil — this project uses Zustand + TanStack Query.
- Do not read game memory or screen pixels — Vanguard anti-cheat risk.
- Do not auto-detect lobbies via LCU API — users input Riot ID manually.
- Do not build a deterministic combat simulator — the AI uses heuristics.
- Do not over-engineer. No feature flags, no DI containers, no abstraction layers for single-use code.

## Project Context

- Full spec: `docs/SPEC.md` (Portuguese)
- Architecture and commands: `CLAUDE.md`
- 4 development phases tracked via GitHub milestones (#1-#24)
- Target: ~50MB RAM, zero FPS impact during TFT gameplay
