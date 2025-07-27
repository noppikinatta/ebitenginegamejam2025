# UIパッケージ修正計画

## 概要
coreパッケージの再設計に伴い、uiパッケージに大量のコンパイルエラーが発生。この計画書では、エラーの分類と修正方針を定める。

## 主要な変更点と対応

### 1. Cards型の削除とCardID型の導入
**変更内容:**
- `core.Cards`型が削除され、`core.CardID`型が新しく導入された
- `CardDeck.Add()`の引数が`*core.Cards`から`core.CardID`に変更

**影響箇所:**
- `ui/battle.go:88,102` - `bv.GameState.CardDeck.Add()`の引数

**修正方針:**
- `&core.Cards{...}`を`core.CardID`の値に変更
- カード生成ロジックを見直し、適切なCardIDを使用

### 2. BattlePointインターフェースのメソッド名変更
**変更内容:**
- `GetEnemy()` → `Enemy()`に変更

**影響箇所:**
- `ui/battle.go:64,226,334,382,389` - `BattlePoint.GetEnemy()`の呼び出し

**修正方針:**
- 全ての`GetEnemy()`を`Enemy()`に変更

### 3. フィールドからメソッドへの変更
**変更内容:**
- `.Cards`が`.Cards()`メソッドに変更
- `.Controlled`が`.Controlled()`メソッドに変更
- `.TerrainType`が非公開フィールドになった

**影響箇所:**
- `ui/territory.go:38,39,83,126,245,259,293` - `Cards`フィールドの使用
- `ui/battle.go:406` - `Controlled`フィールドの使用
- `ui/battle.go:193` - `TerrainType`フィールドの使用

**修正方針:**
- フィールドアクセスをメソッド呼び出しに変更
- 非公開フィールドはアクセサメソッドを使用

### 4. Marketインターフェースの変更
**変更内容:**
- `Nation.GetMarket()`メソッドが削除
- `MarketItem.Price()`がメソッドになった
- `MarketItem.CardPack()`がメソッドになった

**影響箇所:**
- `ui/market.go:98` - `mv.Nation.GetMarket()`の呼び出し
- `ui/market.go:194,202-206` - `price`の使用方法
- `ui/market.go:157` - `item.CardPack.CardPackID`の使用

**修正方針:**
- Marketアクセスの方法を見直し
- `price()`メソッドの戻り値を適切に処理
- `CardPack()`メソッドの戻り値からIDを取得

### 5. 未定義関数エラー
**変更内容:**
- `Input`, `DrawButton`, `DrawBattleCard`, `DrawCardBackground`, `DrawCard`等の関数が見つからない

**影響箇所:**
- 各uiファイルの描画関連処理

**修正方針:**
- 関数の定義場所を確認し、適切なimportを追加
- または代替実装を検討

### 6. インターフェース実装の不備
**変更内容:**
- `BattlePoint`が`Point`インターフェースを実装していない

**影響箇所:**
- `ui/battle.go:385` - `BattlePoint`を`Point`として使用

**修正方針:**
- 適切な型キャストまたはメソッド呼び出しに変更

## 修正作業順序

### Phase 1: 基本的な型とメソッド名の修正
1. `GetEnemy()` → `Enemy()`の一括変更
2. Cards型からCardID型への変更
3. フィールドアクセスからメソッド呼び出しへの変更

### Phase 2: Market関連の修正
1. Market取得方法の見直し
2. Price情報の取得方法修正
3. CardPack関連の修正

### Phase 3: 未定義関数の解決
1. 描画関数の定義場所確認
2. 必要に応じてimport文の修正
3. Input関数の解決

### Phase 4: インターフェース問題の解決
1. Point/BattlePoint関連の型変換修正
2. Territory関連の型変換修正

### Phase 5: 総合テスト
1. 全ファイルのコンパイル確認
2. 個別ファイルの動作確認

## 注意事項
- 修正時は一つずつ段階的に実行し、各段階でコンパイル確認を行う
- coreパッケージのAPIを変更せず、uiパッケージ側で対応する
- 型安全性を保ちながら修正を行う
- 元の機能を損なわないよう注意する 