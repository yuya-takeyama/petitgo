# petitgo Development Guide

## プロジェクト概要

petitgo は小さな Go 実装で、最終的にはセルフホスト（自分自身でコンパイルできる）を目指すプロジェクトです。

## ドキュメント構成

- **@docs/ROADMAP.md**: 全体の開発ロードマップ、9つの Phase の詳細計画
- **@docs/MILESTONE-1.md**: Phase 1 の具体的な実装ステップとSuccess Criteria
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