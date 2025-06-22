# uiの設計

## 基本設計

Widgetインタフェースを満たす構造体。

## UIレイアウト

画面の幅と高さは640x360

(Left, Top, Width, Height)の形式で記載する。

- ResourceView (0,0,300,20)
- CalendarView (520,0,120,20)
- MainView (0,20,520,280)
  - MapGridView,MarketView,BattleView,TerritoryViewも同じ
- InfoView (520,20,120,280)
  - HistoryView,CardInfoView,NationPointView,WildernessPointView,DescriptionViewも同じ
- CardDeckView (0,300,640,60)

各Widgetは配置のための座標を描画時にマージンとして保持しておく。

## Widget詳細

### MainView

mainview.goに実装する。

MapGridView,MarketView,BattleView,TerritoryViewを切り替えるコンテナWidget。
最初に表示するのはMapGridView。

#### MapGridView

mapgrid.goに実装する。

core.MapGridの情報を描画する。

> MapGridのPointは25個で、配置は5x5で固定。

520x280の領域を5で割り、Pointごとの104x56ごとのセルに分ける。
24x24サイズのPoint画像(まずはダミーで良い)をセル中央に描画する。
Point画像の下にPoint名(ダミーとしてPoint0,0~Point4,4でよい)を文字描画する。

> 文字描画のために、drawing.DrawTextを使用できる。fontSizeは12にしてみる。

MyNation(スタート地点)から到達可能なPointまで線を引いていく。線は隣接するPoint同士を繋ぐように引く。

到達可能なPointをクリックすることで、別のViewに遷移する。到達できない場合は遷移しない。
- MyNationPoint: MyNationを描画するMarketView
- OtherNation: OtherNationを描画するMarketView
- WildernessPoint: 未制圧の場合BattleView, 制圧済みの場合TerritoryView
- BossPoint: BattleView

#### MarketView

market.goに実装する。

選択しているMyNationまたはOtherNationのMarketを表示する。

(0,20,480,40)の領域をヘッダとして、「My Nation」またはOtherNationの名前(今はNationIDでいい)を表示する。
(480,20,40,40)の領域に[x]ボタンを作り、クリックでMapGridViewに戻るようにする。
(0,60,260,80)~(260,220,260,80)の6つの領域を、可視状態のCardPack描画に使う。

> ゲームルールとしてCardPackはNationごとに最大6つにする。7つ以上は考慮しなくていい。

CardPack描画領域260x80の領域を分割し、CardPack画像、CardPack名、説明、値段を描画する。

例：(0,60,260,80)の領域の場合
- (0,60,40,40) CardPack画像(24x32)
- (40,60,220,20) CardPack名
- (40,80,220,40) CardPack説明
- (0,120,260,20) CardPackの値段
  - 値段はResourceTypeの画像20x20と価格の数字40x20の組み合わせで、Resource1種類につき60x20の大きさ。
  - Resourceは5種類あるが、ルールとして4種類までしか要求しないことではみ出さないようにする。
  - TreasuryのResourceが足りてない場合、数字を赤文字にする。

> 描画する文字列は適宜ダミーテキストで良い。

CardPackの上をマウスでクリックしたとき、そのCardPackを購入する。

#### BattleView

battle.goに実装する。

Enemyとの戦闘を表現する。core.Battlefieldの情報を表示することになる。

(0,20,520,40)の領域をヘッダとして、「Battle of $POINT_NAME」(今はダミーテキストでいい)を表示する。
(480,20,40,40)の領域に[x]ボタンを作り、クリックでMapGridViewに戻るようにする。

(180,60,160,160)の領域に、敵の画像(今は四角でいい)を描画する。敵を倒せる場合、敵の画像クリックで戦闘に勝利する。

(0,220,480,60)の領域をBattleCard置き場にする。カード1枚あたりのスペースは40x60で、横に12枚置ける。

CardDeckView上でクリックしたBattleCardをBattleViewのBattleCard置き場に置く。
BattleBard置き場でクリックしたカードは、CardDeckViewに戻る。

(480,220,40,60)の領域に、計算されたカード全体のPowerを表示する。


#### TerritoryView

teritorry.goに実装する。

(0,20,520,40)の領域をヘッダとして、「$POINT_NAME」(今はダミーテキストでいい)を表示する。
(480,20,40,40)の領域に[x]ボタンを作り、クリックでMapGridViewに戻るようにする。

(0,60,60,100)の領域に、Territoryの産出量を表示する。60x20の領域5つに分割して、資源ごとの産出量を表示する。

(0,160,520,60)の領域をStructureCard置き場とする。

CardDeckView上でクリックしたBattleCardをTerritoryViewのStructureCard置き場に置く。
StructureCard置き場でクリックしたカードは、CardDeckViewに戻る。

(60,60,460,100)の領域に、置いたStructureCardの効果説明を書く。領域分割は未定。まずはダミーテキストを置く。

### ResourceView

resource.goに実装する。

60x20の領域5つに分け、5種類の資源の在庫と増分を表示する。

ResourceTypeの画像20x20と在庫と増分の数字40x20の組み合わせ。

在庫と増分は、10(+2) のような表示にする。
在庫と増分の文字の色は変えた方がいいかもしれない。

### CalendarView

calendar.goに実装する。

現在のTurnを監視して、その値を年月表示する。Tuen.YearMonthの戻り値の年と月を、「YYYY/MM」のように表示する。

### InfoView

info.goに実装する。

状況によって表示する情報の内容を変える。

#### HistoryView

デフォルトの表示。WildwrnessPointを制圧したり、OtherNationPointに到達できたときのTurnを記録して、一覧表示する。

(520,20,120,280)の領域を、120x20の文字フィールド14行分として表示に使う。

#### CardInfoView

BattleCardやStructureCardにマウスを乗せている場合。カードの詳細な説明を表示する。

高さ280の領域を上から以下のように分割して情報を表示する。

- 20: カード名
- 60: イラスト
- 20: カードの種類(Battle or Structure)
- 180: 説明

説明は、BattleCardかStructureCardかでさらに分ける。

BattleCardの場合
- 20: カードタイプ
- 20: Power
- 20: Skillの名前
- 40: Skillの説明

StructureCardの場合
- 20x9: YieldModifierまたはBattlefieldEffectの説明

> 実際には9個も特殊効果のあるカードはない。

#### NationPointView

MapGridViewでNationPointにマウスが乗っている時に表示する。

高さ280の領域を上から以下のように分割して情報を表示する。

- 20: Point名
- 20: `Card Packs` 固定
- 20x12: CardPack名

#### WildernessPointView

MapGridでWildernessPointにマウスが乗っている時に表示する。

高さ280の領域を上から以下のように分割して情報を表示する。

- 20: Point名
- 20: `Enemy` 固定
- 40: Enemyの画像と名前 (制圧済みの場合バツを付ける)
- 20: EnemyのPower
- 20: `Yields` 固定
- 20x3: Resource種類ごとのYield(1種類につき幅60を使い、1行に2種類表示する。)
- 20: `Structure Cards` 固定
- 20x4: 置かれているStructureCardの名前(5枚以上の場合どうするかは未定。今は4枚まで表示でいい。)

> BossPointでは、Yields以下は表示しない。

#### EnemySkillView

BattleViewでBattleCardにマウスが乗っていないときに表示する。

高さ280の領域のうち上から240を60x4に分けて、20をEnemyのSkillの名前、40をEnemyのSkillの説明に使う。

- 60x4: Enemyのスキルの名前と説明
  - 20: Enemyのスキルの名前
  - 40: Enemyのスキルの説明

#### 情報の切り替え方法

InfoViewの情報切り替えの判断方法。

- CardDeckViewのカードの上にマウスカーソルがある -> CardInfoView

- MainViewで表示しているものが
  - MapGridView
    - MyNationPoint or OtherNationPointの上にマウスカーソルがある -> NationPointView
    - WildernessPoint or BossPointの上にマウスカーソルがある -> WildernessPointView
  - BattleView
    - BattleViewに出したBattleCardの上にマウスカーソルがある -> CardInfoView
    - それ以外 -> EnemySkillView
  - TerritoryView
    - TerritoryViewに出したStructureCardの上にマウスカーソルがある -> CardInfoView

上記のいずれにも当てはまらない場合はHistoryView

### CardDeckView

CardDeckの内容を表示する。

640x60の領域を横に分割する。カード1枚あたりのスペースは40x60で、横に16枚置ける。

表示するカードは以下の条件の場合のみフィルタする。
- MainViewにBattleViewが表示されているとき、BattleCardのみ表示。
- MainViewにTerritoryViewが表示されているとき、StructureCardのみ表示。

MainViewにBattleViewが表示されているとき、カードをクリックするとBattlefieldにBattleCardを追加する。
MainViewにTerritoryViewが表示されているとき、カードをクリックするとTerritoryにStructureCardを追加する。

> カードを破棄するUIを追加するかもしれない。

## 共通の描画方式

### 資源量、資源産出量

60x20の領域を使い、左側の20x20の領域に資源アイコン、右側の40x20に産出量や在庫、増減を表示する。

### カード

1枚当たり40x60の領域を使う。

> 今は四角でいい。

