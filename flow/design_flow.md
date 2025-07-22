# flow package

uiが直接操作するusecaseレイヤー。
Go言語のinterfacedの文法で型定義を書いているが、実際にはstructで設計する。

## BattleFlow

```
type BattleFlow struct {
    RemoveFromBattle(cardIndex int)
    Conquer() bool
    Rollback()
}
```

## CardDeckFlow

```
type CardDeckFlow struct {
    Select(cardIndex int)
}
```

## MapGridFlow

```
type MapGridFlow struct {
    SelectPoint(x,y int)
}
```

## MarketFlow

```
type MarketFlow struct {
    Purchase(marketItemIdx int) bool
}
```

## TerritoryFlow

```
type TerritoryFlow struct {
    RemoveFromPlan(cardIndex int)
    Commit()
    Rollback()
}
```