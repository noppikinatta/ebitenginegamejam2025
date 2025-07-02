# `@/ui` パッケージ構造解説

`@/ui` パッケージは、`@/core` パッケージのデータ構造とロジックをユーザーに視覚的に提示し、インタラクションを受け付ける役割を担います。Ebitengineライブラリを用いて描画されます。

## 1. UIの全体構成 (`gameui.go`)

- **`GameUI`**: UI全体のトップレベルコントローラー。
- `core.GameState` への参照を持ち、ゲームの状態に基づいてUI全体を管理します。
- `ResourceView`, `CalendarView`, `MainView`, `InfoView`, `CardDeckView` といった、画面を構成する主要なUIウィジェットをすべて内包し、それらの `Draw` と `Update` (入力処理) を統括します。
- ウィジェット間の連携（例：カードデッキで選択したカードを情報ビューに表示する）も担当します。

## 2. UIの基本要素

### `widget.go`
- **`Widget` インターフェース**: すべてのUIコンポーネントが実装すべき基本的なインターフェース（`HandleInput`, `Draw`）。
- **共通描画関数**:
    - `DrawResource`, `DrawCard`, `DrawBattleCard`, `DrawButton` など、ゲーム内のオブジェクトやUI部品を描画するためのヘルパー関数群。
    - `core` パッケージのデータ (`core.BattleCard`など) を受け取り、それを画像やテキストとして画面に描画する責務を持ちます。
    - `DrawCardDescriptionTooltip` のように、マウスオーバー時のツールチップ表示なども提供します。

### `input.go`
- **`Input`**: `nyuuryoku` という外部ライブラリをラップし、マウス入力を抽象化します。

## 3. 主要なUIコンポーネント

### `mainview.go`
- **`MainView`**: ゲーム画面中央の最も大きな領域を管理するコンテナウィジェット。
- `ViewType` という状態に応じて、以下の4つのビューを切り替えて表示します。
    - `MapGridView` (マップ)
    - `MarketView` (市場)
    - `BattleView` (戦闘)
    - `TerritoryView` (領地)
- マップ上の地点クリックなど、ユーザーのアクションに応じて適切なビューに切り替えるロジックの中心となります。

#### 3.1. `mapgrid.go`
- **`MapGridView`**: `core.MapGrid` の状態を視覚化します。
- マップ上の各地点 (`Point`) をアイコンで表示し、プレイヤーが到達可能な地点やそこへ至る経路を線で示します。
- 地点クリックを検知し、`MainView` に通知してビュー切り替えをトリガーします。

#### 3.2. `market.go`
- **`MarketView`**: `core.Nation` が持つ市場 (`Market`) の内容を表示します。
- 販売されているカードパックを一覧表示し、価格や購入条件（リソース、市場レベル）を提示します。
- ユーザーの購入操作を処理し、成功すると `core.GameState` を更新してターンを進めます。

#### 3.3. `battle.go`
- **`BattleView`**: `core.Battlefield` での戦闘シーンを描画します。
- 敵の情報、配置された `BattleCard`、それらの合計戦力などをリアルタイムに表示します。
- プレイヤーは `CardDeckView` からカードを選択してこのビューに配置し、戦力が敵を上回ると「征服」ボタンが有効になります。

#### 3.4. `territory.go`
- **`TerritoryView`**: 征服済みの領地 (`Territory`) の管理画面です。
- `StructureCard` を配置するスロットが表示され、プレイヤーは `CardDeckView` からカードを配置できます。
- カード配置による資源産出量の変化をプレビュー表示し、内容を確定させることができます。

### `carddeck.go`
- **`CardDeckView`**: 画面下部に常時表示され、プレイヤーの `CardDeck` の内容を一覧表示します。
- ユーザーがカードをクリックすると、現在の `MainView` の種類（戦闘中か領地管理中か）に応じて適切な処理（カードの配置など）を行うためのコールバックを呼び出します。

### `info.go`
- **`InfoView`**: 画面右側に常時表示される情報パネル。
- `InfoViewMode` の状態によって、表示内容が動的に切り替わります。
    - `InfoModeHistory`: ゲーム内のイベント履歴を表示。
    - `InfoModeCardInfo`: 選択されたカードの詳細情報を表示。
    - `InfoModeNationPoint` / `InfoModeWildernessPoint`: 選択されたマップ地点の詳細情報を表示。

### `resource.go`
- **`ResourceView`**: 画面左上にプレイヤーの現在のリソース保有量と、次のターンの収入予測を常時表示します。

### `calendar.go`
- **`CalendarView`**: 画面右上に現在のゲーム内年月を常時表示します。

## 4. UIコンポーネント間の連携

`GameUI` を頂点として、各UIコンポーネントはコールバック関数や `core.GameState` への参照を通じて密に連携しています。例えば、

1. `MapGridView` で敵のいる地点をクリックする。
2. `MainView` が `BattleView` に切り替える。
3. `CardDeckView` で `BattleCard` をクリックする。
4. `GameUI` を経由して `BattleView` にそのカードが渡され、戦場に配置される。
5. 戦力が更新され、`BattleView` に再描画される。
6. `CardDeckView` でカードにマウスカーソルを合わせる。
7. `GameUI` を経由して `InfoView` の表示内容がそのカードの詳細情報に切り替わる。

といった形で、コンポーネントが協調して動作します。
