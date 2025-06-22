package core_test

import (
	"testing"

	"github.com/noppikinatta/ebitenginegamejam2025/core"
)

func TestEnemy_Basic(t *testing.T) {
	enemy := &core.Enemy{
		EnemyID:        "test_enemy_1",
		EnemyType:      "orc",
		Power:          25.0,
		BattleCardSlot: 3,
		Skills:         []core.EnemySkill{},
	}

	if enemy.EnemyID != "test_enemy_1" {
		t.Errorf("EnemyID = %v, want %v", enemy.EnemyID, "test_enemy_1")
	}
	if enemy.EnemyType != "orc" {
		t.Errorf("EnemyType = %v, want %v", enemy.EnemyType, "orc")
	}
	if enemy.Power != 25.0 {
		t.Errorf("Power = %v, want %v", enemy.Power, 25.0)
	}
	if enemy.BattleCardSlot != 3 {
		t.Errorf("BattleCardSlot = %v, want %v", enemy.BattleCardSlot, 3)
	}
}

func TestEnemySkill_CanAffect(t *testing.T) {
	// テスト用のEnemySkill実装
	warriorDebuffSkill := &mockEnemySkill{
		canAffectFunc: func(card *core.BattleCard) bool {
			return card.Type == "warrior"
		},
		modifyFunc: func(card *core.BattleCard) float64 {
			return float64(card.Power) * 0.5 // 戦士の力を半減
		},
	}

	tests := []struct {
		name     string
		skill    core.EnemySkill
		card     *core.BattleCard
		expected bool
	}{
		{
			name:  "戦士カードに影響する",
			skill: warriorDebuffSkill,
			card: &core.BattleCard{
				CardID: "warrior_1",
				Power:  20.0,
				Type:   "warrior",
			},
			expected: true,
		},
		{
			name:  "魔法使いカードに影響しない",
			skill: warriorDebuffSkill,
			card: &core.BattleCard{
				CardID: "mage_1",
				Power:  15.0,
				Type:   "mage",
			},
			expected: false,
		},
		{
			name:  "動物カードに影響しない",
			skill: warriorDebuffSkill,
			card: &core.BattleCard{
				CardID: "beast_1",
				Power:  12.0,
				Type:   "beast",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.skill.CanAffect(tt.card)
			if result != tt.expected {
				t.Errorf("CanAffect() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestEnemySkill_Modify(t *testing.T) {
	// テスト用のEnemySkill実装
	halfPowerSkill := &mockEnemySkill{
		canAffectFunc: func(card *core.BattleCard) bool {
			return true
		},
		modifyFunc: func(card *core.BattleCard) float64 {
			return float64(card.Power) * 0.5
		},
	}

	tests := []struct {
		name     string
		skill    core.EnemySkill
		card     *core.BattleCard
		expected float64
	}{
		{
			name:  "パワーを半分にする",
			skill: halfPowerSkill,
			card: &core.BattleCard{
				CardID: "test_card",
				Power:  20.0,
				Type:   "warrior",
			},
			expected: 10.0,
		},
		{
			name:  "小数点を含むパワーを半分にする",
			skill: halfPowerSkill,
			card: &core.BattleCard{
				CardID: "test_card_2",
				Power:  15.5,
				Type:   "mage",
			},
			expected: 7.75,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.skill.Modify(tt.card)
			if result != tt.expected {
				t.Errorf("Modify() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestEnemy_WithSkills(t *testing.T) {
	skill1 := &mockEnemySkill{
		canAffectFunc: func(card *core.BattleCard) bool {
			return card.Type == "warrior"
		},
		modifyFunc: func(card *core.BattleCard) float64 {
			return float64(card.Power) * 0.8
		},
	}

	skill2 := &mockEnemySkill{
		canAffectFunc: func(card *core.BattleCard) bool {
			return card.Type == "mage"
		},
		modifyFunc: func(card *core.BattleCard) float64 {
			return float64(card.Power) * 0.6
		},
	}

	enemy := &core.Enemy{
		EnemyID:        "boss_1",
		EnemyType:      "dragon",
		Power:          50.0,
		BattleCardSlot: 2,
		Skills:         []core.EnemySkill{skill1, skill2},
	}

	if len(enemy.Skills) != 2 {
		t.Errorf("Skills length = %v, want %v", len(enemy.Skills), 2)
	}
}

// テスト用のモック実装
type mockEnemySkill struct {
	canAffectFunc func(*core.BattleCard) bool
	modifyFunc    func(*core.BattleCard) float64
}

func (m *mockEnemySkill) CanAffect(battleCard *core.BattleCard) bool {
	return m.canAffectFunc(battleCard)
}

func (m *mockEnemySkill) Modify(battleCard *core.BattleCard) float64 {
	return m.modifyFunc(battleCard)
}
