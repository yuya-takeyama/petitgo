# petitgo Development Guide

## プロジェクト概要

petitgo は小さな Go 実装で、最終的にはセルフホスト（自分自身でコンパイルできる）を目指すプロジェクトです。

## ドキュメント構成

- **@docs/ROADMAP.md**: 全体の開発ロードマップ、9 つの Phase の詳細計画
- **@docs/MILESTONE-1.md**: Phase 1 の具体的な実装ステップと Success Criteria
- **@docs/TESTING-STRATEGY.md**: テスト戦略、Go ライブラリテストから最終的なセルフテストまで
- **@CLAUDE.md**: このファイル、開発ガイドと現在の進捗状況

## Test-Driven Development (TDD) - t-wada style

このプロジェクトでは TDD を実践します。特に t-wada さんが提唱するスタイルで：

1. **まず失敗するテストを書く（RED）**

   - 最小限のテストケースから始める
   - エラーメッセージが分かりやすいように書く

2. **テストをパスする最小限の実装（GREEN）**

   - とにかくテストを通すことを優先
   - ハードコーディングでも OK

3. **リファクタリング（REFACTOR）**
   - テストが通る状態を保ちながら改善
   - 重複を除去、設計を改善

## 開発フロー

1. テストファイルを先に作る（`*_test.go`）
2. `go test` で RED を確認
3. 実装ファイルを作って GREEN に
4. コミット（テストと実装を一緒に）

## コミット方針

- 作業単位ごとに細かくコミット
- 「ポケモンのレポート」のように、こまめにセーブ
- TDD のサイクルごとにコミットするのも OK

## コーディング規約

### コメント

- **全てのコメントは英語で記述する**
- コードの意図や理由を明確に説明する
- TODO や FIXME も英語で記述
- 例: `// Parse the if statement with condition and optional else block`

## 現在の開発状況

### Phase 1: Minimal Calculator ✅ 完了！

- [x] Lexer の実装（完了）
  - [x] Go 本家互換の Token 型定義
  - [x] 数値（INT）のトークン化
  - [x] 演算子（ADD, SUB, MUL, QUO）のトークン化
  - [x] 括弧（LPAREN, RPAREN）のトークン化
  - [x] 空白文字のスキップ処理
  - [x] 包括的なテストスイート
- [x] Parser の実装（完了）
  - [x] AST ノードの定義
  - [x] 再帰下降パーサーの実装
  - [x] 演算子優先順位の処理
  - [x] 括弧サポート
  - [x] 包括的なテストスイート
- [x] Evaluator の実装（完了）
  - [x] AST ノードの評価処理
  - [x] 四則演算のサポート
  - [x] 統合テストの実装
- [x] REPL の実装（完了）
  - [x] 対話型計算機
  - [x] fmt なしのカスタム出力関数
  - [x] exit コマンドサポート
- [x] 全体統合（完了）
  - [x] main.go の実装
  - [x] 動作確認完了

**Success Criteria 達成状況:**

- ✅ `1 + 2` → `3`
- ✅ `10 - 5` → `5`
- ✅ `3 * 4` → `12`
- ✅ `8 / 2` → `4`
- ✅ `2 + 3 * 4` → `14` (演算子優先順位)
- ✅ `(2 + 3) * 4` → `20` (括弧による優先順位制御)

### Phase 2: Variables and Basic Statements ✅ 完了！

- [x] Lexer の拡張（完了）
  - [x] 識別子（IDENT）トークンの実装
  - [x] 代入演算子（ASSIGN: :=, =）の実装
  - [x] 文字と数字の識別機能
- [x] Parser の拡張（完了）
  - [x] Statement インターフェースの定義
  - [x] VarStatement（var x int = 42）の実装
  - [x] AssignStatement（x := 42, x = 42）の実装
  - [x] VariableNode（変数参照）の実装
  - [x] ParseStatement メソッドの実装
- [x] Environment の実装（完了）
  - [x] 変数管理システム（Set/Get）
  - [x] スコープ管理の基盤
  - [x] EvalStatement 関数
  - [x] EvalWithEnvironment 関数
- [x] REPL の拡張（完了）
  - [x] Statement と Expression の自動判定
  - [x] 変数の永続化
  - [x] 対話的な変数操作

**Success Criteria 達成状況:**

- ✅ `var x int = 42` → 変数宣言
- ✅ `x := 10` → 短縮変数宣言
- ✅ `x = 20` → 代入文
- ✅ `x + y` → 変数参照と演算
- ✅ 複数文の実行とスコープ管理

### Phase 3: Control Flow ✅ 完了！

- [x] Lexer の拡張（完了）
  - [x] 比較演算子（==, !=, <, >, <=, >=）の実装
  - [x] 論理演算子（&&, ||, !）の実装
  - [x] 制御構文キーワード（if, else, for, break, continue）の実装
  - [x] ブロック構文（{ }）の実装
- [x] Parser の拡張（完了）
  - [x] IfStatement（if/else 文）の実装
  - [x] ForStatement（condition-only for 文）の実装
  - [x] BlockStatement（ブロック文）の実装
  - [x] ExpressionStatement（式文）の実装
  - [x] CallNode（関数呼び出し）の実装
- [x] Evaluator の拡張（完了）
  - [x] 比較演算の評価（1/0 で true/false）
  - [x] if/else 条件分岐の実装
  - [x] 基本的な for 文（condition-only）の実装
  - [x] print()関数のサポート
- [x] REPL の拡張（完了）
  - [x] 制御構文をステートメントとして認識
  - [x] 新機能の説明を追加

**Success Criteria 達成状況:**

- ✅ `if x > 5 { print(x) }` → 条件分岐と出力
- ✅ `x := 10; x > 5` → 変数を使った比較
- ✅ `5 > 3`, `10 == 10` → 比較演算子
- ✅ REPL での動作確認完了

### Package Refactoring ⚠️ 進行中！

- [x] パッケージ分割の設計（Go 本家スタイル採用）
- [x] token/ パッケージの作成（lexer/token.go から移動）
- [x] scanner/ パッケージの作成（lexer/lexer.go から移動・リネーム）
- [x] ast/ パッケージの作成（parser/ast.go から移動）
- [x] 全てのコメントを英語化
- [ ] parser/ パッケージの修正（進行中・エラー多数）
- [ ] eval/ パッケージの修正（未着手）
- [ ] repl/ パッケージの修正（未着手）
- [ ] main.go の修正（未着手）
- [ ] 全テストファイルの修正（未着手）

**現在の構造（目標）:**

```
petitgo/
├── main.go              # CLI entry point
├── token/              # Token definitions (完了)
│   └── token.go
├── scanner/            # Lexical analysis (完了)
│   ├── scanner.go
│   ├── lexer_test.go (要修正)
│   └── lexer_identifiers_test.go (要修正)
├── ast/                # AST node definitions (完了)
│   └── ast.go
├── parser/             # Syntax analysis (要修正)
│   ├── parser.go (エラー多数)
│   └── parser_test.go (要修正)
├── eval/               # Evaluation (要修正)
│   ├── environment.go
│   └── eval.go
└── repl/               # REPL (要修正)
    └── repl.go
```

**リファクタリング進行状況:**

1. ✅ token パッケージ: Token と TokenInfo の定義完了
2. ✅ scanner パッケージ: Scanner 構造体と NextToken()メソッド完了
3. ✅ ast パッケージ: 全 AST ノード定義完了
4. ⚠️ parser パッケージ: import 修正済みだが、型参照でエラー多数
5. ❌ eval パッケージ: 未着手
6. ❌ repl パッケージ: 未着手
7. ❌ main.go: 未着手
8. ❌ テストファイル群: package 宣言と import 修正が必要

**次回再開時のタスク:**

1. parser/parser.go 内の型参照を ast.前置詞付きに修正
2. scanner.Scanner, token.TokenInfo 等への変更
3. eval パッケージの import 修正
4. repl パッケージの import 修正
5. main.go の import 修正
6. 全テストファイルの修正
7. ビルド・テスト・動作確認

**注意事項:**

- 一括置換は危険（Statement→ast.Statement で関数名まで変わってしまった）
- 段階的に 1 つずつパッケージを修正すること
- 各段階でビルド確認すること
