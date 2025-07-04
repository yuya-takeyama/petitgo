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

### Package Refactoring ✅ 完了！

- [x] パッケージ分割の設計（Go 本家スタイル採用）
- [x] token/ パッケージの作成（lexer/token.go から移動）
- [x] scanner/ パッケージの作成（lexer/lexer.go から移動・リネーム）
- [x] ast/ パッケージの作成（parser/ast.go から移動）
- [x] 全てのコメントを英語化
- [x] parser/ パッケージの修正（完了）
- [x] eval/ パッケージの修正（完了）
- [x] repl/ パッケージの修正（完了）
- [x] main.go の修正（完了）
- [x] 全テストファイルの修正（完了）
- [x] ファイル名の統一（scanner_*.go パターンに統一）

**最終構造:**

```
petitgo/
├── main.go              # CLI entry point
├── token/              # Token definitions
│   └── token.go
├── scanner/            # Lexical analysis
│   ├── scanner.go
│   ├── scanner_test.go
│   └── scanner_identifiers_test.go
├── ast/                # AST node definitions
│   └── ast.go
├── parser/             # Syntax analysis
│   ├── parser.go
│   └── parser_test.go
├── eval/               # Evaluation
│   ├── environment.go
│   └── eval.go
└── repl/               # REPL
    └── repl.go
```

**リファクタリング完了状況:**

1. ✅ token パッケージ: Token と TokenInfo の定義完了
2. ✅ scanner パッケージ: Scanner 構造体と NextToken()メソッド完了
3. ✅ ast パッケージ: 全 AST ノード定義完了
4. ✅ parser パッケージ: 全型参照を ast. プレフィックス付きに修正完了
5. ✅ eval パッケージ: import と型参照修正完了
6. ✅ repl パッケージ: import 修正完了
7. ✅ main.go: 修正不要（既に正しい import）
8. ✅ テストファイル群: package 宣言と import 修正完了
9. ✅ ファイル名統一: scanner_*.go パターンに統一完了

**達成事項:**

- Go 本家スタイルのパッケージ構造採用
- 全 import の正しい修正
- 全コメントの英語化
- ビルド・テスト完全成功（`go build`, `go test ./...` 全て通過）
- REPL 動作確認完了
- ファイル命名規則の統一

### Phase 5: Type System ✅ **完了！**

- [x] Value インターフェース設計（完了）
  - [x] IntValue, StringValue, BoolValue の実装
  - [x] StructValue, SliceValue の基盤実装
  - [x] Type(), String(), IsTruthy() メソッド
- [x] 型チェックシステム（完了）
  - [x] 変数宣言時の型検証（var x int = "string" → zero value）
  - [x] 関数引数の型チェック（引数不足・型不一致対応）
  - [x] 型安全な再代入（既存変数の型保護）
- [x] 型推論システム（完了）
  - [x] := 代入での自動型推定
  - [x] 式評価結果からの型決定
  - [x] 型不一致時のゼロ値フォールバック
- [x] 包括的テストスイート（完了）
  - [x] type_system_test.go - 基本型システムテスト
  - [x] type_checking_test.go - 型チェック機能テスト
  - [x] type_inference_test.go - 型推論機能テスト
  - [x] struct_test.go, slice_test.go - 構造体・スライステスト
- [x] バグ修正（完了）
  - [x] if文条件パース問題（struct literal誤判定）
  - [x] 関数再帰実行でのnilポインタ問題
  - [x] EvalValueWithEnvironment統合
  - [x] Phase 3 統合テストのValue型対応（2025-07-04修正）
  - [x] parseIfCondition修正（数値比較式の正しい解析）

**現在のファイル構造更新:**

```
petitgo/
├── main.go              # CLI entry point
├── token/              # Token definitions
│   └── token.go
├── scanner/            # Lexical analysis
│   ├── scanner.go
│   ├── scanner_test.go
│   └── scanner_identifiers_test.go
├── ast/                # AST node definitions
│   └── ast.go
├── parser/             # Syntax analysis
│   ├── parser.go
│   └── parser_test.go
├── eval/               # Evaluation with type system
│   ├── environment.go
│   ├── eval.go
│   ├── value.go        # Value interface & implementations
│   ├── slice_test.go
│   ├── struct_test.go
│   ├── type_checking_test.go
│   ├── type_inference_test.go
│   └── type_system_test.go
└── repl/               # REPL
    └── repl.go
```

### Phase 8: Compiler ✅ **完了！**

- [x] AST から Go ソースコード生成（完了）
  - [x] 完全な codegen パッケージ実装
  - [x] 全 AST ノード対応（statements, expressions）
  - [x] 演算子、関数呼び出し、制御構文サポート
- [x] build/run コマンド実装（完了）
  - [x] `petitgo build file.pg` - 実行可能バイナリ生成
  - [x] `petitgo run file.pg` - 直接実行
  - [x] Go コンパイラとの統合
- [x] package/import 文サポート（完了）
  - [x] package 宣言のパース
  - [x] import 文のパース
  - [x] 基本的なパッケージ管理
- [x] 組み込み println 関数採用（完了）
  - [x] Go 組み込み println 使用
  - [x] 自動改行付き出力
  - [x] カスタム print 関数削除
- [x] 包括的サンプルコード（完了）
  - [x] examples/ ディレクトリ
  - [x] 8つの実用的サンプル
  - [x] 詳細な README.md

**現在のファイル構造:**

```
petitgo/
├── main.go              # CLI entry point with build/run commands
├── token/              # Token definitions
│   └── token.go
├── scanner/            # Lexical analysis
│   ├── scanner.go
│   ├── scanner_test.go
│   └── scanner_identifiers_test.go
├── ast/                # AST node definitions
│   └── ast.go
├── parser/             # Syntax analysis
│   ├── parser.go
│   └── parser_test.go
├── eval/               # Evaluation with type system
│   ├── environment.go
│   ├── eval.go
│   ├── value.go        # Value interface & implementations
│   ├── slice_test.go
│   ├── struct_test.go
│   ├── type_checking_test.go
│   ├── type_inference_test.go
│   └── type_system_test.go
├── codegen/            # Go source code generation
│   └── codegen.go
├── repl/               # REPL
│   └── repl.go
├── examples/           # Sample programs
│   ├── README.md
│   ├── hello.pg
│   ├── calculator.pg
│   ├── variables.pg
│   ├── conditionals.pg
│   ├── functions.pg
│   ├── fibonacci.pg
│   └── simple_if.pg
└── check_newlines.py   # Development tool
```

**動作実績:**
- ✅ 基本的な petitgo プログラムのコンパイル・実行成功
- ✅ 四則演算、変数、関数、if文など主要機能対応
- ✅ Go 本家互換の println 出力
- ✅ ビルド時間も高速（Go コンパイラ活用）

**既知の課題:**
- ⚠️ 括弧による演算子優先順位制御（`(x+y)*2` が正しく動作しない）
- ⚠️ for 文の完全実装
- ⚠️ 変数再代入（`x = y`）未対応
- ⚠️ コメント（`//`）未対応

### Phase 9: Self-hosting 🚧 **次の目標**

- [ ] petitgo ソースコードの petitgo 対応
- [ ] 必要な Go 機能の追加実装
- [ ] ブートストラップ処理
- [ ] セルフコンパイル検証

## 開発メモリ

### 2025-07-04 作業内容
- Phase 3 統合テストがValue型システムで失敗していた問題を修正
- `parseIfCondition` が数値から始まる比較式（`1 > 0`）を正しく解析できていなかった
- `parseTerm()` から `ParseExpression()` に変更して解決
- 全テストが通ることを確認済み

### 次回の作業
- Phase 9: Self-hosting の実装開始
- petitgo ソースコードを petitgo でコンパイルできるようにする
- 必要な Go 機能の追加実装から着手