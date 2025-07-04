# Milestone 1: Minimal Calculator

## Goal
四則演算ができる最小限のインタープリタを実装する。

## Success Criteria
以下の式が評価できること:
- `1 + 2` → `3`
- `10 - 5` → `5`
- `3 * 4` → `12`
- `8 / 2` → `4`
- `2 + 3 * 4` → `14` (演算子の優先順位)
- `(2 + 3) * 4` → `20` (括弧による優先順位制御)

## Implementation Steps

### Step 1: Project Setup ✅
- [x] Git リポジトリの初期化
- [x] 開発環境の設定（Prettier, aqua）
- [x] ロードマップの作成

### Step 2: Lexer (字句解析器)
- [ ] トークンの定義（NUMBER, PLUS, MINUS, STAR, SLASH, LPAREN, RPAREN, EOF）
- [ ] Lexer 構造体の実装
- [ ] NextToken() メソッドの実装
- [ ] 数値のパース
- [ ] 演算子と括弧の認識
- [ ] 空白文字のスキップ
- [ ] テストの作成

### Step 3: Parser (構文解析器)
- [ ] AST ノードの定義（NumberNode, BinaryOpNode）
- [ ] Parser 構造体の実装
- [ ] 再帰下降パーサーの実装
  - [ ] ParseExpression() - 加減算
  - [ ] ParseTerm() - 乗除算
  - [ ] ParseFactor() - 数値と括弧
- [ ] テストの作成

### Step 4: Evaluator (評価器)
- [ ] Eval() 関数の実装
- [ ] 各 AST ノードの評価
- [ ] エラーハンドリング（ゼロ除算など）
- [ ] テストの作成

### Step 5: REPL
- [ ] 標準入力からの読み込み
- [ ] Lexer → Parser → Evaluator の統合
- [ ] 結果の出力（fmt なしで os.Stdout.Write を使用）
- [ ] エラーメッセージの表示
- [ ] 終了コマンド（exit または Ctrl+D）

### Step 6: Integration & Polish
- [ ] main.go の実装
- [ ] 全体の動作確認
- [ ] README.md の作成
- [ ] サンプル式のドキュメント化

## Technical Notes

- **パッケージ構成**: Phase 1 では全て main パッケージに含める
- **標準ライブラリの使用**: 
  - `os` (ファイル I/O)
  - `bufio` (入力読み取り)
  - `strconv` (数値変換)
  - fmt は使わない（自前で出力実装）
- **エラーハンドリング**: panic ではなく、エラーを返す設計

## Testing Strategy

各コンポーネントに対してテストを書く:
1. `lexer_test.go` - トークン化のテスト
2. `parser_test.go` - AST 構築のテスト
3. `eval_test.go` - 評価結果のテスト

## Next Steps
Milestone 1 完了後、Phase 2（変数とステートメント）に進む。