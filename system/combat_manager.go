package system

import (
	"github.com/noppikinatta/ebitenginegamejam2025/entity"
)

type CombatBattlefield struct {
	frontRow [5]string
	backRow  [5]string
}

type CombatEnemy struct {
	Name    string
	Attack  int
	Defense int
	Health  int
}

type CombatStats struct {
	RoundsCompleted  int
	TotalDamageDealt int
	TotalDamageTaken int
}

type VictoryRewards struct {
	Resources  map[string]int
	Cards      []*entity.Card
	Experience int
}

type CombatManager struct {
	battlefield       *CombatBattlefield
	enemies           map[string]*CombatEnemy
	enemyTemplates    map[string]*entity.Enemy
	bossTemplates     map[string]*entity.Boss
	combatState       string // "ongoing", "victory", "defeat"
	combatStats       *CombatStats
	placedCards       []*entity.Card
	rewards           *VictoryRewards
	playerAttackBonus int
}

func NewCombatManager() *CombatManager {
	return &CombatManager{
		battlefield: &CombatBattlefield{
			frontRow: [5]string{"", "", "", "", ""},
			backRow:  [5]string{"", "", "", "", ""},
		},
		enemies:     make(map[string]*CombatEnemy),
		combatState: "ongoing",
		combatStats: &CombatStats{},
		placedCards: []*entity.Card{},
		rewards: &VictoryRewards{
			Resources:  make(map[string]int),
			Cards:      []*entity.Card{},
			Experience: 0,
		},
		enemyTemplates: make(map[string]*entity.Enemy),
		bossTemplates:  make(map[string]*entity.Boss),
	}
}

func (cb *CombatBattlefield) GetFrontRow() *[5]string {
	return &cb.frontRow
}

func (cb *CombatBattlefield) GetBackRow() *[5]string {
	return &cb.backRow
}

func (cb *CombatBattlefield) PlaceCard(cardName, row string, position int) bool {
	if position < 0 || position >= 5 {
		return false
	}

	if row == "front" {
		if cb.frontRow[position] == "" {
			cb.frontRow[position] = cardName
			return true
		}
	} else if row == "back" {
		if cb.backRow[position] == "" {
			cb.backRow[position] = cardName
			return true
		}
	}

	return false
}

func (cm *CombatManager) GetBattlefield() *CombatBattlefield {
	return cm.battlefield
}

func (cm *CombatManager) PlaceCardInBattle(card *entity.Card, row string, position int) bool {
	success := cm.battlefield.PlaceCard(card.Name, row, position)
	if success {
		cm.placedCards = append(cm.placedCards, card)
	}
	return success
}

func (cm *CombatManager) GetPlacedCards() []*entity.Card {
	return cm.placedCards
}

func (cm *CombatManager) AddEnemy(name string, attack, defense, health int) {
	enemy := &CombatEnemy{
		Name:    name,
		Attack:  attack,
		Defense: defense,
		Health:  health,
	}
	cm.enemies[name] = enemy
}

func (cm *CombatManager) GetEnemyHealth(name string) int {
	if enemy, exists := cm.enemies[name]; exists {
		return enemy.Health
	}
	return 0
}

func (cm *CombatManager) CalculatePlayerDamage() int {
	totalDamage := 0
	for _, card := range cm.placedCards {
		totalDamage += card.Attack
	}
	totalDamage += cm.playerAttackBonus
	return totalDamage
}

func (cm *CombatManager) CalculateEnemyDamage() int {
	totalDamage := 0
	for _, enemy := range cm.enemies {
		if enemy.Health > 0 {
			totalDamage += enemy.Attack
		}
	}
	return totalDamage
}

func (cm *CombatManager) ExecuteCombatRound() {
	cm.combatStats.RoundsCompleted++

	// Player attacks enemies
	playerDamage := cm.CalculatePlayerDamage()
	for _, enemy := range cm.enemies {
		if enemy.Health > 0 {
			damage := playerDamage - enemy.Defense
			if damage < 0 {
				damage = 0
			}
			enemy.Health -= damage
			cm.combatStats.TotalDamageDealt += damage

			if enemy.Health < 0 {
				enemy.Health = 0
			}
		}
	}

	// Check victory condition
	allEnemiesDefeated := true
	for _, enemy := range cm.enemies {
		if enemy.Health > 0 {
			allEnemiesDefeated = false
			break
		}
	}

	if allEnemiesDefeated {
		cm.combatState = "victory"
		cm.calculateVictoryRewards()
		return
	}

	// Enemies attack player (simplified - reduce card health)
	enemyDamage := cm.CalculateEnemyDamage()
	if enemyDamage > 0 {
		cm.combatStats.TotalDamageTaken += enemyDamage
		// Simplified: if enemy damage is too high, defeat
		totalPlayerHealth := len(cm.placedCards) * 3 // Simplified health calculation
		if cm.combatStats.TotalDamageTaken >= totalPlayerHealth {
			cm.combatState = "defeat"
		}
	}
}

func (cm *CombatManager) calculateVictoryRewards() {
	// Calculate rewards based on enemies defeated
	cm.rewards.Resources["Gold"] = len(cm.enemies) * 20
	cm.rewards.Resources["Iron"] = len(cm.enemies) * 5
	cm.rewards.Experience = cm.combatStats.RoundsCompleted * 10
}

func (cm *CombatManager) GetCombatState() string {
	return cm.combatState
}

func (cm *CombatManager) GetCombatStats() *CombatStats {
	return cm.combatStats
}

func (cm *CombatManager) GetVictoryRewards() *VictoryRewards {
	return cm.rewards
}

func (cm *CombatManager) ResetCombat() {
	cm.combatState = "ongoing"
	cm.combatStats = &CombatStats{}
	cm.placedCards = []*entity.Card{}
	cm.enemies = make(map[string]*CombatEnemy)
	cm.rewards = &VictoryRewards{
		Resources:  make(map[string]int),
		Cards:      []*entity.Card{},
		Experience: 0,
	}

	// Reset battlefield
	cm.battlefield = &CombatBattlefield{
		frontRow: [5]string{"", "", "", "", ""},
		backRow:  [5]string{"", "", "", "", ""},
	}
}

func (cm *CombatManager) ClearPlayerCards() {
	cm.placedCards = []*entity.Card{}
	cm.battlefield = &CombatBattlefield{
		frontRow: [5]string{"", "", "", "", ""},
		backRow:  [5]string{"", "", "", "", ""},
	}
}

func (cm *CombatManager) ClearEnemies() {
	cm.enemies = make(map[string]*CombatEnemy)
}

func (cm *CombatManager) GetAllEnemies() []*CombatEnemy {
	enemies := make([]*CombatEnemy, 0, len(cm.enemies))
	for _, e := range cm.enemies {
		enemies = append(enemies, e)
	}
	return enemies
}

func (cm *CombatManager) SetPlayerAttackBonus(bonus int) {
	cm.playerAttackBonus = bonus
}

// NewCombatManagerWithTemplates injects enemy/boss template maps.
func NewCombatManagerWithTemplates(enemies map[string]*entity.Enemy, bosses map[string]*entity.Boss) *CombatManager {
	cm := NewCombatManager()
	if len(enemies) > 0 {
		cm.enemyTemplates = enemies
	}
	if len(bosses) > 0 {
		cm.bossTemplates = bosses
	}
	return cm
}

// AddEnemyByID pulls stats from enemy or boss template maps.
func (cm *CombatManager) AddEnemyByID(id string) bool {
	if e, ok := cm.enemyTemplates[id]; ok {
		cm.AddEnemy(e.Name, e.Attack, e.Defense, e.Health)
		return true
	}
	if b, ok := cm.bossTemplates[id]; ok {
		cm.AddEnemy(b.Name, b.Attack, b.Defense, b.Health)
		return true
	}
	return false
}
