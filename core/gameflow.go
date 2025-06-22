package core

// GameState ゲームの全体状態を管理する
type GameState struct {
	MyNation    *MyNation // プレイヤー国家
	MapGrid     *MapGrid  // マップグリッド
	Treasury    *Treasury // プレイヤーの国庫
	CurrentTurn int       // 現在のターン数
}

// AddYield 制圧済みのTerritoryとMyNationのBasicYieldをTreasuryに加算する
func (gs *GameState) AddYield() {
	// MyNationのBasicYieldを加算
	totalYield := gs.MyNation.BasicYield

	// 制圧済みWildernessPointのTerritoryのYieldを加算
	for _, point := range gs.MapGrid.Points {
		if wilderness, ok := point.(*WildernessPoint); ok && wilderness.Controlled {
			totalYield = totalYield.Add(wilderness.Territory.Yield())
		}
	}

	// Treasuryに加算
	gs.Treasury.Add(totalYield)
}

// NextTurn ターンを進行し、Yieldを加算する
func (gs *GameState) NextTurn() {
	gs.CurrentTurn++
	gs.AddYield()
}

// IsVictory 勝利条件を判定する（全てのBossPointが撃破されているか）
func (gs *GameState) IsVictory() bool {
	for _, point := range gs.MapGrid.Points {
		if bossPoint, ok := point.(*BossPoint); ok {
			if !bossPoint.Defeated {
				return false // 撃破されていないBossが存在する
			}
		}
	}
	return true // 全てのBossが撃破されている（またはBossが存在しない）
}

// CanInteract 指定座標のPointが操作可能かどうかを判定する
func (gs *GameState) CanInteract(x, y int) bool {
	return gs.MapGrid.CanInteract(x, y)
}

// GetPoint 指定座標のPointを取得する
func (gs *GameState) GetPoint(x, y int) Point {
	return gs.MapGrid.GetPoint(x, y)
}
