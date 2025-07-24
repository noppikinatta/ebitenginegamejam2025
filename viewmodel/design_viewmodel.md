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
    CardViewModel(idx int) *BattleCardViewModel
}
```

## CardDeckViewModel

```
type CardDeckViewModel struct {
    NumBattleCards() int
    NumStructureCards() int
    BattleCardViewModel(idx int) *BattleCardViewModel
    StructureCardViewModel(idx int) *StructureCardViewModel
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
    PointViewModel(x,y int) *PointViewModel
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

## TerritoryViewModel