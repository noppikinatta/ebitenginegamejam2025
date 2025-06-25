package load

import "github.com/noppikinatta/ebitenginegamejam2025/core"

// LoadGameState ゲームの初期状態を生成する（ダミーデータ）
func LoadGameState() *core.GameState {
	myNation := &core.MyNation{
		BaseNation: core.BaseNation{
			NationID: "player",
			Market:   &core.Market{},
		},
		BasicYield: core.ResourceQuantity{Food: 10, Wood: 10},
	}
	treasury := &core.Treasury{
		Resources: core.ResourceQuantity{Food: 100, Wood: 100, Money: 50},
	}
	cardDeck := &core.CardDeck{
		Cards: core.Cards{}, // 最初はカードなし
	}

	sizeX, sizeY := 5, 5
	points := make([]core.Point, sizeX*sizeY)

	// (0,0) に自国
	points[0*sizeX+0] = &core.MyNationPoint{MyNation: myNation}

	// (1,0) と (0,1) に未制圧の荒れ地
	points[0*sizeX+1] = &core.WildernessPoint{
		Controlled: false,
		Enemy:      &core.Enemy{Power: 10},
		Territory:  &core.Territory{BaseYield: core.ResourceQuantity{Food: 5}},
	}
	points[1*sizeX+0] = &core.WildernessPoint{
		Controlled: false,
		Enemy:      &core.Enemy{Power: 12},
		Territory:  &core.Territory{BaseYield: core.ResourceQuantity{Food: 3, Wood: 3}},
	}
	// (1,1) に制圧済みの荒れ地
	points[1*sizeX+1] = &core.WildernessPoint{
		Controlled: true,
		Territory:  &core.Territory{BaseYield: core.ResourceQuantity{Money: 2}},
	}

	// (4,4) にボス
	points[4*sizeX+4] = &core.BossPoint{
		Boss:     &core.Enemy{Power: 100},
		Defeated: false,
	}

	// (2,2) に他国
	points[2*sizeX+2] = &core.OtherNationPoint{
		OtherNation: &core.OtherNation{BaseNation: core.BaseNation{NationID: "enemy1"}},
	}

	// 残りは空の荒れ地（操作不可）
	for i := range points {
		if points[i] == nil {
			// WildernessPointのPassableはControlledに依存する。
			// また、Enemyがnilだと戦闘が発生しないので、Territoryもnilにしておく。
			points[i] = &core.WildernessPoint{Controlled: false, Enemy: nil, Territory: nil}
		}
	}

	mapGrid := &core.MapGrid{
		SizeX:  sizeX,
		SizeY:  sizeY,
		Points: points,
	}
	mapGrid.UpdateAccesibles()

	gs := &core.GameState{
		MyNation:    myNation,
		CardDeck:    cardDeck,
		MapGrid:     mapGrid,
		Treasury:    treasury,
		CurrentTurn: 1,
	}

	return gs
}
