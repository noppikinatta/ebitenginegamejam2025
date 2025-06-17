# CSV Design for Game Data

This document defines the **CSV formats** that will be used to drive dynamic content for the game.  
All CSV files live under the `data/` directory and are encoded in **UTF-8 without BOM**.  
The loader package (`github.com/noppikinatta/ebitenginegamejam2025/loader`) will provide strongly typed accessors.

---

## 1. Card List – `data/cards.csv`
| Column | Type | Example | Description |
|--------|------|---------|-------------|
| `id` | string | `CARD_WARRIOR` | Unique identifier (primary key). |
| `name` | string | `Warrior` | Display name. |
| `type` | enum `Unit \| Enchant \| Building` | `Unit` | Card category. |
| `attack` | int | `4` | Base attack value. |
| `defense` | int | `3` | Base defense value. |
| `cost_gold` | int | `50` | Gold cost. |
| `cost_iron` | int | `10` | Iron cost. *(0 if not required)* |
| `cost_wood` | int | `0` | Wood cost. |
| `cost_grain` | int | `0` | Grain cost. |
| `cost_mana` | int | `0` | Mana cost. |
| `description` | string | `Basic melee unit` | One-line explanation.

**Sample Row**
```
CARD_WARRIOR,Warrior,Unit,4,3,50,10,0,0,0,Basic melee unit
```

---

## 2. Nation / NPC List – `data/nations.csv`
| Column | Type | Example | Description |
|--------|------|---------|-------------|
| `id` | string | `NATION_IRON_REPUBLIC` | Primary key. |
| `name` | string | `Iron Republic` | Display name. |
| `initial_relationship` | int [0-100] | `60` | Starting relationship toward the player. |
| `ally_bonus_gold` | int | `10` | Monthly bonus gained while allied. |
| `ally_bonus_attack` | int | `1` | Combat attack bonus while allied. |
| `flavor` | string | `A proud nation of smiths.` | Short description.

---

## 3. Enemy List – `data/enemies.csv`
| Column | Type | Example | Description |
|--------|------|---------|-------------|
| `id` | string | `ENEMY_GOBLIN` | Primary key. |
| `name` | string | `Goblin` | Display name. |
| `attack` | int | `3` | Base attack. |
| `defense` | int | `2` | Base defense. |
| `health` | int | `5` | Hit points. |
| `reward_gold` | int | `15` | Gold reward after defeat. |
| `reward_card_id` | string | `CARD_FIRE_WEAPON` | Optional card granted. *(empty if none)* |

---

## 4. Boss List – `data/bosses.csv`
| Column | Type | Example | Description |
|--------|------|---------|-------------|
| `id` | string | `BOSS_DARK_LORD` | Primary key. |
| `name` | string | `Dark Lord` | Display name. |
| `attack` | int | `20` | Base attack. |
| `defense` | int | `12` | Base defense. |
| `health` | int | `200` | Hit points. |
| `reward_gold` | int | `500` | Gold reward. |
| `reward_card_id` | string | `CARD_LEGENDARY_SWORD` | Guaranteed legendary drop.

---

## 5. Loader Package Plan – `loader`
```
loader/
  cards.go      // LoadCards() []entity.CardTemplate
  nations.go    // LoadNations() []entity.Nation
  enemies.go    // LoadEnemies() []entity.EnemyTemplate
  bosses.go     // LoadBosses() []entity.BossTemplate
  util.go       // shared helpers (openCSV, numeric conversion, etc.)
```

### API Sketch
```go
package loader

func LoadCards(path string) (map[string]*entity.Card, error)
func LoadNations(path string) (map[string]*entity.Nation, error)
func LoadEnemies(path string) (map[string]*entity.Enemy, error)
func LoadBosses(path string) (map[string]*entity.Boss, error)
```
* Each function accepts a **file path** so unit tests can inject fixtures.
* Internally uses `encoding/csv`, performs type conversion and validation.
* Unknown columns are ignored to allow forward compatibility.

### Error Handling
* Malformed lines are collected and returned as `[]error` alongside data.
* Critical format errors (missing mandatory columns) cause function to return an error.

### Hot-Reload Support (Future)
Loader functions are pure; the caller (e.g., `CardManager`) can watch file timestamps and re-invoke loaders to support modding or balancing patches without recompiling.

---

## 6. Next Steps
1. Create the `data/` directory and stub CSV files with headers only.  
2. Implement `loader` package following the API sketch.  
3. Update system managers (CardManager, AllianceManager, CombatManager, etc.) to accept injected data so they no longer rely on hard-coded templates.
