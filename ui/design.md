# uiの設計

## 基本設計

Widgetインタフェースを満たす構造体。

## UIレイアウト

画面の幅と高さは640x360

(Left, Top, Width, Height)の形式で記載する。

- ResourceView (0,0,520,20)
- CalendarView (520,0,120,20)
- MainView (0,20,520,280)
  - MapGridView,MarketView,BattleView,TerritoryViewも同じ
- InfoView (520,20,120,280)
  - HistoryView,CardInfoView,NationPointView,WildernessPointView,DescriptionViewも同じ
- CardDeckView (0,300,640,60)

各Widgetは配置のための座標を描画時にマージンとして保持しておく。

## Widget詳細

### MainView

MapGridView,MarketView,BattleView,TerritoryViewを切り替えるコンテナWidget。
最初に表示するのはMapGridView。

#### MapGridView

core.MapGridの情報を描画する。

> MapGridのPointの配置は5x5で固定。

520x280の領域を5で割り、Pointごとの104x56ごとのセルに分ける。
24x24サイズのPoint画像(まずはダミーで良い)をセル中央に描画する。
Point画像の下にPoint名(ダミーとしてPoint0,0~Point4,4でよい)を文字描画する。

> 文字描画のために、drawing.DrawTextを使用できる。fontSizeは12にしてみる。

MyNationから到達可能なPointまで線を引いていく。線は隣接するPoint同士を繋ぐように引く。

#### MarketView

選択しているMyNationまたはOtherNationのMarketを表示する。

(0,20,520,40)の領域をヘッダとして、「My Nation」またはOtherNationの名前(今はNationIDでいい)を表示する。
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

#### TerritoryView

### ResourceView

### CalendarView

### InfoView

#### HistoryView

#### CardInfoView

#### NationPointView

#### WildernessPointView

#### DescriptionView

### CardDeckView