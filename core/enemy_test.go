package core_test

import (
	"testing"

	"github.com/noppikinatta/ebitenginegamejam2025/core"
)

func TestEnemy_Basic(t *testing.T) {
	enemy := core.NewEnemy("test_enemy_1", "orc", 25.0, []*core.EnemySkill{}, 3)

	if enemy.ID() != "test_enemy_1" {
		t.Errorf("ID() = %v, want %v", enemy.ID(), "test_enemy_1")
	}
	if enemy.Type() != "orc" {
		t.Errorf("Type() = %v, want %v", enemy.Type(), "orc")
	}
	if enemy.Power() != 25.0 {
		t.Errorf("Power() = %v, want %v", enemy.Power(), 25.0)
	}
	if enemy.BattleCardSlot() != 3 {
		t.Errorf("BattleCardSlot() = %v, want %v", enemy.BattleCardSlot(), 3)
	}
}
