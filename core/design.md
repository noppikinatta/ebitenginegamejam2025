# coreパッケージ: ドメインレイヤーのルール表現

coreパッケージには、ゲームのドメインルールを実装する。ゲームのルールは以下の通り。

## MapGridへのインタラクション

MagGrid上の操作可能なPointを操作して、ターンを経過させる。

> 操作はキャンセルできるかもしれない。その場合ターンは経過しない。

### MyNationPointの操作

MyNationPointでは、MarketにあるCardPackを1つ購入できる。

### OtherNationPointの操作

OtherNationPointでは、MarketにあるCardPackを1つ購入できる。

### WildernessPointの操作

WildernessPointでは、Enemyと戦闘して勝つことで、制圧(Controled=true)になる。

Controledな場合、TerritoryからYieldとしてResourceをTreasuryに加算する。
TerritoryにはStructureCardを置くことができる。StructureCardの配置を変更し、変更をコミットした場合もターンが経過する。

### BossPointの操作

BossPointでは、Enemyと戦闘できる。このEnemyを倒すのがゲーム全体の目的。

## ターン経過時の出来事

ターン経過時に、MyNation.BasicYieldと、ControledなWildernessPointのTerritoryの計算済みのYieldを加算して、Treasuryに加算する。

## 戦闘ルール

CardDeckにあるBattleCardのうち一部をBattlefieldに出して、Powerを計算し、EnemyのPower以上であれば勝利。制圧できる。
EnemyのPowerは固定。
BattleCardのPowerは加算する。
Enemyごとに、プレイヤーが戦闘で出せるBattleCardの枚数が決まっている。
BattleCardSkillやEnemySkillは、BattleCardのPowerの計算に影響する。
WildernessPointに隣接する制圧済みWildernessPointに配置したStructureCardは、隣接するWildernessPointの戦闘に影響する。