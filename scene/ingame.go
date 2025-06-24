package scene

import (
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/bamenn"
	"github.com/noppikinatta/ebitenginegamejam2025/core"
	"github.com/noppikinatta/ebitenginegamejam2025/ui"
)

type InGame struct {
	gameState  *core.GameState
	gameUI     *ui.GameUI
	input      *ui.Input
	nextScene  ebiten.Game
	sequence   *bamenn.Sequence
	transition bamenn.Transition
}

func NewInGame(input *ui.Input) *InGame {
	// ダミーGameStateを作成（Game Jam向け簡易実装）
	gameState := createDummyGameState()

	// GameUIを初期化
	gameUI := ui.NewGameUI(gameState)

	return &InGame{
		gameState: gameState,
		gameUI:    gameUI,
		input:     input,
	}
}

func (g *InGame) Init(nextScene ebiten.Game, sequence *bamenn.Sequence, transition bamenn.Transition) {
	g.nextScene = nextScene
	g.sequence = sequence
	g.transition = transition
}

func (g *InGame) Update() error {
	// マウス位置をGameUIに設定
	mouseX, mouseY := ebiten.CursorPosition()
	g.gameUI.SetMousePosition(mouseX, mouseY)

	// GameUIの入力処理
	if err := g.gameUI.HandleInput(g.input); err != nil {
		return err
	}

	// GameUIの更新処理
	if err := g.gameUI.Update(); err != nil {
		return err
	}

	// デバッグ用：スペースキーでテストイベント追加
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.gameUI.AddHistoryEvent("Test event triggered")
	}

	return nil
}

func (g *InGame) Draw(screen *ebiten.Image) {
	// GameUIで全ての描画を行う
	g.gameUI.Draw(screen)
}

func (g *InGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 360
}

// createDummyGameState Game Jam向けのダミーGameStateを作成
func createDummyGameState() *core.GameState {
	// ダミーのMyNation作成
	myNation := &core.MyNation{
		BaseNation: core.BaseNation{
			NationID: "player_nation",
			Market: &core.Market{
				Level: 1.0,
				Items: createDummyMarketItems(),
			},
		},
		BasicYield: core.ResourceQuantity{
			Money: 10,
			Food:  5,
			Wood:  3,
			Iron:  2,
			Mana:  1,
		},
	}

	// ダミーのEnemy作成
	enemy := &core.Enemy{
		EnemyID:        "forest_orc",
		EnemyType:      "orc",
		Power:          20.0,
		BattleCardSlot: 3,
		Skills:         []core.EnemySkill{},
	}

	boss := &core.Enemy{
		EnemyID:        "ancient_dragon",
		EnemyType:      "dragon",
		Power:          100.0,
		BattleCardSlot: 5,
		Skills:         []core.EnemySkill{},
	}

	// 5x5マップのPoint配置
	points := make([]core.Point, 25)

	// 全てのインデックスにPointを配置
	for i := 0; i < 25; i++ {
		if i == 0 {
			// points[0] = MyNationPoint
			points[i] = &core.MyNationPoint{MyNation: myNation}
		} else if i%3 == 0 && i != 24 {
			// インデックスが3の倍数かつ24でない場合 = OtherNationPoint
			otherNation := &core.OtherNation{
				BaseNation: core.BaseNation{
					NationID: core.NationID("ally_nation" + strconv.Itoa(i)),
					Market: &core.Market{
						Level: 0.8,
						Items: createDummyMarketItems(),
					},
				},
			}
			points[i] = &core.OtherNationPoint{OtherNation: otherNation}
		} else if i == 24 {
			// points[24] = BossPoint
			points[i] = &core.BossPoint{
				Boss:     boss,
				Defeated: false,
			}
		} else {
			// その他 = WildernessPoint
			territory := &core.Territory{
				TerritoryID: core.TerritoryID("wilderness" + strconv.Itoa(i)),
				Cards:       []*core.StructureCard{},
				CardSlot:    3,
				BaseYield: core.ResourceQuantity{
					Money: 8,
					Food:  4,
					Wood:  2,
				},
			}
			points[i] = &core.WildernessPoint{
				Controlled: false,
				Enemy:      enemy,
				Territory:  territory,
			}
		}
	}

	// MapGrid作成
	mapGrid := &core.MapGrid{
		SizeX:  5,
		SizeY:  5,
		Points: points,
	}

	// Treasury作成
	treasury := &core.Treasury{
		Resources: core.ResourceQuantity{
			Money: 150,
			Food:  80,
			Wood:  50,
			Iron:  30,
			Mana:  20,
		},
	}

	// GameState作成
	gameState := &core.GameState{
		MyNation:    myNation,
		MapGrid:     mapGrid,
		Treasury:    treasury,
		CardDeck:    createDummyCardDeck(),
		CurrentTurn: 1,
	}

	return gameState
}

func createDummyCardDeck() *core.CardDeck {
	return &core.CardDeck{
		Cards: core.Cards{
			BattleCards: []*core.BattleCard{
				{CardID: "swordsman", Power: 10, Type: "human"},
				{CardID: "archer", Power: 8, Type: "human"},
				{CardID: "cavalry", Power: 15, Type: "human"},
			},
			StructureCards: []*core.StructureCard{
				{CardID: "wooden_wall"},
				{CardID: "watch_tower"},
			},
			ResourceCards: []*core.ResourceCard{
				{CardID: "gold_coin", ResourceQuantity: core.ResourceQuantity{Money: 10}},
				{CardID: "bread", ResourceQuantity: core.ResourceQuantity{Food: 10}},
			},
		},
	}
}

// createDummyMarketItems ダミーのMarketItem配列を作成
func createDummyMarketItems() []*core.MarketItem {
	// ダミーのCardPack作成
	cardPack1 := &core.CardPack{
		CardPackID: "basic_pack",
		Ratios: map[core.CardID]int{
			"warrior_card": 50,
			"archer_card":  30,
			"mage_card":    20,
		},
		NumPerOpen: 3,
	}

	cardPack2 := &core.CardPack{
		CardPackID: "advanced_pack",
		Ratios: map[core.CardID]int{
			"knight_card":   40,
			"wizard_card":   30,
			"dragon_card":   20,
			"artifact_card": 10,
		},
		NumPerOpen: 2,
	}

	return []*core.MarketItem{
		{
			CardPack: cardPack1,
			Price: core.ResourceQuantity{
				Money: 50,
				Food:  20,
			},
			RequiredLevel: 0.5,
		},
		{
			CardPack: cardPack2,
			Price: core.ResourceQuantity{
				Money: 100,
				Food:  30,
				Wood:  20,
			},
			RequiredLevel: 1.0,
		},
	}
}
