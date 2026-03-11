# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

TFT Oracle is an ultra-lightweight desktop app (~50MB RAM) that acts as an AI-powered coach for Teamfight Tactics. It uses AI to translate complex game data into actionable tactical tips in real-time. Currently in early development (documentation phase — no source code yet).

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

**Frontend**: Tauri v2 (Rust) + React + Vite + TailwindCSS + shadcn/ui. State via TanStack Query + Zustand. DnD via @dnd-kit (Phase 3).

**Backend**: Go with gRPC via Connect RPC (Buf). Protobuf contracts auto-generate typed React hooks (@connectrpc/connect-query). Database access via sqlc + pgx (compiled SQL, no ORM).

**Data**: PostgreSQL v16+ for relational storage. Redis for Riot API rate limit caching (Phase 2+). CommunityDragon for game data sync.

**AI**: OpenAI GPT-4o-mini with Structured Outputs (JSON Schema) for typed tactical analysis responses.

## Development Phases

The project follows 4 incremental phases (tracked via GitHub issues #1-#18):

- **Phase 1** (#1-#6): Foundation — Protobuf contracts, CommunityDragon sync worker, Go backend setup, Tauri+React frontend setup, PostgreSQL schema, champion/item display
- **Phase 2** (#7-#11): Analytics — Riot API integration, match history storage, player dashboard, analytics module, Redis caching
- **Phase 3** (#12-#15): AI Simulator — DnD board builder, prompt engineering, OpenAI integration, simulation results UI
- **Phase 4** (#16-#18, #22-#24): Crawler & Launch — Go crawler scraping MetaTFT/TFTactics/Mobalytics daily, consolidation engine for Super Tier List, UX polish, Tauri build (.exe)

## Prerequisites

- Node.js v20+, Go v1.22+, Rust (latest stable), PostgreSQL v16+, Buf CLI

## Key Constraints

- **Performance**: Must stay under ~50MB RAM with zero FPS impact during gameplay
- **Polyglot**: Rust (desktop shell), Go (backend), TypeScript/React (UI), Protobuf (contracts), SQL (database)
- **Out of MVP scope**: No memory/screen reading (Vanguard anti-cheat risk), no auto lobby detection (manual Riot ID), no deterministic combat sim (AI heuristics instead)

## Key Files

- `docs/SPEC.md` — Full technical specification (Portuguese)
- `.github/ISSUE_TEMPLATE/` — Bug report and feature request templates (YAML forms)
- `.github/pull_request_template.md` — PR template with type-of-change checklist
