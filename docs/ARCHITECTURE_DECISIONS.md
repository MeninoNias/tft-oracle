# Decisões de Arquitetura (ADRs) — TFT Oracle

Registro das decisões técnicas fundamentais do projeto e suas justificativas.

---

## ADR-1: Connect RPC em vez de gRPC puro

**Status:** Aceito
**Data:** 2026-03-13

### Contexto

O TFT Oracle usa Tauri v2 como shell desktop. O Tauri renderiza a interface React dentro de uma **webview** (navegador). Precisamos que o React se comunique com o backend Go.

### Problema

O gRPC puro usa **HTTP/2 com protocolo binário** — navegadores não conseguem falar este protocolo nativamente. A solução tradicional é adicionar um proxy `grpc-web`, o que adiciona:
- Mais um processo para gerenciar
- Mais latência (hop extra)
- Mais complexidade de deploy

### Decisão

Usar **Connect RPC** (feito pelo time do Buf) em vez de gRPC puro.

### Como funciona

```
┌─────────────────────────────────┐
│          Tauri v2 (Rust)        │  ← Só janela, sem lógica
│  ┌───────────────────────────┐  │
│  │    React (WebView)        │  │     HTTP/1.1 + JSON
│  │    useQuery(getPatchData) │──────────────────────►  Go Backend :8080
│  └───────────────────────────┘  │     Connect RPC
└─────────────────────────────────┘
```

- React faz uma **requisição HTTP normal** com JSON no body
- Go backend recebe via handler Connect RPC
- Sem proxy, sem HTTP/2 obrigatório

### Comparação

| | gRPC puro | Connect RPC |
|---|---|---|
| Protocolo | HTTP/2, binário | HTTP/1.1 ou HTTP/2, JSON ou binário |
| Browser/WebView | **Não** — precisa de proxy grpc-web | **Sim** — HTTP nativo |
| Usa `.proto`? | Sim | Sim (mesmos contratos) |
| Compatível com gRPC? | — | Sim, fala ambos protocolos |
| Geração de código | protoc + plugins | `buf generate` (tudo integrado) |

### Consequências

- React chama Go diretamente — zero proxies
- Mesmos `.proto` files geram Go + TypeScript
- Auto-generated React hooks via `@connectrpc/connect-query`
- Se no futuro precisarmos de um cliente gRPC nativo (ex: mobile), Connect RPC é compatível

---

## ADR-2: sqlc em vez de ORM (GORM)

**Status:** Aceito
**Data:** 2026-03-13

### Contexto

Precisamos de acesso ao PostgreSQL no backend Go. As duas opções principais são:
- **GORM** — o ORM mais popular de Go (equivalente ao SQLAlchemy/Hibernate)
- **sqlc** — compilador que transforma SQL puro em funções Go tipadas

### Problema

O TFT Oracle tem um target de **~50MB RAM** e precisa rodar sem impacto nos FPS do jogador. As queries são simples (buscar campeões, salvar partidas), não precisam de um ORM complexo.

### Decisão

Usar **sqlc + pgx** em vez de GORM.

### Comparação

| | GORM (ORM) | sqlc (SQL → Go) |
|---|---|---|
| Você escreve | Código Go que gera SQL | SQL puro que gera código Go |
| Performance | Queries ocultas, reflection, N+1 | SQL exato, zero overhead |
| Type safety | Erros em runtime | Erros em compile-time |
| RAM | Reflection, caching de models | Zero — só funções geradas |
| Debugging | "Que SQL o GORM gerou?" | Você escreveu o SQL |

### Como funciona

```sql
-- backend/sqlc/queries/champions.sql
-- name: GetChampionsBySet :many
SELECT * FROM champions WHERE set_number = $1 ORDER BY cost;
```

```bash
sqlc generate
```

```go
// Gerado automaticamente — sem magia, sem reflection:
func (q *Queries) GetChampionsBySet(ctx context.Context, setNumber int32) ([]Champion, error)
```

### Consequências

- SQL é explícito — sem surpresas de N+1 ou queries ocultas
- Erros de tipo são pegos em compile-time
- Zero overhead de runtime — ideal para o target de ~50MB RAM
- Precisa escrever SQL manualmente (trade-off aceito — queries são simples)

---

## ADR-3: Tauri como shell mínimo

**Status:** Aceito
**Data:** 2026-03-13

### Contexto

O Tauri v2 usa Rust para criar a janela desktop e embedda uma webview para a interface React.

### Problema

Existe a tentação de colocar lógica no Rust (Tauri commands, IPC complexo) para intermediar a comunicação React ↔ Go. Isso duplicaria lógica e adicionaria complexidade.

### Decisão

Tauri é **apenas o shell da janela**. Toda lógica de negócio fica no Go backend. React chama Go diretamente via HTTP.

### Fluxo

```
❌ Errado:  React ──IPC──► Tauri (Rust) ──HTTP──► Go Backend
✅ Correto: React ──HTTP──► Go Backend    |    Tauri = só janela
```

### Responsabilidades do Tauri (Rust)

- Abrir e gerenciar a janela desktop
- Configurar a webview (CSP, permissões)
- Lifecycle do processo (iniciar/parar backend Go)
- **Nada mais** — sem lógica de negócio, sem proxy de API

### Consequências

- Rust layer fica mínimo e fácil de manter
- Sem duplicação de lógica entre Rust e Go
- React pode ser desenvolvido/testado independentemente no browser
- Build mais rápido (menos código Rust para compilar)

---

## ADR-4: Protobuf como contrato único

**Status:** Aceito
**Data:** 2026-03-13

### Decisão

Os arquivos `.proto` são a **única fonte de verdade** para a API. Não existem DTOs paralelos, types manuais, ou schemas duplicados.

### Fluxo de geração

```
proto/tft/v1/patch.proto
        │
        ▼  buf generate
        │
        ├──► backend/gen/    → Go structs + Connect handlers
        └──► frontend/src/gen/ → TypeScript types + React hooks
```

### Regras

- Nunca editar código em `backend/gen/` ou `frontend/src/gen/` — é gerado
- Sempre rodar `buf lint` antes de commitar mudanças em `.proto`
- Sempre rodar `buf generate` depois de alterar protos
- O campo `apiName` é a chave universal de junção entre CommunityDragon e Riot API

---

## Referências

- [Connect RPC Documentation](https://connectrpc.com/)
- [Buf Documentation](https://buf.build/docs/)
- [sqlc Documentation](https://docs.sqlc.dev/)
- [GORM Documentation](https://gorm.io/) (referência — não usado no projeto)
- [Tauri v2 Architecture](https://v2.tauri.app/concept/)
