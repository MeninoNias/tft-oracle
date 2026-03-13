# Fontes de Dados — TFT Oracle

Mapeamento completo das fontes de dados externas usadas pelo TFT Oracle para alimentar os contratos Protobuf e a base de dados.

## Visão Geral

| Fonte | Tipo | Uso | Fase |
|-------|------|-----|------|
| CommunityDragon | Dados estáticos (JSON) | Dicionário de campeões, itens, traits, augments | Phase 1 |
| Riot API — Account V1 | REST API | Lookup de jogador por Riot ID → PUUID | Phase 2 |
| Riot API — TFT Summoner V1 | REST API | Perfil do invocador (nível, ícone) | Phase 2 |
| Riot API — TFT League V1 | REST API | Dados de ranked (tier, LP, wins/losses) | Phase 2 |
| Riot API — TFT Match V1 | REST API | Histórico de partidas (composições, colocações) | Phase 2 |
| MetaTFT / TFTactics / Mobalytics | Web Scraping | Tier lists, win rates, meta comps | Phase 4 |

> **Chave universal de junção:** o campo `apiName` conecta todas as fontes — CommunityDragon usa como ID primário, Riot API usa nos campos `characterId`, `itemNames`, `name` (trait) e `augments`.

---

## 1. CommunityDragon — Dados Estáticos do Patch

**URL:** `https://raw.communitydragon.org/latest/cdragon/tft/{locale}.json`
**Tamanho:** ~21MB (en_us) | **Locales:** 28 idiomas disponíveis
**Atualização:** a cada novo patch (~2 semanas)

### Estrutura do JSON

```
{
  "items":   [...],    // 3112 itens (augments + craftáveis + componentes)
  "setData": [...],    // 34 variantes de set (com campeões, traits, itens, augments)
  "sets":    {...}     // Índice simplificado por set (só campeões + traits)
}
```

### Champion

Localizado em `setData[].champions[]`

| Campo | Tipo | Exemplo | Descrição |
|-------|------|---------|-----------|
| `apiName` | string | `"TFT14_Ahri"` | ID único (chave de junção com Riot API) |
| `characterName` | string | `"TFT14_Ahri"` | Nome interno do personagem |
| `name` | string | `"Ahri"` | Nome de exibição (localizado) |
| `cost` | int | `1`–`5` | Custo em ouro / raridade |
| `role` | string | `"APTank"`, `"Carry"` | Função tática |
| `traits` | []string | `["TFT14_Mage"]` | Traits associadas (apiName) |
| `icon` | string | path `.tex` | Splash art |
| `squareIcon` | string | path `.tex` | Ícone quadrado (mobile) |
| `tileIcon` | string | path `.tex` | Ícone HUD |
| `stats` | object | ver abaixo | Atributos base |
| `ability` | object | ver abaixo | Habilidade do campeão |

**Stats:**

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `hp` | float | Vida base |
| `armor` | float | Armadura |
| `magicResist` | float | Resistência mágica |
| `damage` | float | Dano de ataque |
| `attackSpeed` | float | Velocidade de ataque |
| `range` | float | Alcance de ataque (hexes) |
| `mana` | int | Mana máxima |
| `initialMana` | float | Mana inicial |
| `critChance` | float | Chance de crítico (0.25 = 25%) |
| `critMultiplier` | float | Multiplicador de crítico |

**Ability:**

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `name` | string | Nome da habilidade |
| `desc` | string | Descrição com templates `@var@` |
| `icon` | string | Ícone da habilidade |
| `variables` | []object | `{name, value[]}` — valores por star level |

### Trait

Localizado em `setData[].traits[]`

| Campo | Tipo | Exemplo | Descrição |
|-------|------|---------|-----------|
| `apiName` | string | `"TFT14_Divinicorp"` | ID único |
| `name` | string | `"Divinicorp"` | Nome de exibição (localizado) |
| `desc` | string | HTML-like | Descrição com breakpoints |
| `icon` | string | path `.tex` | Ícone do trait |
| `effects` | []object | ver abaixo | Efeitos por nível de ativação |

**Effects (por nível):**

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `minUnits` | int | Mínimo de unidades para ativar |
| `maxUnits` | int | Máximo de unidades neste nível |
| `style` | int | Estilo visual (bronze, silver, gold, chromatic) |
| `variables` | object | Pares chave-valor dos efeitos numéricos |

### Item

Localizado em `items[]` (global, todos os sets)

| Campo | Tipo | Exemplo | Descrição |
|-------|------|---------|-----------|
| `apiName` | string | `"TFT_Item_RecurveBow"` | ID único |
| `name` | string | `"Recurve Bow"` | Nome de exibição (localizado) |
| `desc` | string | template `@var@` | Descrição do efeito |
| `composition` | []string | `[]` base, `["X","Y"]` craft | Receita de crafting |
| `effects` | object | `{AS: 10.0}` | Valores numéricos dos efeitos |
| `icon` | string | path `.tex` | Ícone do item |
| `associatedTraits` | []string | `["TFT14_Mage"]` | Traits vinculadas (emblemas) |
| `incompatibleTraits` | []string | `[]` | Traits conflitantes |
| `tags` | []string | `["component"]` | Categorias |
| `unique` | bool | `false` | Equipamento único |
| `from` | null\|int | `null` | Referência ao item base |
| `id` | null\|int | `null` | ID numérico |

**Categorias de itens (por apiName):**

| Tipo | Contagem | Identificação |
|------|----------|---------------|
| Augments | ~1601 | `"Augment"` no apiName |
| Craftáveis | ~302 | `composition` não-vazio |
| Outros | ~1209 | Componentes, especiais, emblemas |

### Augments

Localizado em `setData[].augments[]` como `[]string` de apiNames.
Dados completos de cada augment estão no array global `items[]`.

### Set Metadata

| Campo | Tipo | Exemplo | Descrição |
|-------|------|---------|-----------|
| `number` | int | `16` | Número do set |
| `name` | string | `"Set16"` | Nome do set |
| `mutator` | string | `"TFTSet16"` | Variante do modo de jogo |

---

## 2. Riot Games API — Dados do Jogador

Base URL regional: `https://{region}.api.riotgames.com`
Autenticação: Header `X-Riot-Token: {API_KEY}`

### Account V1 — Lookup por Riot ID

**Endpoint:** `GET /riot/account/v1/accounts/by-riot-id/{gameName}/{tagLine}`
**Região:** `americas`, `europe`, `asia`, `esports`

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `puuid` | string | UUID do jogador (chave primária universal) |
| `gameName` | string | Nome no jogo |
| `tagLine` | string | Tag (ex: `"BR1"`) |

### TFT Summoner V1

**Endpoint:** `GET /tft/summoner/v1/summoners/by-puuid/{puuid}`
**Região:** `br1`, `na1`, `euw1`, `kr`, etc.

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `id` | string | Summoner ID (encriptado) |
| `accountId` | string | Account ID (encriptado) |
| `puuid` | string | Player UUID |
| `profileIconId` | int32 | ID do ícone de perfil |
| `summonerLevel` | int64 | Nível do invocador |
| `revisionDate` | int64 | Última atualização (unix timestamp) |

### TFT League V1 — Dados de Ranked

**Endpoint:** `GET /tft/league/v1/entries/by-puuid/{puuid}`

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `leagueId` | string | ID da liga |
| `puuid` | string | Player UUID |
| `queueType` | string | `"RANKED_TFT"`, `"RANKED_TFT_TURBO"`, etc. |
| `tier` | string | `"IRON"`, `"BRONZE"`, ..., `"CHALLENGER"` |
| `rank` | string | `"I"`, `"II"`, `"III"`, `"IV"` |
| `leaguePoints` | int32 | LP (League Points) |
| `wins` | int32 | Vitórias |
| `losses` | int32 | Derrotas |
| `hotStreak` | bool | Em sequência de vitórias |
| `veteran` | bool | Jogador veterano na liga |
| `freshBlood` | bool | Recém-promovido |
| `inactive` | bool | Inativo |
| `miniSeries` | object? | Promoção: `{progress, wins, losses, target}` |

### TFT Match V1 — Histórico de Partidas

**Listar IDs:** `GET /tft/match/v1/matches/by-puuid/{puuid}/ids`
**Detalhe:** `GET /tft/match/v1/matches/{matchId}`

#### MatchDTO

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `metadata.dataVersion` | string | Versão dos dados |
| `metadata.matchId` | string | `"BR1_1234567"` |
| `metadata.participants` | []string | Lista de PUUIDs |
| `info.gameDatetime` | int64 | Timestamp unix (ms) |
| `info.gameLength` | float32 | Duração em segundos |
| `info.gameVersion` | string | Versão do patch |
| `info.gameVariation` | string | Modo de jogo |
| `info.queueId` | int32 | ID da fila |
| `info.tftSetNumber` | int32 | Número do set |
| `info.tftSetCoreName` | string | Nome base do set |
| `info.tftGameType` | string | Tipo do jogo |
| `info.endOfGameResult` | string | Resultado |
| `info.participants` | []ParticipantDTO | Jogadores da partida |

#### ParticipantDTO

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `puuid` | string | Player UUID |
| `placement` | int32 | Colocação (1-8) |
| `level` | int32 | Nível do jogador |
| `goldLeft` | int32 | Ouro restante |
| `lastRound` | int32 | Última rodada jogada |
| `timeEliminated` | float32 | Tempo de eliminação (s) |
| `totalDamageToPlayers` | int32 | Dano total a jogadores |
| `playersEliminated` | int32 | Jogadores eliminados |
| `partnerGroupId` | int32 | Parceiro no Double Up |
| `augments` | []string | Augments (apiNames) |
| `companion` | CompanionDTO | Little Legend |
| `traits` | []TraitDTO | Traits ativas |
| `units` | []UnitDTO | Unidades no board |

#### TraitDTO

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `name` | string | apiName do trait (ex: `"TFT14_Mage"`) |
| `numUnits` | int32 | Unidades contribuindo |
| `style` | int32 | 0=inativo, 1=bronze, 2=silver, 3=gold, 4=chromatic |
| `tierCurrent` | int32 | Nível ativo atual |
| `tierTotal` | int32 | Nível máximo possível |

#### UnitDTO

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `characterId` | string | apiName do campeão (ex: `"TFT14_Ahri"`) |
| `name` | string | Nome de exibição |
| `tier` | int32 | Star level (1-3) |
| `rarity` | int32 | Custo/raridade (0-4 → 1-5 gold) |
| `items` | []int32 | IDs numéricos dos itens |
| `itemNames` | []string | apiNames dos itens |

#### CompanionDTO

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `contentId` | string | ID do conteúdo |
| `species` | string | Espécie do Little Legend |
| `itemId` | int32 | ID do item cosmético |
| `skinId` | int32 | ID da skin |

---

## 3. Web Scrapers — Tier Lists (Phase 4)

Fontes: **MetaTFT**, **TFTactics**, **Mobalytics**
Frequência: Crawl diário via Goroutines

### Dados Esperados (a confirmar durante implementação)

| Dado | Descrição |
|------|-----------|
| Tier da composição | S / A / B / C ranking |
| Win rate | Taxa de vitória por comp |
| Play rate | Taxa de uso por comp |
| Avg placement | Colocação média |
| Itens recomendados | BiS (Best in Slot) por campeão |
| Augment tier list | Ranking de augments por tier |
| Unidades da comp | Lista de campeões que formam a comp |

> Detalhes do schema de scraping serão definidos na Phase 4 (issues #22, #23, #24).

---

## Mapeamento: apiName como Chave Universal

```
CommunityDragon                    Riot API
─────────────────                  ────────────────
champion.apiName  ◄──────────────► unit.characterId
item.apiName      ◄──────────────► unit.itemNames[]
trait.apiName     ◄──────────────► trait.name
augment apiName   ◄──────────────► participant.augments[]
```

Este mapeamento é fundamental para o design dos Protobuf contracts — os dados estáticos do CommunityDragon servem como "dicionário" e os dados da Riot API referenciam esse dicionário via `apiName`.

---

## Assets e Imagens

CommunityDragon fornece caminhos de assets no formato:
```
ASSETS/Maps/TFT/Icons/Items/Hexcore/TFT_Item_RecurveBow.TFT_Set13.tex
```

Para construir URLs de imagem:
```
https://raw.communitydragon.org/latest/game/{path_lowercase}
```

Substituir `.tex` por `.png` e converter o path para lowercase.

---

## Referências

- [Riot Developer Portal — TFT](https://developer.riotgames.com/docs/tft)
- [CommunityDragon TFT Data](https://raw.communitydragon.org/latest/cdragon/tft/)
- [CommunityDragon Assets Docs](https://github.com/CommunityDragon/Docs/blob/master/assets.md)
- [Equinox Go TFT Client (structs)](https://pkg.go.dev/github.com/Kyagara/equinox/clients/tft)
- [fightmegg/riot-api TypeScript types](https://github.com/fightmegg/riot-api/blob/master/src/@types/index.ts)
