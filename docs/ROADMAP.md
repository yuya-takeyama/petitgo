# petitgo Development Roadmap

## Overview

petitgo は小さな Go 実装で、最終的にはセルフホスト（自分自身でコンパイルできる）を目指すプロジェクトです。
最初は最小限の機能から始めて、段階的に機能を追加していきます。

## Development Phases

### Phase 1: Minimal Calculator ✅ **完了**
- **目標**: 四則演算ができる最小限のインタープリタ
- **実装内容**:
  - [x] プロジェクトのセットアップ
    - [x] Go 1.24.4 環境
    - [x] Git リポジトリ初期化
    - [x] 開発ツール設定（Prettier、aqua、gofmt hook）
  - [x] 字句解析器（Lexer）- 数値と演算子のトークン化
    - [x] Go 本家互換の Token 型定義
    - [x] 数値（INT）のトークン化
    - [x] 演算子（ADD, SUB, MUL, QUO）のトークン化
    - [x] 括弧（LPAREN, RPAREN）のトークン化
    - [x] 空白文字のスキップ処理
    - [x] 包括的なテストスイート
  - [x] パーサー - 四則演算の式をパース
    - [x] AST ノードの定義
    - [x] 再帰下降パーサーの実装
    - [x] 演算子優先順位の処理
  - [x] 評価器（Evaluator）- 式を評価して結果を返す
  - [x] print 関数 - 結果を出力（fmt パッケージなしで実装）
  - [x] REPL - 対話的に式を評価

**サポートする演算**: `+`, `-`, `*`, `/`, `()` （括弧による優先順位）

### Phase 2: Variables and Basic Statements ✅ **完了**
- **目標**: 変数の代入と参照、基本的な文の実装
- **実装内容**:
  - [x] 変数宣言（`var`）
  - [x] 代入文（`:=`）
  - [x] スコープの実装
  - [x] 複数文の実行

### Phase 3: Control Flow ✅ **完了**
- **目標**: 制御構文の実装
- **実装内容**:
  - [x] if/else 文
  - [x] for ループ
  - [x] break/continue

### Phase 4: Functions ✅ **完了**
- **目標**: 関数の定義と呼び出し
- **実装内容**:
  - [x] 関数定義（`func`）
  - [x] 関数呼び出し
  - [x] 引数と戻り値
  - [ ] クロージャ（未実装）

### Phase 5: Type System ✅ **完了**
- **目標**: Go の型システムの基本実装
- **実装内容**:
  - [x] 基本型（int, string, bool）
  - [x] Value システム（型情報付き評価）
  - [x] 型チェック（変数宣言、関数引数）
  - [x] 型推論（:= 代入文）
  - [x] 型安全性（型不一致時のゼロ値フォールバック）
  - [ ] 構造体（struct）- 部分実装済み
  - [ ] スライス - 部分実装済み
  - [x] 包括的テストスイート（type_system_test.go, type_checking_test.go, type_inference_test.go）

**完了した機能**:
- IntValue, StringValue, BoolValue, StructValue, SliceValue の Value インターフェース
- var 文での型宣言チェック（var x int = "string" → zero value 使用）
- 関数引数の型チェック（引数不足や型不一致も対応）
- := での型推論（x := 42 → int 型として推論）
- 既存変数への型安全な再代入（型不一致時はゼロ値使用）

### Phase 6: Standard Library (Minimal) ⏭️ **スキップ**
- **理由**: Go の標準ライブラリを活用する方針に変更

### Phase 7: Packages and Imports ✅ **基本実装完了**
- **目標**: パッケージシステムの実装
- **実装内容**:
  - [x] package 宣言のパース
  - [x] import 文のパース
  - [ ] 複数ファイルのサポート（将来実装）

### Phase 8: Compiler ✅ **完了**
- **目標**: インタープリタからコンパイラへ
- **実装内容**:
  - [x] AST から Go ソースコードへの変換
  - [x] Go コンパイラとの統合
  - [x] build/run コマンド実装
  - [x] 組み込み println 関数採用
  - [x] 包括的サンプルプログラム作成

**達成事項**:
- 完全な codegen パッケージによる Go ソース生成
- `petitgo build file.pg` でバイナリ生成
- `petitgo run file.pg` で直接実行
- examples/ ディレクトリに8つの実用サンプル
- Go 組み込み println による見やすい出力

### Phase 9: Self-hosting 🚧 **次の目標**
- **目標**: petitgo で petitgo をコンパイル
- **実装内容**:
  - [ ] petitgo ソースコードの petitgo 対応
    - [ ] package/import の複数ファイル対応
    - [ ] struct/interface の完全実装
    - [ ] map, slice の完全サポート
  - [ ] 必要な Go 機能の追加実装
    - [ ] 変数再代入（x = y）
    - [ ] コメント（// と /* */）
    - [ ] より複雑な制御構文
  - [ ] ブートストラップ処理
    - [ ] petitgo 自身のコンパイル検証
    - [ ] セルフコンパイル可能性の確認

## Design Principles

1. **シンプルさ優先**: 最初は正確さよりシンプルさを重視
2. **段階的な成長**: 各フェーズで動くものを作る
3. **テスト駆動**: 各機能にテストを書く
4. **Go らしさ**: Go の設計哲学に従う

## Directory Structure (Phase 1)

```
petitgo/
├── main.go          # エントリーポイント
├── lexer.go         # 字句解析器
├── parser.go        # パーサー
├── eval.go          # 評価器
├── repl.go          # REPL
└── docs/
    └── ROADMAP.md   # このファイル
```

パッケージ分割は Phase 2 以降で検討します。