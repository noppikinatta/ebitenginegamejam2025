# MapGrid計画

## 国

配置
- 0,0: 王国(MyNation)
- 0,2: 森の国
- 2,0: 山の国
- 2,2: 砂漠の国
- 0,4: 侍の国
- 4,0: 魔法の国
- 2,4: 機械の国
- 4,2: お祭りの国

Marketレベルごとのカードパックを示す。

### 王国
- 1: cardpack-free
- 1: cardpack-soldiers
- 2: cardpack-politics
- 3: cardpack-knights
- 5: cardpack-war

### 森の国
- 1: cardpack-forest
- 2: cardpack-politics
- 4: cardpack-magic

### 山の国
- 1: cardpack-mountain
- 2: cardpack-mineral
- 4: cardpack-siege

### 砂漠の国
- 1: cardpack-desert
- 1: cardpack-politics
- 3: cardpack-finance
- 5: cardpack-building

### 侍の国
- 1: cardpack-samurai
- 3: cardpack-mineral
- 4: cardpack-war

### 魔法の国
- 1: cardpack-magic
- 3: cardpack-mystic

### 機械の国
- 1: cardpack-mechanical
- 3: cardpack-siege
- 4: cardpack-building

### お祭りの国
- 1: cardpack-fancy
- 2: cardpack-finance

## カードパック

開封時のカード枚数と価格を記載する。

ratio:card でカード構成と出現率を示す。

### cardpack-free, "無料カードパック"

1 card

no price

10:battlecard-soldier,"兵士"
1:battlecard-knight,"騎士"

### cardpack-soldiers, "兵団カードパック"

3 cards

price
- 2 food
- 2 money

5:battlecard-soldier,"兵士"
5:battlecard-archer,"弓使い"
1:battlecard-knight,"騎士"

### cardpack-knights, "騎士団カードパック"

3 cards

price
- 5 food
- 5 iron
- 10 money

3:battlecard-knight,"騎士"
1:battlecard-general,"将軍"
2:structurecard-catapult,"カタパルト"

### cardpack-politics, "内政カードパック"

2 cards

price
- 5 wood

1:structurecard-farm,"農場"
1:structurecard-woodcutter,"木こり小屋"
1:structurecard-tunnel,"坑道"
1:structurecard-market,"市場"

### cardpack-war, "戦時体制カードパック"

5 cards

price
- 20 iron
- 20 wood
- 20 money

2:structurecard-catapult,"カタパルト"
1:structurecard-ballista,"バリスタ"
1:structurecard-camp,"野営地"

### cardpack-magic, "魔法カードパック"

3 cards

price
- 5 mana
- 5 food

5:battlecard-wizard,"魔法使い"
5:structurecard-mana-node,"マナ・ノード"
1:battlecard-mage,"賢者"

### cardpack-mystic, "神秘のカードパック"

3 cards

price
- 20 mana

2:battlecard-fortune,"占い師"
2:structurecard-mana-node,"マナ・ノード"
2:structurecard-temple,"祖廟"
1:battlecard-mage,"賢者"

### cardpack-mineral, "鉱石カードパック"

2 cards

price
- 10 wood

4:structurecard-tunnel,"坑道"
1:structurecard-smelter,"溶鉱炉"
4:battlecard-blacksmith,"鍛冶屋"


### cardpack-mechanical, "機械カードパック"

2 cards

price
- 30 iron

1:battlecard-golem,"ゴーレム"
2:structurecard-smelter,"溶鉱炉"
2:battlecard-artillery,"砲兵"
1:structurecard-ballista,"バリスタ"

### cardpack-fancy, "派手なカードパック"

3 cards

price
- 50 money
- 10 food

5:battlecard-clown,"道化師"
2:battlecard-wrestler,"レスラー"
5:battlecard-bard,"吟遊詩人"
1:battlecard-fortune,"占い師"

### cardpack-samurai, "侍カードパック"

2 cards

price
- 50 iron
- 10 food

4:battlecard-samurai,"侍"
2:battlecard-ninja,"忍者"
3:battlecard-monk,"武僧"
1:structurecard-camp,"野営地"

### cardpack-siege, "攻城兵器カードパック"

2 cards

price
- 50 iron
- 50 wood

2:structurecard-catapult,"カタパルト"
2:structurecard-ballista,"バリスタ"
1:structurecard-orban-cannon,"ウルバン砲"
2:battlecard-artillery,"砲兵"

### cardpack-finance, "金融カードパック"

2 cards

price
- 30 money
- 10 wood

2:battlecard-blacksmith,"鍛冶屋"
2:structurecard-market,"市場"
1:structurecard-mint,"造幣局"

### cardpack-building, "建物カードパック"

2 cards

price
- 30 money
- 30 wood

1:structurecard-granary,"穀倉"
1:structurecard-sawmill,"製材所"
1:structurecard-smelter,"溶鉱炉"
1:structurecard-mint,"造幣局"
1:structurecard-temple,"祖廟"
1:structurecard-camp,"野営地"

### cardpack-forest, "森のカードパック"

3 cards

price
- 5 food
- 5 wood

2:battlecard-archer,"弓使い"
1:structurecard-farm,"農場"
2:structurecard-woodcutter,"木こり小屋"
1:structurecard-mana-node,"マナ・ノード"

### cardpack-desert, "砂漠のカードパック"

3 cards

price
- 10 money
- 10 food

2:battlecard-fortune,"占い師"
1:battlecard-bard,"吟遊詩人"
2:structurecard-market,"市場"
1:structurecard-mana-node,"マナ・ノード"

### cardpack-mountain, "山のカードパック"

3 cards

price
- 5 food
- 5 wood

2:battlecard-blacksmith,"鍛冶屋"
1:structurecard-tunnel,"坑道"
2:structurecard-woodcutter,"木こり小屋"
1:battlecard-soldier,"兵士"

### Appendix:全カードリスト

battlecard-soldier,"兵士"
battlecard-knight,"騎士"
battlecard-general,"将軍"
battlecard-archer,"弓使い"
battlecard-fortune,"占い師"
battlecard-wizard,"魔法使い"
battlecard-mage,"賢者"
battlecard-blacksmith,"鍛冶屋"
battlecard-golem,"ゴーレム"
battlecard-samurai,"侍"
battlecard-ninja,"忍者"
battlecard-monk,"武僧"
battlecard-bard,"吟遊詩人"
battlecard-artillery,"砲兵"
battlecard-clown,"道化師"
battlecard-wrestler,"レスラー"
structurecard-farm,"農場"
structurecard-woodcutter,"木こり小屋"
structurecard-tunnel,"坑道"
structurecard-market,"市場"
structurecard-mana-node,"マナ・ノード"
structurecard-granary,"穀倉"
structurecard-sawmill,"製材所"
structurecard-smelter,"溶鉱炉"
structurecard-mint,"造幣局"
structurecard-temple,"祖廟"
structurecard-camp,"野営地"
structurecard-catapult,"カタパルト"
structurecard-ballista,"バリスタ"
structurecard-orban-cannon,"ウルバン砲"

