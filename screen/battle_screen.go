package screen

type Battlefield struct {
	frontRow [5]string
	backRow  [5]string
}

type Enemy struct {
	Name    string
	Attack  int
	Defense int
	Health  int
}

type BattleScreen struct {
	battlefield *Battlefield
	enemies     []*Enemy
}

func NewBattleScreen() *BattleScreen {
	battlefield := &Battlefield{
		frontRow: [5]string{"", "", "", "", ""},
		backRow:  [5]string{"", "", "", "", ""},
	}

	enemies := []*Enemy{
		{Name: "Goblin Warrior", Attack: 3, Defense: 2, Health: 5},
		{Name: "Orc Brute", Attack: 5, Defense: 3, Health: 8},
		{Name: "Dark Mage", Attack: 4, Defense: 1, Health: 4},
	}

	return &BattleScreen{
		battlefield: battlefield,
		enemies:     enemies,
	}
}

func (bs *BattleScreen) GetBattlefield() *Battlefield {
	return bs.battlefield
}

func (bs *BattleScreen) GetEnemies() []*Enemy {
	return bs.enemies
}

func (bf *Battlefield) GetFrontRow() *[5]string {
	return &bf.frontRow
}

func (bf *Battlefield) GetBackRow() *[5]string {
	return &bf.backRow
}

func (bf *Battlefield) PlaceCard(cardName, row string, position int) bool {
	if position < 0 || position >= 5 {
		return false
	}

	if row == "front" {
		if bf.frontRow[position] == "" {
			bf.frontRow[position] = cardName
			return true
		}
	} else if row == "back" {
		if bf.backRow[position] == "" {
			bf.backRow[position] = cardName
			return true
		}
	}

	return false
}
