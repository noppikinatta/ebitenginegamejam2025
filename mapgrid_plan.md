# MapGrid計画 - 敵

## 地形

WildernessPoint.TerrainTypeに設定するIDと、Territoryに設定する基本的な地形算出。

- terrain-plain, "平原" : Food +2
- terrain-forest, "森" : Wood +2
- terrain-mountain, "山" : Iron +2
- terrain-desert, "砂漠" : None
- terrain-mana-node, "マナノード" : Mana +3

配置
- (1,0): forest
- (0,1): mountain
- (1,1): plain
- (2,1): desert
- (1,2): desert
- (0,3): mountain
- (3,0): forest
- (1,3): mana-node
- (3,1): mountain
- (1,4): plain
- (4,1): plain
- (2,3): forest
- (3,2): mana-node
- (3,4): forest
- (4,3): mountain
- (3,3): mana-node
- (4,4): boss pointなので産出なし

> TerritoryのCardSlotは全て3にする。

## 敵

(4,4)にはBossPointを配置し、そのほかにはWildernessPointを配置する。

配置
- (1,0): enemy-goblin,"ゴブリン"
- (0,1): enemy-sabrelouse,"サーベルねずみ"
- (1,1): enemy-rattlesnake,"ガラガラヘビ"
- (2,1): enemy-condor,"コンドル"
- (1,2): enemy-slime,"スライム"
- (0,3): enemy-crocodile,"クロコダイル"
- (3,0): enemy-grizzly,"ハイイログマ"
- (1,3): enemy-skeleton,"スケルトン"
- (3,1): enemy-elemental,"エレメンタル"
- (1,4): enemy-dragon,"ドラゴン"
- (4,1): enemy-griffin,"グリフォン"
- (2,3): enemy-vampire,"ヴァンパイア"
- (3,2): enemy-living-armor,"リビングアーマー"
- (3,4): enemy-arc-demon,"アークデーモン"
- (4,3): enemy-durendal,"デュラハン"
- (3,3): enemy-obelisk,"オベリスク"
- (4,4): enemy-final-boss,"魔王"

以下は敵のデータ

### enemy-goblin,"ゴブリン"

Power 3
EnemyType demonic
CardSlot 3

Skill なし

### enemy-sabrelouse,"サーベルねずみ"

Power 4
EnemyType animal
CardSlot 3

Skill なし

### enemy-rattlesnake,"ガラガラヘビ"

Power 6
EnemyType dragon
CardSlot 3

Skill なし

### enemy-condor,"コンドル"

Power 6
EnemyType flying
CardSlot 3

Skill evasion

### enemy-slime,"スライム"

Power 6
EnemyType unknown
CardSlot 3

Skill soft

### enemy-crocodile,"クロコダイル"

Power 10
EnemyType dragon
CardSlot 4

Skill なし

### enemy-grizzly,"ハイイログマ"

Power 12
EnemyType animal
CardSlot 4

Skill なし

### enemy-skeleton,"スケルトン"

Power 12
EnemyType undead
CardSlot 4

Skill longbow

### enemy-elemental,"エレメンタル"

Power 20
EnemyType unknown
CardSlot 5

Skill incorporeality

### enemy-dragon,"ドラゴン"

Power 30
EnemyType dragon
CardSlot 6

Skill pressure

### enemy-griffin,"グリフォン"

Power 25
EnemyType flying
CardSlot 4

Skill evasion

### enemy-vampire,"ヴァンパイア"

Power 30
EnemyType undead
CardSlot 6

Skill charm

### enemy-living-armor,"リビングアーマー"

Power 50
EnemyType unknown
CardSlot 7

Skill なし

### enemy-arc-demon,"アークデーモン"

Power 45
EnemyType demonic
CardSlot 7

Skill magic-barrier

### enemy-durendal,"デュラハン"

Power 45
EnemyType undead
CardSlot 7

Skill side-attack

### enemy-obelisk,"オベリスク"

Power 40
EnemyType unknown
CardSlot 8

Skill laser

### enemy-final-boss,"魔王"

Power 60
EnemyType demonic
CardSlot 8

Skill wave

## Appendix: Enemy Types

enemy-type-animal, "動物"
enemy-type-flying, "飛行"
enemy-type-undead, "アンデッド"
enemy-type-dragon, "ドラゴン"
enemy-type-demonic, "悪魔"
enemy-type-unknown, "不明"

## Appendix: Enemy Skills

enemy-skill-evasion, "回避"
enemy-skill-evasion-desc, "力タイプのカードパワー2"
enemy-skill-soft,"軟体"
enemy-skill-soft-desc, "魔タイプ以外のカードパワー-50%"
enemy-skill-longbow,"ロングボウ"
enemy-skill-longbow-desc, "最も後ろのカードパワー-100%"
enemy-skill-incorporeality,"霊体"
enemy-skill-incorporeality-desc, "魔タイプ以外のカードパワー-100%"
enemy-skill-pressure,"威圧"
enemy-skill-pressure-desc, "全てのカードパワー-1"
enemy-skill-charm,"魅了"
enemy-skill-charm-desc, "先頭から3枚のカードパワー-100%"
enemy-skill-magic-barrier,"魔法障壁"
enemy-skill-magic-barrier-desc, "魔法タイプのカードパワー-100%"
enemy-skill-laser,"レーザー"
enemy-skill-laser-desc, "後方から3枚のカードパワー-100%"
enemy-skill-side-attack,"側面攻撃"
enemy-skill-side-attack-desc, "先頭から5枚のカードパワー-50%"
enemy-skill-wave,"波動"
enemy-skill-wave-desc, "全てのカードパワー-2"
