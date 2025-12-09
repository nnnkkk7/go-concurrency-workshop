# Goで学ぶ並行処理：アクセスログ解析チャレンジ

実務でよくある「大量のログファイルを急いで解析してほしい」という状況を題材に、Goの並行処理を学ぶハンズオンワークショップです。

##  このワークショップで学べること

- **goroutineとchannelの実践的な使い方**
- **段階的な改善アプローチ**：逐次処理 → 並行化 → 最適化
- **ワーカープールパターン**の実装
- **Go 1.25の新機能** `WaitGroup.Go()`の活用
- **パフォーマンスチューニング**の手法

##  必要な環境

- **Go 1.25以降**（WaitGroup.Go()を使用）
- エディタ（VS Code推奨）
- ターミナル

##  クイックスタート

### 1. リポジトリをクローン

```bash
git clone https://github.com/nnnkkk7/go-concurrency-workshop.git
cd go-concurrency-workshop
```

### 2. ログファイルを生成

```bash
go run cmd/loggen/main.go
```

これで `logs/` ディレクトリに50個のログファイル（合計約500MB）が生成されます。

### 3. ワークショップを開始

```bash
# まずは workshop/ で挑戦
go run ./workshop/phase1/main.go

# 詰まったら hint/ を参照
go run ./hint/phase1/main.go
```

##  ワークショップの進め方

このワークショップでは **2つの難易度レベル** から選べます：

###  まずは workshop/ で挑戦

各Phaseを自分で実装していきます。シンプルなTODOコメントのみが付いています。

```bash
go run ./workshop/phase1/main.go
go run ./workshop/phase2/main.go
go run ./workshop/phase3/main.go
go run ./workshop/phase4/main.go  # 自由課題
```

詳しくは [workshop/README.md](workshop/README.md) を参照。

###  詰まったら hint/ を参照

`workshop/` で詰まったら、`hint/` ディレクトリに詳細なヒント付きコードがあります。

```bash
go run ./hint/phase1/main.go
go run ./hint/phase2/main.go
go run ./hint/phase3/main.go
```

詳しくは [hint/README.md](hint/README.md) を参照。

---

### Phase 1: 逐次処理版（15分）

まずは並行処理を使わず、シンプルなfor文で実装します。
これが基準値になります。

### Phase 2: 基本並行処理版（20分）

goroutineとchannelを使って並行処理化します。
目標：5〜10倍の改善

### Phase 3: ワーカープール版（17分）

ワーカープールパターンで最適化します。
**Go 1.25の新機能 `WaitGroup.Go()` を活用！**

### Phase 4: さらなる高速化（自由課題）

Phase 3を超える最適化に挑戦します。
**制約なし - あらゆる手法を試してみましょう！**

##  成果物

各Phaseの模範解答は `solutions/` ディレクトリにあります：

- [solutions/phase1/main.go](solutions/phase1/main.go) - 逐次処理版
- [solutions/phase2/main.go](solutions/phase2/main.go) - 基本並行処理版
- [solutions/phase3/main.go](solutions/phase3/main.go) - ワーカープール版
- [solutions/phase4/main.go](solutions/phase4/main.go) - 最適化版

##  ドキュメント

### 参加者向け
- [workshop/README.md](workshop/README.md) - 実装ガイド
- [hint/README.md](hint/README.md) - ヒント付きコードの使い方
- [docs/HINTS.md](docs/HINTS.md) - レベル別ヒント集

### 講師向け
- [docs/FACILITATOR_GUIDE.md](docs/FACILITATOR_GUIDE.md) - 講師用ガイド
- [docs/SLIDES.md](docs/SLIDES.md) - 講義スライド
- [docs/SETUP_CHECK.md](docs/SETUP_CHECK.md) - 環境確認チェックリスト

##  プロジェクト構成

```
go-concurrency-workshop/
├── cmd/
│   └── loggen/          # ログ生成ツール
├── pkg/
│   └── logparser/       # ログパース共通処理
├── workshop/            # 実装用: シンプルなTODOのみ
│   ├── phase1/
│   ├── phase2/
│   ├── phase3/
│   └── phase4/
├── hint/                # ヒント: 詳細なコメント付き
│   ├── phase1/
│   ├── phase2/
│   └── phase3/
├── solutions/           # 模範解答
│   ├── phase1/
│   ├── phase2/
│   ├── phase3/
│   └── phase4/
├── logs/                # 生成されたログファイル
└── docs/                # ドキュメント
```

##  よくある質問

### Q: Go 1.25未満でも動きますか？

A: Phase 3の `WaitGroup.Go()` はGo 1.25の新機能です。それ以前のバージョンでは、従来の `wg.Add(1)` + `go func()` パターンを使ってください。

### Q: 処理時間が期待より短い/長い

A: マシンスペックによって絶対値は大きく変わります。重要なのは「改善率」です。Phase 1の時間を基準に、何倍速くなったかを測定しましょう。

### Q: ログファイルのサイズを変更したい

A: 以下のオプションが使えます：

```bash
go run cmd/loggen/main.go --files=100 --lines=50000
```

##  学習リソース

- [Go Concurrency Patterns](https://antonz.org/go-concurrency/goroutines/) - goroutinesの詳しい解説
- [Go 1.25 Release Notes](https://tip.golang.org/doc/go1.25) - 最新機能の公式ドキュメント

##  ライセンス

MIT License

##  質問・フィードバック

Issuesでお気軽にどうぞ！