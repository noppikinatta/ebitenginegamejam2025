package core_test

import (
	"testing"

	"github.com/noppikinatta/ebitenginegamejam2025/core"
)

func TestBattlefield_CanBeat(t *testing.T) {
	// Test enemy
	weakEnemy := &core.Enemy{
		EnemyID:        "weak_orc",
		EnemyType:      "orc",
		Power:          10.0,
		BattleCardSlot: 2,
		Skills:         []core.EnemySkill{},
	}

	strongEnemy := &core.Enemy{
		EnemyID:        "strong_dragon",
		EnemyType:      "dragon",
		Power:          50.0,
		BattleCardSlot: 3,
		Skills:         []core.EnemySkill{},
	}

	// Test battle card
	weakCard := &core.BattleCard{
		CardID:    "weak_warrior",
		BasePower: 5.0,
		Type:      "warrior",
	}

	strongCard := &core.BattleCard{
		CardID:    "strong_mage",
		BasePower: 30.0,
		Type:      "mage",
	}

	tests := []struct {
		name        string
		battlefield *core.Battlefield
		expected    bool
	}{
		{
			name: "When power is insufficient",
			battlefield: &core.Battlefield{
				Enemy:            weakEnemy,
				BattleCards:      []*core.BattleCard{weakCard},
				BaseSupportPower: 0.0,
			},
			expected: false,
		},
		{
			name: "When power is exactly equal",
			battlefield: &core.Battlefield{
				Enemy:            weakEnemy,
				BattleCards:      []*core.BattleCard{weakCard, weakCard}, // 5 + 5 = 10
				BaseSupportPower: 0.0,
			},
			expected: true,
		},
		{
			name: "When power is sufficient",
			battlefield: &core.Battlefield{
				Enemy:            weakEnemy,
				BattleCards:      []*core.BattleCard{strongCard},
				BaseSupportPower: 0.0,
			},
			expected: true,
		},
		{
			name: "Win with support power",
			battlefield: &core.Battlefield{
				Enemy:            weakEnemy,
				BattleCards:      []*core.BattleCard{weakCard}, // 5 + 6 = 11 > 10
				BaseSupportPower: 6.0,
			},
			expected: true,
		},
		{
			name: "Insufficient power against strong enemy",
			battlefield: &core.Battlefield{
				Enemy:            strongEnemy,
				BattleCards:      []*core.BattleCard{strongCard}, // 30 < 50
				BaseSupportPower: 0.0,
			},
			expected: false,
		},
		{
			name: "Defeat strong enemy with multiple cards and support power",
			battlefield: &core.Battlefield{
				Enemy:            strongEnemy,
				BattleCards:      []*core.BattleCard{strongCard, strongCard}, // 30 + 30 + 5 = 65 > 50
				BaseSupportPower: 5.0,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.battlefield.CanBeat()
			if result != tt.expected {
				t.Errorf("CanBeat() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestBattlefield_PowerCalculation(t *testing.T) {
	enemy := &core.Enemy{
		EnemyID:        "calc_test_enemy",
		EnemyType:      "goblin",
		Power:          25.0,
		BattleCardSlot: 3,
		Skills:         []core.EnemySkill{},
	}

	card1 := &core.BattleCard{
		CardID:    "card1",
		BasePower: 10.0,
		Type:      "warrior",
	}

	card2 := &core.BattleCard{
		CardID:    "card2",
		BasePower: 8.0,
		Type:      "mage",
	}

	card3 := &core.BattleCard{
		CardID:    "card3",
		BasePower: 7.5,
		Type:      "beast",
	}

	tests := []struct {
		name         string
		cards        []*core.BattleCard
		supportPower float64
		canBeat      bool
	}{
		{
			name:         "Insufficient with a single card",
			cards:        []*core.BattleCard{card1}, // 10 < 25
			supportPower: 0.0,
			canBeat:      false,
		},
		{
			name:         "Insufficient with two cards",
			cards:        []*core.BattleCard{card1, card2}, // 10 + 8 = 18 < 25
			supportPower: 0.0,
			canBeat:      false,
		},
		{
			name:         "Win with three cards",
			cards:        []*core.BattleCard{card1, card2, card3}, // 10 + 8 + 7.5 = 25.5 >= 25
			supportPower: 0.0,
			canBeat:      true,
		},
		{
			name:         "Win with support power",
			cards:        []*core.BattleCard{card1, card2}, // 10 + 8 + 10 = 28 >= 25
			supportPower: 10.0,
			canBeat:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			battlefield := &core.Battlefield{
				Enemy:            enemy,
				BattleCards:      tt.cards,
				BaseSupportPower: tt.supportPower,
			}

			result := battlefield.CanBeat()
			if result != tt.canBeat {
				t.Errorf("CanBeat() = %v, want %v", result, tt.canBeat)
			}
		})
	}
}
