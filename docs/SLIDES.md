---
marp: true
theme: default
paginate: true
---

# Goで学ぶ並行処理

ログ解析で goroutine と channel を体験
処理速度を改善する

---

<!-- omit in toc -->
## 目次

1. [はじめに](#はじめに)
   - [タイムスケジュール](#タイムスケジュール)
   - [お題](#お題)
   - [逐次処理だと](#逐次処理だと)
   - [なぜ遅いのか](#なぜ遅いのか)
   - [並行処理なら](#並行処理なら)
2. [goroutine](#goroutine)
   - [goroutine とは](#goroutine-とは)
   - [go キーワード](#go-キーワード)
   - [何が起きているのか](#何が起きているのか)
   - [goroutine が軽い理由](#goroutine-が軽い理由)
   - [問題: main が先に終わる](#問題-main-が先に終わる)
   - [図で見ると](#図で見ると)
   - [解決: sync.WaitGroup](#解決-syncwaitgroup)
   - [WaitGroup の使い方](#waitgroup-の使い方)
   - [カウンタの動きを追う](#カウンタの動きを追う)
   - [複数の goroutine を待つ](#複数の-goroutine-を待つ)
   - [なぜ defer を使うのか](#なぜ-defer-を使うのか)
   - [Go 1.25: WaitGroup.Go()](#go-125-waitgroupgo)
   - [WaitGroup.Go() の利点](#waitgroupgo-の利点)
   - [ロジックは分けておく](#ロジックは分けておく)
3. [channel](#channel)
   - [新しい問題](#新しい問題)
   - [channel とは](#channel-とは)
   - [channel の作り方と使い方](#channel-の作り方と使い方)
   - [送信と受信の対応](#送信と受信の対応)
   - [ブロックとは?](#ブロックとは)
   - [送信時のブロック](#送信時のブロック)
   - [受信時のブロック](#受信時のブロック)
   - [なぜブロックするのか](#なぜブロックするのか)
   - [結果を集める](#結果を集める)
   - [流れを図で見る](#流れを図で見る)
   - [このパターンのポイント](#このパターンのポイント)
   - [ハマりどころ: デッドロック](#ハマりどころ-デッドロック)
   - [なぜデッドロックになるのか](#なぜデッドロックになるのか)
   - [解決策: 送信と受信を別の goroutine で](#解決策-送信と受信を別の-goroutine-で)
   - [デッドロックを防ぐコツ](#デッドロックを防ぐコツ)
4. [ワーカープール](#ワーカープール)
   - [Phase 2 の方法の問題点](#phase-2-の方法の問題点)
   - [大量の goroutine の問題](#大量の-goroutine-の問題)
   - [ワーカープールの考え方](#ワーカープールの考え方)
   - [ワーカープールの実装(前半)](#ワーカープールの実装前半)
   - [for range channel の動き](#for-range-channel-の動き)
   - [ワーカープールの実装(後半)](#ワーカープールの実装後半)
   - [close の役割](#close-の役割)
   - [close すると何が起きる?](#close-すると何が起きる)
   - [close を忘れるとどうなる?](#close-を忘れるとどうなる)
   - [バッファなし channel の動き](#バッファなし-channel-の動き)
   - [バッファ付き channel の動き](#バッファ付き-channel-の動き)
   - [バッファあり/なし の使い分け](#バッファありなし-の使い分け)
   - [ワーカープールでバッファを使う理由](#ワーカープールでバッファを使う理由)
   - [ワーカープールのポイント](#ワーカープールのポイント)
5. [並行処理パターン集](#並行処理パターン集)
   - [なぜパターンを学ぶのか](#なぜパターンを学ぶのか)
   - [パターン1: Generator](#パターン1-generator)
   - [パターン2: Pipeline](#パターン2-pipeline)
   - [パターン3: Fan-out / Fan-in](#パターン3-fan-out--fan-in)
   - [パターン4: select](#パターン4-select)
   - [パターン5: Done Channel(キャンセル)](#パターン5-done-channelキャンセル)
   - [パターン6: Timeout](#パターン6-timeout)
   - [パターン7: Semaphore](#パターン7-semaphore)
   - [パターン8: Rate Limiting](#パターン8-rate-limiting)
   - [パターン9: context.Context](#パターン9-contextcontext)
   - [パターンの選び方](#パターンの選び方)
   - [やりたいこと → パターン対応表](#やりたいこと--パターン対応表)
   - [実務での組み合わせ例](#実務での組み合わせ例)
   - [パターン選択のコツ](#パターン選択のコツ)
   - [今日のハンズオンで使うパターン](#今日のハンズオンで使うパターン)
6. [ハンズオン](#ハンズオン)
   - [4つの Phase](#4つの-phase)
   - [ルール](#ルール)

---

# はじめに

---

## タイムスケジュール

| 時間 | 内容 |
|------|------|
| 00:00-00:32 | この講義 |
| 00:32-00:37 | Phase 1(逐次処理) |
| 00:37-00:52 | Phase 2(並行処理1) |
| 00:52-01:05 | Phase 3(並行処理2) |
| 01:05-01:20 | Phase 4(さらなる高速化) |
| 01:20-01:25 | 最終計測 |
| 01:25-01:30 | 結果発表 |

---

## お題

「このログ、急ぎで解析して」

- 200ファイル × 50,000行 ≈ 1,000万行
- ステータスコード別にカウントしたい

---

## 逐次処理

```go
for _, file := range files {
    result := processFile(file)
}
```

200ファイルを1つずつ処理。

```
時間 →
[file1を処理]→[file2を処理]→[file3を処理]→ ...
              ↑
         file1が終わるまで
         file2は待っている
```

CPUは暇な時間が多い。ファイルI/Oの待ち時間がもったいない。

---

## 並行処理なら

```
時間 →
[file1を処理]→
[file2を処理]→
[file3を処理]→
    ...
```

複数のファイルを同時に処理

この仕組みを理解して実装する。

ただし、

- **並行処理は常に速くなるわけではない**
- goroutine の作成・管理にもコストがかかる
- 処理が小さすぎると、オーバーヘッドの方が大きくなる
- goroutineを大量に作りすぎると、メモリやCPUの負荷が増える
- 適切な並行化の設計が重要

 参考: [Goroutines in Go - GetStream](https://getstream.io/blog/goroutines-go-concurrency-guide/) | [Go Concurrency Patterns](https://ggbaker.ca/prog-langs/content/go-concurrency.html)

---

# goroutine

---

## goroutine とは

Goランタイムが管理する軽量なスレッド。

普通の関数呼び出しは、その関数が終わるまで次に進めないが、(同期実行)。
goroutine を使うと、関数を別の実行フローで動かし、呼び出し元は待たずに次の処理に進める(非同期実行)。

参考: [A Tour of Go - Goroutines](https://go.dev/tour/concurrency/1)

---

## go キーワード

```go
// 普通に呼ぶ → processFile が終わるまで待つ
processFile("access_001.json")

// goroutine として呼ぶ → 待たずに次へ進む
go processFile("access_001.json")
```

`go` を付けるだけで、その関数は別の流れで実行される。

 参考: [Go Spec - Go statements](https://go.dev/ref/spec#Go_statements) | [Effective Go - Goroutines](https://go.dev/doc/effective_go#goroutines)

---

## 何が起きているのか

```go
func main() {
    go processFile("file1.json")  // ①
    go processFile("file2.json")  // ②
    fmt.Println("done")          // ③
}
```

①と②は「処理を開始しろ」という指示を出すだけ。
実際の処理完了を待たずに、すぐ③に進む。

---

## goroutine が軽い理由

- 初期スタックサイズ: わずか2KB(Go 1.4以降)
- Goランタイムが管理、必要に応じて動的に拡張・縮小
- 最大: 64bit 1GB / 32bit 250MB
- 数千〜数万個でも問題なく動く

 参考: [What is a goroutine? And what is their size?](https://tpaschalis.me/goroutines-size/) | [Cloudflare: How Stacks are Handled in Go](https://blog.cloudflare.com/how-stacks-are-handled-in-go/)

---

## 問題: main が先に終わる

```go
func main() {
    go processFile("file1.json")
    go processFile("file2.json")
    // ここで main が終わる
}
```

出力: 何も表示されない

main関数が終わると、プログラム全体が終了する。
goroutine が処理中でも、終了する。

---

## 図で見ると

```
main        [開始]──────────────────[終了] ← プログラム終了
                ↓         ↓
file1           [処理中...] ← 途中で強制終了
file2             [処理中...] ← 途中で強制終了
```

main は goroutine の完了を待っていない。
「待つ仕組み」が必要。

---

## sync.WaitGroup

「全部終わるまで待つ」ための仕組み。
内部にカウンタを持っていて、0になるまで待機できる。

```go
var wg sync.WaitGroup  // カウンタは0で始まる
```

---

## WaitGroup の使い方

```go
var wg sync.WaitGroup

wg.Add(1)  // カウンタを1増やす(1になる)

go func() {
    defer wg.Done()  // 終了時にカウンタを1減らす
    processFile("file1.json")
}()

wg.Wait()  // カウンタが0になるまでここで待つ
```

 参考: [sync.WaitGroup - pkg.go.dev](https://pkg.go.dev/sync#WaitGroup) | [WaitGroup.Go - pkg.go.dev](https://pkg.go.dev/sync#WaitGroup.Go)

---

## カウンタの動き

```go
var wg sync.WaitGroup        // カウンタ: 0

wg.Add(1)                    // カウンタ: 1
go func() {
    defer wg.Done()          // (終了時に) カウンタ: 0
    processFile("file1.json")
}()

wg.Wait()                    // カウンタが0になるまで待機
fmt.Println("完了")
```

---

## 複数の goroutine を待つ

```go
var wg sync.WaitGroup

for _, file := range files {
    wg.Add(1)  // ループごとにカウンタ+1
    go func() {
        defer wg.Done()  // 終わったらカウンタ-1
        processFile(file)
    }()
}

wg.Wait()  // 全部終わるまで待つ
```

200ファイルなら、カウンタは 0→1→2→...→200→199→...→0 と動く。

---

## なぜ defer を使うのか

```go
go func() {
    defer wg.Done()  // ← これ
    processFile(f)
}()
```

`defer` は「この関数を抜ける直前に（正常終了でも panic でも）登録した呼び出しを実行する」キーワード。引数の評価は `defer` を書いた時点で行われる。

processFile でエラーが起きても、Done() は必ず呼ばれる。
カウンタが減らないまま残る事故を防げる。

参考: [Go Blog - Defer, Panic, and Recover](https://go.dev/blog/defer-panic-and-recover)

---

## Go 1.25: WaitGroup.Go()

Go 1.25 から、もっと簡単に書ける新しいメソッドが追加された。

従来のパターン:

```go
wg.Add(1)
go func() {
    defer wg.Done()
    processFile(f)
}()
```

新しいパターン (Go 1.25+):

```go
wg.Go(func() {
    processFile(f)
})
```

`WaitGroup.Go()` は内部で `Add(1)` と `defer Done()` を実行している。

---

## WaitGroup.Go() の利点

1. Add/Done の書き忘れを防ぐ

- 手動で Add(1) を書く必要がない
- defer wg.Done() も不要

2. コードが簡潔になる

- 読みやすく、ミスも減る

 参考: [WaitGroup.Go - pkg.go.dev](https://pkg.go.dev/sync#WaitGroup.Go) | [Go 1.25 Release Notes](https://go.dev/doc/go1.25)

---

## ロジックは分けておく

```go
// ビジネスロジック(並行処理を知らない)
func processFile(name string) Result {
    // ファイルを開いて処理して返す
}

// 並行処理は呼び出し側で制御
go func() {
    defer wg.Done()
    result := processFile(name)
}()
```

processFile は「自分が goroutine で呼ばれているか」を知らなくていい。
テストも書きやすいし、逐次処理でも並行処理でも使える。

---

# channel

---

## 新しい問題

goroutine で処理を並行化できた。
でも、各 goroutine の結果をどうやって集める?

```go
for _, file := range files {
    go func() {
        result := processFile(file)
        // この result をどこに返す?
    }()
}
// ここで全ファイルの結果を集計したい
```

---

## channel とは

goroutine 同士がデータをやり取りするための「通り道」。

```
┌────────────┐             ┌────────────┐
│ goroutine  │ ─── 値 ───→ │ goroutine  │
│     A      │   channel   │     B      │
└────────────┘             └────────────┘
```

一方が値を送り、もう一方が値を受け取る。

Go 言語仕様では「並行実行される関数が、指定された型の値を送受信することで通信するための仕組み」と定義されている。

参考: [Go言語仕様 - Channel types](https://go.dev/ref/spec#Channel_types) | [Effective Go - Channels](https://go.dev/doc/effective_go#channels) | [Go by Example - Channels](https://gobyexample.com/channels)

---

## channel の作り方と使い方

```go
// channel を作る(int型の値を流せる)
ch := make(chan int)

// 値を送る
ch <- 42

// 値を受け取る
value := <-ch
```

`<-` はデータの流れる向きを表している。

 参考: [Go Spec - Channel types](https://go.dev/ref/spec#Channel_types) | [Go Spec - Send statements](https://go.dev/ref/spec#Send_statements) | [Go Spec - Receive operator](https://go.dev/ref/spec#Receive_operator)

---

## 送信と受信の対応

```go
ch := make(chan int)

// 送る側(別の goroutine で)
go func() {
    ch <- 42  // 42 を channel に送る
}()

// 受け取る側(main で)
value := <-ch  // channel から値を受け取る
fmt.Println(value)  // 42
```

---

## ブロックとは?

「ブロック」= goroutine が一時停止する状態

プログラムの実行が「その場で止まる」こと

- その行で待機する
- 次の行には進めない
- 条件が満たされるまで待ち続ける


---

## 送信時のブロック

バッファなしの channel で送信すると、受信側が現れるまで**ブロック**

送信側 goroutine の状態変化

1. 実行中: `ch <- 42` を実行しようとする
2. **ブロック開始**: 受信側がいない → この場で停止
3. 待機中: 他の goroutine は動いている(自分だけ止まる)
4. **ブロック解除**: 受信側が現れた!
5. 実行再開: 値を渡して次の行へ進む

```go
go func() {
    fmt.Println("送信前")
    ch <- 42  // ← ここでブロック(止まる)
    fmt.Println("送信後")  // ← ブロック解除後に実行
}()
```

goroutine は止まっているが、プログラム全体は動いている

---

## 受信時のブロック

受信も同様にブロックする(送信側が現れるまで)

受信側 goroutine の状態変化

1. 実行中: `<-ch` を実行しようとする
2. **ブロック開始**: 送信側がいない → この場で停止
3. 待機中: 値が来るまで待つ
4. **ブロック解除**: 送信側が値を送ってきた!
5. 実行再開: 値を受け取って次の行へ進む

```go
func main() {
    value := <-ch  // ← ここでブロック(止まる)
    fmt.Println(value)  // ← ブロック解除後に実行
}
```

---

## なぜブロックするのか

channel は「データの受け渡し場所」ではなく「待ち合わせ場所」

ブロックの3つの役割

1つ目: 同期を取る

- 送信側と受信側のタイミングを合わせる
- データが確実に受け渡される保証

2つ目: 安全性を保つ

- ロックや条件変数が不要
- Go ランタイムが管理してくれる

3つ目: データ競合を防ぐ

- 送り手と受け手が揃うまで待つ
- 中途半端な状態を防ぐ

```text
送信側: 「42を渡したい」 → ブロック(待機)
受信側:                    → 「受け取る準備OK」
        ← 42 が渡される → (ブロック解除)
送信側: 「渡せた、次へ進む」
```

ブロックがあるから、安全な並行処理ができる

---

## 補足: channelのブロック挙動

| 状態 | いつブロックするか | いつブロックしないか |
| --- | --- | --- |
| バッファなし | 送り手と受け手のどちらかが不在 | 双方が揃った瞬間に進む |
| バッファあり | 送信: 満杯 / 受信: 空 | 送信: 空きあり / 受信: データあり |
| nil channel | 送受信どちらも永遠にブロック | なし（初期化漏れ注意） |
| close済み | 送信はpanic | 受信は即時にゼロ値・ok=false、rangeは即終了 |

 参考: [Go Spec - Channel types](https://go.dev/ref/spec#Channel_types) | [Go Memory Model](https://go.dev/ref/mem) | [Gist of Go: Channels](https://antonz.org/go-concurrency/channels/)

---

## 結果を集める

```go
results := make(chan Result)

// 各ファイルを goroutine で処理
for _, file := range files {
    go func() {
        results <- processFile(file)  // 結果を送信
    }()
}

// 結果を受け取って集計
for i := 0; i < len(files); i++ {
    r := <-results
    // 集計処理
}
```

---

## 流れを図で見る

```
main          files を回して goroutine を起動
                ↓     ↓     ↓
goroutine1    [処理] → results に送信
goroutine2    [処理] → results に送信
goroutine3    [処理] → results に送信

main          results から len(files) 回受信して集計
```

送った数と受け取る数を合わせるのがポイント。

---

## このパターンのポイント

- goroutine は結果を channel に送るだけ(集計は知らない)
- main 側で必要な回数だけ受信して集計
- 送信回数と受信回数を一致させる(これ重要)

今日の Phase 2 はこれを使う

---

## ハマりどころ: デッドロック

全ての goroutine が何かを待っていて、誰も先に進めない状態。

```go
ch := make(chan int)
ch <- 42  // 受け取る人がいない → 永遠に待つ
```

```
fatal error: all goroutines are asleep - deadlock!
```

Go ランタイムがこれを検出してプログラムを止めてくれる。

---

## なぜデッドロックになるのか

```go
func main() {
    ch := make(chan int)
    ch <- 42  // ここで止まる
    fmt.Println(<-ch)  // ここには来ない
}
```

1. main が `ch <- 42` で送信しようとする
2. 受信者がいないので main は待機
3. main 以外に goroutine がいない
4. 誰も受信できない → 永遠に待つ → デッドロック

---

## 解決策: 送信と受信を別の goroutine で

```go
func main() {
    ch := make(chan int)

    go func() {
        ch <- 42  // 別の goroutine で送信
    }()

    fmt.Println(<-ch)  // main で受信
}
```

送信と受信が別々の goroutine にいるので、お互いを待てる。

---

## デッドロックを防ぐコツ

1. 送信回数と受信回数を一致させる。
   50個送るなら、50回受け取る

2. 送信と受信を別の goroutine で行う。
   同じ goroutine 内で両方やると詰まりやすい

3. 「誰が受け取るのか」を常に意識。
   送る前に、受け取る側が存在するか確認

 参考: [Go Memory Model](https://go.dev/ref/mem) | [Effective Go - Channels](https://go.dev/doc/effective_go#channels)

---

# ワーカープール

---

## Phase 2 の方法の問題点

```go
for _, file := range files {
    go func() {
        results <- processFile(file)
    }()
}
```

200ファイルなら200個の goroutine が同時に動く。
これは問題ないが、5000ファイルだったら?

---

## 大量の goroutine の問題

goroutine は軽量だが、無制限に作ると問題が起きる

- メモリ消費
    - 1個あたり約2.7KB使う(最初は2KBだが、実際は少し増える)
    - 大量に作ると、メモリが足りなくなる可能性


- ファイルを同時に開ける数に上限がある
    - OS には「一度に開けるファイル数」の制限がある

- CPU で同時に動けるのは限られている
    - 8コアのマシンで5000個の goroutine を起動しても実際に CPU 上で同時に実行されるのは最大8個だけ
    - 残りは順番待ち(切り替えながら実行)
    - 切り替えの処理にもコストがかかる

結論: goroutine の数を適切に制限した方が効率的

 参考: [What is a goroutine? And what is their size?](https://tpaschalis.me/goroutines-size/) | [How Many Goroutines Can Go Run?](https://leapcell.io/blog/how-many-goroutines-can-go-run)

---

## ワーカープールの考え方

「仕事をするワーカー」を固定数だけ先に用意しておく。
ワーカーは「仕事キュー」から仕事を取って処理する。

```
         jobs channel
仕事 → [file1][file2][file3]...
              ↓
        ┌─────┼─────┐
     worker1  worker2  worker3  (固定数)
        │      │       │
        └──────┼───────┘
               ↓
         results channel
```

- 5000ファイルでも、ワーカーは8個だけ(CPU コア数分)
- 各ワーカーは jobs から順番に仕事を取る
- 全員が同じ jobs channel を見ている

このパターンは Go の並行処理パターンで **Fan-out** と呼ばれる。複数の関数(ワーカー)が同じ channel から読み取ることで、作業を分散し CPU と I/O を並列化できる。

参考: [Go Blog: Pipelines and cancellation](https://go.dev/blog/pipelines) | [Go by Example: Worker Pools](https://gobyexample.com/worker-pools)

---

## ワーカープールの実装(前半)

```go
numWorkers := runtime.NumCPU()  // CPUコア数
jobs := make(chan string, 100)
results := make(chan Result, 100)

// ワーカーを先に起動
for i := 0; i < numWorkers; i++ {
    go func() {
        for file := range jobs {  // jobs が閉じるまでループ
            results <- processFile(file)
        }
    }()
}
```

---

## for range channel の動き

```go
for file := range jobs {
    // ...
}
```

この構文の動作:

1. jobs から値を1つ受け取る
2. ループ本体を実行
3. また jobs から受け取る
4. **jobs が close されるまで繰り返す**

通常の for ループとの違い:

- 通常: `for i := 0; i < 10; i++` → 回数が決まっている
- channel: `for file := range jobs` → **終了条件 = channel が閉じる**

close されると、バッファ内の残りの値を全て処理してから、ループを抜ける。

---

## ワーカープールの実装(後半)

```go
// 仕事を投げる
for _, file := range files {
    jobs <- file
}
close(jobs)  // 「もう仕事はない」と伝える

// 結果を集める
for i := 0; i < len(files); i++ {
    r := <-results
    // 集計処理
}
```

---

## close の役割

close は「もう値を送らない」という合図

```go
close(jobs)  // jobsチャネルを閉じる
```

- ✓ 送信側(書き込む側)だけがcloseすべき
- ✗ 受信側(読み込む側)はcloseしてはいけない
- ✓ channelは1回だけcloseできる(2回目はpanic)

---

## close すると何が起きる?

1. 送信側(main)

- `close(jobs)` を呼ぶ
- これ以降、送信できなくなる
- 送信すると → panic 発生

2. 受信側(ワーカー)

- `for file := range jobs` が終了条件を検知
- バッファに残っている値は全て処理できる
- 全て処理したらループを抜ける
- その後の受信 → ゼロ値が返る(ブロックしない)

---

## close を忘れるとどうなる?

```go
for _, file := range files {
    jobs <- file
}
// close(jobs)  ← これを忘れると...
```

結果: **デッドロック発生**

なぜデッドロックになるのか:

1. ワーカーは `for file := range jobs` で待ち続ける
2. main は結果を待ち続ける
3. jobs は閉じられないので、ワーカーは永遠に待つ
4. → 誰も進めない


 参考: [Go builtin: close](https://pkg.go.dev/builtin#close) | [Go by Example: Closing Channels](https://gobyexample.com/closing-channels) | [Gist of Go: Channels](https://antonz.org/go-concurrency/channels/)

---

## バッファなし channel の動き

```go
ch := make(chan int)  // バッファサイズ: 0(デフォルト)
```

送信側の動き:

1. 値を送ろうとする `ch <- 42`
2. 受信側が準備できるまで **待つ**(ブロック) ← **送信側は止まる**
3. 受信側が受け取ったら、次に進める

受信側の動き:

1. 値を受け取ろうとする `<-ch`
2. 送信側が送るまで **待つ**(ブロック) ← **受信側も止まる**
3. 送信側が送ったら、値を受け取って次に進める


---

## バッファ付き channel の動き

```go
ch := make(chan int, 3)  // バッファサイズ: 3
```

送信側の動き:

1. 値を送る `ch <- 42`
2. **バッファに空きがあれば**、すぐに次に進める ← **送信側は待たない**
3. バッファが満杯なら、空くまで待つ(ブロック)

受信側の動き:

1. 値を受け取ろうとする `<-ch`
2. **バッファに値があれば**、すぐに受け取って次に進める ← **受信側も待たない**
3. バッファが空なら、値が来るまで待つ(ブロック)


---

## バッファあり/なし の使い分け

バッファなし(デフォルト)の使いどころ:

- 送信側と受信側を厳密に同期させたい
- 「確実に受け取られた」ことを確認したい
- シンプルな通知やシグナル

バッファありの使いどころ:

- 送信側と受信側の速度が異なる
- まとめて送信してから処理したい

---

## ワーカープールでバッファを使う理由

```go
numWorkers := runtime.NumCPU()
jobs := make(chan string, len(files))      // バッファあり
results := make(chan Result, len(files))   // バッファあり
```

バッファがないと:

- 送信のたびにワーカーが受け取るまで待つ
- 200ファイルなら200回ブロック
- 非効率

バッファがあると:

- 全ファイルをまとめて投入できる
- ワーカーは自分のペースで処理
- main も待たずに次へ進める

 参考: [Go by Example: Channel Buffering](https://gobyexample.com/channel-buffering) | [Go Tour: Buffered Channels](https://go.dev/tour/concurrency/3) | [Gist of Go: Channels](https://antonz.org/go-concurrency/channels/)

---

## ワーカープールのポイント

- 固定数のワーカーを先に起動(CPU コア数など)
- jobs channel から仕事を取り出して処理
- `close(jobs)` でワーカーに「もう仕事はない」と伝える
- 同時実行数をコントロールできる

今日の Phase 3 で使う

---

# 並行処理パターン集

---

## なぜパターンを学ぶのか

並行処理には「よくある問題」と「定番の解決策」がある。

パターンを知っていれば:
- 車輪の再発明を避けられる
- バグを踏みにくくなる
- コードを読むときに意図がわかる

ここからは実務でよく使うパターンを紹介する。

---

## パターン1: Generator

---

## Generator とは

channel を返す関数。データソースを抽象化できる。

```go
func generateNumbers(n int) <-chan int {
    ch := make(chan int)
    go func() {
        for i := 0; i < n; i++ {
            ch <- i
        }
        close(ch)
    }()
    return ch
}
```

---

## Generator の使い方

```go
// 使う側はデータの生成方法を知らなくていい
for num := range generateNumbers(10) {
    fmt.Println(num)
}
```

ポイント
- 関数内で goroutine を起動し、channel を返す
- 呼び出し側は for range で受け取るだけ
- 生成側と消費側が疎結合になる

---

## Generator の応用例

```go
// ファイルから1行ずつ読む Generator
func readLines(filename string) <-chan string {
    ch := make(chan string)
    go func() {
        defer close(ch)
        file, _ := os.Open(filename)
        defer file.Close()
        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
            ch <- scanner.Text()
        }
    }()
    return ch
}
```

巨大なファイルでもメモリを食わない。必要な分だけ読める。

---

## パターン2: Pipeline

---

## Pipeline とは

処理を段階に分けて、channel で繋ぐパターン。

```
入力 → [Stage1] → [Stage2] → [Stage3] → 出力
         ↓           ↓           ↓
      channel     channel     channel
```

各ステージは独立した goroutine で動く。

---

## Pipeline の実装例

```go
// Stage 1: 数値を生成
func generate(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()
    return out
}
```

---

## Pipeline の実装例(続き)

```go
// Stage 2: 2倍にする
func double(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * 2
        }
        close(out)
    }()
    return out
}
```

---

## Pipeline を繋げる

```go
// パイプラインを構築
nums := generate(1, 2, 3, 4, 5)
doubled := double(nums)
quadrupled := double(doubled)

// 結果を受け取る
for n := range quadrupled {
    fmt.Println(n)  // 4, 8, 12, 16, 20
}
```

ポイント: 各ステージは入力の channel を受け取り、出力の channel を返す

---

## Pipeline のメリット

- 関心の分離: 各ステージは自分の仕事だけに集中
- 再利用性: ステージを組み替えて別のパイプラインを作れる
- 並行性: 各ステージが同時に動く(Stage1が次を出力している間にStage2が処理)

---

## パターン3: Fan-out / Fan-in

---

## Fan-out とは

1つの入力を複数のワーカーに分散させること。

```
              ┌→ [worker1] →┐
入力 → channel ─→ [worker2] →─ 結果
              └→ [worker3] →┘
```

重い処理を並列化したいときに使う。

---

## Fan-out の実装

```go
// 同じ channel から複数のワーカーが読む
jobs := make(chan Job)

for i := 0; i < numWorkers; i++ {
    go func() {
        for job := range jobs {
            process(job)
        }
    }()
}

// 仕事を投入
for _, job := range allJobs {
    jobs <- job
}
```

これは実はワーカープールと同じ。

---

## Fan-in とは

複数の channel を1つにまとめること。

```
[source1] →┐
[source2] →┼→ 1つの channel → 消費者
[source3] →┘
```

複数のデータソースを統合したいときに使う。

---

## Fan-in の実装

```go
func fanIn(channels ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup

    for _, ch := range channels {
        wg.Add(1)
        go func(c <-chan int) {
            defer wg.Done()
            for v := range c {
                out <- v
            }
        }(ch)
    }

    go func() {
        wg.Wait()
        close(out)
    }()
    return out
}
```

---

## Fan-in の使い方

```go
ch1 := generateNumbers(5)
ch2 := generateNumbers(5)
ch3 := generateNumbers(5)

// 3つの channel を1つにまとめる
merged := fanIn(ch1, ch2, ch3)

for n := range merged {
    fmt.Println(n)  // 順序は不定
}
```

注意: 出力の順序は保証されない(先に来たものから出る)

---

## パターン4: select

---

## select とは

複数の channel を同時に待ち受ける構文。

```go
select {
case msg := <-ch1:
    fmt.Println("ch1から:", msg)
case msg := <-ch2:
    fmt.Println("ch2から:", msg)
}
```

どちらか先に来た方を処理する。

---

## select の動作

```go
select {
case v := <-ch1:
    // ch1 から受信できたらここ
case v := <-ch2:
    // ch2 から受信できたらここ
case ch3 <- value:
    // ch3 に送信できたらここ
default:
    // どれもすぐに実行できないならここ
}
```

- 複数が同時に ready なら、ランダムに1つ選ばれる
- default があると、どれも ready でなくてもブロックしない

 参考: [Go Spec - Select statements](https://go.dev/ref/spec#Select_statements) | [Effective Go - Select](https://go.dev/doc/effective_go#select)

---

## select の使い所

1. 複数のデータソースから受信
2. タイムアウトの実装
3. キャンセル処理
4. ノンブロッキング送受信

これらを見ていく。

---

## パターン5: Done Channel(キャンセル)

---

## なぜキャンセルが必要か

goroutine は起動したら勝手に終わらない。

```go
go func() {
    for {
        // 永遠に動き続ける...
    }
}()
```

「もう結果は要らない」と伝える仕組みが必要。

---

## Done Channel パターン

```go
done := make(chan struct{})  // 空の struct(メモリ0バイト)

go func() {
    for {
        select {
        case <-done:
            fmt.Println("キャンセルされた")
            return
        default:
            // 通常処理
        }
    }
}()

// キャンセルしたいとき
close(done)
```

---

## なぜ struct{} を使うのか

```go
done := make(chan struct{})
```

- `struct{}` はサイズ0バイトの型
- 「値を送る」のではなく「シグナルを送る」目的
- close すると、全ての受信側が即座に起きる

```go
close(done)  // 全ての <-done が解除される
```

---

## 複数の goroutine を止める

```go
done := make(chan struct{})

// 10個の goroutine を起動
for i := 0; i < 10; i++ {
    go func(id int) {
        for {
            select {
            case <-done:
                fmt.Printf("worker %d: 終了\n", id)
                return
            default:
                // 作業
            }
        }
    }(i)
}

// 全員を一斉に止める
close(done)
```

---

## パターン6: Timeout

---

## タイムアウトの必要性

外部APIの呼び出しなど、いつまでも待てない処理がある。

「3秒待って返事がなければ諦める」を実装したい。

---

## time.After を使ったタイムアウト

```go
select {
case result := <-ch:
    fmt.Println("結果:", result)
case <-time.After(3 * time.Second):
    fmt.Println("タイムアウト")
}
```

`time.After(d)` は、時間 d が経過すると値を送る channel を返す。

---

## 処理全体にタイムアウトをかける

```go
func fetchWithTimeout(url string) (string, error) {
    result := make(chan string)

    go func() {
        // 時間のかかる処理
        body := fetch(url)
        result <- body
    }()

    select {
    case body := <-result:
        return body, nil
    case <-time.After(5 * time.Second):
        return "", errors.New("timeout")
    }
}
```

---

## パターン7: Semaphore

---

## Semaphore とは

同時実行数を制限する仕組み。

ワーカープールと似ているが、「トークン」を使って制御する。

---

## バッファ付き channel で Semaphore

```go
// 同時に3つまで
sem := make(chan struct{}, 3)

for _, task := range tasks {
    sem <- struct{}{}  // トークンを取得(空きがなければ待つ)

    go func(t Task) {
        defer func() { <-sem }()  // 終わったらトークンを返す
        process(t)
    }(task)
}
```

---

## Semaphore の動き

```
バッファサイズ: 3

task1 開始 → sem: [●][_][_]
task2 開始 → sem: [●][●][_]
task3 開始 → sem: [●][●][●]
task4 開始 → 待機...(空きがない)
task1 終了 → sem: [_][●][●]
task4 開始 → sem: [●][●][●]
```

バッファが「許可証」の役割を果たす。

---

## ワーカープールとの違い

ワーカープール
- 固定数のワーカーを先に起動
- ワーカーが仕事を取りに行く

Semaphore
- goroutine は都度起動
- 起動前に許可を取る

Semaphore の方がシンプルだが、goroutine の起動コストがかかる。

---

## パターン8: Rate Limiting

---

## Rate Limiting とは

単位時間あたりの処理数を制限すること。

例:「1秒に10リクエストまで」

APIのレート制限を守るときなどに使う。

---

## time.Tick を使った Rate Limiting

```go
// 100ms ごとに1つ処理 = 1秒に10個
rate := time.Tick(100 * time.Millisecond)

for _, req := range requests {
    <-rate  // 100ms 経つまで待つ
    go process(req)
}
```

`time.Tick(d)` は、一定間隔で値を送り続ける channel を返す。

---

## バースト対応の Rate Limiting

最初の数個は即座に処理し、その後は制限をかけたい場合。

```go
// バースト: 最初の3個は即座に処理可能
burstyLimiter := make(chan time.Time, 3)
for i := 0; i < 3; i++ {
    burstyLimiter <- time.Now()
}

// その後は 200ms ごとに補充
go func() {
    for t := range time.Tick(200 * time.Millisecond) {
        burstyLimiter <- t
    }
}()
```

---

## バースト対応の使い方

```go
for _, req := range requests {
    <-burstyLimiter  // トークンを取得
    go process(req)
}
```

最初の3個は即座に処理され、4個目以降は 200ms 間隔になる。

---

## パターン9: context.Context

---

## context.Context とは

Go 1.7 で導入された、キャンセル・タイムアウト・値の受け渡しを統合した仕組み。

Done channel + Timeout + 値の受け渡しをまとめたもの。

---

## context の基本

```go
// 空の context を作る
ctx := context.Background()

// キャンセル可能な context
ctx, cancel := context.WithCancel(context.Background())

// タイムアウト付き context
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

// デッドライン付き context
ctx, cancel := context.WithDeadline(context.Background(), deadline)
```

---

## context を使ったキャンセル

```go
func worker(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            fmt.Println("キャンセル:", ctx.Err())
            return
        default:
            // 作業
        }
    }
}
```

`ctx.Done()` は、キャンセルされると close される channel を返す。

 参考: [context package - pkg.go.dev](https://pkg.go.dev/context) | [Go blog - Context](https://go.dev/blog/context)

---

## context の伝播

```go
func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    result, err := fetchData(ctx)
}

func fetchData(ctx context.Context) (Data, error) {
    // 子の処理にも ctx を渡す
    return callAPI(ctx, url)
}
```

慣習: context は関数の第1引数に渡す。

---

## context を使うべき場面

- HTTP ハンドラ(リクエストごとにキャンセル可能に)
- データベースクエリ
- 外部 API 呼び出し
- 長時間実行されるバックグラウンド処理

標準ライブラリの多くの関数が context を受け取る設計になっている。

---

## context 使用時の注意

```go
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()  // 必ず呼ぶ
```

- cancel は必ず呼ぶ: リソースリークを防ぐ
- context を struct に入れない: 関数の引数で渡す
- nil context を渡さない: context.TODO() を使う

---

## パターンの選び方

---

## やりたいこと → パターン対応表

| やりたいこと | パターン |
|------------|---------|
| データソースを抽象化 | Generator |
| 処理を段階に分ける | Pipeline |
| 重い処理を並列化 | Fan-out / Worker Pool |
| 複数ソースを統合 | Fan-in |
| 複数 channel を待つ | select |
| 処理をキャンセル | Done channel / context |
| 時間制限を設ける | Timeout / context |
| 同時実行数を制限 | Semaphore / Worker Pool |
| リクエスト頻度を制限 | Rate Limiting |

---

## 実務での組み合わせ例

Web クローラー
```
URL Generator → Worker Pool(Fan-out) → 結果収集(Fan-in)
                    ↑
            context でタイムアウト
            Semaphore で同時接続数制限
```

ログ処理パイプライン
```
ファイル読み込み → パース → フィルタ → 集計
    Generator      Pipeline stages
         ↑
    context でキャンセル可能に
```

---

## パターン選択のコツ

1. まず問題を明確に: 何を並行化したいのか
2. シンプルに始める: いきなり複雑なパターンを使わない
3. 組み合わせる: 1つのパターンで解決しなくていい
4. context を活用: キャンセルとタイムアウトは context で統一

---

# ハンズオン

---

## 4つの Phase

Phase 1 逐次処理
　goroutineを使わずに、まず動くものを作る → 基準タイム

Phase 2 並行処理
　goroutine + channel を使う。

Phase 3 ワーカープール
　固定数のgoroutineで処理。(Go 1.25の WaitGroup.Go() を活用)

Phase 4 さらなる高速化
　制約なし

 参考: [Go 1.25 Release Notes](https://go.dev/doc/go1.25) | [WaitGroup.Go - pkg.go.dev](https://pkg.go.dev/sync#WaitGroup.Go)

---

## ルール

- 2人1組で進める
- 改善率で競う(PCスペック差を吸収)
- 隣に聞いてOK、教え合い推奨
- 困ったら手を挙げて

---
