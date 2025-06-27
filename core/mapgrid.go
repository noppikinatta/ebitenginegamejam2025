package core

// MapGrid関連の実装（旧仕様コメントを削除し、実装に置き換え済み）

// Point は、MapGrid上のPointを表すインターフェース
type Point interface {
	Passable() bool
	IsMyNation() bool
}

type BattlePoint interface {
	Point
	GetEnemy() *Enemy
	SetControlled(bool)
}

// MyNationPoint プレイヤー国家のPoint
type MyNationPoint struct {
	MyNation *MyNation
}

func (p *MyNationPoint) Passable() bool {
	return true
}

func (p *MyNationPoint) IsMyNation() bool {
	return true
}

// OtherNationPoint NPC国家のPoint
type OtherNationPoint struct {
	OtherNation *OtherNation
}

func (p *OtherNationPoint) Passable() bool {
	return true
}

func (p *OtherNationPoint) IsMyNation() bool {
	return false
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

func (p *WildernessPoint) IsMyNation() bool {
	return false
}

func (p *WildernessPoint) GetEnemy() *Enemy {
	return p.Enemy
}

func (p *WildernessPoint) SetControlled(controlled bool) {
	p.Controlled = controlled
}

// BossPoint ボスのPoint
type BossPoint struct {
	Boss     *Enemy
	Defeated bool // ボスが撃破されているかどうか
}

func (p *BossPoint) Passable() bool {
	return false
}

func (p *BossPoint) IsMyNation() bool {
	return false
}

func (p *BossPoint) GetEnemy() *Enemy {
	return p.Boss
}

func (p *BossPoint) SetControlled(controlled bool) {
	p.Defeated = controlled
}

// MapGrid ゲームのマップグリッド
type MapGrid struct {
	Size       MapGridSize
	Points     []Point // Pointの一覧。インデックスは y*SizeX + x で計算
	accesibles []bool
}

// GetPoint 指定座標のPointを取得する
func (m *MapGrid) GetPoint(x, y int) Point {
	index, ok := m.IndexFromXY(x, y)
	if !ok {
		return nil
	}
	return m.Points[index]
}

func (m *MapGrid) XYOfPoint(p Point) (int, int, bool) {
	for i, pp := range m.Points {
		if pp == p {
			return m.XYFromIndex(i)
		}
	}
	return 0, 0, false
}

func (m *MapGrid) IndexFromXY(x, y int) (int, bool) {
	if x < 0 || x >= m.Size.X || y < 0 || y >= m.Size.Y {
		return 0, false
	}
	return m.Size.Index(x, y), true
}

func (m *MapGrid) XYFromIndex(index int) (int, int, bool) {
	x, y := m.Size.XY(index)
	if x < 0 || x >= m.Size.X || y < 0 || y >= m.Size.Y {
		return 0, 0, false
	}
	return x, y, true
}

func (m *MapGrid) UpdateAccesibles() {
	alreadySet := make(map[int]struct{})
	remainingIdxs := make([]int, 0, len(m.Points))
	var ri int

	if m.accesibles == nil {
		m.accesibles = make([]bool, len(m.Points))
	}

	for i, p := range m.Points {
		m.accesibles[i] = p.IsMyNation()
		if p.IsMyNation() {
			remainingIdxs = append(remainingIdxs, i)
		}
	}

	for ri < len(remainingIdxs) {
		idx := remainingIdxs[ri]
		alreadySet[idx] = struct{}{}
		ri++

		x, y, ok := m.XYFromIndex(idx)
		if !ok {
			continue
		}

		set := func(x, y int) {
			idx, ok := m.IndexFromXY(x, y)
			if !ok {
				return
			}
			if _, ok := alreadySet[idx]; ok {
				return
			}
			p := m.Points[idx]
			m.accesibles[idx] = true
			if p != nil && p.Passable() {
				remainingIdxs = append(remainingIdxs, idx)
			}
		}

		set(x+1, y)
		set(x-1, y)
		set(x, y+1)
		set(x, y-1)
	}
}

// CanInteract 指定座標のPointが操作可能かどうかを判定する
// MyNationPointから制圧済みのWildernessPointやOtherNationPoint、BossPointへの
// 連続した制圧済みのルートが存在する場合のみ操作可能
func (m *MapGrid) CanInteract(x, y int) bool {
	if m.accesibles == nil {
		m.UpdateAccesibles()
	}

	idx, ok := m.IndexFromXY(x, y)
	if !ok {
		return false
	}

	return m.accesibles[idx]
}

type MapGridSize struct {
	X int
	Y int
}

func (s MapGridSize) Index(x, y int) int {
	return y*s.X + x
}

func (s MapGridSize) XY(index int) (int, int) {
	return index % s.X, index / s.X
}

func (s MapGridSize) Length() int {
	return s.X * s.Y
}
