# 設計文書

## 概要

この設計は、coreパッケージの包括的なリファクタリングを行い、カプセル化の改善、アーキテクチャの明確化、および複数カードインスタンスのサポートを実現します。設計では、既存のコードベースを段階的に変更し、新しいviewmodelとflowパッケージを導入してクリーンアーキテクチャを実現します。

## アーキテクチャ

### レイヤー構造

```
UI Layer (ui package)
    ↓
Presentation Layer (viewmodel package) 
    ↓
Use Case Layer (flow package)
    ↓
Domain Layer (core package)
    ↓
Infrastructure Layer (load package)
```

### パッケージ責務

- **core**: ゲームロジックとドメインモデル（不変オブジェクト、ビジネスルール）
- **viewmodel**: UI表示用のデータ変換とプレゼンテーション
- **flow**: ユースケース操作とアプリケーションロジック
- **ui**: ユーザーインターフェースとイベント処理
- **load**: データ永続化とゲーム状態の読み込み

## コンポーネントと インターフェース

### 1. CardDeck の再設計

#### 現在の問題
- BattleCardへの直接参照を保持
- Experience概念による複雑な重複管理
- カプセル化不足

#### 新設計

```go
type CardDeck struct {
    hand map[CardID]int
}

func NewCardDeck() *CardDeck {
    return &CardDeck{
        hand: make(map[CardID]int),
    }
}

func (cd *CardDeck) Add(cardID CardID) {
    cd.hand[cardID]++
}

func (cd *CardDeck) Remove(cardID CardID) bool {
    if cd.hand[cardID] <= 0 {
        return false
    }
    cd.hand[cardID]--
    if cd.hand[cardID] == 0 {
        delete(cd.hand, cardID)
    }
    return true
}

func (cd *CardDeck) Count(cardID CardID) int {
    return cd.hand[cardID]
}

func (cd *CardDeck) GetAllCardIDs() []CardID {
    var cardIDs []CardID
    for cardID, count := range cd.hand {
        for i := 0; i < count; i++ {
            cardIDs = append(cardIDs, cardID)
        }
    }
    return cardIDs
}
```

### 2. BattleCard の不変化

#### 現在の問題
- Experienceフィールドによる可変性
- 公開フィールドによるカプセル化不足

#### 新設計

```go
type BattleCard struct {
    id       CardID
    power    float64
    skill    *BattleCardSkill
    cardType BattleCardType
}

func NewBattleCard(id CardID, power float64, skill *BattleCardSkill, cardType BattleCardType) *BattleCard {
    return &BattleCard{
        id:       id,
        power:    power,
        skill:    skill,
        cardType: cardType,
    }
}

func (c *BattleCard) ID() CardID { return c.id }
func (c *BattleCard) Power() float64 { return c.power }
func (c *BattleCard) Skill() *BattleCardSkill { return c.skill }
func (c *BattleCard) Type() BattleCardType { return c.cardType }
```

### 3. BattleCardSkill の簡素化

#### 現在の問題
- DescriptionKeyによる表示関心事の混入

#### 新設計

```go
type BattleCardSkill struct {
    id         BattleCardSkillID
    calculator BattleCardSkillCalculator
}

func NewBattleCardSkill(id BattleCardSkillID, calculator BattleCardSkillCalculator) *BattleCardSkill {
    return &BattleCardSkill{
        id:         id,
        calculator: calculator,
    }
}

func (s *BattleCardSkill) ID() BattleCardSkillID { return s.id }
func (s *BattleCardSkill) Calculate(options *BattleCardSkillCalculationOptions) {
    s.calculator.Calculate(options)
}
```

### 4. StructureCard の再設計

#### 現在の問題
- DescriptionKeyによる表示関心事の混入
- 複雑な抽象化

#### 新設計

```go
type StructureCard struct {
    cardID              CardID
    yieldAdditiveValue  ResourceQuantity
    yieldModifier       ResourceModifier
    supportPower        float64
    supportCardSlot     int
}

func NewStructureCard(cardID CardID, yieldAdditiveValue ResourceQuantity, yieldModifier ResourceModifier, supportPower float64, supportCardSlot int) *StructureCard {
    return &StructureCard{
        cardID:              cardID,
        yieldAdditiveValue:  yieldAdditiveValue,
        yieldModifier:       yieldModifier,
        supportPower:        supportPower,
        supportCardSlot:     supportCardSlot,
    }
}

func (c *StructureCard) ID() CardID { return c.cardID }
func (c *StructureCard) YieldAdditiveValue() ResourceQuantity { return c.yieldAdditiveValue }
func (c *StructureCard) YieldModifier() ResourceModifier { return c.yieldModifier }
func (c *StructureCard) SupportPower() float64 { return c.supportPower }
func (c *StructureCard) SupportCardSlot() int { return c.supportCardSlot }
```

### 5. Territory と Terrain の分離

#### 現在の問題
- 可変な建設計画と不変な地形データの混在

#### 新設計

```go
type TerrainID string

type Terrain struct {
    id        TerrainID
    baseYield ResourceQuantity
    cardSlot  int
}

func NewTerrain(id TerrainID, baseYield ResourceQuantity, cardSlot int) *Terrain {
    return &Terrain{
        id:        id,
        baseYield: baseYield,
        cardSlot:  cardSlot,
    }
}

func (t *Terrain) ID() TerrainID { return t.id }
func (t *Terrain) BaseYield() ResourceQuantity { return t.baseYield }
func (t *Terrain) CardSlot() int { return t.cardSlot }

type Territory struct {
    id      TerritoryID
    terrain *Terrain
    cards   []*StructureCard
}

func NewTerritory(id TerritoryID, terrain *Terrain) *Territory {
    return &Territory{
        id:      id,
        terrain: terrain,
        cards:   make([]*StructureCard, 0, terrain.CardSlot()),
    }
}

func (t *Territory) ID() TerritoryID { return t.id }
func (t *Territory) Terrain() *Terrain { return t.terrain }
func (t *Territory) Cards() []*StructureCard { 
    // 防御的コピーを返す
    result := make([]*StructureCard, len(t.cards))
    copy(result, t.cards)
    return result
}

type ConstructionPlan struct {
    cards []*StructureCard
}

func NewConstructionPlan(territory *Territory) *ConstructionPlan {
    // 既存のカードを防御的コピー
    cards := make([]*StructureCard, len(territory.cards))
    copy(cards, territory.cards)
    return &ConstructionPlan{cards: cards}
}

func (cp *ConstructionPlan) Cards() []*StructureCard {
    result := make([]*StructureCard, len(cp.cards))
    copy(result, cp.cards)
    return result
}

func (cp *ConstructionPlan) AddCard(card *StructureCard) bool {
    // カードスロット制限チェックは呼び出し側で行う
    cp.cards = append(cp.cards, card)
    return true
}

func (cp *ConstructionPlan) RemoveCard(index int) (*StructureCard, bool) {
    if index < 0 || index >= len(cp.cards) {
        return nil, false
    }
    card := cp.cards[index]
    cp.cards = append(cp.cards[:index], cp.cards[index+1:]...)
    return card, true
}

func (t *Territory) ApplyConstructionPlan(plan *ConstructionPlan) {
    // メモリ共有を避けるため防御的コピー
    t.cards = make([]*StructureCard, len(plan.cards))
    copy(t.cards, plan.cards)
}
```

### 6. Point システムの改善

#### 現在の問題
- UIでの直接型キャスト
- 型安全性の不足

#### 新設計

```go
type PointType int

const (
    PointTypeUnknown PointType = iota
    PointTypeMyNation
    PointTypeOtherNation
    PointTypeWilderness
    PointTypeBoss
)

type Point interface {
    PointType() PointType
    Passable() bool
    AsBattlePoint() (BattlePoint, bool)
    AsTerritoryPoint() (TerritoryPoint, bool)
    AsMarketPoint() (MarketPoint, bool)
}

type BattlePoint interface {
    Enemy() *Enemy
    Conquer()
}

type TerritoryPoint interface {
    Yield() ResourceQuantity
    Terrain() *Terrain
    CardSlot() int
    Cards() []*StructureCard
}

type MarketPoint interface {
    Nation() Nation
}
```

### 7. Enemy の簡素化

#### 現在の問題
- Questionフィールドによる表示関心事の混入
- EnemySkillの過度な抽象化

#### 新設計

```go
type Enemy struct {
    id             EnemyID
    enemyType      EnemyType
    power          float64
    skills         []*EnemySkill
    battleCardSlot int
}

func NewEnemy(id EnemyID, enemyType EnemyType, power float64, skills []*EnemySkill, battleCardSlot int) *Enemy {
    return &Enemy{
        id:             id,
        enemyType:      enemyType,
        power:          power,
        skills:         skills,
        battleCardSlot: battleCardSlot,
    }
}

type EnemySkill struct {
    id        EnemySkillID
    condition func(idx int, options *EnemySkillCalculationOptions) bool
    modifier  *BattleCardPowerModifier
}

func NewEnemySkill(id EnemySkillID, condition func(idx int, options *EnemySkillCalculationOptions) bool, modifier *BattleCardPowerModifier) *EnemySkill {
    return &EnemySkill{
        id:        id,
        condition: condition,
        modifier:  modifier,
    }
}
```

### 8. Market システムの強化

#### 現在の問題
- レベル効果の手動管理
- 投資型アイテムの未サポート

#### 新設計

```go
type MarketItem struct {
    cardPack      *CardPack        // nilを許可（投資型アイテム用）
    price         ResourceQuantity
    requiredLevel MarketLevel
    levelEffect   MarketLevel      // 購入時のレベル効果
    resourceQuantity *ResourceQuantity // リソース取引用（nilを許可）
}

func NewMarketItem(cardPack *CardPack, price ResourceQuantity, requiredLevel MarketLevel, levelEffect MarketLevel) *MarketItem {
    return &MarketItem{
        cardPack:      cardPack,
        price:         price,
        requiredLevel: requiredLevel,
        levelEffect:   levelEffect,
    }
}

func (m *Market) Purchase(index int, treasury *Treasury) (*CardPack, bool) {
    if !m.CanPurchase(index, treasury) {
        return nil, false
    }

    item := m.Items[index]
    if !treasury.Sub(item.Price) {
        return nil, false
    }

    // レベル効果を自動適用
    m.Level += item.levelEffect

    // リソース取引の処理
    if item.resourceQuantity != nil {
        treasury.Add(*item.resourceQuantity)
    }

    return item.cardPack, true
}
```

## データモデル

### ViewModelパッケージ

ViewModelは描画用の情報を直接保持せず、core.GameStateを参照し、langパッケージやdrawingパッケージの機能を使って描画に必要な情報を都度計算して返します。これにより、UIパッケージで直接GameStateを参照していた処理をViewModelのメソッドに移行できます。

```go
// BattleViewModel - バトル画面用
type BattleViewModel struct {
    gameState   *core.GameState
    battlefield *core.Battlefield
}

func NewBattleViewModel(gameState *core.GameState, battlefield *core.Battlefield) *BattleViewModel {
    return &BattleViewModel{
        gameState:   gameState,
        battlefield: battlefield,
    }
}

func (vm *BattleViewModel) Title() string {
    // langパッケージを使用してローカライズされたタイトルを取得
    return lang.GetText("battle_title")
}

func (vm *BattleViewModel) EnemyImage() *ebiten.Image {
    // drawingパッケージを使用して敵の画像を取得
    return drawing.GetEnemyImage(vm.battlefield.Enemy().ID())
}

func (vm *BattleViewModel) EnemyType() string {
    // langパッケージを使用してローカライズされた敵タイプを取得
    return lang.GetEnemyTypeName(vm.battlefield.Enemy().Type())
}

func (vm *BattleViewModel) EnemyPower() float64 {
    return vm.battlefield.Enemy().Power()
}

func (vm *BattleViewModel) EnemyTalk() string {
    // langパッケージを使用して敵の台詞を取得
    return lang.GetEnemyTalk(vm.battlefield.Enemy().ID())
}

func (vm *BattleViewModel) EnemySkillNames() []string {
    skills := vm.battlefield.Enemy().Skills()
    names := make([]string, len(skills))
    for i, skill := range skills {
        names[i] = lang.GetEnemySkillName(skill.ID())
    }
    return names
}

func (vm *BattleViewModel) EnemySkillDescriptions() []string {
    skills := vm.battlefield.Enemy().Skills()
    descriptions := make([]string, len(skills))
    for i, skill := range skills {
        descriptions[i] = lang.GetEnemySkillDescription(skill.ID())
    }
    return descriptions
}

func (vm *BattleViewModel) CardSlot() int {
    return vm.battlefield.Enemy().BattleCardSlot()
}

func (vm *BattleViewModel) CanBeat() bool {
    return vm.battlefield.CanBeat()
}

func (vm *BattleViewModel) TotalPower() float64 {
    return vm.battlefield.CalculateTotalPower()
}

func (vm *BattleViewModel) NumCards() int {
    return len(vm.battlefield.BattleCards())
}

func (vm *BattleViewModel) Card(idx int) *BattleCardViewModel {
    cards := vm.battlefield.BattleCards()
    if idx < 0 || idx >= len(cards) {
        return nil
    }
    return NewBattleCardViewModel(vm.gameState, cards[idx], vm.battlefield.CalculateCardPower(idx))
}

// CardDeckViewModel - カードデッキ表示用
type CardDeckViewModel struct {
    gameState *core.GameState
}

func NewCardDeckViewModel(gameState *core.GameState) *CardDeckViewModel {
    return &CardDeckViewModel{gameState: gameState}
}

func (vm *CardDeckViewModel) NumBattleCards() int {
    count := 0
    for cardID := range vm.gameState.CardDeck.GetAllCardCounts() {
        if vm.gameState.GetBattleCard(cardID) != nil {
            count++
        }
    }
    return count
}

func (vm *CardDeckViewModel) NumStructureCards() int {
    count := 0
    for cardID := range vm.gameState.CardDeck.GetAllCardCounts() {
        if vm.gameState.GetStructureCard(cardID) != nil {
            count++
        }
    }
    return count
}

func (vm *CardDeckViewModel) BattleCard(idx int) *BattleCardViewModel {
    battleCardIDs := vm.getBattleCardIDs()
    if idx < 0 || idx >= len(battleCardIDs) {
        return nil
    }
    card := vm.gameState.GetBattleCard(battleCardIDs[idx])
    return NewBattleCardViewModel(vm.gameState, card, card.Power())
}

func (vm *CardDeckViewModel) StructureCard(idx int) *StructureCardViewModel {
    structureCardIDs := vm.getStructureCardIDs()
    if idx < 0 || idx >= len(structureCardIDs) {
        return nil
    }
    card := vm.gameState.GetStructureCard(structureCardIDs[idx])
    return NewStructureCardViewModel(vm.gameState, card)
}

func (vm *CardDeckViewModel) getBattleCardIDs() []CardID {
    var cardIDs []CardID
    for cardID := range vm.gameState.CardDeck.GetAllCardCounts() {
        if vm.gameState.GetBattleCard(cardID) != nil {
            cardIDs = append(cardIDs, cardID)
        }
    }
    return cardIDs
}

func (vm *CardDeckViewModel) getStructureCardIDs() []CardID {
    var cardIDs []CardID
    for cardID := range vm.gameState.CardDeck.GetAllCardCounts() {
        if vm.gameState.GetStructureCard(cardID) != nil {
            cardIDs = append(cardIDs, cardID)
        }
    }
    return cardIDs
}

// BattleCardViewModel - バトルカード表示用
type BattleCardViewModel struct {
    gameState     *core.GameState
    card          *core.BattleCard
    calculatedPower float64
}

func NewBattleCardViewModel(gameState *core.GameState, card *core.BattleCard, calculatedPower float64) *BattleCardViewModel {
    return &BattleCardViewModel{
        gameState:       gameState,
        card:            card,
        calculatedPower: calculatedPower,
    }
}

func (vm *BattleCardViewModel) Image() *ebiten.Image {
    // drawingパッケージを使用してカード画像を取得
    return drawing.GetBattleCardImage(vm.card.ID())
}

func (vm *BattleCardViewModel) Name() string {
    // langパッケージを使用してローカライズされたカード名を取得
    return lang.GetBattleCardName(vm.card.ID())
}

func (vm *BattleCardViewModel) Duplicates() int {
    return vm.gameState.CardDeck.Count(vm.card.ID())
}

func (vm *BattleCardViewModel) CardTypeImage() *ebiten.Image {
    // drawingパッケージを使用してカードタイプ画像を取得
    return drawing.GetBattleCardTypeImage(vm.card.Type())
}

func (vm *BattleCardViewModel) CardTypeName() string {
    // langパッケージを使用してローカライズされたカードタイプ名を取得
    return lang.GetBattleCardTypeName(vm.card.Type())
}

func (vm *BattleCardViewModel) Power() float64 {
    return vm.calculatedPower
}

func (vm *BattleCardViewModel) SkillName() string {
    if vm.card.Skill() == nil {
        return ""
    }
    // langパッケージを使用してローカライズされたスキル名を取得
    return lang.GetBattleCardSkillName(vm.card.Skill().ID())
}

func (vm *BattleCardViewModel) SkillDescription() string {
    if vm.card.Skill() == nil {
        return ""
    }
    // langパッケージを使用してローカライズされたスキル説明を取得
    return lang.GetBattleCardSkillDescription(vm.card.Skill().ID())
}

// StructureCardViewModel - 建設カード表示用
type StructureCardViewModel struct {
    gameState *core.GameState
    card      *core.StructureCard
}

func NewStructureCardViewModel(gameState *core.GameState, card *core.StructureCard) *StructureCardViewModel {
    return &StructureCardViewModel{
        gameState: gameState,
        card:      card,
    }
}

func (vm *StructureCardViewModel) Image() *ebiten.Image {
    // drawingパッケージを使用してカード画像を取得
    return drawing.GetStructureCardImage(vm.card.ID())
}

func (vm *StructureCardViewModel) Name() string {
    // langパッケージを使用してローカライズされたカード名を取得
    return lang.GetStructureCardName(vm.card.ID())
}

func (vm *StructureCardViewModel) Duplicates() int {
    return vm.gameState.CardDeck.Count(vm.card.ID())
}

// MapGridViewModel - マップグリッド表示用
type MapGridViewModel struct {
    gameState *core.GameState
}

func NewMapGridViewModel(gameState *core.GameState) *MapGridViewModel {
    return &MapGridViewModel{gameState: gameState}
}

func (vm *MapGridViewModel) Size() core.MapGridSize {
    return vm.gameState.MapGrid.Size()
}

func (vm *MapGridViewModel) Point(x, y int) *PointViewModel {
    point := vm.gameState.MapGrid.GetPoint(x, y)
    if point == nil {
        return nil
    }
    return NewPointViewModel(vm.gameState, point)
}

func (vm *MapGridViewModel) ShouldDrawLineToRight(x, y int) bool {
    // マップグリッドの接続情報を使用して線を描画するかを判定
    return vm.gameState.MapGrid.HasConnectionToRight(x, y)
}

func (vm *MapGridViewModel) ShouldDrawLineToUpper(x, y int) bool {
    // マップグリッドの接続情報を使用して線を描画するかを判定
    return vm.gameState.MapGrid.HasConnectionToUpper(x, y)
}

// PointViewModel - ポイント表示用
type PointViewModel struct {
    gameState *core.GameState
    point     core.Point
}

func NewPointViewModel(gameState *core.GameState, point core.Point) *PointViewModel {
    return &PointViewModel{
        gameState: gameState,
        point:     point,
    }
}

func (vm *PointViewModel) Image() *ebiten.Image {
    // drawingパッケージを使用してポイント画像を取得
    switch vm.point.PointType() {
    case core.PointTypeMyNation:
        return drawing.GetMyNationPointImage()
    case core.PointTypeOtherNation:
        return drawing.GetOtherNationPointImage()
    case core.PointTypeWilderness:
        return drawing.GetWildernessPointImage()
    case core.PointTypeBoss:
        return drawing.GetBossPointImage()
    default:
        return drawing.GetUnknownPointImage()
    }
}

func (vm *PointViewModel) Name() string {
    // langパッケージを使用してローカライズされたポイント名を取得
    return lang.GetPointName(vm.point)
}

func (vm *PointViewModel) HasEnemy() bool {
    if battlePoint, ok := vm.point.AsBattlePoint(); ok {
        return battlePoint.Enemy() != nil
    }
    return false
}

func (vm *PointViewModel) EnemyPower() float64 {
    if battlePoint, ok := vm.point.AsBattlePoint(); ok {
        if enemy := battlePoint.Enemy(); enemy != nil {
            return enemy.Power()
        }
    }
    return 0.0
}

// MarketViewModel - マーケット表示用
type MarketViewModel struct {
    gameState *core.GameState
    market    *core.Market
}

func NewMarketViewModel(gameState *core.GameState, market *core.Market) *MarketViewModel {
    return &MarketViewModel{
        gameState: gameState,
        market:    market,
    }
}

func (vm *MarketViewModel) Title() string {
    // langパッケージを使用してローカライズされたマーケットタイトルを取得
    return lang.GetMarketTitle(vm.market.ID())
}

func (vm *MarketViewModel) Level() float64 {
    return vm.market.Level()
}

func (vm *MarketViewModel) NumItems() int {
    return len(vm.market.Items())
}

func (vm *MarketViewModel) Item(idx int) *MarketItemViewModel {
    items := vm.market.Items()
    if idx < 0 || idx >= len(items) {
        return nil
    }
    return NewMarketItemViewModel(vm.gameState, items[idx])
}

// MarketItemViewModel - マーケットアイテム表示用
type MarketItemViewModel struct {
    gameState *core.GameState
    item      *core.MarketItem
}

func NewMarketItemViewModel(gameState *core.GameState, item *core.MarketItem) *MarketItemViewModel {
    return &MarketItemViewModel{
        gameState: gameState,
        item:      item,
    }
}

func (vm *MarketItemViewModel) ItemName() string {
    // langパッケージを使用してローカライズされたアイテム名を取得
    return lang.GetMarketItemName(vm.item.ID())
}

func (vm *MarketItemViewModel) RequiredLevel() int {
    return int(vm.item.RequiredLevel())
}

func (vm *MarketItemViewModel) Unlocked() bool {
    return vm.item.RequiredLevel() <= vm.gameState.GetCurrentMarket().Level()
}

func (vm *MarketItemViewModel) CanPurchase() bool {
    return vm.Unlocked() && vm.gameState.Treasury.CanAfford(vm.item.Price())
}

func (vm *MarketItemViewModel) Price() core.ResourceQuantity {
    return vm.item.Price()
}

func (vm *MarketItemViewModel) ResourceSufficiency() *ResourceSufficiency {
    treasury := vm.gameState.Treasury
    price := vm.item.Price()
    
    return &ResourceSufficiency{
        Money: treasury.Money() >= price.Money,
        Food:  treasury.Food() >= price.Food,
        Wood:  treasury.Wood() >= price.Wood,
        Iron:  treasury.Iron() >= price.Iron,
        Mana:  treasury.Mana() >= price.Mana,
    }
}

// ResourceSufficiency - リソース充足状況
type ResourceSufficiency struct {
    Money bool
    Food  bool
    Wood  bool
    Iron  bool
    Mana  bool
}

// TerritoryViewModel - テリトリー表示用
type TerritoryViewModel struct {
    gameState *core.GameState
    territory *core.Territory
    plan      *core.ConstructionPlan
}

func NewTerritoryViewModel(gameState *core.GameState, territory *core.Territory, plan *core.ConstructionPlan) *TerritoryViewModel {
    return &TerritoryViewModel{
        gameState: gameState,
        territory: territory,
        plan:      plan,
    }
}

func (vm *TerritoryViewModel) Title() string {
    // langパッケージを使用してローカライズされたテリトリータイトルを取得
    return lang.GetTerritoryTitle(vm.territory.ID())
}

func (vm *TerritoryViewModel) CardSlot() int {
    return vm.territory.Terrain().CardSlot()
}

func (vm *TerritoryViewModel) NumCards() int {
    if vm.plan != nil {
        return len(vm.plan.Cards())
    }
    return len(vm.territory.Cards())
}

func (vm *TerritoryViewModel) Card(idx int) *StructureCardViewModel {
    var cards []*core.StructureCard
    if vm.plan != nil {
        cards = vm.plan.Cards()
    } else {
        cards = vm.territory.Cards()
    }
    
    if idx < 0 || idx >= len(cards) {
        return nil
    }
    return NewStructureCardViewModel(vm.gameState, cards[idx])
}

func (vm *TerritoryViewModel) Yield() core.ResourceQuantity {
    // テリトリーの基本収益と建設カードの効果を計算
    baseYield := vm.territory.Terrain().BaseYield()
    
    var cards []*core.StructureCard
    if vm.plan != nil {
        cards = vm.plan.Cards()
    } else {
        cards = vm.territory.Cards()
    }
    
    totalYield := baseYield
    for _, card := range cards {
        // カードの収益効果を適用
        totalYield = totalYield.Add(card.YieldAdditiveValue())
        totalYield = card.YieldModifier().Apply(totalYield)
    }
    
    return totalYield
}

func (vm *TerritoryViewModel) SupportPower() float64 {
    var cards []*core.StructureCard
    if vm.plan != nil {
        cards = vm.plan.Cards()
    } else {
        cards = vm.territory.Cards()
    }
    
    totalSupportPower := 0.0
    for _, card := range cards {
        totalSupportPower += card.SupportPower()
    }
    
    return totalSupportPower
}

func (vm *TerritoryViewModel) SupportCardSlot() int {
    var cards []*core.StructureCard
    if vm.plan != nil {
        cards = vm.plan.Cards()
    } else {
        cards = vm.territory.Cards()
    }
    
    totalSupportCardSlot := 0
    for _, card := range cards {
        totalSupportCardSlot += card.SupportCardSlot()
    }
    
    return totalSupportCardSlot
}

// ResourceViewModel - リソース表示用
type ResourceViewModel struct {
    gameState *core.GameState
}

func NewResourceViewModel(gameState *core.GameState) *ResourceViewModel {
    return &ResourceViewModel{gameState: gameState}
}

func (vm *ResourceViewModel) Quantity() core.ResourceQuantity {
    return vm.gameState.Treasury.GetResourceQuantity()
}

// CalendarViewModel - カレンダー表示用
type CalendarViewModel struct {
    gameState *core.GameState
}

func NewCalendarViewModel(gameState *core.GameState) *CalendarViewModel {
    return &CalendarViewModel{gameState: gameState}
}

func (vm *CalendarViewModel) YearMonth() string {
    // langパッケージを使用してローカライズされた年月表示を取得
    return lang.FormatYearMonth(vm.gameState.Calendar.Year(), vm.gameState.Calendar.Month())
}

// HistoryViewModel - 履歴表示用
type HistoryViewModel struct {
    gameState *core.GameState
}

func NewHistoryViewModel(gameState *core.GameState) *HistoryViewModel {
    return &HistoryViewModel{gameState: gameState}
}

func (vm *HistoryViewModel) Title() string {
    // langパッケージを使用してローカライズされた履歴タイトルを取得
    return lang.GetText("history_title")
}

func (vm *HistoryViewModel) NumEvents() int {
    return len(vm.gameState.History.Events())
}

func (vm *HistoryViewModel) Event(idx int) *HistoryEventViewModel {
    events := vm.gameState.History.Events()
    if idx < 0 || idx >= len(events) {
        return nil
    }
    return NewHistoryEventViewModel(vm.gameState, events[idx])
}

// HistoryEventViewModel - 履歴イベント表示用
type HistoryEventViewModel struct {
    gameState *core.GameState
    event     *core.HistoryEvent
}

func NewHistoryEventViewModel(gameState *core.GameState, event *core.HistoryEvent) *HistoryEventViewModel {
    return &HistoryEventViewModel{
        gameState: gameState,
        event:     event,
    }
}

func (vm *HistoryEventViewModel) YearMonth() string {
    // langパッケージを使用してローカライズされた年月表示を取得
    return lang.FormatYearMonth(vm.event.Year(), vm.event.Month())
}

func (vm *HistoryEventViewModel) Text() string {
    // langパッケージを使用してローカライズされたイベントテキストを取得
    return lang.GetHistoryEventText(vm.event.ID(), vm.event.Parameters())
}
```

### Flowパッケージ

```go
// BattleFlow - バトル操作用
type BattleFlow struct {
    gameState *GameState
    battlefield *Battlefield
}

func (bf *BattleFlow) RemoveFromBattle(cardIndex int) {
    card, ok := bf.battlefield.RemoveBattleCard(cardIndex)
    if ok {
        // カードをデッキに戻す
        bf.gameState.CardDeck.Add(card.ID())
    }
}

func (bf *BattleFlow) Conquer() bool {
    if bf.battlefield.CanBeat() {
        // 征服処理
        return true
    }
    return false
}

func (bf *BattleFlow) Rollback() {
    // バトル状態をリセット
}

// TerritoryFlow - テリトリー操作用
type TerritoryFlow struct {
    territory *Territory
    plan      *ConstructionPlan
}

func (tf *TerritoryFlow) RemoveFromPlan(cardIndex int) {
    tf.plan.RemoveCard(cardIndex)
}

func (tf *TerritoryFlow) Commit() {
    tf.territory.ApplyConstructionPlan(tf.plan)
}

func (tf *TerritoryFlow) Rollback() {
    tf.plan = NewConstructionPlan(tf.territory)
}
```

## エラーハンドリング

### 設計原則
1. **防御的プログラミング**: 不正な入力に対する適切な処理
2. **早期リターン**: エラー条件の早期検出と処理
3. **型安全性**: コンパイル時エラーの最大化

### エラー処理パターン

```go
// ブール戻り値による成功/失敗の表現
func (cd *CardDeck) Remove(cardID CardID) bool {
    if cd.hand[cardID] <= 0 {
        return false
    }
    // 処理続行
    return true
}

// 値とブールの組み合わせによる安全な変換
func (p Point) AsBattlePoint() (BattlePoint, bool) {
    if bp, ok := p.(BattlePoint); ok {
        return bp, true
    }
    return nil, false
}
```

## テスト戦略

### 単体テスト
- **不変オブジェクト**: 状態変更がないことの確認
- **カプセル化**: 公開APIのみを使用したテスト
- **ビジネスロジック**: ドメインルールの検証

### 統合テスト
- **レイヤー間連携**: viewmodel ↔ core の連携
- **データ整合性**: メモリ共有の回避確認
- **ユースケース**: flow パッケージの操作検証

### テストデータ
- **モックオブジェクト**: 外部依存の分離
- **テストビルダー**: 複雑なオブジェクト構築の簡素化
- **テーブル駆動テスト**: 多様な入力パターンの検証

```go
func TestCardDeck_Add(t *testing.T) {
    tests := []struct {
        name     string
        initial  map[CardID]int
        cardID   CardID
        expected map[CardID]int
    }{
        {
            name:     "新しいカードを追加",
            initial:  map[CardID]int{},
            cardID:   "card1",
            expected: map[CardID]int{"card1": 1},
        },
        {
            name:     "既存カードの枚数増加",
            initial:  map[CardID]int{"card1": 1},
            cardID:   "card1", 
            expected: map[CardID]int{"card1": 2},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            deck := &CardDeck{hand: tt.initial}
            deck.Add(tt.cardID)
            
            if !reflect.DeepEqual(deck.hand, tt.expected) {
                t.Errorf("expected %v, got %v", tt.expected, deck.hand)
            }
        })
    }
}
```