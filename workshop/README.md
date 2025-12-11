# Workshop 実装ディレクトリ

このディレクトリは、ワークショップ参加者が実際にコードを書いて学ぶための作業スペースです。

## 前提条件

ワークショップを始める前に、[docs/SETUP_CHECK.md](../docs/SETUP_CHECK.md) のチェックリストを完了してください。

特に以下を確認：
- Go 1.25以降がインストールされている
- ログファイル（50個、約500MB）が生成済み
- スターターコードが正常に動作する

##  ディレクトリ構成

```
workshop/
├── phase1/
│   └── main.go    # Phase 1: 逐次処理版を実装
├── phase2/
│   └── main.go    # Phase 2: 並行処理版を実装
├── phase3/
│   └── main.go    # Phase 3: ワーカープール版を実装
└── phase4/
    └── main.go    # Phase 4: さらなる高速化に挑戦
```


##  使い方

### Phase 1: 逐次処理版

```bash
# リポジトリのルートから実行
go run ./workshop/phase1/main.go
```

### Phase 2: 並行処理版

```bash
# リポジトリのルートから実行
go run ./workshop/phase2/main.go
```


### Phase 3: ワーカープール版

```bash
# リポジトリのルートから実行
go run ./workshop/phase3/main.go
```


### Phase 4: さらなる高速化

```bash
# リポジトリのルートから実行
go run ./workshop/phase4/main.go
```


Phase 3よりもさらに高速化する。あらゆる最適化手法を試してください。

##  ヒントが必要な場合

詰まったときは以下を参照してください。

- **[docs/HINTS.md](../docs/HINTS.md)** - ヒント
- **[solutions/](../solutions/)** - 模範解答


##  パフォーマンス測定

各Phaseの処理時間を記録しましょう。

```
Phase 1: _____ 秒（基準値）
Phase 2: _____ 秒（改善率: _____倍）
Phase 3: _____ 秒（改善率: _____倍）
Phase 4: _____ 秒（改善率: _____倍）
```

頑張ってください！ 
