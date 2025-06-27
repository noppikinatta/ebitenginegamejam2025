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
		ResourceCards:  createResourceCards(),
	}
}

func createBattleCards() map[core.CardID]*core.BattleCard {
	cards := []*core.BattleCard{
		{
			CardID: "battle_card_1",
			Power:  3,
			Skill:  nil,
			Type:   "warrior",
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
			CardID: "structure_card_1",
			YieldModifier: &core.MultiplyYieldModifier{
				Multiply: 1.5,
			},
			BattlefieldModifier: &core.CardSlotBattlefieldModifier{
				Value: 1,
			},
		},
	}

	cardMap := make(map[core.CardID]*core.StructureCard)
	for _, card := range cards {
		cardMap[card.CardID] = card
	}
	return cardMap
}

func createResourceCards() map[core.CardID]*core.ResourceCard {
	cards := []*core.ResourceCard{
		{
			CardID: "resource_card_1",
			ResourceQuantity: core.ResourceQuantity{
				Food: 1,
			},
		},
	}

	cardMap := make(map[core.CardID]*core.ResourceCard)
	for _, card := range cards {
		cardMap[card.CardID] = card
	}
	return cardMap
}
