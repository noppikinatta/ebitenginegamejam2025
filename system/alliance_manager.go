package system

import "github.com/noppikinatta/ebitenginegamejam2025/entity"

type AllianceBonus struct {
	ResourceBonus map[string]int
	MilitaryBonus int
	TradingBonus  int
}

type AllianceManager struct {
	relationships map[string]int           // NPC name -> relationship (0-100)
	allies        []string                 // List of allied NPCs
	bonuses       *AllianceBonus
	nationData    map[string]*entity.Nation // Optional: loaded from CSV, keyed by NPC name
}

func NewAllianceManager() *AllianceManager {
	return &AllianceManager{
		relationships: map[string]int{
			"Iron Republic":   50,
			"Forest Alliance": 60,
			"Desert Emirate":  40,
			"Mountain Clans":  35,
		},
		allies:     []string{},
		bonuses:    &AllianceBonus{ResourceBonus: make(map[string]int)},
		nationData: make(map[string]*entity.Nation),
	}
}

func (am *AllianceManager) GetAllies() []string {
	return am.allies
}

func (am *AllianceManager) GetRelationship(npcName string) int {
	if relationship, exists := am.relationships[npcName]; exists {
		return relationship
	}
	return 0
}

func (am *AllianceManager) ImproveRelationship(npcName string, amount int) {
	if _, exists := am.relationships[npcName]; exists {
		am.relationships[npcName] += amount
		if am.relationships[npcName] > 100 {
			am.relationships[npcName] = 100
		}
	}
}

func (am *AllianceManager) FormAlliance(npcName string) bool {
	// Check if relationship is high enough
	if am.GetRelationship(npcName) < 70 {
		return false
	}

	// Check if already allied
	for _, ally := range am.allies {
		if ally == npcName {
			return false
		}
	}

	// Form alliance
	am.allies = append(am.allies, npcName)
	am.updateAllianceBonuses()
	return true
}

func (am *AllianceManager) updateAllianceBonuses() {
	// Reset bonuses
	am.bonuses.ResourceBonus = make(map[string]int)
	am.bonuses.MilitaryBonus = 0
	am.bonuses.TradingBonus = 0

	// Calculate bonuses based on allies. Prefer CSV-loaded data if available,
	// otherwise fall back to the original hard-coded values so existing
	// behaviour (and tests) stay unchanged.
	for _, ally := range am.allies {
		if n, ok := am.nationData[ally]; ok {
			if n.AllyBonusGold != 0 {
				am.bonuses.ResourceBonus["Gold"] += n.AllyBonusGold
			}
			if n.AllyBonusAttack != 0 {
				am.bonuses.MilitaryBonus += n.AllyBonusAttack
			}
			// Extend easily later with more resource types
			continue
		}

		// Fallback
		switch ally {
		case "Iron Republic":
			am.bonuses.ResourceBonus["Iron"] += 5
			am.bonuses.MilitaryBonus += 10
		case "Forest Alliance":
			am.bonuses.ResourceBonus["Wood"] += 8
			am.bonuses.ResourceBonus["Grain"] += 3
		case "Desert Emirate":
			am.bonuses.ResourceBonus["Gold"] += 7
			am.bonuses.TradingBonus += 15
		case "Mountain Clans":
			am.bonuses.ResourceBonus["Iron"] += 3
			am.bonuses.ResourceBonus["Mana"] += 5
			am.bonuses.MilitaryBonus += 5
		}
	}
}

func (am *AllianceManager) GetAllianceBonuses() *AllianceBonus {
	return am.bonuses
}

func (am *AllianceManager) IsAllied(npcName string) bool {
	for _, ally := range am.allies {
		if ally == npcName {
			return true
		}
	}
	return false
}

func (am *AllianceManager) BreakAlliance(npcName string) {
	for i, ally := range am.allies {
		if ally == npcName {
			am.allies = append(am.allies[:i], am.allies[i+1:]...)
			am.updateAllianceBonuses()
			break
		}
	}
}

func (am *AllianceManager) GetAllRelationships() map[string]int {
	result := make(map[string]int)
	for k, v := range am.relationships {
		result[k] = v
	}
	return result
}

// NewAllianceManagerFromData initialises relationships from CSV-loaded nations map.
func NewAllianceManagerFromData(nations map[string]*entity.Nation) *AllianceManager {
	am := NewAllianceManager()
	if len(nations) == 0 {
		return am
	}
	for _, n := range nations {
		am.relationships[n.Name] = n.InitialRelationship
		am.nationData[n.Name] = n
	}

	// Ensure bonuses map is initialised (if caller constructs via FromData)
	if am.bonuses == nil {
		am.bonuses = &AllianceBonus{ResourceBonus: make(map[string]int)}
	}
	return am
}
