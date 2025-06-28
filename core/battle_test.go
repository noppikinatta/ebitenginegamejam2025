package core_test

import (
	"testing"

	"github.com/noppikinatta/ebitenginegamejam2025/core"
)

func TestBattlefield_CanBeat(t *testing.T) {
	// テスト用の敵
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

	// テスト用のバトルカード
	weakCard := &core.BattleCard{
		CardID: "weak_warrior",
		Power:  5.0,
		Type:   "warrior",
	}

	strongCard := &core.BattleCard{
		CardID: "strong_mage",
		Power:  30.0,
		Type:   "mage",
	}

	tests := []struct {
		name        string
		battlefield *core.Battlefield
		expected    bool
	}{
		{
			name: "戦闘力が足りない場合",
			battlefield: &core.Battlefield{
				Enemy:            weakEnemy,
				BattleCards:      []*core.BattleCard{weakCard},
				BaseSupportPower: 0.0,
			},
			expected: false,
		},
		{
			name: "戦闘力がちょうど等しい場合",
			battlefield: &core.Battlefield{
				Enemy:            weakEnemy,
				BattleCards:      []*core.BattleCard{weakCard, weakCard}, // 5 + 5 = 10
				BaseSupportPower: 0.0,
			},
			expected: true,
		},
		{
			name: "戦闘力が十分な場合",
			battlefield: &core.Battlefield{
				Enemy:            weakEnemy,
				BattleCards:      []*core.BattleCard{strongCard},
				BaseSupportPower: 0.0,
			},
			expected: true,
		},
		{
			name: "サポートパワーで勝利",
			battlefield: &core.Battlefield{
				Enemy:            weakEnemy,
				BattleCards:      []*core.BattleCard{weakCard}, // 5 + 6 = 11 > 10
				BaseSupportPower: 6.0,
			},
			expected: true,
		},
		{
			name: "強敵に対して戦闘力不足",
			battlefield: &core.Battlefield{
				Enemy:            strongEnemy,
				BattleCards:      []*core.BattleCard{strongCard}, // 30 < 50
				BaseSupportPower: 0.0,
			},
			expected: false,
		},
		{
			name: "複数カードとサポートパワーで強敵に勝利",
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

func TestBattlefield_Beat(t *testing.T) {
	enemy := &core.Enemy{
		EnemyID:        "test_enemy",
		EnemyType:      "orc",
		Power:          20.0,
		BattleCardSlot: 2,
		Skills:         []core.EnemySkill{},
	}

	card := &core.BattleCard{
		CardID: "test_card",
		Power:  25.0,
		Type:   "warrior",
	}

	battlefield := &core.Battlefield{
		Enemy:            enemy,
		BattleCards:      []*core.BattleCard{card},
		BaseSupportPower: 0.0,
	}

	// 勝利可能かチェック
	if !battlefield.CanBeat() {
		t.Fatal("Expected to be able to beat enemy")
	}

	// Beat()メソッドを呼び出す
	battlefield.Beat()

	// Beat()は成功時に何も返さないので、エラーが発生しないことを確認
	// 将来的に戦闘結果やログを返すようになるかもしれない
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
		CardID: "card1",
		Power:  10.0,
		Type:   "warrior",
	}

	card2 := &core.BattleCard{
		CardID: "card2",
		Power:  8.0,
		Type:   "mage",
	}

	card3 := &core.BattleCard{
		CardID: "card3",
		Power:  7.5,
		Type:   "beast",
	}

	tests := []struct {
		name         string
		cards        []*core.BattleCard
		supportPower float64
		canBeat      bool
	}{
		{
			name:         "単体カードで不足",
			cards:        []*core.BattleCard{card1}, // 10 < 25
			supportPower: 0.0,
			canBeat:      false,
		},
		{
			name:         "2枚で不足",
			cards:        []*core.BattleCard{card1, card2}, // 10 + 8 = 18 < 25
			supportPower: 0.0,
			canBeat:      false,
		},
		{
			name:         "3枚で勝利",
			cards:        []*core.BattleCard{card1, card2, card3}, // 10 + 8 + 7.5 = 25.5 >= 25
			supportPower: 0.0,
			canBeat:      true,
		},
		{
			name:         "サポートパワーで勝利",
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
