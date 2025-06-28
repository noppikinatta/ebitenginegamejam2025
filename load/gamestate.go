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
		},
		BasicYield: core.ResourceQuantity{
			Food:  1,
			Wood:  1,
			Money: 1,
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

	cardPacks := map[string]*core.CardPack{
		"cardpack-free": {
			NumPerOpen: 1,
			Ratios: map[core.CardID]int{
				"battlecard-soldier": 10,
				"battlecard-knight":  1,
			},
		},
		"cardpack-soldiers": {
			NumPerOpen: 3,
			Ratios: map[core.CardID]int{
				"battlecard-soldier": 5,
				"battlecard-archer":  5,
				"battlecard-knight":  1,
			},
		},
		"cardpack-knights": {
			NumPerOpen: 3,
			Ratios: map[core.CardID]int{
				"battlecard-knight":      3,
				"battlecard-general":     1,
				"structurecard-catapult": 2,
			},
		},
		"cardpack-politics": {
			NumPerOpen: 2,
			Ratios: map[core.CardID]int{
				"structurecard-farm":       1,
				"structurecard-woodcutter": 1,
				"structurecard-tunnel":     1,
				"structurecard-market":     1,
			},
		},
		"cardpack-war": {
			NumPerOpen: 5,
			Ratios: map[core.CardID]int{
				"structurecard-catapult": 2,
				"structurecard-ballista": 1,
				"structurecard-camp":     1,
			},
		},
		"cardpack-magic": {
			NumPerOpen: 3,
			Ratios: map[core.CardID]int{
				"battlecard-wizard":       5,
				"structurecard-mana-node": 5,
				"battlecard-mage":         1,
			},
		},
		"cardpack-mystic": {
			NumPerOpen: 3,
			Ratios: map[core.CardID]int{
				"battlecard-fortune":      2,
				"structurecard-mana-node": 2,
				"structurecard-temple":    2,
				"battlecard-mage":         1,
			},
		},
		"cardpack-mineral": {
			NumPerOpen: 2,
			Ratios: map[core.CardID]int{
				"structurecard-tunnel":  4,
				"structurecard-smelter": 1,
				"battlecard-blacksmith": 4,
			},
		},
		"cardpack-mechanical": {
			NumPerOpen: 2,
			Ratios: map[core.CardID]int{
				"battlecard-golem":       1,
				"structurecard-smelter":  2,
				"battlecard-artillery":   2,
				"structurecard-ballista": 1,
			},
		},
		"cardpack-fancy": {
			NumPerOpen: 3,
			Ratios: map[core.CardID]int{
				"battlecard-clown":    5,
				"battlecard-wrestler": 2,
				"battlecard-bard":     5,
				"battlecard-fortune":  1,
			},
		},
		"cardpack-samurai": {
			NumPerOpen: 2,
			Ratios: map[core.CardID]int{
				"battlecard-samurai": 4,
				"battlecard-ninja":   2,
				"battlecard-monk":    3,
				"structurecard-camp": 1,
			},
		},
		"cardpack-siege": {
			NumPerOpen: 2,
			Ratios: map[core.CardID]int{
				"structurecard-catapult":     2,
				"structurecard-ballista":     2,
				"structurecard-orban-cannon": 1,
				"battlecard-artillery":       2,
			},
		},
		"cardpack-finance": {
			NumPerOpen: 2,
			Ratios: map[core.CardID]int{
				"battlecard-blacksmith": 2,
				"structurecard-market":  2,
				"structurecard-mint":    1,
			},
		},
		"cardpack-building": {
			NumPerOpen: 2,
			Ratios: map[core.CardID]int{
				"structurecard-granary": 1,
				"structurecard-sawmill": 1,
				"structurecard-smelter": 1,
				"structurecard-mint":    1,
				"structurecard-temple":  1,
				"structurecard-camp":    1,
			},
		},
		"cardpack-forest": {
			NumPerOpen: 3,
			Ratios: map[core.CardID]int{
				"battlecard-archer":        2,
				"structurecard-farm":       1,
				"structurecard-woodcutter": 2,
				"structurecard-mana-node":  1,
			},
		},
		"cardpack-desert": {
			NumPerOpen: 3,
			Ratios: map[core.CardID]int{
				"battlecard-fortune":      2,
				"battlecard-bard":         1,
				"structurecard-market":    2,
				"structurecard-mana-node": 1,
			},
		},
		"cardpack-mountain": {
			NumPerOpen: 3,
			Ratios: map[core.CardID]int{
				"battlecard-blacksmith":    2,
				"structurecard-tunnel":     1,
				"structurecard-woodcutter": 2,
				"battlecard-soldier":       1,
			},
		},
	}
	cardPackPrices := map[string]core.ResourceQuantity{
		"cardpack-free":       {},
		"cardpack-soldiers":   {Food: 2, Money: 2},
		"cardpack-knights":    {Food: 5, Iron: 5, Money: 10},
		"cardpack-politics":   {Wood: 5},
		"cardpack-war":        {Iron: 20, Wood: 20, Money: 20},
		"cardpack-magic":      {Mana: 5, Food: 5},
		"cardpack-mystic":     {Mana: 20},
		"cardpack-mineral":    {Wood: 10},
		"cardpack-mechanical": {Iron: 30},
		"cardpack-fancy":      {Money: 50, Food: 10},
		"cardpack-samurai":    {Iron: 50, Food: 10},
		"cardpack-siege":      {Iron: 50, Wood: 50},
		"cardpack-finance":    {Money: 30, Wood: 10},
		"cardpack-building":   {Money: 30, Wood: 30},
		"cardpack-forest":     {Food: 5, Wood: 5},
		"cardpack-desert":     {Money: 10, Food: 10},
		"cardpack-mountain":   {Food: 5, Wood: 5},
	}

	myNation.Market = &core.Market{
		Level: 1.0,
		Items: []*core.MarketItem{
			{
				CardPack: cardPacks["cardpack-free"],
				Price:    cardPackPrices["cardpack-free"],
			},
			{
				CardPack: cardPacks["cardpack-soldiers"],
				Price:    cardPackPrices["cardpack-soldiers"],
			},
		},
	}
	points[size.Index(0, 0)] = &core.MyNationPoint{MyNation: myNation}

	// 各国
	points[size.Index(0, 2)] = &core.OtherNationPoint{ // 森の国
		OtherNation: &core.OtherNation{
			BaseNation: core.BaseNation{
				NationID: "nation-forest",
				Market: &core.Market{
					Level: 1.0,
					Items: []*core.MarketItem{
						{
							CardPack: cardPacks["cardpack-forest"],
							Price:    cardPackPrices["cardpack-forest"],
						},
					},
				},
			},
		},
	}
	points[size.Index(2, 0)] = &core.OtherNationPoint{ // 山の国
		OtherNation: &core.OtherNation{
			BaseNation: core.BaseNation{
				NationID: "nation-mountain",
				Market: &core.Market{
					Level: 1.0,
					Items: []*core.MarketItem{
						{
							CardPack: cardPacks["cardpack-mountain"],
							Price:    cardPackPrices["cardpack-mountain"],
						},
					},
				},
			},
		},
	}
	points[size.Index(2, 2)] = &core.OtherNationPoint{ // 砂漠の国
		OtherNation: &core.OtherNation{
			BaseNation: core.BaseNation{
				NationID: "nation-desert",
				Market: &core.Market{
					Level: 1.0,
					Items: []*core.MarketItem{
						{
							CardPack: cardPacks["cardpack-desert"],
							Price:    cardPackPrices["cardpack-desert"],
						},
						{
							CardPack: cardPacks["cardpack-politics"],
							Price:    cardPackPrices["cardpack-politics"],
						},
					},
				},
			},
		},
	}
	points[size.Index(0, 4)] = &core.OtherNationPoint{ // 侍の国
		OtherNation: &core.OtherNation{
			BaseNation: core.BaseNation{
				NationID: "nation-samurai",
				Market: &core.Market{
					Level: 1.0,
					Items: []*core.MarketItem{
						{
							CardPack: cardPacks["cardpack-samurai"],
							Price:    cardPackPrices["cardpack-samurai"],
						},
					},
				},
			},
		},
	}
	points[size.Index(4, 0)] = &core.OtherNationPoint{ // 魔法の国
		OtherNation: &core.OtherNation{
			BaseNation: core.BaseNation{
				NationID: "nation-magical",
				Market: &core.Market{
					Level: 1.0,
					Items: []*core.MarketItem{
						{
							CardPack: cardPacks["cardpack-magic"],
							Price:    cardPackPrices["cardpack-magic"],
						},
					},
				},
			},
		},
	}
	points[size.Index(2, 4)] = &core.OtherNationPoint{ // 機械の国
		OtherNation: &core.OtherNation{
			BaseNation: core.BaseNation{
				NationID: "nation-mechanical",
				Market: &core.Market{
					Level: 1.0,
					Items: []*core.MarketItem{
						{
							CardPack: cardPacks["cardpack-mechanical"],
							Price:    cardPackPrices["cardpack-mechanical"],
						},
					},
				},
			},
		},
	}
	points[size.Index(4, 2)] = &core.OtherNationPoint{ // お祭りの国
		OtherNation: &core.OtherNation{
			BaseNation: core.BaseNation{
				NationID: "nation-carnival",
				Market: &core.Market{
					Level: 1.0,
					Items: []*core.MarketItem{
						{
							CardPack: cardPacks["cardpack-fancy"],
							Price:    cardPackPrices["cardpack-fancy"],
						},
					},
				},
			},
		},
	}

	for i := range points {
		if points[i] != nil {
			continue
		}

		x, y := size.XY(i)

		if x == 4 && y == 4 {
			points[i] = &core.BossPoint{
				Boss:     &core.Enemy{Power: 30, BattleCardSlot: 10},
				Defeated: false,
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
