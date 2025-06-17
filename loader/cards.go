package loader

import (
	"errors"
	"fmt"

	"github.com/noppikinatta/ebitenginegamejam2025/entity"
)

// LoadCards reads data/cards.csv and returns a map[id]*entity.Card
func LoadCards(path string) (map[string]*entity.Card, error) {
    r, closer, err := openCSV(path)
    if err != nil {
        return nil, err
    }
    defer closer.Close()

    records, err := r.ReadAll()
    if err != nil {
        return nil, err
    }
    if len(records) == 0 {
        return nil, errors.New("cards csv: empty file")
    }
    cards := make(map[string]*entity.Card)
    // assume header row exists -> skip first row
    for i, rec := range records[1:] {
        if len(rec) < 11 {
            return nil, fmt.Errorf("cards csv: row %d expect 11 columns, got %d", i+2, len(rec))
        }
        id := rec[0]
        name := rec[1]
        ctype := rec[2]
        attack, err := parseInt(rec[3])
        if err != nil { return nil, err }
        defense, err := parseInt(rec[4])
        if err != nil { return nil, err }
        costGold, _ := parseInt(rec[5])
        costIron, _ := parseInt(rec[6])
        costWood, _ := parseInt(rec[7])
        costGrain, _ := parseInt(rec[8])
        costMana, _ := parseInt(rec[9])
        desc := rec[10]

        cost := map[string]int{"Gold": costGold, "Iron": costIron, "Wood": costWood, "Grain": costGrain, "Mana": costMana}
        card := &entity.Card{
            Name: name,
            Type: ctype,
            Cost: cost,
            Description: desc,
            Attack: attack,
            Defense: defense,
        }
        cards[id] = card
    }
    return cards, nil
} 