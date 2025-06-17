package loader

import (
	"errors"
	"fmt"

	"github.com/noppikinatta/ebitenginegamejam2025/entity"
)

func LoadNations(path string) (map[string]*entity.Nation, error) {
    r, closer, err := openCSV(path)
    if err != nil { return nil, err }
    defer closer.Close()

    recs, err := r.ReadAll(); if err != nil { return nil, err }
    if len(recs)==0 { return nil, errors.New("nations csv empty") }
    nations := make(map[string]*entity.Nation)
    for i, rec := range recs[1:] {
        if len(rec)<6 { return nil, fmt.Errorf("nations csv row %d col mismatch", i+2) }
        id:=rec[0]; name:=rec[1]
        initRel, _ := parseInt(rec[2])
        bonusGold,_ := parseInt(rec[3])
        bonusAtk,_ := parseInt(rec[4])
        flavor := rec[5]
        nations[id]=&entity.Nation{ID:id,Name:name,InitialRelationship:initRel,AllyBonusGold:bonusGold,AllyBonusAttack:bonusAtk,Flavor:flavor}
    }
    return nations,nil
} 