# Hint ディレクトリ

このディレクトリには、**詳細なTODOコメント付き**のコードが含まれています。

`workshop/` でチャレンジしてみて、詰まったときにこちらを参照してください。

##  ディレクトリ構成

```
hint/
├── phase1/
│   └── main.go    # Phase 1: 詳細なヒント付き
├── phase2/
│   └── main.go    # Phase 2: 詳細なヒント付き
└── phase3/
    └── main.go    # Phase 3: 詳細なヒント付き
```

##  使い方

### まずは workshop/ で挑戦

1. `workshop/phase1/` で自力で実装を試みる
2. 詰まったらこちら `hint/phase1/` を見る
3. より詳細なTODOコメントを参考に実装を進める

### 実行方法

```bash
# リポジトリのルートから実行
go run ./hint/phase1/main.go
go run ./hint/phase2/main.go
go run ./hint/phase3/main.go
```

### hint/ の特徴

- **段階的なヒント**: 実装の各ステップにコメントがある
- **コード例**: 具体的な実装パターンが示されている
- **比較説明**: 複数のアプローチの違いが説明されている

##  推奨フロー

```
workshop/ で挑戦
    ↓
詰まった？
    ↓
hint/ を参照
    ↓
それでも分からない？
    ↓
docs/HINTS.md を参照
    ↓
最後の手段
    ↓
solutions/ を参照
```

##  その他の参考資料

- **[docs/HINTS.md](../docs/HINTS.md)** - レベル別のヒント集（コンセプトレベル）
- **[solutions/](../solutions/)** - 模範解答（完全な実装）
- **[docs/FACILITATOR_GUIDE.md](../docs/FACILITATOR_GUIDE.md)** - よくあるバグと対処法

## Phase 別のヒント概要

### Phase 1 のヒント

- ファイルのループ処理
- エラーハンドリング
- スライスへの追加

### Phase 2 のヒント

- channelの作成（バッファ付き）
- goroutineの起動
- 変数キャプチャの注意点
- channelからの結果収集
- 代替実装（WaitGroup + Mutex）

### Phase 3 のヒント

- ワーカー数の決定
- ジョブchannelの作成
- 結果channelの作成
- ワーカーgoroutineの起動
- Go 1.25の `WaitGroup.Go()` の使い方
- Go 1.24以前との比較

頑張ってください！ 
