# Backend — Go + Connect RPC

API gateway and orchestrator for all services. Handles Riot API calls, AI prompt orchestration, analytics computation, and CommunityDragon sync.

## Structure

```
backend/
├── cmd/server/     # Entry point (main.go)
├── internal/       # Business logic (unexported packages)
├── gen/            # Generated Protobuf/Connect code (gitignored)
└── sqlc/
    └── queries/    # Raw SQL files for sqlc generation
```

## Commands

```bash
task dev:backend    # Run with hot reload
task test:go        # Run tests
task lint:go        # golangci-lint
task fmt:go         # gofmt
task build:backend  # Build binary → backend/bin/server
```

See issue #6 for setup details.
