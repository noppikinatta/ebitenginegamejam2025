# リファクタリング計画(coreパッケージ)

ゲームルールとして、同じカードを複数持っても良いことにする。
現在はBattleCardが1種類1枚になっている。

## 全体的に

FieldがPublicになっているものが多い。カプセル化の試みとして、可能な限りFieldをpackage private (小文字始まり) にしていく。

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
	skillID  BattleCardSkillID
	type     BattleCardType
}
```

BattleCardはimmutableになった。IDなどのフィールドの情報を取得するメソッドを追加すべき。
BattleCardSkillのポインタを持つのをやめ、BattleCardSkillIDを持たせる。これは将来的に外部ファイルから読み込みを行うとき、各オブジェクトを独立させておいた方が処理しやすいため。

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

そのまま。

## BattleCardSkillCalculatorBuff

新設。BattleCardSkillCalculatorインタフェースを実装する具体型を元に、範囲/条件/効果の3つを指定する構造体を定義する。

```
type BattleCardSkillCalculatorBuff struct {
	target    BattleCardSkillCalculatorBuffTarget
	condition BattleCardSkillCalculatorBuffCondition
	effect    BattleCardSkillCalculatorBuffEffect
}
```

範囲/条件/効果の3つはインタフェース型。

## BattleCardSkillCalculatorBuffTarget

BattleCardSkillCalculatorBuffの対象。

```
type BattleCardSkillCalculatorBuffTarget interface {
	IsTarget(idx int, options *BattleCardSkillCalculationOptions) bool
}
```

BattleCardSkillCalculationOptions.BattleCardsをループするときに、idxを渡して対象かどうかを判定する。

実装は次の２種類。

```
type BattleCardSkillCalculatorBuffTargetIndices struct {
	idxs []int
}
```

BattleCardSkillCalculationOptions.BattleCardIndexとidxsを加算し、どれかに当てはまればtrueを返す。

例
* BattleCardSkillCalculatorBuffTargetIndices.idxs == []int{0} の場合は自分自身。
* BattleCardSkillCalculatorBuffTargetIndices.idxs == []int{-1,1} の場合は自分の前後のカード。


```
type BattleCardSkillCalculatorBuffTargetAll struct {
}
```

常にtrueを返す。

## BattleCardSkillCalculatorBuffCondition

Skillの効果対象を決めるinterface。

```
type BattleCardSkillCalculatorBuffCondition interface {
    Apply(opt *BattleCardSkillCalculationOptions, effect BattleCardSkillEffect)
}
```

## BattleCardSkillCalculatorBuffEffect

```
type BattleCardSkillEffect interface {
    Apply(opt *BattleCardSkillCalculationOptions, effect BattleCardSkillEffect)
}
```
