package core

// MapGrid関連の実装（旧仕様コメントを削除し、実装に置き換え済み）

// Point は、MapGrid上のPointを表すインターフェース
type Point interface{}

// MyNationPoint プレイヤー国家のPoint
type MyNationPoint struct {
	MyNation *MyNation
}

// OtherNationPoint NPC国家のPoint
type OtherNationPoint struct {
	OtherNation *OtherNation
}

// WildernessPoint 制圧可能な野生のPoint
type WildernessPoint struct {
	Controlled bool        // 制圧済みかどうか
	Enemy      *Enemy      // 守っているEnemy
	Territory  *Territory  // 制圧後のTerritory
}

// BossPoint ボスのPoint
type BossPoint struct {
	Boss     *Enemy
	Defeated bool // ボスが撃破されているかどうか
}

// MapGrid ゲームのマップグリッド
type MapGrid struct {
	SizeX  int     // X方向のサイズ
	SizeY  int     // Y方向のサイズ
	Points []Point // Pointの一覧。インデックスは y*SizeX + x で計算
}

// GetPoint 指定座標のPointを取得する
func (mg *MapGrid) GetPoint(x, y int) Point {
	if x < 0 || x >= mg.SizeX || y < 0 || y >= mg.SizeY {
		return nil
	}
	
	index := y*mg.SizeX + x
	if index >= len(mg.Points) {
		return nil
	}
	
	return mg.Points[index]
}

// CanInteract 指定座標のPointが操作可能かどうかを判定する
// MyNationPointから制圧済みのWildernessPointやOtherNationPoint、BossPointへの
// 連続した制圧済みのルートが存在する場合のみ操作可能
func (mg *MapGrid) CanInteract(x, y int) bool {
	// 範囲外チェック
	if x < 0 || x >= mg.SizeX || y < 0 || y >= mg.SizeY {
		return false
	}

	// BFS (幅優先探索) で到達可能性を判定
	visited := make([]bool, mg.SizeX*mg.SizeY)
	queue := make([][2]int, 0)

	// MyNationPointを探して開始点に設定
	for startY := 0; startY < mg.SizeY; startY++ {
		for startX := 0; startX < mg.SizeX; startX++ {
			point := mg.GetPoint(startX, startY)
			if _, ok := point.(*MyNationPoint); ok {
				queue = append(queue, [2]int{startX, startY})
				visited[startY*mg.SizeX+startX] = true
				break
			}
		}
		if len(queue) > 0 {
			break // MyNationPointが見つかったら探索を開始
		}
	}

	// MyNationPointが見つからない場合は操作不可
	if len(queue) == 0 {
		return false
	}

	// BFS実行
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		
		currentX, currentY := current[0], current[1]
		
		// 目標座標に到達した場合
		if currentX == x && currentY == y {
			return true
		}

		// 隣接する4方向をチェック
		directions := [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
		for _, dir := range directions {
			nextX, nextY := currentX+dir[0], currentY+dir[1]
			
			// 範囲内チェック
			if nextX < 0 || nextX >= mg.SizeX || nextY < 0 || nextY >= mg.SizeY {
				continue
			}
			
			nextIndex := nextY*mg.SizeX + nextX
			if visited[nextIndex] {
				continue
			}
			
			point := mg.GetPoint(nextX, nextY)
			if point == nil {
				continue
			}
			
			// 操作可能なPointかどうかチェック
			canPass := false
			switch p := point.(type) {
			case *MyNationPoint:
				canPass = true
			case *OtherNationPoint:
				canPass = true
			case *BossPoint:
				canPass = true
			case *WildernessPoint:
				canPass = p.Controlled // 制圧済みのみ通過可能
			}
			
			if canPass {
				visited[nextIndex] = true
				queue = append(queue, [2]int{nextX, nextY})
			}
		}
	}

	return false
}
