# Frontend — React + Vite + TailwindCSS

Desktop UI for TFT Oracle. Dark mode gamer aesthetic with shadcn/ui components.

## Structure

```
frontend/
└── src/
    ├── components/   # UI components (shadcn/ui based)
    ├── hooks/        # Custom React hooks
    ├── stores/       # Zustand client state
    └── gen/          # Generated Connect-Query hooks (gitignored)
```

## Commands

```bash
task dev:frontend     # Vite dev server
task test:frontend    # Run tests
task lint:frontend    # ESLint
task fmt:frontend     # Prettier
task build:frontend   # Production build → frontend/dist/
```

See issue #4 for setup details.
