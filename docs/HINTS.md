# ヒント集

このファイルは、workshop各Phaseの実装で詰まったときに参照してください。

---

## Phase 1: 逐次処理版

### 何をするか

ログファイルを1つずつ順番に処理して、全体の結果を集計します。

### 必要な要素

- 結果を格納するスライス
- ファイルを1つずつ処理するループ
- エラーハンドリング（エラーが出ても処理を続ける）
- 結果の収集

### 実装の流れ

1. 結果を入れる空のスライスを作る
2. ファイル一覧をループで回す
3. 各ファイルを processFile() で処理
4. エラーが出たら標準エラーに出力して次へ
5. 成功したら結果をスライスに追加
6. 全部終わったらスライスを返す

---

## Phase 2: 並行処理版

### 何をするか

Phase 1では1つずつ順番に処理しましたが、Phase 2では全ファイルを同時に処理します。
goroutineを使って各ファイルを別々に処理し、結果をchannelで集めます。

### 必要な要素

- 結果を送受信するchannel
- 全goroutineの完了を待つWaitGroup
- 各ファイルに対するgoroutine
- channelのclose処理

### 実装の流れ


1. バッファ付きchannelを作る（ファイル数分）
2. WaitGroupを作る
3. 各ファイルに対してgoroutineを起動
   - WaitGroupにAdd(1)
   - goroutine内でdefer Done()
   - processFile()を呼び出し
   - 結果をchannelに送信
4. 別のgoroutineでWaitGroupを待ち、完了したらchannelをclose
5. channelから結果を全部受け取る
`

---

## Phase 3: ワーカープール版

### 何をするか

Phase 2では全ファイルに対してgoroutineを作りましたが、Phase 3では固定数のワーカー（goroutine）を作り、
それらがファイルを順番に処理していきます。これにより大量のファイルでもgoroutine数を制御できます。

Go 1.25の新機能 `WaitGroup.Go()` を使います。

### 必要な要素

- ジョブ（ファイル名）を配布するchannel
- 結果を集めるchannel
- 固定数のワーカーgoroutine
- ワーカー数の決定（`runtime.NumCPU()`）
- ジョブのchannel投入とclose
- 結果のchannel受信

### 実装の流れ


1. jobsとresultsの2つのchannelを作る
2. WaitGroupを作る
3. ワーカー数分（CPU数分）だけgoroutineを起動
   - 各ワーカーはjobsからファイル名を取り出して処理
   - 結果をresultsに送信
4. 全ファイル名をjobsに投入してclose
5. 別のgoroutineで全ワーカーの完了を待ち、resultsをclose
6. resultsから結果を全部受け取る


## 参考資料

### Go公式ドキュメント

- [Effective Go - Concurrency](https://go.dev/doc/effective_go#concurrency)
- [A Tour of Go - Concurrency](https://go.dev/tour/concurrency/1)
- [Go 1.25 Release Notes](https://tip.golang.org/doc/go1.25)

### 並行処理パターン

- [Go Concurrency Patterns](https://go.dev/blog/pipelines)
- [Advanced Go Concurrency Patterns](https://go.dev/blog/io2013-talk-concurrency)

### ワーカープール

- [Worker Pool Pattern](https://gobyexample.com/worker-pools)
- [Go by Example: Worker Pools](https://gobyexample.com/worker-pools)

---

## 最後に
詰まったら気軽に `solutions/` を参照してください。学習が目的なので、完璧を目指す必要はありません！
