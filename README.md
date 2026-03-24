# TFT Oracle

Your pocket coach for Teamfight Tactics. An ultra-lightweight desktop app that uses AI to translate complex game data into actionable tactical tips in real-time.

## What Makes It Different

Unlike MetaTFT, TFTactics, or Mobalytics, TFT Oracle is designed with **zero impact on your FPS** — running invisibly during matches with ~50MB of RAM usage. It leverages AI to provide real-time tactical advice, not just raw stats.

## Tech Stack

| Layer | Technology | Why |
|---|---|---|
| **Desktop** | Tauri v2 (Rust) | Native binary, ~50MB RAM vs 300MB+ Electron/Overwolf |
| **UI** | React 19 + Vite | Rich ecosystem, Drag-and-Drop via `@dnd-kit` |
| **Styling** | TailwindCSS v4 | Dark/light theme, lofi monospace aesthetic |
| **State** | TanStack Query + Zustand | Efficient caching and client state management |
| **Backend** | Go + Connect RPC (Buf) | Protobuf contracts, auto-generated React hooks |
| **Database** | PostgreSQL v16+ + sqlc/pgx | Compiled SQL, no ORM overhead |
| **Cache** | Redis | Riot API rate limit management |
| **AI Engine** | OpenAI GPT-4o-mini | Structured Outputs for typed tactical analysis |
| **Data** | CommunityDragon + Riot API | Champions, items, traits, augments, match history |

## Architecture

```
┌─────────────────────────┐     HTTP+JSON (Connect RPC)     ┌─────────────────────────┐
│      Tauri v2 App       │◄──────────────────────────────►│      Go Backend :8080    │
│  React + Vite + Zustand │                                 │                         │
│  TailwindCSS + @dnd-kit │                                 │  PatchService           │
│                         │                                 │  PlayerService          │
│  Pages:                 │                                 │  AuthService            │
│  - Champions            │                                 │  SimulationService      │
│  - Items                │                                 ├─────────────────────────┤
│  - Augments             │                                 │  PostgreSQL             │
│  - Simulator (DnD)      │                                 │  Redis                  │
│  - Profile              │                                 │  OpenAI GPT-4o-mini     │
│  - Player Search        │                                 │  Riot API               │
│  - Settings             │                                 │  CommunityDragon        │
└─────────────────────────┘                                 └─────────────────────────┘
```

## Features

### Phase 1 — Foundation & Data Dictionary (complete)
- [x] Protobuf contracts (PatchService, PlayerService, AuthService)
- [x] CommunityDragon sync (champions, items, traits, augments)
- [x] Frontend: champion, item, and augment pages with search + filters
- [x] Boot splash screen with terminal-style animation
- [x] Error pages (404, 500, 401) with lofi terminal aesthetic
- [x] Light/dark theme support

### Phase 2 — Match History & Personal Analytics (complete)
- [x] Riot API client with retry + exponential backoff on 429
- [x] Player profile: Riot ID lookup, ranked stats, match history
- [x] Three-tier caching: Redis -> DB -> Riot API
- [x] JWT authentication with access key system
- [x] Onboarding flow with 3-step wizard (identify -> connect -> augment)
- [x] Friendly error messages (Connect RPC error parsing)
- [x] Settings page with preferences management

### Phase 3 — AI Simulator (complete)
- [x] Drag-and-drop board builder (4x7 grid, player + opponent boards)
- [x] Champion palette with search + cost filter (draggable)
- [x] Item selector + star level toggle on placed champions
- [x] Real-time synergy panel (trait activations with tier coloring)
- [x] AI prompt engineering (board enrichment + trait computation)
- [x] OpenAI GPT-4o-mini integration with Structured Outputs
- [x] SimulationService (Connect RPC endpoint)
- [x] Win probability gauge (color-coded) + tactical analysis display
- [x] Positioning tips, key factors, suggested changes

### Phase 4 — Crawler, Super Tier List & Launch (next)
- [ ] Go Crawler scraping MetaTFT, TFTactics, and Mobalytics daily
- [ ] Consolidation engine: normalize and cross-rank data from 3 sources
- [ ] Super Tier List UI with source attribution and confidence scores
- [ ] Tauri build (.exe for Windows)

## Getting Started

### Prerequisites

- [Node.js](https://nodejs.org/) (v20+)
- [Go](https://go.dev/) (v1.26+)
- [Rust](https://www.rust-lang.org/) (latest stable, for Tauri)
- [PostgreSQL](https://www.postgresql.org/) (v16+)
- [Buf CLI](https://buf.build/) (for Protobuf/Connect)
- [Redis](https://redis.io/) (optional, for caching)

### Setup

```bash
# Clone
git clone https://github.com/MeninoNias/tft-oracle.git
cd tft-oracle

# Environment
cp .env.example .env
# Edit .env with your DATABASE_URL, RIOT_API_KEY, OPENAI_API_KEY, JWT_SECRET

# Database
createdb tft_oracle

# Backend
cd backend
go run ./cmd/server       # Runs migrations + CommunityDragon sync + starts server on :8080

# Frontend (new terminal)
cd frontend
npm install
npm run dev               # Vite dev server on :5173

# Code generation (after proto changes)
buf generate              # Go stubs + TypeScript hooks
sqlc generate             # Go from SQL queries
```

### Running Tests

```bash
cd backend
go test ./...             # All backend tests
go test ./... -cover      # With coverage
```

## Out of Scope (MVP)

- Memory/screen reading (avoids Vanguard anti-cheat bans)
- Automatic lobby detection via LCU API (manual Riot ID input)
- Deterministic frame-by-frame combat calculation (AI uses heuristics)

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes following conventional commits (`feat:`, `fix:`, `docs:`, etc.)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request using the provided template

## License

This project is under development. License TBD.
