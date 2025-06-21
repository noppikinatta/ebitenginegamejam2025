package core

// struct Battlefield　は未制圧のWildernessで戦闘開始時に生成される戦場オブジェクトです。
// - Enemy *Enemy 戦闘相手のEnemy。
// - Effects []BattlefieldEffect
// - SupportPower float64 周囲のTerritoryに置いたStructureCardの影響で増加したPower。
// - BattleCards []*BattleCard 戦闘中に出すBattleCardの集合。
// - CardSlot int BattleCardを置くことができる枚数。
// Battlefieldのインスタンスを作る関数がどこかに必要。
