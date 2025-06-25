package core

// GameState ゲームの全体状態を管理する
type GameState struct {
	MyNation    *MyNation // プレイヤー国家
	CardDeck    *CardDeck // プレイヤーの所持カード
	MapGrid     *MapGrid  // マップグリッド
	Treasury    *Treasury // プレイヤーの国庫
	CurrentTurn int       // 現在のターン数
}

// AddYield 制圧済みのTerritoryとMyNationのBasicYieldをTreasuryに加算する
func (g *GameState) AddYield() {
	// MyNationのBasicYieldを加算
	totalYield := g.MyNation.BasicYield

	// 制圧済みWildernessPointのTerritoryのYieldを加算
	for _, point := range g.MapGrid.Points {
		if wilderness, ok := point.(*WildernessPoint); ok && wilderness.Controlled {
			totalYield = totalYield.Add(wilderness.Territory.Yield())
		}
	}

	// Treasuryに加算
	g.Treasury.Add(totalYield)
}

// NextTurn ターンを進行し、Yieldを加算する
func (g *GameState) NextTurn() {
	g.CurrentTurn++
	g.AddYield()
}

// IsVictory 勝利条件を判定する（全てのBossPointが撃破されているか）
func (g *GameState) IsVictory() bool {
	for _, point := range g.MapGrid.Points {
		if bossPoint, ok := point.(*BossPoint); ok {
			if !bossPoint.Defeated {
				return false // 撃破されていないBossが存在する
			}
		}
	}
	return true // 全てのBossが撃破されている（またはBossが存在しない）
}

// CanInteract 指定座標のPointが操作可能かどうかを判定する
func (g *GameState) CanInteract(x, y int) bool {
	return g.MapGrid.CanInteract(x, y)
}

// GetPoint 指定座標のPointを取得する
func (g *GameState) GetPoint(x, y int) Point {
	return g.MapGrid.GetPoint(x, y)
}
