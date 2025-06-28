package load

import "github.com/noppikinatta/ebitenginegamejam2025/core"

// LoadGameState ゲームの初期状態を生成する（ダミーデータ）
func LoadGameState() *core.GameState {
	myNation := createMyNation()
	treasury := createTreasury()
	cardDeck := createCardDeck()
	mapGrid := createMapGrid(myNation)
	cardGenerator := createCardGenerator()

	gs := &core.GameState{
		MyNation:      myNation,
		CardDeck:      cardDeck,
		MapGrid:       mapGrid,
		Treasury:      treasury,
		CurrentTurn:   0,
		CardGenerator: cardGenerator,
	}

	return gs
}

func createMyNation() *core.MyNation {
	return &core.MyNation{
		BaseNation: core.BaseNation{
			NationID: "player",
			Market: &core.Market{
				Level: 1.0,
				Items: []*core.MarketItem{
					{
						CardPack: &core.CardPack{
							Ratios: map[core.CardID]int{
								"battle_card_1": 1,
							},
							NumPerOpen: 2,
						},
						Price: core.ResourceQuantity{Money: 0},
					},
				},
			},
		},
		BasicYield: core.ResourceQuantity{
			Food:  1,
			Wood:  1,
			Money: 3,
			Iron:  0,
			Mana:  0,
		},
	}
}

func createTreasury() *core.Treasury {
	return &core.Treasury{}
}

func createCardDeck() *core.CardDeck {
	return &core.CardDeck{}
}

func createMapGrid(myNation *core.MyNation) *core.MapGrid {
	size := core.MapGridSize{X: 5, Y: 5}
	points := make([]core.Point, size.Length())

	for i := range points {
		x, y := size.XY(i)

		if x == 0 && y == 0 {
			points[i] = &core.MyNationPoint{MyNation: myNation}
			continue
		}

		if x == 4 && y == 4 {
			points[i] = &core.BossPoint{
				Boss:     &core.Enemy{Power: 30, BattleCardSlot: 10},
				Defeated: false,
			}
			continue
		}

		if x%2 == 0 && y%2 == 0 {
			points[i] = &core.OtherNationPoint{
				OtherNation: &core.OtherNation{
					BaseNation: core.BaseNation{
						NationID: "other",
						Market: &core.Market{
							Level: 1.0,
							Items: []*core.MarketItem{
								{
									CardPack: &core.CardPack{
										Ratios: map[core.CardID]int{
											"battle_card_1":    1,
											"structure_card_1": 1,
										},
										NumPerOpen: 2,
									},
									Price: core.ResourceQuantity{Money: 3},
								},
							},
						},
					},
				},
			}
			continue
		}

		points[i] = &core.WildernessPoint{
			Controlled: false,
			Enemy:      &core.Enemy{Power: 10, BattleCardSlot: 5},
			Territory:  &core.Territory{BaseYield: core.ResourceQuantity{Food: 5}, CardSlot: 3},
		}
	}

	mapGrid := &core.MapGrid{
		Size:   size,
		Points: points,
	}
	mapGrid.UpdateAccesibles()

	return mapGrid
}

func createCardGenerator() *core.CardGenerator {
	return &core.CardGenerator{
		BattleCards:    createBattleCards(),
		StructureCards: createStructureCards(),
	}
}

func createBattleCards() map[core.CardID]*core.BattleCard {
	cards := []*core.BattleCard{
		{
			CardID: "battlecard-soldier",
			Power:  3,
			Skill: &core.BattleCardSkill{
				BattleCardSkillID: "battlecardskill-cooperation",
				DescriptionKey:    "battlecardskill-cooperation-desc",
				Calculator: &core.BattleCardSkillCalculatorBoostBuff{
					BoostBuff: 2.0,
				},
			},
			Type: "cardtype-str",
		},
		{
			CardID: "battlecard-knight",
			Power:  3,
			Skill: &core.BattleCardSkill{
				BattleCardSkillID: "battlecardskill-dragon-killer",
				DescriptionKey:    "battlecardskill-dragon-killer-desc",
				Calculator: &core.BattleCardSkillCalculatorEnemyType{
					EnemyType:  "enemy-type-dragon",
					Multiplier: 2.0,
				},
			},
			Type: "cardtype-str",
		},
		{
			CardID: "battlecard-general",
			Power:  3,
			Skill: &core.BattleCardSkill{
				BattleCardSkillID: "battlecardskill-command",
				DescriptionKey:    "battlecardskill-command-desc",
				Calculator: &core.BattleCardSkillCalculatorTrailings{
					Multiplier: 0.2,
				},
			},
			Type: "cardtype-str",
		},
		{
			CardID: "battlecard-archer",
			Power:  3,
			Skill: &core.BattleCardSkill{
				BattleCardSkillID: "battlecardskill-sniper",
				DescriptionKey:    "battlecardskill-sniper-desc",
				Calculator: &core.BattleCardSkillCalculatorComposite{
					Calculators: []core.BattleCardSkillCalculator{
						&core.BattleCardSkillCalculatorEnemyType{
							EnemyType:  "enemy-type-animal",
							Multiplier: 2.0,
						},
						&core.BattleCardSkillCalculatorEnemyType{
							EnemyType:  "enemy-type-flying",
							Multiplier: 2.0,
						},
					},
				},
			},
			Type: "cardtype-agi",
		},
		{
			CardID: "battlecard-fortune",
			Power:  3,
			Skill: &core.BattleCardSkill{
				BattleCardSkillID: "battlecardskill-forecast",
				DescriptionKey:    "battlecardskill-forecast-desc",
				Calculator: &core.BattleCardSkillCalculatorProofBuff{
					Value: 0.2,
				},
			},
			Type: "cardtype-mag",
		},
		{
			CardID: "battlecard-wizard",
			Power:  3,
			Skill: &core.BattleCardSkill{
				BattleCardSkillID: "battlecardskill-long-spell",
				DescriptionKey:    "battlecardskill-long-spell-desc",
				Calculator:        core.AddingByIndexBattleCardSkillCalculator,
			},
			Type: "cardtype-mag",
		},
		{
			CardID: "battlecard-mage",
			Power:  3,
			Skill: &core.BattleCardSkill{
				BattleCardSkillID: "battlecardskill-magic-amplifier",
				DescriptionKey:    "battlecardskill-magic-amplifier-desc",
				Calculator: &core.BattleCardSkillCalculatorAllByCardType{
					CardType:   "cardtype-mag",
					Multiplier: 0.5,
				},
			},
			Type: "cardtype-mag",
		},
		{
			CardID: "battlecard-blacksmith",
			Power:  3,
			Skill: &core.BattleCardSkill{
				BattleCardSkillID: "battlecardskill-weapon-enhancement",
				DescriptionKey:    "battlecardskill-weapon-enhancement-desc",
				Calculator: &core.BattleCardSkillCalculatorAllByCardType{
					CardType:   "cardtype-str",
					Multiplier: 0.5,
				},
			},
			Type: "cardtype-str",
		},
		{
			CardID: "battlecard-samurai",
			Power:  3,
			Skill: &core.BattleCardSkill{
				BattleCardSkillID: "battlecardskill-bushido",
				DescriptionKey:    "battlecardskill-bushido-desc",
				Calculator: &core.BattleCardSkillCalculatorByIdx{
					Index:      0,
					Multiplier: 1.0,
				},
			},
			Type: "cardtype-str",
		},
		{
			CardID: "battlecard-ninja",
			Power:  3,
			Skill: &core.BattleCardSkill{
				BattleCardSkillID: "battlecardskill-stealth",
				DescriptionKey:    "battlecardskill-stealth-desc",
				Calculator: &core.BattleCardSkillCalculatorProofBuff{
					Value: 1.0,
				},
			},
			Type: "cardtype-agi",
		},
		{
			CardID: "battlecard-monk",
			Power:  3,
			Skill: &core.BattleCardSkill{
				BattleCardSkillID: "battlecardskill-ki",
				DescriptionKey:    "battlecardskill-ki-desc",
				Calculator: &core.BattleCardSkillCalculatorEnemyType{
					EnemyType:  "enemy-type-undead",
					Multiplier: 2.0,
				},
			},
			Type: "cardtype-str",
		},
		{
			CardID: "battlecard-bard",
			Power:  3,
			Skill: &core.BattleCardSkill{
				BattleCardSkillID: "battlecardskill-support",
				DescriptionKey:    "battlecardskill-support-desc",
				Calculator: &core.BattleCardSkillCalculatorAll{
					ModifierFunc: func(modifier *core.BattleCardPowerModifier) {
						modifier.AdditiveBuff += 1
					},
				},
			},
			Type: "cardtype-agi",
		},
		{
			CardID: "battlecard-artillery",
			Power:  3,
			Skill: &core.BattleCardSkill{
				BattleCardSkillID: "battlecardskill-shooting-observation",
				DescriptionKey:    "battlecardskill-shooting-observation-desc",
				Calculator: &core.BattleCardSkillCalculatorSupportPowerMultiplier{
					Multiplier: 1.0,
				},
			},
			Type: "cardtype-str",
		},
		{
			CardID: "battlecard-clown",
			Power:  3,
			Skill: &core.BattleCardSkill{
				BattleCardSkillID: "battlecardskill-viper-master",
				DescriptionKey:    "battlecardskill-viper-master-desc",
				Calculator: &core.BattleCardSkillCalculatorProofDebufNeighboring{
					Value: 1.0,
				},
			},
			Type: "cardtype-agi",
		},
		{
			CardID: "battlecard-wrestler",
			Power:  3,
			Skill: &core.BattleCardSkill{
				BattleCardSkillID: "battlecardskill-two-platoon",
				DescriptionKey:    "battlecardskill-two-platoon-desc",
				Calculator: &core.BattleCardSkillCalculatorTwoPlatoon{
					Multiplier: 1.0,
					CardType:   "cardtype-str",
				},
			},
			Type: "cardtype-str",
		},
		{
			CardID: "battlecard-golem",
			Power:  5,
			Skill:  nil,
			Type:   "cardtype-str",
		},
	}

	cardMap := make(map[core.CardID]*core.BattleCard)
	for _, card := range cards {
		cardMap[card.CardID] = card
	}
	return cardMap
}

func createStructureCards() map[core.CardID]*core.StructureCard {
	cards := []*core.StructureCard{
		{
			CardID:         "structurecard-farm",
			DescriptionKey: "structurecard-farm-desc",
			YieldModifier: &core.AddYieldModifier{
				ResourceQuantity: core.ResourceQuantity{
					Food: 2,
				},
			},
		},
		{
			CardID:         "structurecard-woodcutter",
			DescriptionKey: "structurecard-woodcutter-desc",
			YieldModifier: &core.AddYieldModifier{
				ResourceQuantity: core.ResourceQuantity{
					Wood: 2,
				},
			},
		},
		{
			CardID:         "structurecard-tunnel",
			DescriptionKey: "structurecard-tunnel-desc",
			YieldModifier: &core.AddYieldModifier{
				ResourceQuantity: core.ResourceQuantity{
					Iron: 2,
				},
			},
		},
		{
			CardID:         "structurecard-market",
			DescriptionKey: "structurecard-market-desc",
			YieldModifier: &core.AddYieldModifier{
				ResourceQuantity: core.ResourceQuantity{
					Money: 2,
				},
			},
		},
		{
			CardID:         "structurecard-mana-node",
			DescriptionKey: "structurecard-mana-node-desc",
			YieldModifier: &core.AddYieldModifier{
				ResourceQuantity: core.ResourceQuantity{
					Mana: 2,
				},
			},
		},
		{
			CardID:         "structurecard-granary",
			DescriptionKey: "structurecard-granary-desc",
			YieldModifier: &core.MultiplyYieldModifier{
				FoodMultiply: 0.5,
			},
		},
		{
			CardID:         "structurecard-sawmill",
			DescriptionKey: "structurecard-sawmill-desc",
			YieldModifier: &core.MultiplyYieldModifier{
				WoodMultiply: 0.5,
			},
		},
		{
			CardID:         "structurecard-smelter",
			DescriptionKey: "structurecard-smelter-desc",
			YieldModifier: &core.MultiplyYieldModifier{
				IronMultiply: 0.5,
			},
		},
		{
			CardID:         "structurecard-mint",
			DescriptionKey: "structurecard-mint-desc",
			YieldModifier: &core.MultiplyYieldModifier{
				MoneyMultiply: 0.5,
			},
		},
		{
			CardID:         "structurecard-temple",
			DescriptionKey: "structurecard-temple-desc",
			YieldModifier: &core.MultiplyYieldModifier{
				ManaMultiply: 0.5,
			},
		},
		{
			CardID:              "structurecard-camp",
			DescriptionKey:      "structurecard-camp-desc",
			BattlefieldModifier: &core.CardSlotBattlefieldModifier{Value: 1},
		},
		{
			CardID:              "structurecard-catapult",
			DescriptionKey:      "structurecard-catapult-desc",
			BattlefieldModifier: &core.SupportPowerBattlefieldModifier{Value: 3},
		},
		{
			CardID:              "structurecard-ballista",
			DescriptionKey:      "structurecard-ballista-desc",
			BattlefieldModifier: &core.SupportPowerBattlefieldModifier{Value: 5},
		},
		{
			CardID:              "structurecard-orban-cannon",
			DescriptionKey:      "structurecard-orban-cannon-desc",
			BattlefieldModifier: &core.SupportPowerBattlefieldModifier{Value: 8},
		},
	}

	cardMap := make(map[core.CardID]*core.StructureCard)
	for _, card := range cards {
		cardMap[card.CardID] = card
	}
	return cardMap
}
