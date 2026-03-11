# Documento de Especificacao Tecnica e Produto: TFT Oracle (MVP)

## 1. Visao Geral do Produto

O TFT Oracle e um aplicativo desktop ultraleve projetado para ser o "coach de bolso" dos jogadores de Teamfight Tactics. O seu grande diferencial em relacao aos concorrentes (MetaTFT, TFTactics, Mobalytics) e o baixissimo consumo de memoria (rodando de forma invisivel durante as partidas) e a utilizacao de Inteligencia Artificial para traduzir dados matematicos complexos em dicas taticas acionaveis em tempo real.

## 2. Arquitetura Tecnica (A Stack Definitiva)

A arquitetura foi desenhada para suportar alta concorrencia, processamento pesado de JSON (Riot API) e impacto zero nos FPS do jogador, utilizando as melhores praticas da industria.

### 2.1. Frontend (Cliente Desktop)

| Componente | Tecnologia | Justificativa |
|---|---|---|
| Core Desktop | Tauri v2 (Rust) | Binario nativo minusculo, ~50MB RAM (vs 300MB+ Electron/Overwolf) |
| Framework UI | React + Vite | Ecossistema vasto, Drag-and-Drop via `@dnd-kit` |
| Estilizacao | TailwindCSS + shadcn/ui | Interface gamer (Dark Mode), responsiva |
| Estado e Cache | TanStack Query + Zustand | Gerenciamento eficiente de estado e cache |

### 2.2. Backend (API Gateway & Orquestrador)

| Componente | Tecnologia | Justificativa |
|---|---|---|
| Linguagem Core | Go (Golang) | Concorrencia nativa (Goroutines), processamento pesado de JSON |
| Protocolo | gRPC via Connect RPC (Buf) | Contratos tipados, hooks React auto-gerados |
| Banco de Dados | sqlc + pgx | SQL puro compilado para Go (sem ORMs pesados) |

### 2.3. Banco de Dados e Caching

| Componente | Tecnologia | Justificativa |
|---|---|---|
| Relacional | PostgreSQL | Modelagem estruturada (Utilizador -> Partidas -> Composicoes -> Itens) |
| Cache | Redis | Poupar requisicoes a Riot API (obrigatorio na Fase 2) |

### 2.4. Motor de Inteligencia Artificial (O Simulador)

| Componente | Tecnologia | Justificativa |
|---|---|---|
| LLM Engine | OpenAI GPT-4o-mini | Custo-eficiente, Structured Outputs nativo |
| Paradigma | Structured Outputs (JSON Schema) | Respostas tipadas e validaveis |

## 3. Escopo de Funcionalidades e Ordem Cronologica

### FASE 1: A Fundacao e Dicionario de Dados

**Objetivo**: Estabelecer a comunicacao base entre o App, o Servidor e a Riot API.

- Modelagem Protobuf (.proto): Contratos gRPC para `GetPatchData` e `GetPlayerProfile`
- Sincronizacao do CommunityDragon: Worker Go para download diario dos recursos do patch atual
- Frontend Base: Setup do Tauri com React, Tailwind e Connect-Query
- Exibicao de lista de campeoes e itens na interface

**Issues**: #1, #2, #3, #4, #5, #6

### FASE 2: Historico e Analytics Pessoal

**Objetivo**: Trazer os dados do jogador e gerar valor imediato.

- Integracao Riot API: Consulta de historico via Riot ID
- Armazenamento Relacional: Composicoes, colocacoes (Top 1-8), delta de PDL
- Dashboard UI: Historico detalhado do jogador
- Modulo de Analytics: Taxas de vitoria por trait/composicao

**Issues**: #7, #8, #9, #10, #11

### FASE 3: O Simulador IA (O Diferencial de Mercado)

**Objetivo**: Criar o "Sandbox" onde a magica acontece.

- Interface de Drag-and-Drop: Tabuleiro virtual com insercao de campeoes e itens
- Engenharia de Prompt (Go): Traducao da configuracao para prompt claro
- Integracao OpenAI: Requisicao com Response Format JSON Schema
- Exibicao Tatica: Probabilidade de Vitoria e Dica de Posicionamento

**Issues**: #12, #13, #14, #15

### FASE 4: Polimento, Tier Lists e Lancamento

**Objetivo**: Finalizar as pontas soltas para a versao 1.0.

- Tier List Estatica: Melhores composicoes do patch
- Design e UX: Animacoes, skeletons, tratamento de erros
- Build do Tauri: Binario .exe para Windows

**Issues**: #16, #17, #18

## 4. Fora do Escopo (MVP)

- Leitura de memoria/tela (Visao Computacional) — evitar bans do Vanguard
- Detecao automatica do lobby via LCU API — usuario digitara o Riot ID
- Calculo deterministico frame-a-frame de combate — IA usara heuristica e probabilidade
