# Goで学ぶ並行処理ワークショップ

実務でよくある「大量のログファイルを急いで解析してほしい」という状況を題材に、Goの並行処理を学ぶハンズオンワークショップです。

## 1. 必要な環境

環境構築は [docs/SETUP_CHECK.md](docs/SETUP_CHECK.md) のチェックリストを完了してください。


##  2. ワークショップの進め方

**[workshop/README.md](workshop/README.md)** を参照してください。


##  ドキュメント

- [workshop/README.md](workshop/README.md) - ワークショップ実施ガイド（各Phaseの詳細な説明）
- [docs/SETUP_CHECK.md](docs/SETUP_CHECK.md) - 環境構築方法とチェックリスト
- [docs/HINTS.md](docs/HINTS.md) - 実装のヒント集


##  プロジェクト構成

```
go-concurrency-workshop/
├── cmd/loggen/          # ログ生成ツール
├── pkg/logparser/       # ログパース共通処理
├── workshop/            # 実装用
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
└── logs/                # 生成されたログファイル
```


### ログファイルのサイズを変更したい

以下のオプションが使えます。

```bash
go run cmd/loggen/main.go --files=100 --lines=50000
```

### Make コマンド

```bash
make help           # コマンド一覧を表示
make gen            # ログファイルを生成
make w1 w2 w3 w4    # Workshop Phase 1-4 を実行
make s1 s2 s3 s4    # Solution Phase 1-4 を実行
```

##  ライセンス

MIT License
