# ヒント集

詰まったときに参照してください。

---

## Phase 1: 逐次処理版

### やること

50個のログファイルを1つずつ順番に処理して、ステータスコード別に集計する。

まずファイル一覧を取得（`filepath.Glob`を使う）。次に各ファイルをループで処理。最後に結果を統合して表示。

### ファイル処理の基本

1つのファイルを読む関数を作る。ファイル名を受け取って、結果とエラーを返す。

やることは単純：
- ファイルを開く（`os.Open`）
- 読み終わったら閉じる（`defer file.Close()`で確実に）
- JSON を1行ずつ読んで集計
- 結果を返す

### JSONの読み方

`encoding/json` の `NewDecoder` を使うと、ファイル全体を読み込まなくても1行ずつ処理できる。

```go
decoder := json.NewDecoder(file)
for decoder.More() {
    // 1行読んでresult.AddEntry()
}
```

不正な行があってもスキップして続ける。

### 結果の集め方

スライスを作って、ループで各ファイルの結果を追加。全部終わったら `logparser.MergeResults()` で統合。

エラーが出たファイルは標準エラーに出力して次へ。

---

## Phase 2: 並行処理版

### 考え方

Phase 1は1つずつ順番に処理した。Phase 2では全ファイルを同時に処理する。

goroutineを使って各ファイルを別々に処理。結果はchannelで集める。

### channelの準備

結果を集めるchannelを作る。バッファ付きにすると、goroutineが結果を送るときに待たなくていい。

```go
results := make(chan *logparser.Result, len(files))
```

バッファサイズはファイル数と同じにしておけば、デッドロックの心配がない。

### goroutineを起動

各ファイルに対してgoroutineを起動。

```go
for _, filename := range files {
    go func(name string) {
        // processFile(name)
        // results <- result
    }(filename)
}
```

**重要:** ループ変数をそのまま使うと全部同じファイルを処理してしまう。引数で渡す。

### 結果の回収

ファイル数だけループして受信。

```go
for i := 0; i < len(files); i++ {
    result := <-results
    // 集める
}
```

または、WaitGroupで全goroutineの終了を待ってからchannelをclose、`for range`で受信する方法もある。

---

## Phase 3: ワーカープール版

### なぜワーカープール？

Phase 2は50ファイルに50個のgoroutineを作った。でも5000ファイルだったら？

goroutineは軽いけど、数万個作るとメモリを使う。ファイルも同時に開きすぎると上限に引っかかる。

ワーカープールは固定数のgoroutine（例えば8個）を作って、仕事を順番に処理させる。

### 構造

2つのchannelを使う：
- `jobs`: ファイル名を入れる（仕事の待ち行列）
- `results`: 処理結果を入れる

固定数のワーカーが `jobs` から仕事を取り出して処理、結果を `results` に送る。

### ワーカー数

`runtime.NumCPU()` でCPU数を取得。ファイル読み込みとJSON解析が主な処理なので、CPU数くらいが妥当。

### ワーカーの動き

各ワーカーはループで仕事を取り出して処理。

```go
for filename := range jobs {
    // 処理してresultsに送る
}
```

`jobs` がcloseされたらループが終わる。

Go 1.25なら `wg.Go()` が便利。従来の `wg.Add(1)` + `go func()` + `defer wg.Done()` を1行で書ける。

### 仕事の投入

全ファイル名を `jobs` に送ったら、必ず `close(jobs)` する。これで「もう仕事はない」とワーカーに伝わる。

### 結果を集める

別のgoroutineで `wg.Wait()` して、全ワーカーが終わったら `close(results)` する。

メインは `for result := range results` で全部受け取る。

---

## よくあるエラー

### デッドロック

「all goroutines are asleep」が出たら、全員が待ち状態で進めない。

原因：

- channelに送ったけど誰も受け取らない
- channelから受け取ろうとしたけど誰も送らない

対策：

- 送信数と受信数を合わせる
- バッファ付きchannelを使う
- 送信と受信を別のgoroutineで

### 同じファイルばかり処理される

ループ変数をgoroutineで直接使うと、全部最後の値になる。

```go
// ダメ
for _, file := range files {
    go func() { process(file) }()
}

// OK
for _, file := range files {
    go func(f string) { process(f) }(file)
}
```

### プログラムが終わらない

原因1: channelをcloseし忘れ。`for range channel` は閉じられるまで待つ。

原因2: WaitGroupの `Done()` 忘れ。`defer wg.Done()` を書くか、`wg.Go()` を使う。

### panic: send on closed channel

closeしたchannelに送信しようとした。

原則：送信側がcloseする。受信側はcloseしない。

---

## さらに速くするには

### バッファサイズを変えてみる

channelのバッファを変えると速度が変わるかも。いろいろ試してみる。

### ワーカー数を調整

CPU数の半分、2倍、4倍など試して、どれが速いか測る。

### 環境変数で制御

`NUM_WORKERS` 環境変数があればその値を使う、なければCPU数。実行時に調整できて便利。

### contextでキャンセル

Ctrl+Cで処理を中断できるようにする。`context.WithCancel()` と `signal.Notify()` を調べてみる。

---

## 参考資料

- [Effective Go - Concurrency](https://go.dev/doc/effective_go#concurrency)
- [Go Concurrency Patterns](https://antonz.org/go-concurrency/goroutines/)
- [Go 1.25 Release Notes](https://tip.golang.org/doc/go1.25)

詰まったら `solutions/` を見てもOK。頑張ってください！
