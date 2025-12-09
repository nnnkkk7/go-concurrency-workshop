# Workshop 実装ディレクトリ

このディレクトリは、ワークショップ参加者が実際にコードを書いて学ぶための作業スペースです。

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

各Phaseは独立したディレクトリになっているため、**前のPhaseのコードを残したまま次のPhaseに進める**ようになっています。

##  使い方

### Phase 1: 逐次処理版（15分）

```bash
# リポジトリのルートから実行
go run ./workshop/phase1/main.go
```

**実装する関数：** `processFilesSequentially()`

### Phase 2: 並行処理版（20分）

```bash
# リポジトリのルートから実行
go run ./workshop/phase2/main.go
```

**実装する関数：** `processFilesConcurrently()`

### Phase 3: ワーカープール版（17分）

```bash
# リポジトリのルートから実行
go run ./workshop/phase3/main.go
```

**実装する関数：** `processFilesWithWorkerPool()`

### Phase 4: さらなる高速化（自由課題）

```bash
# リポジトリのルートから実行
go run ./workshop/phase4/main.go
```

**実装する関数：** `processFiles()`

**目標：** Phase 3よりもさらに高速化する
**制約：** なし（あらゆる最適化手法を試してください）

##  ヒントが必要な場合

詰まったときは以下を参照してください：

- **[hint/](../hint/)** - より詳細なTODOコメント付きのコード
- **[docs/HINTS.md](../docs/HINTS.md)** - レベル別のヒント集
- **[solutions/](../solutions/)** - 模範解答

##  パフォーマンス測定

各Phaseの処理時間を記録しましょう：

```
Phase 1: _____ 秒（基準値）
Phase 2: _____ 秒（改善率: _____倍）
Phase 3: _____ 秒（改善率: _____倍）
Phase 4: _____ 秒（改善率: _____倍）
```

##  目標

- Phase 1をベースラインとして処理時間を計測
- Phase 2で5〜10倍の高速化を実現
- Phase 3でより安定した性能を実現
- Phase 4でPhase 3を超える最適化に挑戦

頑張ってください！ 
