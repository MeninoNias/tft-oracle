# Proto — Protobuf Contracts

Single source of truth for all API contracts between frontend and backend.

## Structure

```
proto/
└── tft/
    └── v1/
        ├── patch.proto        # GetPatchData — champions, items, traits
        └── player.proto       # GetPlayerProfile — match history, stats
```

## Usage

```bash
task generate:proto   # Generates Go stubs → backend/gen/
                      #          TS hooks  → frontend/src/gen/
```

See issue #1 for contract definitions.
