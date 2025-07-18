package load

import (
	"fmt"

	"github.com/noppikinatta/ebitenginegamejam2025/core"
)

// LoadGameState generates the initial game state (dummy data)
func LoadGameState() *core.GameState {
	myNation := createMyNation()
	treasury := createTreasury()
	cardGenerator := createCardGenerator()
	cardDeck := createCardDeck(cardGenerator)
	mapGrid := createMapGrid(myNation)

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
			NationID: "nation-mynation",
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

func createCardDeck(cardGenerator *core.CardGenerator) *core.CardDeck {
	firstCards, _ := cardGenerator.Generate([]core.CardID{
		"battlecard-soldier",
		"battlecard-archer",
		//"battlecard-debug",
	})

	deck := core.CardDeck{}
	deck.Add(firstCards)

	return &deck
}

func createMapGrid(myNation *core.MyNation) *core.MapGrid {
	size := core.MapGridSize{X: 5, Y: 5}
	points := make([]core.Point, size.Length())

	cardPacks := map[string]*core.CardPack{
		"cardpack-free": {
			CardPackID: "cardpack-free",
			NumPerOpen: 1,
			Ratios: map[core.CardID]int{
				"battlecard-soldier": 10,
				"battlecard-knight":  1,
			},
		},
		"cardpack-soldiers": {
			CardPackID: "cardpack-soldiers",
			NumPerOpen: 3,
			Ratios: map[core.CardID]int{
				"battlecard-soldier": 5,
				"battlecard-archer":  5,
				"battlecard-knight":  1,
			},
		},
		"cardpack-knights": {
			CardPackID: "cardpack-knights",
			NumPerOpen: 3,
			Ratios: map[core.CardID]int{
				"battlecard-knight":      3,
				"battlecard-general":     1,
				"structurecard-catapult": 2,
			},
		},
		"cardpack-politics": {
			CardPackID: "cardpack-politics",
			NumPerOpen: 2,
			Ratios: map[core.CardID]int{
				"structurecard-farm":       1,
				"structurecard-woodcutter": 1,
				"structurecard-tunnel":     1,
				"structurecard-market":     1,
			},
		},
		"cardpack-war": {
			CardPackID: "cardpack-war",
			NumPerOpen: 5,
			Ratios: map[core.CardID]int{
				"structurecard-catapult": 2,
				"structurecard-ballista": 1,
				"structurecard-camp":     1,
			},
		},
		"cardpack-magic": {
			CardPackID: "cardpack-magic",
			NumPerOpen: 3,
			Ratios: map[core.CardID]int{
				"battlecard-wizard":    5,
				"structurecard-shrine": 5,
				"battlecard-mage":      1,
			},
		},
		"cardpack-mystic": {
			CardPackID: "cardpack-mystic",
			NumPerOpen: 3,
			Ratios: map[core.CardID]int{
				"battlecard-fortune":   2,
				"structurecard-shrine": 2,
				"structurecard-temple": 2,
				"battlecard-mage":      1,
			},
		},
		"cardpack-mineral": {
			CardPackID: "cardpack-mineral",
			NumPerOpen: 2,
			Ratios: map[core.CardID]int{
				"structurecard-tunnel":  4,
				"structurecard-smelter": 1,
				"battlecard-blacksmith": 4,
			},
		},
		"cardpack-mechanical": {
			CardPackID: "cardpack-mechanical",
			NumPerOpen: 2,
			Ratios: map[core.CardID]int{
				"battlecard-golem":       1,
				"structurecard-smelter":  2,
				"battlecard-artillery":   2,
				"structurecard-ballista": 1,
			},
		},
		"cardpack-fancy": {
			CardPackID: "cardpack-fancy",
			NumPerOpen: 3,
			Ratios: map[core.CardID]int{
				"battlecard-clown":    5,
				"battlecard-wrestler": 2,
				"battlecard-bard":     5,
				"battlecard-fortune":  1,
			},
		},
		"cardpack-samurai": {
			CardPackID: "cardpack-samurai",
			NumPerOpen: 2,
			Ratios: map[core.CardID]int{
				"battlecard-samurai": 4,
				"battlecard-ninja":   2,
				"battlecard-monk":    3,
				"structurecard-camp": 1,
			},
		},
		"cardpack-siege": {
			CardPackID: "cardpack-siege",
			NumPerOpen: 2,
			Ratios: map[core.CardID]int{
				"structurecard-catapult":     2,
				"structurecard-ballista":     2,
				"structurecard-orban-cannon": 1,
				"battlecard-artillery":       2,
			},
		},
		"cardpack-finance": {
			CardPackID: "cardpack-finance",
			NumPerOpen: 2,
			Ratios: map[core.CardID]int{
				"battlecard-blacksmith": 2,
				"structurecard-market":  2,
				"structurecard-mint":    1,
			},
		},
		"cardpack-building": {
			CardPackID: "cardpack-building",
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
			CardPackID: "cardpack-forest",
			NumPerOpen: 3,
			Ratios: map[core.CardID]int{
				"battlecard-archer":        2,
				"structurecard-farm":       1,
				"structurecard-woodcutter": 2,
				"structurecard-shrine":     1,
			},
		},
		"cardpack-desert": {
			CardPackID: "cardpack-desert",
			NumPerOpen: 3,
			Ratios: map[core.CardID]int{
				"battlecard-fortune":   2,
				"battlecard-bard":      1,
				"structurecard-market": 2,
				"structurecard-shrine": 1,
			},
		},
		"cardpack-mountain": {
			CardPackID: "cardpack-mountain",
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
				CardPack:      cardPacks["cardpack-free"],
				Price:         cardPackPrices["cardpack-free"],
				RequiredLevel: 1,
			},
			{
				CardPack:      cardPacks["cardpack-soldiers"],
				Price:         cardPackPrices["cardpack-soldiers"],
				RequiredLevel: 1,
			},
			{
				CardPack:      cardPacks["cardpack-politics"],
				Price:         cardPackPrices["cardpack-politics"],
				RequiredLevel: 2,
			},
			{
				CardPack:      cardPacks["cardpack-knights"],
				Price:         cardPackPrices["cardpack-knights"],
				RequiredLevel: 3,
			},
			{
				CardPack:      cardPacks["cardpack-war"],
				Price:         cardPackPrices["cardpack-war"],
				RequiredLevel: 5,
			},
		},
	}
	points[size.Index(0, 0)] = &core.MyNationPoint{MyNation: myNation}

	// Nations
	points[size.Index(0, 2)] = &core.OtherNationPoint{ // Forest Nation
		OtherNation: &core.OtherNation{
			BaseNation: core.BaseNation{
				NationID: "nation-forest",
				Market: &core.Market{
					Level: 1.0,
					Items: []*core.MarketItem{
						{
							CardPack:      cardPacks["cardpack-forest"],
							Price:         cardPackPrices["cardpack-forest"],
							RequiredLevel: 1,
						},
						{
							CardPack:      cardPacks["cardpack-politics"],
							Price:         cardPackPrices["cardpack-politics"],
							RequiredLevel: 2,
						},
						{
							CardPack:      cardPacks["cardpack-magic"],
							Price:         cardPackPrices["cardpack-magic"],
							RequiredLevel: 4,
						},
					},
				},
			},
		},
	}
	points[size.Index(2, 0)] = &core.OtherNationPoint{ // Mountain Nation
		OtherNation: &core.OtherNation{
			BaseNation: core.BaseNation{
				NationID: "nation-mountain",
				Market: &core.Market{
					Level: 1.0,
					Items: []*core.MarketItem{
						{
							CardPack:      cardPacks["cardpack-mountain"],
							Price:         cardPackPrices["cardpack-mountain"],
							RequiredLevel: 1,
						},
						{
							CardPack:      cardPacks["cardpack-mineral"],
							Price:         cardPackPrices["cardpack-mineral"],
							RequiredLevel: 2,
						},
						{
							CardPack:      cardPacks["cardpack-siege"],
							Price:         cardPackPrices["cardpack-siege"],
							RequiredLevel: 4,
						},
					},
				},
			},
		},
	}
	points[size.Index(2, 2)] = &core.OtherNationPoint{ // Desert Nation
		OtherNation: &core.OtherNation{
			BaseNation: core.BaseNation{
				NationID: "nation-desert",
				Market: &core.Market{
					Level: 1.0,
					Items: []*core.MarketItem{
						{
							CardPack:      cardPacks["cardpack-desert"],
							Price:         cardPackPrices["cardpack-desert"],
							RequiredLevel: 1,
						},
						{
							CardPack:      cardPacks["cardpack-politics"],
							Price:         cardPackPrices["cardpack-politics"],
							RequiredLevel: 1,
						},
						{
							CardPack:      cardPacks["cardpack-finance"],
							Price:         cardPackPrices["cardpack-finance"],
							RequiredLevel: 3,
						},
						{
							CardPack:      cardPacks["cardpack-building"],
							Price:         cardPackPrices["cardpack-building"],
							RequiredLevel: 5,
						},
					},
				},
			},
		},
	}
	points[size.Index(0, 4)] = &core.OtherNationPoint{ // Samurai Nation
		OtherNation: &core.OtherNation{
			BaseNation: core.BaseNation{
				NationID: "nation-samurai",
				Market: &core.Market{
					Level: 1.0,
					Items: []*core.MarketItem{
						{
							CardPack:      cardPacks["cardpack-samurai"],
							Price:         cardPackPrices["cardpack-samurai"],
							RequiredLevel: 1,
						},
						{
							CardPack:      cardPacks["cardpack-mineral"],
							Price:         cardPackPrices["cardpack-mineral"],
							RequiredLevel: 3,
						},
						{
							CardPack:      cardPacks["cardpack-war"],
							Price:         cardPackPrices["cardpack-war"],
							RequiredLevel: 4,
						},
					},
				},
			},
		},
	}
	points[size.Index(4, 0)] = &core.OtherNationPoint{ // Magical Nation
		OtherNation: &core.OtherNation{
			BaseNation: core.BaseNation{
				NationID: "nation-magical",
				Market: &core.Market{
					Level: 1.0,
					Items: []*core.MarketItem{
						{
							CardPack:      cardPacks["cardpack-magic"],
							Price:         cardPackPrices["cardpack-magic"],
							RequiredLevel: 1,
						},
						{
							CardPack:      cardPacks["cardpack-mystic"],
							Price:         cardPackPrices["cardpack-mystic"],
							RequiredLevel: 3,
						},
					},
				},
			},
		},
	}
	points[size.Index(2, 4)] = &core.OtherNationPoint{ // Mechanical Nation
		OtherNation: &core.OtherNation{
			BaseNation: core.BaseNation{
				NationID: "nation-mechanical",
				Market: &core.Market{
					Level: 1.0,
					Items: []*core.MarketItem{
						{
							CardPack:      cardPacks["cardpack-mechanical"],
							Price:         cardPackPrices["cardpack-mechanical"],
							RequiredLevel: 1,
						},
						{
							CardPack:      cardPacks["cardpack-siege"],
							Price:         cardPackPrices["cardpack-siege"],
							RequiredLevel: 3,
						},
						{
							CardPack:      cardPacks["cardpack-building"],
							Price:         cardPackPrices["cardpack-building"],
							RequiredLevel: 4,
						},
					},
				},
			},
		},
	}
	points[size.Index(4, 2)] = &core.OtherNationPoint{ // Carnival Nation
		OtherNation: &core.OtherNation{
			BaseNation: core.BaseNation{
				NationID: "nation-carnival",
				Market: &core.Market{
					Level: 1.0,
					Items: []*core.MarketItem{
						{
							CardPack:      cardPacks["cardpack-fancy"],
							Price:         cardPackPrices["cardpack-fancy"],
							RequiredLevel: 1,
						},
						{
							CardPack:      cardPacks["cardpack-finance"],
							Price:         cardPackPrices["cardpack-finance"],
							RequiredLevel: 2,
						},
					},
				},
			},
		},
	}

	// Place WildernessPoint and BossPoint
	wildernessConfigs := []struct {
		x, y        int
		enemyID     core.EnemyID
		enemyType   core.EnemyType
		power       float64
		cardSlot    int
		skills      []core.EnemySkill
		terrainType string
		baseYield   core.ResourceQuantity
	}{
		{1, 0, "enemy-goblin", "enemy-type-demonic", 3, 3, []core.EnemySkill{}, "terrain-forest", core.ResourceQuantity{Wood: 2}},
		{0, 1, "enemy-sabrelouse", "enemy-type-animal", 4, 3, []core.EnemySkill{}, "terrain-mountain", core.ResourceQuantity{Iron: 2}},
		{1, 1, "enemy-rattlesnake", "enemy-type-dragon", 6, 3, []core.EnemySkill{}, "terrain-plain", core.ResourceQuantity{Food: 2}},
		{2, 1, "enemy-condor", "enemy-type-flying", 6, 3, []core.EnemySkill{createEvasionSkill()}, "terrain-desert", core.ResourceQuantity{}},
		{1, 2, "enemy-slime", "enemy-type-unknown", 6, 3, []core.EnemySkill{createSoftSkill()}, "terrain-desert", core.ResourceQuantity{}},
		{0, 3, "enemy-crocodile", "enemy-type-dragon", 10, 4, []core.EnemySkill{}, "terrain-mountain", core.ResourceQuantity{Iron: 2}},
		{3, 0, "enemy-grizzly", "enemy-type-animal", 12, 4, []core.EnemySkill{}, "terrain-forest", core.ResourceQuantity{Wood: 2}},
		{1, 3, "enemy-skeleton", "enemy-type-undead", 12, 4, []core.EnemySkill{createLongbowSkill()}, "terrain-mana-node", core.ResourceQuantity{Mana: 3}},
		{3, 1, "enemy-elemental", "enemy-type-unknown", 20, 5, []core.EnemySkill{createIncorporealitySkill()}, "terrain-mountain", core.ResourceQuantity{Iron: 2}},
		{1, 4, "enemy-dragon", "enemy-type-dragon", 30, 6, []core.EnemySkill{createPressureSkill()}, "terrain-plain", core.ResourceQuantity{Food: 2}},
		{4, 1, "enemy-griffin", "enemy-type-flying", 25, 5, []core.EnemySkill{createEvasionSkill()}, "terrain-plain", core.ResourceQuantity{Food: 2}},
		{2, 3, "enemy-vampire", "enemy-type-undead", 30, 7, []core.EnemySkill{createCharmSkill()}, "terrain-forest", core.ResourceQuantity{Wood: 2}},
		{3, 2, "enemy-living-armor", "enemy-type-unknown", 50, 7, []core.EnemySkill{}, "terrain-mana-node", core.ResourceQuantity{Mana: 3}},
		{3, 4, "enemy-arc-demon", "enemy-type-demonic", 40, 8, []core.EnemySkill{createMagicBarrierSkill()}, "terrain-forest", core.ResourceQuantity{Wood: 2}},
		{4, 3, "enemy-durendal", "enemy-type-undead", 40, 8, []core.EnemySkill{createSideAttackSkill()}, "terrain-mountain", core.ResourceQuantity{Iron: 2}},
		{3, 3, "enemy-obelisk", "enemy-type-unknown", 40, 8, []core.EnemySkill{createLaserSkill()}, "terrain-mana-node", core.ResourceQuantity{Mana: 3}},
	}

	for i, config := range wildernessConfigs {
		enemy := &core.Enemy{
			EnemyID:        config.enemyID,
			EnemyType:      config.enemyType,
			Power:          config.power,
			Skills:         config.skills,
			BattleCardSlot: config.cardSlot,
			Question:       fmt.Sprintf("question-%d", i+1),
		}

		territory := &core.Territory{
			BaseYield: config.baseYield,
			CardSlot:  3, // Set all to 3
		}

		points[size.Index(config.x, config.y)] = &core.WildernessPoint{
			TerrainType: config.terrainType,
			Controlled:  false,
			Enemy:       enemy,
			Territory:   territory,
		}
	}

	// Boss Point (4,4)
	boss := &core.Enemy{
		EnemyID:        "enemy-final-boss",
		EnemyType:      "enemy-type-demonic",
		Power:          60,
		Skills:         []core.EnemySkill{createWaveSkill()},
		BattleCardSlot: 9,
		Question:       "question-boss",
	}
	points[size.Index(4, 4)] = &core.BossPoint{
		Boss:     boss,
		Defeated: false,
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
			CardID:    "battlecard-debug",
			BasePower: 999,
			Skill: &core.BattleCardSkill{
				BattleCardSkillID: "battlecardskill-debug",
				DescriptionKey:    "battlecardskill-debug-desc",
				Calculator: &core.BattleCardSkillCalculatorEffectSelf{
					Effect: &core.BattleCardSkillEffect{
						Modifier: &core.BattleCardPowerModifier{
							ProtectionFromDebuff: 999.0,
						},
					},
				},
			},
			Type: "cardtype-str",
		},
		{
			CardID:    "battlecard-soldier",
			BasePower: 3,
			Skill: &core.BattleCardSkill{
				BattleCardSkillID: "battlecardskill-cooperation",
				DescriptionKey:    "battlecardskill-cooperation-desc",
				Calculator: &core.BattleCardSkillCalculatorEffectSelf{
					Effect: &core.BattleCardSkillEffect{
						Modifier: &core.BattleCardPowerModifier{
							BuffBoostedPower: 0.5,
						},
					},
				},
			},
			Type: "cardtype-str",
		},
		{
			CardID:    "battlecard-knight",
			BasePower: 4,
			Skill: &core.BattleCardSkill{
				BattleCardSkillID: "battlecardskill-dragon-killer",
				DescriptionKey:    "battlecardskill-dragon-killer-desc",
				Calculator: &core.BattleCardSkillCalculatorCondition{
					Condition: func(options *core.BattleCardSkillCalculationOptions) bool {
						return options.Enemy.EnemyType == "enemy-type-dragon"
					},
					Calculator: &core.BattleCardSkillCalculatorEffectSelf{
						Effect: &core.BattleCardSkillEffect{
							Modifier: &core.BattleCardPowerModifier{
								MultiplicativeBuff: 1.0,
							},
						},
					},
				},
			},
			Type: "cardtype-str",
		},
		{
			CardID:    "battlecard-general",
			BasePower: 4,
			Skill: &core.BattleCardSkill{
				BattleCardSkillID: "battlecardskill-command",
				DescriptionKey:    "battlecardskill-command-desc",
				Calculator: &core.BattleCardSkillCalculatorEffectAllCondition{
					Condition: func(idx int, options *core.BattleCardSkillCalculationOptions) bool {
						return idx > options.BattleCardIndex
					},
					Effect: &core.BattleCardSkillEffect{
						Modifier: &core.BattleCardPowerModifier{
							MultiplicativeBuff: 0.2,
						},
					},
				},
			},
			Type: "cardtype-str",
		},
		{
			CardID:    "battlecard-archer",
			BasePower: 3,
			Skill: &core.BattleCardSkill{
				BattleCardSkillID: "battlecardskill-sniper",
				DescriptionKey:    "battlecardskill-sniper-desc",
				Calculator: &core.BattleCardSkillCalculatorCondition{
					Condition: func(options *core.BattleCardSkillCalculationOptions) bool {
						if options.Enemy.EnemyType == "enemy-type-animal" {
							return true
						}
						if options.Enemy.EnemyType == "enemy-type-flying" {
							return true
						}
						return false
					},
					Calculator: &core.BattleCardSkillCalculatorEffectSelf{
						Effect: &core.BattleCardSkillEffect{
							Modifier: &core.BattleCardPowerModifier{
								MultiplicativeBuff: 1.0,
							},
						},
					},
				},
			},
			Type: "cardtype-agi",
		},
		{
			CardID:    "battlecard-fortune",
			BasePower: 1,
			Skill: &core.BattleCardSkill{
				BattleCardSkillID: "battlecardskill-forecast",
				DescriptionKey:    "battlecardskill-forecast-desc",
				Calculator: &core.BattleCardSkillCalculatorEffectAll{
					Effect: &core.BattleCardSkillEffect{
						Modifier: &core.BattleCardPowerModifier{
							ProtectionFromDebuff: 0.2,
						},
					},
				},
			},
			Type: "cardtype-mag",
		},
		{
			CardID:    "battlecard-wizard",
			BasePower: 2,
			Skill: &core.BattleCardSkill{
				BattleCardSkillID: "battlecardskill-long-spell",
				DescriptionKey:    "battlecardskill-long-spell-desc",
				Calculator: core.BattleCardSkillCalculationFunc(func(options *core.BattleCardSkillCalculationOptions) {
					options.BattleCardPowerModifiers[options.BattleCardIndex].AdditiveBuff = float64(options.BattleCardIndex)
				}),
			},
			Type: "cardtype-mag",
		},
		{
			CardID:    "battlecard-mage",
			BasePower: 2,
			Skill: &core.BattleCardSkill{
				BattleCardSkillID: "battlecardskill-magic-amplifier",
				DescriptionKey:    "battlecardskill-magic-amplifier-desc",
				Calculator: &core.BattleCardSkillCalculatorEffectAllCondition{
					Condition: func(idx int, options *core.BattleCardSkillCalculationOptions) bool {
						card := options.BattleCards[idx]
						return card.Type == "cardtype-mag"
					},
					Effect: &core.BattleCardSkillEffect{
						Modifier: &core.BattleCardPowerModifier{
							MultiplicativeBuff: 0.3,
						},
					},
				},
			},
			Type: "cardtype-mag",
		},
		{
			CardID:    "battlecard-blacksmith",
			BasePower: 2,
			Skill: &core.BattleCardSkill{
				BattleCardSkillID: "battlecardskill-weapon-enhancement",
				DescriptionKey:    "battlecardskill-weapon-enhancement-desc",
				Calculator: &core.BattleCardSkillCalculatorEffectAllCondition{
					Condition: func(idx int, options *core.BattleCardSkillCalculationOptions) bool {
						card := options.BattleCards[idx]
						return card.Type == "cardtype-str"
					},
					Effect: &core.BattleCardSkillEffect{
						Modifier: &core.BattleCardPowerModifier{
							MultiplicativeBuff: 0.3,
						},
					},
				},
			},
			Type: "cardtype-str",
		},
		{
			CardID:    "battlecard-samurai",
			BasePower: 5,
			Skill: &core.BattleCardSkill{
				BattleCardSkillID: "battlecardskill-bushido",
				DescriptionKey:    "battlecardskill-bushido-desc",
				Calculator: &core.BattleCardSkillCalculatorCondition{
					Condition: func(options *core.BattleCardSkillCalculationOptions) bool {
						return options.BattleCardIndex == 0
					},
					Calculator: &core.BattleCardSkillCalculatorEffectSelf{
						Effect: &core.BattleCardSkillEffect{
							Modifier: &core.BattleCardPowerModifier{
								MultiplicativeBuff: 1.0,
							},
						},
					},
				},
			},
			Type: "cardtype-str",
		},
		{
			CardID:    "battlecard-ninja",
			BasePower: 5,
			Skill: &core.BattleCardSkill{
				BattleCardSkillID: "battlecardskill-stealth",
				DescriptionKey:    "battlecardskill-stealth-desc",
				Calculator: &core.BattleCardSkillCalculatorEffectSelf{
					Effect: &core.BattleCardSkillEffect{
						Modifier: &core.BattleCardPowerModifier{
							ProtectionFromDebuff: 1.0,
						},
					},
				},
			},
			Type: "cardtype-agi",
		},
		{
			CardID:    "battlecard-monk",
			BasePower: 4,
			Skill: &core.BattleCardSkill{
				BattleCardSkillID: "battlecardskill-ki",
				DescriptionKey:    "battlecardskill-ki-desc",
				Calculator: &core.BattleCardSkillCalculatorCondition{
					Condition: func(options *core.BattleCardSkillCalculationOptions) bool {
						return options.Enemy.EnemyType == "enemy-type-undead"
					},
					Calculator: &core.BattleCardSkillCalculatorEffectSelf{
						Effect: &core.BattleCardSkillEffect{
							Modifier: &core.BattleCardPowerModifier{
								MultiplicativeBuff: 1.0,
							},
						},
					},
				},
			},
			Type: "cardtype-mag",
		},
		{
			CardID:    "battlecard-bard",
			BasePower: 1,
			Skill: &core.BattleCardSkill{
				BattleCardSkillID: "battlecardskill-support",
				DescriptionKey:    "battlecardskill-support-desc",
				Calculator: &core.BattleCardSkillCalculatorEffectAll{
					Effect: &core.BattleCardSkillEffect{
						Modifier: &core.BattleCardPowerModifier{
							AdditiveBuff: 1,
						},
					},
				},
			},
			Type: "cardtype-agi",
		},
		{
			CardID:    "battlecard-artillery",
			BasePower: 2,
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
			CardID:    "battlecard-clown",
			BasePower: 1,
			Skill: &core.BattleCardSkill{
				BattleCardSkillID: "battlecardskill-viper-master",
				DescriptionKey:    "battlecardskill-viper-master-desc",
				Calculator: &core.BattleCardSkillCalculatorEffectIdxs{
					IdxDeltas: []int{-1, 1},
					Effect: &core.BattleCardSkillEffect{
						Modifier: &core.BattleCardPowerModifier{
							ProtectionFromDebuff: 1.0,
						},
					},
				},
			},
			Type: "cardtype-agi",
		},
		{
			CardID:    "battlecard-wrestler",
			BasePower: 7,
			Skill: &core.BattleCardSkill{
				BattleCardSkillID: "battlecardskill-two-platoon",
				DescriptionKey:    "battlecardskill-two-platoon-desc",
				Calculator: &core.BattleCardSkillCalculatorCondition{
					Condition: func(options *core.BattleCardSkillCalculationOptions) bool {
						idx := options.BattleCardIndex + 1
						if idx < 0 || idx >= len(options.BattleCards) {
							return false
						}

						return options.BattleCards[idx].Type == "cardtype-str"
					},
					Calculator: &core.BattleCardSkillCalculatorEffectIdxs{
						IdxDeltas: []int{0, 1},
						Effect: &core.BattleCardSkillEffect{
							Modifier: &core.BattleCardPowerModifier{
								MultiplicativeBuff: 1.0,
							},
						},
					},
				},
			},
			Type: "cardtype-str",
		},
		{
			CardID:    "battlecard-golem",
			BasePower: 9,
			Skill:     nil,
			Type:      "cardtype-str",
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
					Money: 5,
				},
			},
		},
		{
			CardID:         "structurecard-shrine",
			DescriptionKey: "structurecard-shrine-desc",
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

// Helper functions for generating enemy skills

func createEvasionSkill() core.EnemySkill {
	// Strength type card power -2
	return &core.EnemySkillImpl{
		IDField: "enemy-skill-evasion",
		Condition: func(idx int, options *core.EnemySkillCalculationOptions) bool {
			card := options.BattleCards[idx]
			return card.Type == "cardtype-str"
		},
		Modifier: &core.BattleCardPowerModifier{
			AdditiveDebuff: 2.0,
		},
	}
}

func createSoftSkill() core.EnemySkill {
	// Non-Magic type card power -50%
	return &core.EnemySkillImpl{
		IDField: "enemy-skill-soft",
		Condition: func(idx int, options *core.EnemySkillCalculationOptions) bool {
			card := options.BattleCards[idx]
			return card.Type != "cardtype-mag"
		},
		Modifier: &core.BattleCardPowerModifier{
			MultiplicativeDebuff: 0.5,
		},
	}
}

func createLongbowSkill() core.EnemySkill {
	// Rearmost card power -100%
	return &core.EnemySkillImpl{
		IDField: "enemy-skill-longbow",
		Condition: func(idx int, options *core.EnemySkillCalculationOptions) bool {
			return idx == len(options.BattleCards)-1
		},
		Modifier: &core.BattleCardPowerModifier{
			MultiplicativeDebuff: 1.0,
		},
	}
}

func createIncorporealitySkill() core.EnemySkill {
	// Non-Magic type card power -100%
	return &core.EnemySkillImpl{
		IDField: "enemy-skill-incorporeality",
		Condition: func(idx int, options *core.EnemySkillCalculationOptions) bool {
			card := options.BattleCards[idx]
			return card.Type != "cardtype-mag"
		},
		Modifier: &core.BattleCardPowerModifier{
			MultiplicativeDebuff: 1.0,
		},
	}
}

func createPressureSkill() core.EnemySkill {
	// All card power -1
	return &core.EnemySkillImpl{
		IDField: "enemy-skill-pressure",
		Condition: func(idx int, options *core.EnemySkillCalculationOptions) bool {
			return true
		},
		Modifier: &core.BattleCardPowerModifier{
			AdditiveDebuff: 1.0,
		},
	}
}

func createCharmSkill() core.EnemySkill {
	// First 3 cards power -100%
	return &core.EnemySkillImpl{
		IDField: "enemy-skill-charm",
		Condition: func(idx int, options *core.EnemySkillCalculationOptions) bool {
			return idx < 3
		},
		Modifier: &core.BattleCardPowerModifier{
			MultiplicativeDebuff: 1.0,
		},
	}
}

func createMagicBarrierSkill() core.EnemySkill {
	// Magic type card power -100%
	return &core.EnemySkillImpl{
		IDField: "enemy-skill-magic-barrier",
		Condition: func(idx int, options *core.EnemySkillCalculationOptions) bool {
			card := options.BattleCards[idx]
			return card.Type == "cardtype-mag"
		},
		Modifier: &core.BattleCardPowerModifier{
			MultiplicativeDebuff: 1.0,
		},
	}
}

func createLaserSkill() core.EnemySkill {
	// Last 3 cards power -100%
	return &core.EnemySkillImpl{
		IDField: "enemy-skill-laser",
		Condition: func(idx int, options *core.EnemySkillCalculationOptions) bool {
			return idx >= len(options.BattleCards)-3
		},
		Modifier: &core.BattleCardPowerModifier{
			MultiplicativeDebuff: 1.0,
		},
	}
}

func createSideAttackSkill() core.EnemySkill {
	// First 5 cards power -50%
	return &core.EnemySkillImpl{
		IDField: "enemy-skill-side-attack",
		Condition: func(idx int, options *core.EnemySkillCalculationOptions) bool {
			return idx < 5
		},
		Modifier: &core.BattleCardPowerModifier{
			MultiplicativeDebuff: 0.5,
		},
	}
}

func createWaveSkill() core.EnemySkill {
	// All card power -2
	return &core.EnemySkillImpl{
		IDField: "enemy-skill-wave",
		Condition: func(idx int, options *core.EnemySkillCalculationOptions) bool {
			return true
		},
		Modifier: &core.BattleCardPowerModifier{
			AdditiveDebuff: 2.0,
		},
	}
}
