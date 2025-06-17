# Game for Ebitengine game jam

## Theme of the game jam

UNION

## Overview

国家連合を作っていく戦略ゲーム

- ワールドマップにスタート地点となる1つの自国、複数のNPC国と敵のいる戦闘フィールド、最終地点となる1つのボス敵拠点が存在する。
- 自国、NPC国、敵のいるポイントが線で繋がっている。
- 自国とNPC国の線を塞ぐ敵を倒しながら同盟を作り、シナジーでデッキを強化して最終的にボス敵の打倒を目指す。

- 画面の大きさは1280x720。ただしLayout関数で2倍にするので、ゲーム内では640x360になる。

## Scenes

### Title

最初に表示されるタイトル画面と簡単なストーリー紹介。

### InGame

メインとなるゲーム部分。画面を構成する要素は、GameMain,ResourceView,Calendar,History,CardDeckの5つ。

{Left,Top,Width,Height}の形式で、画面を構成する要素の位置と大きさを示す。

- GameMain:{0,20,520,280}
- ResourceView:{0,0,520,20}
- Calendar:{520,0,120,40}
- History:{520,80,120,320}
- CardDeck:{0,300,520,60}


- 資源ビュー:画面上部に横に細長く資源の種類と個数を表す一覧がある。
- メイン部分:画面中央左にメイン部分がある。状況によってこの画面は切り替わる。
- カレンダーと年表:画面中央右に現在の年月と年表の表示部分がある。
- カードデッキ:画面下部に自分の手持ちのカードを表示するデッキ部分がある。

以降で画面の表示内容を説明する。GameMainは状況によって変わるので分けて記載する。

#### GameMain1:Map

クリックできるポイントとポイント同士を繋ぐ線で構成されたマップ。
GameMainの領域が520x280のサイズなので、13x7個のポイントを格子状に配置する。左下が自国で右上がボス。残りはランダム。

ポイント一覧
- 自国:1つだけあるスタート地点。
- 野外:敵のいるポイント。複数ある。クリックするとBattle画面になる。倒すとその土地を制圧したことになり、資源を算出するようになる。
- NPC国:同盟国のいるポイント。複数ある。クリックするとDiplomacy画面になる。
- ボス:最終ボスのいるポイント。１つ。クリックするとBattle画面になる。倒すとゲームクリア。

ポイントをクリックして外交や戦闘を行うと、ターンが経過し翌月になる仕組み。何もしたくない場合、自国をクリックすることで内政を行ったことにし、基本的なカードを１つ引く。

自国と他のポイントの間に、敵を倒していない野外ポイントがあると、クリックしても利用できない。

#### GameMain2:Diplomacy

カードを資源で買うことができる。
NPC国の特色によって販売されるカードの種類が変わる。ラインナップはNPC国のカードプールの中からランダムに選ばれる。
カードと価格となる資源の一覧が表示される。

#### GameMain3:Battle

カードのうち、戦闘に出せるものを選んで自陣とし、敵と戦う。敵味方ともに、前衛と後衛をそれぞれ５枚、合計10枚まで出せる。

具体的な戦闘ルールは後で作るので、カードを出して敵を倒すという遷移の部分だけまず作る。

#### ResourceView

資源のアイコンと個数が表示される。資源は５種類。

#### Calendar

現在の年と月を表示する。1ターンで1月経過する。王国暦1000年4月からスタートする。

#### History

年月ごとの戦闘のログを表示して、歴史物っぽい雰囲気を出したい。

#### CardDeck

自国のカード一覧。カードを横一列に並べる。最大枚数はまだ考えていない。

### Result

最終ボスを倒した時点で表示される。年表を表示する。

### Rules

#### Resources

金(コイン),鉄,木材,穀物,マナの５種類。どのような地形が何を算出するかは後で考える。

#### Cards

以下の種類のカードがある。

- ユニットカード:先頭に参加するカード
- エンチャントカード:ユニットに合体させて強化する
- 建物カード:自国や制圧した野外ポイントの性能を高める

## CSV-Driven Data Integration (Phase 8)

To decouple hard-coded values we will switch core managers to rely on the CSV data described in `csv_design.md`.

### Loader Initialisation Flow
1. `app/main.go` start-up creates a `loader.Config{DataDir: "data"}` (or simply pass paths).
2. Call loaders in order:
   ```go
   cards   , _ := loader.LoadCards(filepath.Join(dataDir, "cards.csv"))
   nations , _ := loader.LoadNations(filepath.Join(dataDir, "nations.csv"))
   enemies , _ := loader.LoadEnemies(filepath.Join(dataDir, "enemies.csv"))
   bosses  , _ := loader.LoadBosses(filepath.Join(dataDir, "bosses.csv"))
   ```
3. Inject the maps into system managers on construction:
   * `system.NewCardManager(cards)` – replaces internal template generation.
   * `system.NewAllianceManager(nations)` – sets initial relationship values & bonuses.
   * `system.NewCombatManager(enemies, bosses)` – provides factory methods for enemy/boss instances.

### Manager Changes
| Manager | New Responsibility |
|---------|--------------------|
| **CardManager** | Store `map[string]*entity.Card` templates from CSV. `CreateCard(id)` clones template by id. |
| **AllianceManager** | Holds `map[string]*entity.Nation`. Uses `InitialRelationship`, `AllyBonusGold`, `AllyBonusAttack`. |
| **CombatManager** | Hold enemy & boss templates. `AddEnemyByID(id)` pulls stats from template. Victory rewards use `RewardGold` / `RewardCardID`. |
| **MapGrid / Point** | Replace random enemy stats: each Wild point stores an `EnemyID`. Boss point stores a `BossID`. During click, GameMain asks CombatManager to spawn using ID. |

### Data Consistency Rules
* IDs in `enemies.csv` & `bosses.csv` must be unique and referenced by MapGrid.
* `reward_card_id` must exist in `cards.csv`.
* Loader silently trims unknown columns → forward compatible.

### Hot Reload (Optional)
`loader.Watch(files…)` will use `fsnotify` to reload CSV and push updates to observers (`CardManager.Refresh`). Not required for jam deadline but API left open.

### Acceptance Tests
* Phase 8 tests listed in `.tdd.md` (`T8.1‒T8.6`) will use small fixture CSV files placed under `testdata/`.
* Existing gameplay tests must still pass using fixture data.

---
