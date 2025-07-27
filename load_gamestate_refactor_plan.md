# load/gamestate.go エラー解消 実装計画

## 📋 **概要**

`load/gamestate.go`はゲーム初期状態のロード処理を担当しているが、coreパッケージの再設計により以下の問題が発生：

- 構造体フィールドへの直接値設定が不可能
- privateフィールドへのアクセス不可
- 削除されたstructやinterfaceの参照
- 新しいconstructor APIへの対応が必要

## 🚨 **主要エラー分類**

### 1. **Nation関連エラー**
- `core.BaseNation` → 削除済み
- `MyNation.BasicYield` → 削除済み
- `myNation.Market` → privateフィールド化

### 2. **MarketItem関連エラー**
- 構造体リテラル → `core.NewMarketItem()` 必要
- privateフィールド (`cardPack`, `price`, `requiredLevel`)

### 3. **Enemy/EnemySkill関連エラー**
- 構造体リテラル → `core.NewEnemy()` 必要
- `core.EnemySkillImpl` → 削除済み
- privateフィールド (`enemyType`, `power`, etc.)

### 4. **StructureCard関連エラー**
- 構造体リテラル → `core.NewStructureCard()` 必要
- `core.AddYieldModifier`, `core.MultiplyYieldModifier` → 削除済み
- privateフィールド (`cardID`)

### 5. **Territory/Point関連エラー**
- 構造体リテラル → `core.NewTerritory()` + `core.NewTerrain()` 必要
- `WildernessPoint`, `BossPoint` privateフィールド
- `Territory.BaseYield`, `Territory.CardSlot` → 削除済み

---

## 🚀 **Phase 1: 即座に対応可能な項目**

### 1.1 **EnemySkill関連**
**ステータス**: ✅ 対応可能

**現在の状況**:
- ✅ `core.NewEnemySkill`コンストラクタが存在
- ✅ `Enemy.Type()`ゲッターメソッドが存在

**対応方法**:
```go
// 修正前
return &core.EnemySkillImpl{
    IDField: "enemy-skill-evasion",
    Condition: func(idx int, options *core.EnemySkillCalculationOptions) bool {
        card := options.BattleCards[idx]
        return card.Type == "cardtype-str"
    },
    Modifier: &core.BattleCardPowerModifier{
        AdditiveDebuff: 2.0,
    },
}

// 修正後
return core.NewEnemySkill(
    "enemy-skill-evasion",
    func(idx int, options *core.EnemySkillCalculationOptions) bool {
        card := options.BattleCards[idx]
        return card.Type == "cardtype-str"
    },
    &core.BattleCardPowerModifier{
        AdditiveDebuff: 2.0,
    },
)
```

**Enemy.EnemyType アクセス修正**:
```go
// 修正前
options.Enemy.EnemyType == "enemy-type-dragon"

// 修正後  
options.Enemy.Type() == "enemy-type-dragon"
```

### 1.2 **StructureCard作成**
**ステータス**: ✅ 対応可能

**現在のAPI**: `NewStructureCard(cardID, yieldAdditiveValue, yieldModifier, supportPower, supportCardSlot)`

**対応パターン**:

```go
// Yield Additive系 (farm, woodcutter, tunnel, market, shrine)
core.NewStructureCard(
    "structurecard-farm", 
    core.ResourceQuantity{Food: 2}, 
    core.NewResourceModifier(), 
    0.0, 0)

// Yield Multiplicative系 (granary, sawmill, smelter, mint, temple) 
core.NewStructureCard(
    "structurecard-granary", 
    core.ResourceQuantity{}, 
    core.ResourceModifier{Food: 0.5}, 
    0.0, 0)

// Support Power系 (catapult, ballista, orban-cannon)
core.NewStructureCard(
    "structurecard-catapult", 
    core.ResourceQuantity{}, 
    core.NewResourceModifier(), 
    3.0, 0)

// Support CardSlot系 (camp)
core.NewStructureCard(
    "structurecard-camp", 
    core.ResourceQuantity{}, 
    core.NewResourceModifier(), 
    0.0, 1)
```

### 1.3 **Enemy作成**
**ステータス**: ✅ 対応可能

```go
// 修正前
enemy := &core.Enemy{
    EnemyID:        config.enemyID,
    EnemyType:      config.enemyType,
    Power:          config.power,
    Skills:         config.skills,
    BattleCardSlot: config.cardSlot,
    Question:       fmt.Sprintf("question-%d", i+1), // 削除
}

// 修正後
enemy := core.NewEnemy(
    config.enemyID,
    config.enemyType,
    config.power,
    config.skills,
    config.cardSlot,
)
// Note: Questionフィールドは削除済みのため対応不要
```

### 1.4 **Territory/Terrain作成**
**ステータス**: ✅ 対応可能

```go
// 修正前
territory := &core.Territory{
    BaseYield: config.baseYield,
    CardSlot:  3,
}

// 修正後
terrain := core.NewTerrain(
    core.TerrainID(config.terrainType), 
    config.baseYield, 
    3)
territory := core.NewTerritory("territory_"+config.enemyID, terrain)
```

---

## 🔄 **Phase 2: アーキテクチャ課題の解決**

### 2.1 **Market管理問題**
**ステータス**: 🔧 アーキテクチャ変更必要

**現状の問題**:
```go
// これが不可能
myNation.Market = &core.Market{...}
```

**解決方針A: GameStateでMarket管理**
```go
type GameState struct {
    // ... existing fields ...
    Markets map[NationID]*Market // Nation別Market管理
}

// 使用例
gameState.Markets[myNation.ID()] = &core.Market{...}
```

**解決方針B: 独立したMarket管理**
```go
// createMarkets() 関数を新設
func createMarkets() map[NationID]*Market {
    return map[NationID]*Market{
        "nation-mynation": createMyNationMarket(),
        "nation-forest":   createForestMarket(),
        // ...
    }
}
```

### 2.2 **MarketItem作成**
**ステータス**: ✅ 対応可能

```go
// 修正前
{
    CardPack:      cardPacks["cardpack-free"],
    Price:         cardPackPrices["cardpack-free"],
    RequiredLevel: 1,
}

// 修正後
core.NewMarketItem(
    cardPacks["cardpack-free"],
    cardPackPrices["cardpack-free"],
    1.0,  // requiredLevel
    0.0,  // levelEffect
)
```

### 2.3 **Point作成対応**
**ステータス**: 🔧 部分的対応必要

**WildernessPoint**:
```go
// 修正前
&core.WildernessPoint{
    TerrainType: config.terrainType,
    Controlled:  false,
    Enemy:       enemy,
    Territory:   territory,
}

// 修正後 (test setter利用)
wilderness := &core.WildernessPoint{}
wilderness.SetTerrainTypeForTest(config.terrainType)  // 要追加
wilderness.SetControlledForTest(false)
wilderness.SetEnemyForTest(enemy)
wilderness.SetTerritoryForTest(territory)
```

**BossPoint**:
```go
// 修正前
&core.BossPoint{
    Boss:     boss,
    Defeated: false,
}

// 修正後
bossPoint := &core.BossPoint{}
bossPoint.SetBossForTest(boss)
bossPoint.SetDefeatedForTest(false)
```

---

## 🛠️ **Phase 3: 追加実装が必要な項目**

### 3.1 **core パッケージに追加必要**

#### WildernessPoint関連
```go
// core/mapgrid.go に追加
func (p *WildernessPoint) TerrainType() string { 
    return p.terrainType 
}

func (p *WildernessPoint) SetTerrainTypeForTest(terrainType string) { 
    p.terrainType = terrainType 
}
```

#### GameState Market管理 (方針Aの場合)
```go
// core/gameflow.go に追加
type GameState struct {
    // ... existing fields ...
    Markets map[NationID]*Market
}

func (g *GameState) GetMarket(nationID NationID) *Market {
    return g.Markets[nationID]
}

func (g *GameState) SetMarket(nationID NationID, market *Market) {
    if g.Markets == nil {
        g.Markets = make(map[NationID]*Market)
    }
    g.Markets[nationID] = market
}
```

### 3.2 **load/gamestate.go 修正項目一覧**

#### 関数別修正内容
- **createMyNation()** → `core.NewMyNation`使用
- **createCardDeck()** → 変更不要
- **createStructureCards()** → 全構造体を新パラメータ形式に対応
- **createBattleCards()** → Enemy.Type()アクセス修正
- **createEnemySkill系関数(12個)** → `core.NewEnemySkill`使用
- **createMapGrid()** → 各Point作成時に新API使用、Market管理方式変更

---

## 📅 **実装スケジュール**

### Step 1: Phase 1対応 (即座対応可能)
- [ ] EnemySkill作成関数群の修正 (12関数)
- [ ] StructureCard作成の修正 (14カード)
- [ ] Enemy作成の修正
- [ ] Territory/Terrain作成の修正

### Step 2: Phase 3.1対応 (core追加実装)
- [ ] WildernessPoint.TerrainType()追加
- [ ] WildernessPoint.SetTerrainTypeForTest()追加
- [ ] Market管理アーキテクチャ実装

### Step 3: Phase 2対応 (アーキテクチャ変更)
- [ ] Market管理方式の適用
- [ ] MarketItem作成の修正
- [ ] Point作成の修正

### Step 4: 統合テスト
- [ ] load.LoadGameState()動作確認
- [ ] ゲーム起動確認
- [ ] 各機能の動作確認

---

## 🎯 **成功基準**

1. **コンパイルエラー解消**: load/gamestate.goのlinter errorsが0になる
2. **機能保持**: 既存のゲームロード機能が正常動作する
3. **アーキテクチャ準拠**: 新しいcoreパッケージ設計に準拠する
4. **テスト通過**: 関連するテストが全てPASSする

---

## 📝 **実装時の注意点**

1. **段階的実装**: Phase順に進めて各段階でテスト確認
2. **互換性**: 既存のUI/flow層への影響を最小限に
3. **エラーハンドリング**: 新しいAPIでのエラーケース対応
4. **パフォーマンス**: Market管理方式での性能劣化を避ける

---

*Created: 2024* | *Status: Planning Phase* 