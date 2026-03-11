# TFT Oracle

Your pocket coach for Teamfight Tactics. A ultra-lightweight desktop app that uses AI to translate complex game data into actionable tactical tips in real-time.

## What Makes It Different

Unlike MetaTFT, TFTactics, or Mobalytics, TFT Oracle is designed with **zero impact on your FPS** — running invisibly during matches with ~50MB of RAM usage. It leverages AI to provide real-time tactical advice, not just raw stats.

## Tech Stack

| Layer | Technology | Why |
|---|---|---|
| **Desktop** | Tauri v2 (Rust) | Native binary, ~50MB RAM vs 300MB+ Electron/Overwolf |
| **UI** | React + Vite | Rich ecosystem, Drag-and-Drop via `@dnd-kit` |
| **Styling** | TailwindCSS + shadcn/ui | Dark mode gamer aesthetic, responsive |
| **State** | TanStack Query + Zustand | Efficient caching and state management |
| **Backend** | Go (Golang) | Native concurrency (Goroutines) for Riot API processing + web crawling |
| **API Protocol** | gRPC via Connect RPC (Buf) | End-to-end typed contracts, auto-generated React hooks |
| **Database** | PostgreSQL + sqlc/pgx | Structured relational data, compiled SQL (no ORM overhead) |
| **Cache** | Redis | Save Riot API rate limits (Phase 2+) |
| **AI Engine** | OpenAI GPT-4o-mini | Structured Outputs for typed tactical analysis |

## Architecture Overview

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

## Development Roadmap

### Phase 1 — Foundation & Data Dictionary
- Protobuf contracts (`GetPatchData`, `GetPlayerProfile`)
- CommunityDragon sync worker (daily patch data)
- Frontend base setup (Tauri + React + Tailwind + Connect-Query)
- Champion and item listing from Go server

### Phase 2 — Match History & Personal Analytics
- Riot API integration (match history via Riot ID)
- Relational storage (compositions, placements, LP delta)
- Player dashboard UI
- Analytics module (win rates per trait/comp)

### Phase 3 — AI Simulator (Market Differentiator)
- Drag-and-Drop board builder (carry/tank + items)
- Prompt engineering in Go backend
- OpenAI integration with JSON Schema response
- Win probability and positioning tips display

### Phase 4 — Crawler, Consolidated Super Tier List & Launch
- Go Crawler scraping MetaTFT, TFTactics, and Mobalytics daily
- Consolidation engine: normalize and cross-rank data from all 3 sources
- Super Tier List UI with source attribution and confidence scores
- UX refinement (animations, skeletons, error handling)
- Tauri build (.exe for Windows)

## Out of Scope (MVP)

- Memory/screen reading (avoids Vanguard bans)
- Automatic lobby detection via LCU API (manual Riot ID input)
- Deterministic frame-by-frame combat calculation (AI uses heuristics)

## Getting Started

> Coming soon — the project is in early development.

### Prerequisites

- [Node.js](https://nodejs.org/) (v20+)
- [Go](https://go.dev/) (v1.22+)
- [Rust](https://www.rust-lang.org/) (latest stable)
- [PostgreSQL](https://www.postgresql.org/) (v16+)
- [Buf CLI](https://buf.build/) (for Protobuf/Connect)

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes following the project conventions
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request using the provided template

## License

This project is under development. License TBD.
