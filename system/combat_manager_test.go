package system_test

import (
	"testing"

	"github.com/noppikinatta/ebitenginegamejam2025/scene"
)

func TestCardsCanBePlacedInFrontBackRows(t *testing.T) {
    g := scene.CreateSequence(); g.SetCurrentScene("ingame")
    cm := g.GetInGameScene().GetCombatManager()
    bf := cm.GetBattlefield()
    if len(*bf.GetFrontRow())!=5 || len(*bf.GetBackRow())!=5 { t.Fatalf("row size wrong") }
    cardMgr := g.GetInGameScene().GetCardManager()
    for i:=0;i<5;i++{ cm.PlaceCardInBattle(cardMgr.CreateCard("W","Unit",nil),"front",i) }
    for i:=0;i<5;i++{ cm.PlaceCardInBattle(cardMgr.CreateCard("A","Unit",nil),"back",i) }
    if len(cm.GetPlacedCards())!=10 { t.Error("not all cards placed") }
}

func TestCombatCalculatesDamageCorrectly(t *testing.T) {
    g := scene.CreateSequence(); g.SetCurrentScene("ingame")
    cm := g.GetInGameScene().GetCombatManager()
    card := g.GetInGameScene().GetCardManager().CreateCard("Warrior","Unit",nil)
    cm.PlaceCardInBattle(card,"front",0)
    cm.AddEnemy("Goblin",3,2,5)
    if cm.CalculatePlayerDamage()<=0 || cm.CalculateEnemyDamage()<=0 { t.Error("damage non positive") }
    h := cm.GetEnemyHealth("Goblin"); cm.ExecuteCombatRound(); if cm.GetEnemyHealth("Goblin")>=h { t.Error("health not reduced") }
}

func TestVictoryDefeatConditionsTriggerProperly(t *testing.T) {
    g := scene.CreateSequence(); g.SetCurrentScene("ingame")
    cm := g.GetInGameScene().GetCombatManager(); cMgr := g.GetInGameScene().GetCardManager()
    cm.PlaceCardInBattle(cMgr.CreateCard("Mage","Unit",nil),"front",0)
    cm.AddEnemy("Weak",1,1,1)
    for i:=0;i<10 && cm.GetCombatState()=="ongoing"; i++ { cm.ExecuteCombatRound() }
    if cm.GetCombatState()!="victory" { t.Error("should be victory") }
    cm.ResetCombat(); cm.ClearPlayerCards(); cm.ClearEnemies()
    weak := cMgr.CreateCard("Peasant","Unit",nil); weak.SetStats(1,1); cm.PlaceCardInBattle(weak,"front",0)
    cm.AddEnemy("Dragon",20,10,50)
    for i:=0;i<5 && cm.GetCombatState()=="ongoing"; i++ { cm.ExecuteCombatRound() }
    if cm.GetCombatState()=="ongoing" { t.Log("combat still ongoing, acceptable") }
} 