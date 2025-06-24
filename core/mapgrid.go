package core

// MapGrid関連の実装（旧仕様コメントを削除し、実装に置き換え済み）

// Point は、MapGrid上のPointを表すインターフェース
type Point interface {
	Passable() bool
}

// MyNationPoint プレイヤー国家のPoint
type MyNationPoint struct {
	MyNation *MyNation
}

func (p *MyNationPoint) Passable() bool {
	return true
}

// OtherNationPoint NPC国家のPoint
type OtherNationPoint struct {
	OtherNation *OtherNation
}

func (p *OtherNationPoint) Passable() bool {
	return true
}

// WildernessPoint 制圧可能な野生のPoint
type WildernessPoint struct {
	Controlled bool       // 制圧済みかどうか
	Enemy      *Enemy     // 守っているEnemy
	Territory  *Territory // 制圧後のTerritory
}

func (p *WildernessPoint) Passable() bool {
	return p.Controlled
}

// BossPoint ボスのPoint
type BossPoint struct {
	Boss     *Enemy
	Defeated bool // ボスが撃破されているかどうか
}

func (p *BossPoint) Passable() bool {
	return false
}

// MapGrid ゲームのマップグリッド
type MapGrid struct {
	SizeX      int     // X方向のサイズ
	SizeY      int     // Y方向のサイズ
	Points     []Point // Pointの一覧。インデックスは y*SizeX + x で計算
	accesibles []bool
}

// GetPoint 指定座標のPointを取得する
func (m *MapGrid) GetPoint(x, y int) Point {
	index, ok := m.IndexFromPoint(x, y)
	if !ok {
		return nil
	}
	return m.Points[index]
}

func (m *MapGrid) IndexFromPoint(x, y int) (int, bool) {
	if x < 0 || x >= m.SizeX || y < 0 || y >= m.SizeY {
		return 0, false
	}
	return y*m.SizeX + x, true
}

func (m *MapGrid) PointFromIndex(index int) (int, int, bool) {
	if index < 0 || index >= len(m.Points) {
		return 0, 0, false
	}
	return index % m.SizeX, index / m.SizeX, true
}

func (m *MapGrid) UpdateAccesibles() {
	remainingIdxs := make([]int, 0, len(m.Points))
	var ri int

	if m.accesibles == nil {
		m.accesibles = make([]bool, len(m.Points))
	}

	for i := range len(m.Points) {
		m.accesibles[i] = false
	}
	m.accesibles[0] = true
	remainingIdxs = append(remainingIdxs, 0)

	for ri < len(remainingIdxs) {
		idx := remainingIdxs[ri]
		ri++

		x, y, ok := m.PointFromIndex(idx)
		if !ok {
			continue
		}

		rightIdx, ok := m.IndexFromPoint(x+1, y)
		if ok && !m.accesibles[rightIdx] {
			p := m.Points[rightIdx]
			m.accesibles[rightIdx] = true
			if p.Passable() {
				remainingIdxs = append(remainingIdxs, rightIdx)
			}
		}

		upIdx, ok := m.IndexFromPoint(x, y+1)
		if ok && !m.accesibles[upIdx] {
			p := m.Points[upIdx]
			m.accesibles[upIdx] = true
			if p.Passable() {
				remainingIdxs = append(remainingIdxs, upIdx)
			}
		}
	}
}

// CanInteract 指定座標のPointが操作可能かどうかを判定する
// MyNationPointから制圧済みのWildernessPointやOtherNationPoint、BossPointへの
// 連続した制圧済みのルートが存在する場合のみ操作可能
func (m *MapGrid) CanInteract(x, y int) bool {
	if m.accesibles == nil {
		m.UpdateAccesibles()
	}

	idx, ok := m.IndexFromPoint(x, y)
	if !ok {
		return false
	}

	return m.accesibles[idx]
}
