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

> **📋 環境構築の詳細確認:** ワークショップ当日までに [docs/SETUP_CHECK.md](docs/SETUP_CHECK.md) のチェックリストを完了してください。

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

> **✅ 環境確認:** 正常に動作するか確認したい場合は [docs/SETUP_CHECK.md](docs/SETUP_CHECK.md) を参照してください。

### 3. ワークショップを開始

```bash
# workshop/ で実装に挑戦
go run ./workshop/phase1/main.go

# 詰まったら docs/HINTS.md を参照
```

##  ワークショップの進め方

### workshop/ で実装に挑戦

各Phaseを自分で実装していきます。シンプルなTODOコメントのみが付いています。

```bash
go run ./workshop/phase1/main.go
go run ./workshop/phase2/main.go
go run ./workshop/phase3/main.go
go run ./workshop/phase4/main.go  # 自由課題
```

詳しくは [workshop/README.md](workshop/README.md) を参照。

### 詰まったら docs/HINTS.md を参照

実装に詰まったら、[docs/HINTS.md](docs/HINTS.md) にレベル別のヒントがあります：

- **Level 1-2**: 全体の流れと基本的なアプローチ
- **Level 3-4**: 具体的なコード例とパターン
- **Level 5-6**: 代替実装や詳細な解説
- **Level 7**: 完全な実装例へのリファレンス

---

### Phase 1: 逐次処理版（15分）

まずは並行処理を使わず、シンプルなfor文で実装します。
これが基準値になります。

### Phase 2: 基本並行処理版（20分）

goroutineとchannelを使って並行処理化します。

### Phase 3: ワーカープール版（17分）

ワーカープールパターンで最適化します。

### Phase 4: さらなる高速化（自由課題）

Phase 3を超える最適化に挑戦します。


各Phaseの模範解答は `solutions/` ディレクトリにあります：

- [solutions/phase1/main.go](solutions/phase1/main.go) - 逐次処理版
- [solutions/phase2/main.go](solutions/phase2/main.go) - 基本並行処理版
- [solutions/phase3/main.go](solutions/phase3/main.go) - ワーカープール版


##  ドキュメント

- [workshop/README.md](workshop/README.md) - 実装ガイド
- [docs/HINTS.md](docs/HINTS.md) - レベル別ヒント集


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
├── solutions/           # 模範解答
│   ├── phase1/
│   ├── phase2/
│   ├── phase3/
│   └── phase4/
├── docs/                # ドキュメント
│   ├── HINTS.md         # レベル別ヒント集
│   ├── FACILITATOR_GUIDE.md
│   └── SLIDES.md
└── logs/                # 生成されたログファイル
```


### ログファイルのサイズを変更したい

以下のオプションが使えます。

```bash
go run cmd/loggen/main.go --files=100 --lines=50000
```

##  ライセンス

MIT License
