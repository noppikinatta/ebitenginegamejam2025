package load

import (
	"github.com/noppikinatta/ebitenginegamejam2025/core"
)

// LoadGameState generates the initial game state (dummy data)
func LoadGameState() *core.GameState {
	myNation := createMyNation()
	treasury := createTreasury()
	cardGenerator, cardDisplayOrder := createCardGenerator()
	cardDeck := createCardDeck(cardGenerator)
	cardPacks, cardPackPrices := createCardPacksAndPrices()
	markets := createMarkets(cardPacks, cardPackPrices)
	mapGrid := createMapGrid(myNation, cardPacks, cardPackPrices)

	gs := &core.GameState{
		MyNation:         myNation,
		CardDeck:         cardDeck,
		MapGrid:          mapGrid,
		Treasury:         treasury,
		CurrentTurn:      0,
		CardGenerator:    cardGenerator,
		Markets:          markets,
		CardDisplayOrder: cardDisplayOrder,
	}

	return gs
}

func createMyNation() *core.MyNation {
	return core.NewMyNation("nation-mynation", "My Nation")
}

func createTreasury() *core.Treasury {
	return &core.Treasury{}
}

func createCardDeck(cardGenerator *core.CardGenerator) *core.CardDeck {
	deck := core.NewCardDeck()

	// Add initial cards
	deck.Add("battlecard-soldier")
	deck.Add("battlecard-archer")
	// deck.Add("battlecard-debug") // commented out

	return deck
}

func createCardPacksAndPrices() (map[string]*core.CardPack, map[string]core.ResourceQuantity) {
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

	return cardPacks, cardPackPrices
}

func createMapGrid(myNation *core.MyNation, cardPacks map[string]*core.CardPack, cardPackPrices map[string]core.ResourceQuantity) *core.MapGrid {
	size := core.MapGridSize{X: 5, Y: 5}
	points := make([]core.Point, size.Length())

	points[size.Index(0, 0)] = &core.MyNationPoint{MyNation: myNation}

	// Nations
	points[size.Index(0, 2)] = &core.OtherNationPoint{ // Forest Nation
		OtherNation: core.NewOtherNation("nation-forest", "Forest Nation"),
	}

	points[size.Index(2, 0)] = &core.OtherNationPoint{ // Mountain Nation
		OtherNation: core.NewOtherNation("nation-mountain", "Mountain Nation"),
	}

	points[size.Index(2, 2)] = &core.OtherNationPoint{ // Desert Nation
		OtherNation: core.NewOtherNation("nation-desert", "Desert Nation"),
	}

	points[size.Index(0, 4)] = &core.OtherNationPoint{ // Samurai Nation
		OtherNation: core.NewOtherNation("nation-samurai", "Samurai Nation"),
	}

	points[size.Index(4, 0)] = &core.OtherNationPoint{ // Magical Nation
		OtherNation: core.NewOtherNation("nation-magical", "Magical Nation"),
	}

	points[size.Index(2, 4)] = &core.OtherNationPoint{ // Mechanical Nation
		OtherNation: core.NewOtherNation("nation-mechanical", "Mechanical Nation"),
	}

	points[size.Index(4, 2)] = &core.OtherNationPoint{ // Carnival Nation
		OtherNation: core.NewOtherNation("nation-carnival", "Carnival Nation"),
	}

	// Place WildernessPoint and BossPoint
	wildernessConfigs := []struct {
		x, y        int
		enemyID     core.EnemyID
		enemyType   core.EnemyType
		power       float64
		cardSlot    int
		skills      []*core.EnemySkill
		terrainType string
		baseYield   core.ResourceQuantity
	}{
		{1, 0, "enemy-goblin", "enemy-type-demonic", 3, 3, []*core.EnemySkill{}, "terrain-forest", core.ResourceQuantity{Wood: 2}},
		{0, 1, "enemy-sabrelouse", "enemy-type-animal", 4, 3, []*core.EnemySkill{}, "terrain-mountain", core.ResourceQuantity{Iron: 2}},
		{1, 1, "enemy-rattlesnake", "enemy-type-dragon", 6, 3, []*core.EnemySkill{}, "terrain-plain", core.ResourceQuantity{Food: 2}},
		{2, 1, "enemy-condor", "enemy-type-flying", 6, 3, []*core.EnemySkill{createEvasionSkill()}, "terrain-desert", core.ResourceQuantity{}},
		{1, 2, "enemy-slime", "enemy-type-unknown", 6, 3, []*core.EnemySkill{createSoftSkill()}, "terrain-desert", core.ResourceQuantity{}},
		{0, 3, "enemy-crocodile", "enemy-type-dragon", 10, 4, []*core.EnemySkill{}, "terrain-mountain", core.ResourceQuantity{Iron: 2}},
		{3, 0, "enemy-grizzly", "enemy-type-animal", 12, 4, []*core.EnemySkill{}, "terrain-forest", core.ResourceQuantity{Wood: 2}},
		{1, 3, "enemy-skeleton", "enemy-type-undead", 12, 4, []*core.EnemySkill{createLongbowSkill()}, "terrain-mana-node", core.ResourceQuantity{Mana: 3}},
		{3, 1, "enemy-elemental", "enemy-type-unknown", 20, 5, []*core.EnemySkill{createIncorporealitySkill()}, "terrain-mountain", core.ResourceQuantity{Iron: 2}},
		{1, 4, "enemy-dragon", "enemy-type-dragon", 30, 6, []*core.EnemySkill{createPressureSkill()}, "terrain-plain", core.ResourceQuantity{Food: 2}},
		{4, 1, "enemy-griffin", "enemy-type-flying", 25, 5, []*core.EnemySkill{createEvasionSkill()}, "terrain-plain", core.ResourceQuantity{Food: 2}},
		{2, 3, "enemy-vampire", "enemy-type-undead", 30, 7, []*core.EnemySkill{createCharmSkill()}, "terrain-forest", core.ResourceQuantity{Wood: 2}},
		{3, 2, "enemy-living-armor", "enemy-type-unknown", 50, 7, []*core.EnemySkill{}, "terrain-mana-node", core.ResourceQuantity{Mana: 3}},
		{3, 4, "enemy-arc-demon", "enemy-type-demonic", 40, 8, []*core.EnemySkill{createMagicBarrierSkill()}, "terrain-forest", core.ResourceQuantity{Wood: 2}},
		{4, 3, "enemy-durendal", "enemy-type-undead", 40, 8, []*core.EnemySkill{createSideAttackSkill()}, "terrain-mountain", core.ResourceQuantity{Iron: 2}},
		{3, 3, "enemy-obelisk", "enemy-type-unknown", 40, 8, []*core.EnemySkill{createLaserSkill()}, "terrain-mana-node", core.ResourceQuantity{Mana: 3}},
	}

	for _, config := range wildernessConfigs {
		enemy := core.NewEnemy(
			core.EnemyID(config.enemyID),
			core.EnemyType(config.enemyType),
			config.power,
			config.skills,
			config.cardSlot,
		)

		terrain := core.NewTerrain(
			core.TerrainID(config.terrainType),
			config.baseYield,
			3, // Set all to 3
		)
		territory := core.NewTerritory(
			core.TerritoryID("territory-"+config.enemyID),
			terrain,
		)

		wilderness := &core.WildernessPoint{}
		wilderness.SetControlledForTest(false)
		wilderness.SetEnemyForTest(enemy)
		wilderness.SetTerritoryForTest(territory)

		points[size.Index(config.x, config.y)] = wilderness
	}

	// Boss Point (4,4)
	boss := core.NewEnemy(
		"enemy-final-boss",
		"enemy-type-demonic",
		60,
		[]*core.EnemySkill{createWaveSkill()},
		9,
	)
	bossPoint := &core.BossPoint{}
	bossPoint.SetBossForTest(boss)
	bossPoint.SetDefeatedForTest(false)
	points[size.Index(4, 4)] = bossPoint

	mapGrid := &core.MapGrid{
		Size:   size,
		Points: points,
	}
	mapGrid.UpdateAccesibles()

	return mapGrid
}

func createMarkets(cardPacks map[string]*core.CardPack, cardPackPrices map[string]core.ResourceQuantity) map[core.NationID]*core.Market {
	markets := make(map[core.NationID]*core.Market)

	// MyNation Market
	markets["nation-mynation"] = &core.Market{
		Level: 1.0,
		Items: []*core.MarketItem{
			core.NewMarketItem(cardPacks["cardpack-free"], cardPackPrices["cardpack-free"], 1, 0),
			core.NewMarketItem(cardPacks["cardpack-soldiers"], cardPackPrices["cardpack-soldiers"], 1, 0),
			core.NewMarketItem(cardPacks["cardpack-politics"], cardPackPrices["cardpack-politics"], 2, 0),
			core.NewMarketItem(cardPacks["cardpack-knights"], cardPackPrices["cardpack-knights"], 3, 0),
			core.NewMarketItem(cardPacks["cardpack-war"], cardPackPrices["cardpack-war"], 5, 0),
		},
	}

	// Forest Nation Market
	markets["nation-forest"] = &core.Market{
		Level: 1.0,
		Items: []*core.MarketItem{
			core.NewMarketItem(cardPacks["cardpack-forest"], cardPackPrices["cardpack-forest"], 1, 0),
			core.NewMarketItem(cardPacks["cardpack-politics"], cardPackPrices["cardpack-politics"], 2, 0),
			core.NewMarketItem(cardPacks["cardpack-magic"], cardPackPrices["cardpack-magic"], 4, 0),
		},
	}

	// Mountain Nation Market
	markets["nation-mountain"] = &core.Market{
		Level: 1.0,
		Items: []*core.MarketItem{
			core.NewMarketItem(cardPacks["cardpack-mountain"], cardPackPrices["cardpack-mountain"], 1, 0),
			core.NewMarketItem(cardPacks["cardpack-mineral"], cardPackPrices["cardpack-mineral"], 2, 0),
			core.NewMarketItem(cardPacks["cardpack-siege"], cardPackPrices["cardpack-siege"], 4, 0),
		},
	}

	// Desert Nation Market
	markets["nation-desert"] = &core.Market{
		Level: 1.0,
		Items: []*core.MarketItem{
			core.NewMarketItem(cardPacks["cardpack-desert"], cardPackPrices["cardpack-desert"], 1, 0),
			core.NewMarketItem(cardPacks["cardpack-politics"], cardPackPrices["cardpack-politics"], 1, 0),
			core.NewMarketItem(cardPacks["cardpack-finance"], cardPackPrices["cardpack-finance"], 3, 0),
			core.NewMarketItem(cardPacks["cardpack-building"], cardPackPrices["cardpack-building"], 5, 0),
		},
	}

	// Samurai Nation Market
	markets["nation-samurai"] = &core.Market{
		Level: 1.0,
		Items: []*core.MarketItem{
			core.NewMarketItem(cardPacks["cardpack-samurai"], cardPackPrices["cardpack-samurai"], 1, 0),
			core.NewMarketItem(cardPacks["cardpack-mineral"], cardPackPrices["cardpack-mineral"], 3, 0),
			core.NewMarketItem(cardPacks["cardpack-war"], cardPackPrices["cardpack-war"], 4, 0),
		},
	}

	// Magical Nation Market
	markets["nation-magical"] = &core.Market{
		Level: 1.0,
		Items: []*core.MarketItem{
			core.NewMarketItem(cardPacks["cardpack-magic"], cardPackPrices["cardpack-magic"], 1, 0),
			core.NewMarketItem(cardPacks["cardpack-mystic"], cardPackPrices["cardpack-mystic"], 3, 0),
		},
	}

	// Mechanical Nation Market
	markets["nation-mechanical"] = &core.Market{
		Level: 1.0,
		Items: []*core.MarketItem{
			core.NewMarketItem(cardPacks["cardpack-mechanical"], cardPackPrices["cardpack-mechanical"], 1, 0),
			core.NewMarketItem(cardPacks["cardpack-siege"], cardPackPrices["cardpack-siege"], 3, 0),
			core.NewMarketItem(cardPacks["cardpack-building"], cardPackPrices["cardpack-building"], 4, 0),
		},
	}

	// Carnival Nation Market
	markets["nation-carnival"] = &core.Market{
		Level: 1.0,
		Items: []*core.MarketItem{
			core.NewMarketItem(cardPacks["cardpack-fancy"], cardPackPrices["cardpack-fancy"], 1, 0),
			core.NewMarketItem(cardPacks["cardpack-finance"], cardPackPrices["cardpack-finance"], 2, 0),
		},
	}

	return markets
}

func createCardGenerator() (*core.CardGenerator, []core.CardID) {
	battleCards := createBattleCards()
	structureCards := createStructureCards()

	// 表示順序を作成（BattleCard → StructureCard の順）
	var displayOrder []core.CardID

	// createBattleCardsで定義されている順序で追加
	battleCardOrder := []core.CardID{
		"battlecard-debug",
		"battlecard-soldier",
		"battlecard-knight",
		"battlecard-general",
		"battlecard-archer",
		"battlecard-fortune",
		"battlecard-wizard",
		"battlecard-mage",
		"battlecard-blacksmith",
		"battlecard-samurai",
		"battlecard-ninja",
		"battlecard-monk",
		"battlecard-bard",
		"battlecard-artillery",
		"battlecard-clown",
		"battlecard-wrestler",
		"battlecard-golem",
	}
	displayOrder = append(displayOrder, battleCardOrder...)

	// createStructureCardsで定義されている順序で追加
	structureCardOrder := []core.CardID{
		"structurecard-farm",
		"structurecard-woodcutter",
		"structurecard-tunnel",
		"structurecard-market",
		"structurecard-shrine",
		"structurecard-granary",
		"structurecard-sawmill",
		"structurecard-smelter",
		"structurecard-mint",
		"structurecard-temple",
		"structurecard-camp",
		"structurecard-catapult",
		"structurecard-ballista",
		"structurecard-orban-cannon",
	}
	displayOrder = append(displayOrder, structureCardOrder...)

	cardGenerator := &core.CardGenerator{
		BattleCards:    battleCards,
		StructureCards: structureCards,
	}

	return cardGenerator, displayOrder
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
						return options.Enemy.Type() == "enemy-type-dragon"
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
						if options.Enemy.Type() == "enemy-type-animal" {
							return true
						}
						if options.Enemy.Type() == "enemy-type-flying" {
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
						return options.Enemy.Type() == "enemy-type-undead"
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
		// Yield Additive系
		core.NewStructureCard(
			"structurecard-farm",
			core.ResourceQuantity{Food: 2},
			core.NewResourceModifier(),
			0.0, 0),
		core.NewStructureCard(
			"structurecard-woodcutter",
			core.ResourceQuantity{Wood: 2},
			core.NewResourceModifier(),
			0.0, 0),
		core.NewStructureCard(
			"structurecard-tunnel",
			core.ResourceQuantity{Iron: 2},
			core.NewResourceModifier(),
			0.0, 0),
		core.NewStructureCard(
			"structurecard-market",
			core.ResourceQuantity{Money: 5},
			core.NewResourceModifier(),
			0.0, 0),
		core.NewStructureCard(
			"structurecard-shrine",
			core.ResourceQuantity{Mana: 2},
			core.NewResourceModifier(),
			0.0, 0),

		// Yield Multiplicative系
		core.NewStructureCard(
			"structurecard-granary",
			core.ResourceQuantity{},
			core.ResourceModifier{Food: 0.5},
			0.0, 0),
		core.NewStructureCard(
			"structurecard-sawmill",
			core.ResourceQuantity{},
			core.ResourceModifier{Wood: 0.5},
			0.0, 0),
		core.NewStructureCard(
			"structurecard-smelter",
			core.ResourceQuantity{},
			core.ResourceModifier{Iron: 0.5},
			0.0, 0),
		core.NewStructureCard(
			"structurecard-mint",
			core.ResourceQuantity{},
			core.ResourceModifier{Money: 0.5},
			0.0, 0),
		core.NewStructureCard(
			"structurecard-temple",
			core.ResourceQuantity{},
			core.ResourceModifier{Mana: 0.5},
			0.0, 0),

		// Support CardSlot系
		core.NewStructureCard(
			"structurecard-camp",
			core.ResourceQuantity{},
			core.NewResourceModifier(),
			0.0, 1),

		// Support Power系
		core.NewStructureCard(
			"structurecard-catapult",
			core.ResourceQuantity{},
			core.NewResourceModifier(),
			3.0, 0),
		core.NewStructureCard(
			"structurecard-ballista",
			core.ResourceQuantity{},
			core.NewResourceModifier(),
			5.0, 0),
		core.NewStructureCard(
			"structurecard-orban-cannon",
			core.ResourceQuantity{},
			core.NewResourceModifier(),
			8.0, 0),
	}

	cardMap := make(map[core.CardID]*core.StructureCard)
	for _, card := range cards {
		cardMap[card.ID()] = card
	}
	return cardMap
}

// Helper functions for generating enemy skills

func createEvasionSkill() *core.EnemySkill {
	// Strength type card power -2
	return core.NewEnemySkill(
		"enemy-skill-evasion",
		func(idx int, options *core.EnemySkillCalculationOptions) bool {
			card := options.BattleCards[idx]
			return card.Type == "cardtype-str"
		},
		&core.BattleCardPowerModifier{
			AdditiveDebuff: 2.0,
		},
	)
}

func createSoftSkill() *core.EnemySkill {
	// Non-Magic type card power -50%
	return core.NewEnemySkill(
		"enemy-skill-soft",
		func(idx int, options *core.EnemySkillCalculationOptions) bool {
			card := options.BattleCards[idx]
			return card.Type != "cardtype-mag"
		},
		&core.BattleCardPowerModifier{
			MultiplicativeDebuff: 0.5,
		},
	)
}

func createLongbowSkill() *core.EnemySkill {
	// Rearmost card power -100%
	return core.NewEnemySkill(
		"enemy-skill-longbow",
		func(idx int, options *core.EnemySkillCalculationOptions) bool {
			return idx == len(options.BattleCards)-1
		},
		&core.BattleCardPowerModifier{
			MultiplicativeDebuff: 1.0,
		},
	)
}

func createIncorporealitySkill() *core.EnemySkill {
	// Non-Magic type card power -100%
	return core.NewEnemySkill(
		"enemy-skill-incorporeality",
		func(idx int, options *core.EnemySkillCalculationOptions) bool {
			card := options.BattleCards[idx]
			return card.Type != "cardtype-mag"
		},
		&core.BattleCardPowerModifier{
			MultiplicativeDebuff: 1.0,
		},
	)
}

func createPressureSkill() *core.EnemySkill {
	// All card power -1
	return core.NewEnemySkill(
		"enemy-skill-pressure",
		func(idx int, options *core.EnemySkillCalculationOptions) bool {
			return true
		},
		&core.BattleCardPowerModifier{
			AdditiveDebuff: 1.0,
		},
	)
}

func createCharmSkill() *core.EnemySkill {
	// First 3 cards power -100%
	return core.NewEnemySkill(
		"enemy-skill-charm",
		func(idx int, options *core.EnemySkillCalculationOptions) bool {
			return idx < 3
		},
		&core.BattleCardPowerModifier{
			MultiplicativeDebuff: 1.0,
		},
	)
}

func createMagicBarrierSkill() *core.EnemySkill {
	// Magic type card power -100%
	return core.NewEnemySkill(
		"enemy-skill-magic-barrier",
		func(idx int, options *core.EnemySkillCalculationOptions) bool {
			card := options.BattleCards[idx]
			return card.Type == "cardtype-mag"
		},
		&core.BattleCardPowerModifier{
			MultiplicativeDebuff: 1.0,
		},
	)
}

func createLaserSkill() *core.EnemySkill {
	// Last 3 cards power -100%
	return core.NewEnemySkill(
		"enemy-skill-laser",
		func(idx int, options *core.EnemySkillCalculationOptions) bool {
			return idx >= len(options.BattleCards)-3
		},
		&core.BattleCardPowerModifier{
			MultiplicativeDebuff: 1.0,
		},
	)
}

func createSideAttackSkill() *core.EnemySkill {
	// First 5 cards power -50%
	return core.NewEnemySkill(
		"enemy-skill-side-attack",
		func(idx int, options *core.EnemySkillCalculationOptions) bool {
			return idx < 5
		},
		&core.BattleCardPowerModifier{
			MultiplicativeDebuff: 0.5,
		},
	)
}

func createWaveSkill() *core.EnemySkill {
	// All card power -2
	return core.NewEnemySkill(
		"enemy-skill-wave",
		func(idx int, options *core.EnemySkillCalculationOptions) bool {
			return true
		},
		&core.BattleCardPowerModifier{
			AdditiveDebuff: 2.0,
		},
	)
}
