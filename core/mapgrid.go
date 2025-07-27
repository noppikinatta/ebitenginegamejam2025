package core

// Implementation related to MapGrid (old specification comments have been deleted and replaced with implementation)

// PointType represents the type of a point on the MapGrid.
type PointType int

const (
	PointTypeUnknown PointType = iota
	PointTypeMyNation
	PointTypeOtherNation
	PointTypeWilderness
	PointTypeBoss
)

// Point is an interface representing a point on the MapGrid.
type Point interface {
	PointType() PointType
	Passable() bool
	AsBattlePoint() (BattlePoint, bool)
	AsTerritoryPoint() (TerritoryPoint, bool)
	AsMarketPoint() (MarketPoint, bool)
}

// BattlePoint represents a point where battles can occur.
type BattlePoint interface {
	Enemy() *Enemy
	Conquer()
}

// TerritoryPoint represents a point that provides territory functionality.
type TerritoryPoint interface {
	Yield() ResourceQuantity
	Terrain() *Terrain
	CardSlot() int
	Cards() []*StructureCard
}

// MarketPoint represents a point that provides market functionality.
type MarketPoint interface {
	Nation() Nation
}

// MyNationPoint is a point of the player's nation.
type MyNationPoint struct {
	MyNation *MyNation
}

func (p *MyNationPoint) PointType() PointType {
	return PointTypeMyNation
}

func (p *MyNationPoint) Passable() bool {
	return true
}

func (p *MyNationPoint) AsBattlePoint() (BattlePoint, bool) {
	return nil, false
}

func (p *MyNationPoint) AsTerritoryPoint() (TerritoryPoint, bool) {
	return nil, false
}

func (p *MyNationPoint) AsMarketPoint() (MarketPoint, bool) {
	return p, true
}

func (p *MyNationPoint) Nation() Nation {
	return p.MyNation
}

// OtherNationPoint is a point of an NPC nation.
type OtherNationPoint struct {
	OtherNation *OtherNation
}

func (p *OtherNationPoint) PointType() PointType {
	return PointTypeOtherNation
}

func (p *OtherNationPoint) Passable() bool {
	return true
}

func (p *OtherNationPoint) AsBattlePoint() (BattlePoint, bool) {
	return nil, false
}

func (p *OtherNationPoint) AsTerritoryPoint() (TerritoryPoint, bool) {
	return nil, false
}

func (p *OtherNationPoint) AsMarketPoint() (MarketPoint, bool) {
	return p, true
}

func (p *OtherNationPoint) Nation() Nation {
	return p.OtherNation
}

// WildernessPoint is a conquerable wild point.
type WildernessPoint struct {
	terrainType string
	controlled  bool       // Whether it is controlled
	enemy       *Enemy     // The Enemy guarding it
	territory   *Territory // The Territory after conquest
}

func (p *WildernessPoint) PointType() PointType {
	return PointTypeWilderness
}

func (p *WildernessPoint) Passable() bool {
	return p.controlled
}

func (p *WildernessPoint) AsBattlePoint() (BattlePoint, bool) {
	if !p.controlled && p.enemy != nil {
		return p, true
	}
	return nil, false
}

func (p *WildernessPoint) AsTerritoryPoint() (TerritoryPoint, bool) {
	if p.controlled && p.territory != nil {
		return p, true
	}
	return nil, false
}

func (p *WildernessPoint) AsMarketPoint() (MarketPoint, bool) {
	return nil, false
}

// BattlePoint interface implementation
func (p *WildernessPoint) Enemy() *Enemy {
	return p.enemy
}

func (p *WildernessPoint) Conquer() {
	p.controlled = true
}

// TerritoryPoint interface implementation
func (p *WildernessPoint) Yield() ResourceQuantity {
	if p.territory != nil {
		return p.territory.Yield()
	}
	return ResourceQuantity{}
}

func (p *WildernessPoint) Terrain() *Terrain {
	if p.territory != nil {
		return p.territory.Terrain()
	}
	return nil
}

func (p *WildernessPoint) CardSlot() int {
	if p.territory != nil {
		return p.territory.Terrain().CardSlot()
	}
	return 0
}

func (p *WildernessPoint) Cards() []*StructureCard {
	if p.territory != nil {
		return p.territory.Cards()
	}
	return []*StructureCard{}
}

// SetControlledForTest sets the controlled status for testing purposes.
func (p *WildernessPoint) SetControlledForTest(controlled bool) {
	p.controlled = controlled
}

// SetTerritoryForTest sets the territory for testing purposes.
func (p *WildernessPoint) SetTerritoryForTest(territory *Territory) {
	p.territory = territory
}

// SetEnemyForTest sets the enemy for testing purposes.
func (p *WildernessPoint) SetEnemyForTest(enemy *Enemy) {
	p.enemy = enemy
}

// Controlled returns the controlled status.
func (p *WildernessPoint) Controlled() bool {
	return p.controlled
}

// Territory returns the territory.
func (p *WildernessPoint) Territory() *Territory {
	return p.territory
}

// BossPoint is a point of a boss.
type BossPoint struct {
	boss     *Enemy
	defeated bool // Whether the boss has been defeated
}

func (p *BossPoint) PointType() PointType {
	return PointTypeBoss
}

func (p *BossPoint) Passable() bool {
	return false
}

func (p *BossPoint) AsBattlePoint() (BattlePoint, bool) {
	if !p.defeated && p.boss != nil {
		return p, true
	}
	return nil, false
}

func (p *BossPoint) AsTerritoryPoint() (TerritoryPoint, bool) {
	return nil, false
}

func (p *BossPoint) AsMarketPoint() (MarketPoint, bool) {
	return nil, false
}

// BattlePoint interface implementation
func (p *BossPoint) Enemy() *Enemy {
	return p.boss
}

func (p *BossPoint) Conquer() {
	p.defeated = true
}

// SetBossForTest sets the boss for testing purposes.
func (p *BossPoint) SetBossForTest(boss *Enemy) {
	p.boss = boss
}

// SetDefeatedForTest sets the defeated status for testing purposes.
func (p *BossPoint) SetDefeatedForTest(defeated bool) {
	p.defeated = defeated
}

// Boss returns the boss enemy.
func (p *BossPoint) Boss() *Enemy {
	return p.boss
}

// MapGrid is the game's map grid.
type MapGrid struct {
	Size       MapGridSize
	Points     []Point // List of Points. The index is calculated by y*SizeX + x.
	accesibles []bool
}

// GetPoint gets the Point at the specified coordinates.
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
		m.accesibles[i] = p.PointType() == PointTypeMyNation
		if p.PointType() == PointTypeMyNation {
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

// CanInteract determines whether the Point at the specified coordinates can be interacted with.
// From MyNationPoint to a controlled WildernessPoint, OtherNationPoint, or BossPoint
// It can be interacted with only if a continuous controlled route exists.
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
