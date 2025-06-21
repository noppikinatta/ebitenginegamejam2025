package core

// struct MapGrid
// - SizeX, SizeY int
// - Points []Point
// - func CanInteract(x, y int) bool 引数x, yの座標にあるPointを操作できるかどうかを返す。
//   - MyNationPointは常に操作できる。
//   - MyNationPointに隣接するPointは操作できる。
//   - MyNationPointから座標を辿っていくとき、WildernessPoint(Controled=false)を通った場合は操作できない。
//     * MyNationPointから遠回りしていくパターンをテストすべき。
//     * WildernessPointのControledは制圧時に切り替わるので、その時点で操作できるかどうかのキャッシュを作る方がいい。

// interface Point
// - func Location() (x, y int) 座標を返す。

// struct MyNationPoint
// - MyNation *MyNation

// struct OtherNationPoint
// - OtherNation *OtherNation

// struct WildernessPoint
// - Controled bool 未制圧ならfalse、制圧されていればtrue。Enemyを倒すとtrueになる。
// - Enemy *Enemy 制圧のために倒さなければならない敵。
// - Territory *Territory 制圧後に算出を得たり、カードを置くための場所。

// struct BossPoint
// - Boss *Enemy ボスのEnemy。これを倒すとゲームクリア。
