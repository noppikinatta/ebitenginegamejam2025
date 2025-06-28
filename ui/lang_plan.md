# UI多言語化計画

## 調査概要

`/ui`ディレクトリのコンポーネントを調査し、`lang.Text()`で多言語化すべき箇所を特定しました。現在多くの文字列がハードコードされており、`asset/lang/`の言語データファイルに定義されたキーを使用することで多言語対応できます。

## 多言語化が必要な箇所

### 1. battle.go - バトル画面
- **タイトル**: `"Battle of %s"` → テンプレート形式に変更が必要
- **ボタン**: `"CONQUER"`, `"Need Power"` → 新しいキーが必要
- **敵情報**: `"Enemy: %s"`, `"Power: %.1f"`, `"Card Limit: %d"` → 新しいキーが必要
- **勝敗メッセージ**: `"CLICK TO WIN!"`, `"Need more power"` → 新しいキーが必要
- **敵名**: `enemy.EnemyID` → `lang.Text("enemy-" + enemy.EnemyID)`で既存キー利用可能

### 2. market.go - マーケット画面
- **カードパック名**: `item.CardPack.CardPackID` → `lang.Text("cardpack-" + cardPackID)`で既存キー利用可能
- **説明**: `"Required Level: %.1f"`, `"Card pack with various cards"` → 新しいキーが必要
- **リソース名**: `"Money"`, `"Food"`, `"Wood"`, `"Iron"`, `"Mana"` → 新しいキーが必要

### 3. carddeck.go - カードデッキ画面
- **メッセージ**: `"No card deck"`, `"No cards in deck"` → 新しいキーが必要
- **カード情報**: `"Type: Battle"`, `"Type: Structure"`, `"STR"` → 新しいキーが必要
- **BattleCard名**: `card.CardID` → `lang.Text("battlecard-" + cardID)`で既存キー利用可能
- **StructureCard名**: `card.CardID` → `lang.Text("structurecard-" + cardID)`で既存キー利用可能
- **カードタイプ**: `card.Type` → `lang.Text("cardtype-" + cardType)`で既存キー利用可能

### 4. territory.go - 領土画面
- **ボタン**: `"CONFIRM"`, `"No Changes"` → 新しいキーが必要
- **効果説明**: `"Structure Effects:"`, `"No structure cards placed."`, `"Place cards to get bonuses!"` → 新しいキーが必要
- **カード効果**: `"Boosts resource production"` → 新しいキーが必要
- **StructureCard名**: `card.CardID` → `lang.Text("structurecard-" + cardID)`で既存キー利用可能

### 5. info.go - 情報画面
- **タイトル**: `"History"`, `"Card:"`, `"Type: Battle"`, `"Type: Structure"` → 新しいキーが必要
- **カード情報**: `"Class: %s"`, `"Power: %.1f"`, `"Skill: Active"`, `"Skill: None"` → 新しいキーが必要
- **国名**: `"My Nation"` → `lang.Text("nation-mynation")` で既存キー利用可能
- **地点名**: `"Wilderness"`, `"Boss Point"`, `"Enemy:"`, `"Boss:"` → 新しいキーが必要
- **リソース**: `"Yields:"`, `"Structures:"` → 新しいキーが必要
- **敵名**: `enemy.EnemyID` → `lang.Text("enemy-" + enemyID)`で既存キー利用可能

### 6. mapgrid.go - マップグリッド画面
- **地点名**: 
  - `"My Nation"` → `lang.Text("nation-mynation")`
  - `"Nation %s"` → 国名の多言語化が必要
  - `"Area %d,%d"`, `"Wild %d,%d"`, `"Point %d,%d"` → 新しいキーが必要
  - `"Boss"` → 新しいキーが必要

### 7. resource.go - リソース表示
- **リソース名**: `"Money"`, `"Food"`, `"Wood"`, `"Iron"`, `"Mana"` → 新しいキーが必要

## 必要な新キーの提案

以下のキーを言語ファイルに追加する必要があります：

```csv
# UI基本
ui-back, "Back"
ui-confirm, "Confirm"
ui-no-changes, "No Changes"
ui-conquer, "Conquer"
ui-need-power, "Need Power"
ui-click-to-win, "Click to Win!"
ui-need-more-power, "Need more power"
ui-history, "History"
ui-no-events, "No events yet."

# バトル関連
battle-title, "Battle of {{.location}}"
battle-enemy, "Enemy:"
battle-power, "Power:"
battle-card-limit, "Card Limit:"

# カード関連
card-type-battle, "Battle"
card-type-structure, "Structure"
card-class, "Class:"
card-skill-active, "Skill: Active"
card-skill-none, "Skill: None"
card-no-deck, "No card deck"
card-no-cards, "No cards in deck"

# 領土関連
territory-structure-effects, "Structure Effects:"
territory-no-structures, "No structure cards placed."
territory-place-cards, "Place cards to get bonuses!"
territory-yield-effect, "Yield Effect:"
territory-boosts-production, "Boosts resource production"
territory-no-yield-effect, "No yield effect"

# リソース名
resource-money, "Money"
resource-food, "Food"
resource-wood, "Wood"
resource-iron, "Iron"
resource-mana, "Mana"

# 地点関連
point-wilderness, "Wilderness"
point-boss, "Boss Point"
point-boss, "Boss"
point-area, "Area {{.x}},{{.y}}"
point-wild, "Wild {{.x}},{{.y}}"
point-yields, "Yields:"
point-structures, "Structures:"

# マーケット関連
market-required-level, "Required Level: {{.level}}"
market-card-pack-desc, "Card pack with various cards"
```

## 実装の優先順位

1. **高優先度**: 敵名、カード名、国名（既存キーが利用可能）
2. **中優先度**: UI基本文字列（戻るボタン、確認ボタンなど）
3. **低優先度**: 説明文、ヘルプテキスト

## 実装方法

各コンポーネントで以下のように変更：

```go
// Before
drawing.DrawText(screen, "CONQUER", 14, opt)

// After  
drawing.DrawText(screen, lang.Text("ui-conquer"), 14, opt)
```

```go
// Before
enemyName := fmt.Sprintf("Enemy: %s", enemy.EnemyID)

// After
enemyName := fmt.Sprintf("%s: %s", lang.Text("battle-enemy"), lang.Text("enemy-" + string(enemy.EnemyID)))
```
