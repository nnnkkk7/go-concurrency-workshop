# 環境確認チェックリスト

## 1. Go のインストール確認

### 1.1 Go 1.25 以降がインストールされているか

```bash
go version
```

**期待される出力:**
```
go version go1.25.0 darwin/arm64
```

または

```
go version go1.25.x ...
```

**確認ポイント:**
- [ ] `go version`が表示される
- [ ] バージョンが1.25.0以降である

**問題がある場合:**
- [Go公式サイト](https://go.dev/dl/)からインストール


## 2. リポジトリのクローン確認

### 2.1 リポジトリをクローン

```bash
git clone https://github.com/nnnkkk7/go-concurrency-workshop.git
cd go-concurrency-workshop
```

**確認ポイント:**
- [ ] クローンが成功した
- [ ] ディレクトリに移動できた


---

## 3. ログファイルの生成

### 3.1 ログ生成ツールの実行

```bash
go run cmd/loggen/main.go
```

**期待される出力:**
```
Generating log files...
  [1/50] access_001.log (67000 lines, 10.2MB)
  [2/50] access_002.log (67000 lines, 10.1MB)
  ...
  [50/50] access_050.log (67000 lines, 10.3MB)

Done! Generated 50 files (512MB total) in 5.2s
```


**確認ポイント:**
- [ ] エラーなく完了した

### 3.2 ログファイルの確認

```bash
ls logs/*.log | wc -l
```

**期待される出力:**
```
50
```

**確認ポイント:**
- [ ] 50ファイルが存在する

---

## 4. 記録

### 4.1 処理時間の記録

Phase 1の基準値として、Phase1の実装が終わったあとに、処理時間を記録しておきましょう。

```
Phase 1の処理時間: _______ 秒
```
