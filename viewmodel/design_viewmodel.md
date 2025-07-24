# viewmodel package

ui packageのUI部品が情報を参照するための型。
Go言語のinterfacedの文法で型定義を書いているが、実際にはstructで設計する。

## BattleViewModel

```
type BattleViewModel struct {
    Title() string
    EnemyImage() *ebiten.Image
    EnemyType() string
    EnemyPower() float64
    EnemyTalk() string
    EnemySkillNames() []string
    EnemySkillDescriptions() []string
    CardSlot() int
    CanBeat() bool
    TotalPower() float64
    NumCards() int
    Card(idx int) *BattleCardViewModel
}
```

## CardDeckViewModel

```
type CardDeckViewModel struct {
    NumBattleCards() int
    NumStructureCards() int
    BattleCard(idx int) *BattleCardViewModel
    StructureCard(idx int) *StructureCardViewModel
}
```

## CardViewModel

```
type BattleCardViewModel struct {
    Image() *ebiten.Image
    Name() string
    Duplicates() int
    CardTypeImage() *ebiten.Image
    CardTypeName() string
    Power() float64
    SkillName() string
    SkillDescription() string
}
```

```
type StructureCardViewModel struct {
    Image() *ebiten.Image
    Name() string
    Duplicates() int
}
```

## MapGridViewModel

```
type MapGridViewModel struct {
    Size() core.MapGridSize
    Point(x,y int) *PointViewModel
    ShouldDrawLineToRight(x,y int) bool
    ShouldDrawLineToUpper(x,y int) bool
}
```

```
type PointViewModel struct {
    Image() *ebiten.Image
    Name() string
    HasEnemy() bool
    EnemyPower() float64
}
```

## MarketViewModel

```
type MarketViewModel struct {
    Title() string
    Level() float64
    NumItems() int
    Item(idx int) *MarketItemViewModel
}
```

```
type MarketItemViewModel struct {
    ItemName() string
    RequiredLevel() int
    Unlocked() bool
    CanPurchase() bool
    Price() core.ResourceQuantity
}
```

```
type ResouceSufficiency struct {
	Money bool
	Food  bool
	Wood  bool
	Iron  bool
	Mana  bool
}
```

## TerritoryViewModel

```
type TerriroryViewModel struct {
    Title() string
    CardSlot() int
    NumCards() int
    Card(idx int) *StructureCardViewModel
    Yield() ResourceQuantity
    SupportPower() float64
    SupportCardSlot() int
}
```

## ResourceViewModel

```
type ResourceViewModel struct {
    Quantity() core.ResourceQuantity
}
```

## CalendarViewModel

```
type CalendarViewModel struct {
    YearMonth() string
}
```

## HistoryViewModel

```
type HistoryViewModel struct {
    Title() string
    NumEvents() int
    Events(idx int) *HistoryEventViewModel
}
```

```
type HistoryEventViewModel struct {
    YearMonth() string
    Text() string
}
```