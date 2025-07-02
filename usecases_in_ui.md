### 1. `ui/market.go` より: カードパック購入処理

`MarketView.PurchaseCardPack` メソッドは、カードパックの購入という単一のユーザー操作のために、複数のドメイン層のメソッドを順番に呼び出しています。

- **該当コード (`ui/market.go`)**:
  ```go
  // ui/market.go
  func (mv *MarketView) PurchaseCardPack(item *core.MarketItem) error {
      // ... (事前チェック) ...

      // 1. 購入処理の実行
      cardPack, ok := mv.Nation.Purchase(itemIndex, mv.GameState.Treasury)
      if !ok { /* ... */ }

      // 2. カードパックの開封
      rng := newSimpleRand()
      cardIDs := cardPack.Open(rng)

      // 3. カードインスタンスの生成
      cards, ok := mv.GameState.CardGenerator.Generate(cardIDs)
      if !ok { /* ... */ }

      // 4. 生成されたカードをデッキに追加
      mv.GameState.CardDeck.Add(cards)

      // 5. ターンを進行
      mv.GameState.NextTurn()

      return nil
  }
  ```

- **問題点**:
  この「購入 → 開封 → 生成 → デッキに追加 → ターン進行」という一連の処理フローは、UIではなくドメイン層が責務を持つべき「カードパックを購入する」というユースケースそのものです。UI層がこの手順を組み立てるのではなく、`core.GameState` に `PurchaseCardPack(...)` のような単一のメソッドを定義し、UIからはそれを呼び出すだけ、というのが理想的な形です。

### 2. `ui/battle.go` より: 戦闘開始時の初期化処理

`BattleView.createBattlefield` メソッドは、戦闘を開始するにあたり、隣接する領地の建造物カードの効果を調べて `Battlefield` を初期化するという、複雑なロジックを実装しています。

- **該当コード (`ui/battle.go`)**:
  ```go
  // ui/battle.go
  func (bv *BattleView) createBattlefield(point core.BattlePoint) *core.Battlefield {
      // ...
      battlefield := core.NewBattlefield(enemy, 0.0)

      // 周囲4方向のPointを調査
      mapGrid := bv.GameState.MapGrid
      // ... (for loop for directions) ...
          p := mapGrid.GetPoint(checkX, checkY)

          if wildernessPoint, ok := p.(*core.WildernessPoint); ok {
              if wildernessPoint.Controlled && wildernessPoint.Territory != nil {
                  territory := wildernessPoint.Territory

                  // 領地の建造物カードの効果を適用
                  for _, card := range territory.Cards {
                      if card.BattlefieldModifier != nil {
                          card.BattlefieldModifier.Modify(battlefield)
                      }
                  }
              }
          }
      // ...
      return battlefield
  }
  ```

- **問題点**:
  「ある地点での戦闘開始時に、周囲の状況からどのような支援効果が発生するか」というルールは、まさしくドメインロジックです。この処理はUIから分離し、`core`パッケージ内に `core.NewBattlefieldForPoint(point, mapGrid)` のようなファクトリ関数として実装されるべきです。

### 3. `ui/battle.go` より: 敵拠点征服時の処理

`BattleView.Conquer` メソッドは、戦闘に勝利した際に行われるべき一連のドメイン状態の更新を直接実行しています。

- **該当コード (`ui/battle.go`)**:
  ```go
  // ui/battle.go
  func (bv *BattleView) Conquer() bool {
      // ... (事前チェック) ...

      // 1. 戦闘勝利処理
      bv.Battlefield.Beat()

      // 2. 使用したカードをデッキに戻す
      bv.GameState.CardDeck.Add(&core.Cards{BattleCards: bv.Battlefield.BattleCards})

      // 3. 地点の支配状態を変更
      bv.BattlePoint.SetControlled(true)
      bv.GameState.MapGrid.UpdateAccesibles()

      // 4. プレイヤー国家の市場レベルを上昇
      bv.GameState.MyNation.AppendLevel(0.5)

      // 5. ターンを進行
      bv.GameState.NextTurn()

      return true
  }
  ```

- **問題点**:
  これも「戦闘に勝利する」というユースケースです。カードの返却、領地状態の更新、報酬の付与（市場レベル上昇）、ターン進行といった一連のドメイン操作は、UIが個別に呼び出すのではなく、`core.GameState` が `WinBattle(...)` のような単一のメソッドで責務を持つべきです。特に、「拠点を征服すると市場レベルが0.5上がる」というルールは、UIではなくドメインが知っているべき情報です。

### 4. `ui/territory.go` より: 建造物カードのデッキへの返却

`TerritoryView.RemoveCard` メソッドは、領地に配置したカードをデッキに戻す際、`core.Territory` の内部データである `Cards` スライスを直接操作しています。

- **該当コード (`ui/territory.go`)**:
  ```go
  // ui/territory.go
  func (tv *TerritoryView) RemoveCard(card *core.StructureCard) bool {
      // ... (インデックス検索) ...

      // core.Territoryが持つスライスをUI層が直接変更している
      tv.Territory.Cards = append(tv.Territory.Cards[:cardIndex], tv.Territory.Cards[cardIndex+1:]...)

      // デッキにカードを追加
      if tv.GameState != nil {
          cards := &core.Cards{StructureCards: []*core.StructureCard{card}}
          tv.GameState.CardDeck.Add(cards)
      }

      return true
  }
  ```

- **問題点**:
  これはカプセル化の破壊にあたります。`core.Territory` は自身の `Cards` スライスを管理する `RemoveCard(index int)` メソッドをすでに持っています。UI層はそれを呼び出すべきであり、内部のスライス構造を直接変更してはいけません。また、「領地からカードを取り除き、デッキに戻す」という一連の操作も、ドメイン層のユースケースとしてカプセル化することが望ましいです。
