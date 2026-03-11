# Contributing to TFT Oracle

Thanks for your interest in contributing! This guide will help you get started.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/<your-user>/tft-oracle.git`
3. Create a feature branch: `git checkout -b feature/your-feature`
4. Make your changes
5. Push and open a PR against `main`

## Prerequisites

- Node.js v20+
- Go v1.22+
- Rust (latest stable)
- PostgreSQL v16+
- Buf CLI (for Protobuf)

## Development Workflow

### Branch Naming

| Type | Pattern | Example |
|------|---------|---------|
| Feature | `feature/short-description` | `feature/champion-grid` |
| Bug fix | `fix/short-description` | `fix/api-timeout` |
| Infra | `infra/short-description` | `infra/ci-pipeline` |
| Docs | `docs/short-description` | `docs/api-reference` |

### Commit Messages

Use conventional-style commits:

```
feat: add champion grid component
fix: resolve API timeout on match history
docs: update README with setup instructions
chore: update Go dependencies
refactor: extract shared types to proto
```

### Pull Requests

- Fill out the PR template completely
- Link related issues with `Closes #N`
- Keep PRs focused — one concern per PR
- Ensure all checks pass before requesting review

## Code Standards

- **Go**: Follow `gofmt` and `golangci-lint`
- **TypeScript/React**: ESLint + Prettier
- **Rust**: `cargo fmt` + `cargo clippy`
- **SQL**: Use sqlc — no raw query strings
- **Protobuf**: Lint with `buf lint`

## Reporting Issues

- Use the issue templates (bug report or feature request)
- For questions, use [GitHub Discussions](https://github.com/MeninoNias/tft-oracle/discussions)

## License

By contributing, you agree that your contributions will be licensed under the same license as the project.
