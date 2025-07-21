# リファクタリング計画(coreパッケージ)

ゲームルールとして、同じカードを複数持っても良いことにする。
現在はBattleCardが1種類1枚になっている。

## 全体的に

FieldがPublicになっているものが多い。カプセル化の試みとして、可能な限りFieldをpackage private (小文字始まり) にしていく。
他のパッケージから設定しているものは、コンストラクタ(New+型名の関数)を追加する。おそらくloadパッケージのみ使うだろう。一部は新設するflowパッケージ(Usecaseレイヤー)で使う。
uiパッケージから参照しているものは、新設するviewmodelパッケージ(Presenterレイヤー)にプロキシ型を通じて取得する予定。

## CardDeck

CardDeckの責務はどのカードを何枚持っているかの管理だけ。なので、BattleCardやStructrueCardへのポインタは必要なく、CardIDだけでよい。

```
type CardDeck struct {
	hand map[CardID]int
}
```

のようになるはず。APIを以下の通りにする。

* Add(CardID) カードを1枚追加
* Remove(CardID) bool カードを1毎削除(削除できない場合falseを返す)
* Count(CardID) int カードの枚数を返す

> AddCards([]CardID) のようなAPIが必要になるかもしれない。

## BattleCard

Experienceの概念をやめる。元々これはBattleCardの枚数が多すぎて表示しきれない問題を締め切り当日に解決しようとしたもの。
BattleCardを複数枚購入したとき、Experienceが増えるのではなくBattleCardの枚数が増えてよい。

```
type BattleCard struct {
	id       CardID
	power    float64  
	skill    *BattleCardSkill
	type     BattleCardType
}
```

BattleCardはimmutableになった。IDなどのフィールドの情報を取得するメソッドを追加すべき。

## BattleCardSkill

DescriptionKeyは表示にしか関係ない。coreの責務ではないので削除する。
BattleCardSkillCalculatorはBattleCardSkillTargetとBattleCardSkillEffectに分ける。

```
type BattleCardSkill struct {
	id         BattleCardSkillID
	calculator BattleCardSkillCalculator
}
```

## BattleCardSkillCalculator

そのまま。具象型群もそのままでいい。FieldだけPackage Privateを目指す。
そのため、具象型のコンストラクタ関数が大量にできるかもしれない。

## StructureCard

YieldModifierやBattlefieldModifierで抽象化していたが、過度の抽象化なのでやめる。
DescriptionKeyは表示にしか関係ない。coreの責務ではないので削除する。

```
type StructureCard struct {
	cardID              CardID
	yieldAdditiveValue  ResourceQuantity
	yieldModifier       ResourceModifier
	supportPower        float64
	supportCardSlot     int
}
```

## Enemy

Questionは削除する。表示にしか関係ないため。

## EnemySkill, EnemySkillImpl

EnemySkillをinterfaceにしたのは過度の抽象化だった。
EnemySkillImplをEnemySkillとして、structで再定義すべき。

## Battlefield

CalculateTotalPowerとは別に、個々のBattleCardの計算されたPowerを取得できるようにする。表示に使うため。
同じように、計算されたSupportPowerも表示のために必要になる。

## Nation

VisibleMarketItemsは意味がなかった。なぜなら、利用できないMarketItemが見えないよりも、見えるようにしてアンロック条件を表示した方がプレイヤーにとっては有益なためである。

結局各所ではGetMarketを使っている。

## OtherNation

購入時にMarket.Levelを0.5上げているが、これはやめたい。Marketに購入時の変化を持たせるように変更する。

## Market, MarketItem

MarketItemにLevelEffectの項目を設け、購入時にMarket.Levelに加算する。

MarketItem.CardPackはnilを許容するかもしれない。それは、LevelEffectを多く持たせ、投資のような性質のMarketItemを作るためである。

MarketItem.ResourceQuantityを追加するかもしれない。これは、資源の取引を表す。

## Treasury

国庫を表す概念なのでMyNationの一部のような気がする。

## Territory

一部の情報を新設するTerrainにうつす。

```
type Territory struct {
	id      TerritoryID
	terrain *Terrain
	cards   []*StructureCard
}
```

## ConstructionPlan

ui.TerritoryView表示時に作る。またはcardsをクリアする。
ui.BattleView表示時に作られるBattlefieldのような役割。
Territoryのcardsの変更計画であって、ui.TerritoryViewで計画をコミットするとき、cardsの内容をTerriroryに複写するようにしたい。

**TerritoryとConstructionPlanでcardsのSliceのメモリを共有してしまわないように注意する。**

```
type ConstructionPlan struct {
	cards   []*StructureCard
}
```

## Terrain

地形の情報。immutable。

```
type Terrain struct {
	id        TerrainID
	baseYield ResourceQuantity
	cardSlot  int
}
```

## MapGrid



## Point

Pointの具体型はMyNationPoint,OtherNationPoint,WildernessPoint,BossPointの4種。
ui packageでPointをこの具体型にキャストして処理を分岐しているが、これはあまり良くない。

この処理を分岐している原因は、Pointの種類によって、MapGridView画面から開く画面の種類や必要なデータが異なるからである。

Pointは、IsMyNationの代わりにPointTypeを持ち、次のように定義するPointType型の値を返す。

```
type PointType int
const (
	PointTypeUnknown PointType = iota
	PointTypeMyNation
	PointTypeOtherNation
	PointTypeWilderness
	PointTypeBoss
)
```

Pointは、MapGridViewから開く子画面(BattleView, TerritotyView, MarketView)に対応する情報を返せる必要がある。

```
type Point interface {
	PointType() PointType
	Passible() bool
	AsBattlePoint() BattlePoint

}
```

```
type BattlePoint interface {
	Point
	Enemy() *Enemy
	Conquer()
}
```

```
type TerriroryPoint interface {
	Yield() ResourceQuantity
	Terrain() *Terrain
	CardSlot() int
	Cards() []*StructureCard
}
```

```
type MarketPoint interface {
	Nation() Nation
}
```