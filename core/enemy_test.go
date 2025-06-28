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
