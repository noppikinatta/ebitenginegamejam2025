# 実装計画

- [ ] 1. CardDeckの再設計と実装
  - CardDeckをmap[CardID]intベースに変更し、Add、Remove、Countメソッドを実装
  - 既存のExperience概念を削除し、カード枚数管理に簡素化
  - _要件: 1.1, 1.2, 3.1, 3.2, 3.3, 3.4_

- [ ] 2. BattleCardの不変化
  - BattleCardのフィールドをpackage privateに変更
  - Experienceフィールドを削除し、不変オブジェクトとして再設計
  - ID、Power、Skill、Typeのgetterメソッドを実装
  - _要件: 2.1, 2.2, 4.1, 4.2, 5.1, 5.2, 5.3_

- [ ] 3. BattleCardSkillの簡素化
  - DescriptionKeyフィールドを削除
  - package privateフィールドとgetterメソッドを実装
  - _要件: 6.1, 6.4_

- [ ] 4. StructureCardの再設計
  - DescriptionKeyフィールドを削除
  - YieldModifierとBattlefieldModifierの抽象化を簡素化
  - package privateフィールドとgetterメソッドを実装
  - _要件: 6.2, 6.4_

- [ ] 5. Enemyの簡素化
  - Questionフィールドを削除
  - EnemySkillをinterfaceからstructに変更
  - package privateフィールドとコンストラクタを実装
  - _要件: 6.3, 6.4_

- [ ] 6. TerritoryとTerrainの分離
- [ ] 6.1 Terrain型の実装
  - 不変なTerrain型をid、baseYield、cardSlotフィールドで作成
  - Terrainのgetterメソッドを実装
  - _要件: 7.2_

- [ ] 6.2 Territoryの再設計
  - TerritoryをTerrain参照とカードリストで再構築
  - 防御的コピーを使用するCardsメソッドを実装
  - _要件: 7.1_

- [ ] 6.3 ConstructionPlanの実装
  - 建設計画管理のためのConstructionPlan型を作成
  - AddCard、RemoveCard、Cardsメソッドを実装
  - メモリ共有を避ける防御的コピーを実装
  - _要件: 7.3, 7.4_

- [ ] 7. Pointシステムの改善
- [ ] 7.1 PointType列挙とPointインターフェースの実装
  - PointType列挙値を定義
  - PointType、Passable、変換メソッドを持つPointインターフェースを作成
  - _要件: 8.1, 8.4_

- [ ] 7.2 Point変換メソッドの実装
  - AsBattlePoint、AsTerritoryPoint、AsMarketPointメソッドを実装
  - 型安全な変換とブール戻り値を提供
  - _要件: 8.2, 8.3_

- [ ] 8. Marketシステムの強化
  - MarketItemにLevelEffectとResourceQuantityフィールドを追加
  - 購入時の自動レベル効果適用を実装
  - nil CardPackサポートを追加（投資型アイテム用）
  - _要件: 11.1, 11.2, 11.3, 11.4_

- [ ] 9. ViewModelパッケージの実装
- [ ] 9.1 BattleViewModelの実装
  - core.GameStateとcore.Battlefieldを参照するBattleViewModelを作成
  - langとdrawingパッケージを使用する表示メソッドを実装
  - 敵情報、カードデータ、バトル状態のメソッドを実装
  - _要件: 9.1, 9.2_

- [ ] 9.2 CardDeckViewModelの実装
  - core.GameStateを参照するCardDeckViewModelを作成
  - バトルカードと建設カードの動的カウントを実装
  - カードアクセスメソッドを実装
  - _要件: 9.1_

- [ ] 9.3 BattleCardViewModelとStructureCardViewModelの実装
  - langとdrawingパッケージを使用するカード表示ViewModelを作成
  - 画像、名前、重複枚数、スキル情報のメソッドを実装
  - _要件: 9.1_

- [ ] 9.4 MapGridViewModelとPointViewModelの実装
  - マップグリッドとポイント表示のためのViewModelを作成
  - drawingパッケージを使用する画像取得を実装
  - ポイント情報と接続データのメソッドを実装
  - _要件: 9.1_

- [ ] 9.5 MarketViewModelとMarketItemViewModelの実装
  - マーケット表示のためのViewModelを作成
  - アイテム情報、価格、購入可能性のメソッドを実装
  - ResourceSufficiency型を実装
  - _要件: 9.1, 9.3_

- [ ] 9.6 TerritoryViewModelの実装
  - テリトリー表示のためのViewModelを作成
  - 収穫量、サポートパワー、カードスロット計算を実装
  - ConstructionPlan統合を実装
  - _要件: 9.1, 9.4_

- [ ] 9.7 その他のViewModelの実装
  - ResourceViewModel、CalendarViewModel、HistoryViewModelを実装
  - langパッケージを使用するローカライゼーション機能を実装
  - _要件: 9.1_

- [ ] 10. Flowパッケージの実装
- [ ] 10.1 BattleFlowの実装
  - バトル操作のためのBattleFlowを作成
  - RemoveFromBattle、Conquer、Rollbackメソッドを実装
  - _要件: 10.1, 10.2_

- [ ] 10.2 TerritoryFlowの実装
  - テリトリー操作のためのTerritoryFlowを作成
  - RemoveFromPlan、Commit、Rollbackメソッドを実装
  - _要件: 10.1, 10.4_

- [ ] 10.3 その他のFlowの実装
  - CardDeckFlow、MapGridFlow、MarketFlowを実装
  - 適切なユースケース操作メソッドを実装
  - _要件: 10.1, 10.3_

- [ ] 11. 既存コードの統合とテスト
- [ ] 11.1 loadパッケージの更新
  - 新しいコンストラクタ関数を使用するようにloadパッケージを更新
  - データ読み込みロジックを新しい型構造に適応
  - _要件: 2.2_

- [ ] 11.2 単体テストの実装
  - 各core型の単体テストを作成
  - 不変性、カプセル化、ビジネスロジックをテスト
  - テーブル駆動テストパターンを使用
  - _要件: 2.1, 5.1, 5.3_

- [ ] 11.3 統合テストの実装
  - viewmodel ↔ core連携の統合テストを作成
  - メモリ共有回避の確認テストを実装
  - flowパッケージ操作の検証テストを作成
  - _要件: 7.4, 9.1, 10.1_

- [ ] 11.4 UIパッケージのリファクタリング
  - UIコードを新しいviewmodelとflowパッケージを使用するように更新
  - 直接core依存を削除し、適切なレイヤーを通じてアクセス
  - 型キャストを変換メソッドに置き換え
  - _要件: 8.3, 9.1, 10.1_