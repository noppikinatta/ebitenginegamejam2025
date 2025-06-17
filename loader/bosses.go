package loader

import (
	"errors"
	"fmt"

	"github.com/noppikinatta/ebitenginegamejam2025/entity"
)

func LoadBosses(path string) (map[string]*entity.Boss, error) {
    r,c,err := openCSV(path); if err!=nil {return nil,err}; defer c.Close()
    recs,err := r.ReadAll(); if err!=nil {return nil,err}
    if len(recs)==0 {return nil, errors.New("bosses csv empty")}
    m := make(map[string]*entity.Boss)
    for i, rec := range recs[1:] {
        if len(rec)<7 {return nil, fmt.Errorf("bosses csv row %d col mismatch", i+2)}
        id,name:=rec[0],rec[1]
        atk,_ := parseInt(rec[2]); def,_ := parseInt(rec[3]); hp,_ := parseInt(rec[4])
        rewardGold,_ := parseInt(rec[5]); rewardCard := rec[6]
        m[id]=&entity.Boss{ID:id,Name:name,Attack:atk,Defense:def,Health:hp,RewardGold:rewardGold,RewardCardID:rewardCard}
    }
    return m,nil
} 