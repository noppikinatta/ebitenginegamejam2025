package screen

import (
	"github.com/noppikinatta/ebitenginegamejam2025/system"
)

type Enemy struct {
	Name    string
	Attack  int
	Defense int
	Health  int
}

type BattleScreen struct {
	combatManager *system.CombatManager
	enemies       []*Enemy
}

func NewBattleScreen(cm *system.CombatManager) *BattleScreen {
	defaultEnemies := []*Enemy{
		{Name: "Goblin Warrior", Attack: 3, Defense: 2, Health: 5},
		{Name: "Orc Brute", Attack: 5, Defense: 3, Health: 8},
		{Name: "Dark Mage", Attack: 4, Defense: 1, Health: 4},
	}

	return &BattleScreen{
		combatManager: cm,
		enemies:       defaultEnemies,
	}
}

// Proxy methods ------------------------------------------------------------

func (bs *BattleScreen) GetBattlefield() *system.CombatBattlefield {
	return bs.combatManager.GetBattlefield()
}

func (bs *BattleScreen) GetEnemies() []*Enemy {
	return bs.enemies
}

// Backward compatible helpers (for tests that place cards directly)

func (bs *BattleScreen) PlaceCard(cardName, row string, position int) bool {
	bf := bs.combatManager.GetBattlefield()
	return bf.PlaceCard(cardName, row, position)
}
