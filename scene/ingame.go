package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/core"
	"github.com/noppikinatta/ebitenginegamejam2025/ui"
	"github.com/noppikinatta/nyuuryoku"
)

type InGame struct {
	gameState *core.GameState
	gameUI    *ui.GameUI
	mouse     *nyuuryoku.Mouse
}

func NewInGame() *InGame {
	// ダミーGameStateを作成（Game Jam向け簡易実装）
	gameState := createDummyGameState()

	// GameUIを初期化
	gameUI := ui.NewGameUI(gameState)

	// nyuuryoku Mouseを初期化
	mouse := nyuuryoku.NewMouse()

	return &InGame{
		gameState: gameState,
		gameUI:    gameUI,
		mouse:     mouse,
	}
}

func (g *InGame) Update() error {
	// マウス位置をGameUIに設定
	mouseX, mouseY := ebiten.CursorPosition()
	g.gameUI.SetMousePosition(mouseX, mouseY)

	// Input構造体を作成
	input := &ui.Input{
		Mouse: g.mouse,
	}

	// GameUIの入力処理
	if err := g.gameUI.HandleInput(input); err != nil {
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
		Nation: core.Nation{
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

	// ダミーのOtherNation作成
	otherNation := &core.OtherNation{
		Nation: core.Nation{
			NationID: "ally_nation",
			Market: &core.Market{
				Level: 0.8,
				Items: createDummyMarketItems(),
			},
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

	// ダミーのTerritory作成
	territory := &core.Territory{
		TerritoryID: "forest_territory",
		Cards:       []*core.StructureCard{},
		CardSlot:    3,
		BaseYield: core.ResourceQuantity{
			Money: 8,
			Food:  4,
			Wood:  2,
		},
	}

	// 5x5マップのPoint配置
	points := make([]core.Point, 25)

	// 中央にMyNationPoint
	points[12] = &core.MyNationPoint{MyNation: myNation} // (2, 2)

	// 隣接位置にOtherNationPoint
	points[7] = &core.OtherNationPoint{OtherNation: otherNation} // (2, 1)

	// 制圧済みWildernessPoint
	points[11] = &core.WildernessPoint{ // (1, 2)
		Controlled: true,
		Enemy:      enemy,
		Territory:  territory,
	}

	// 未制圧WildernessPoint
	points[13] = &core.WildernessPoint{ // (3, 2)
		Controlled: false,
		Enemy:      enemy,
		Territory:  territory,
	}

	// BossPoint
	points[24] = &core.BossPoint{ // (4, 4)
		Boss:     boss,
		Defeated: false,
	}

	// 残りはnilのまま（何もないマス）

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
		CurrentTurn: 1,
	}

	return gameState
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
