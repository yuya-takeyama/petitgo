# petitgo Development Guide

## プロジェクト概要

petitgo は小さな Go 実装で、最終的にはセルフホスト（自分自身でコンパイルできる）を目指すプロジェクトです。

## ドキュメント構成

- **@docs/ROADMAP.md**: 全体の開発ロードマップ、9 つの Phase の詳細計画
- **@docs/MILESTONE-1.md**: Phase 1 の具体的な実装ステップと Success Criteria
- **@docs/TESTING-STRATEGY.md**: テスト戦略、Go ライブラリテストから最終的なセルフテストまで
- **@CLAUDE.md**: このファイル、開発ガイドと現在の進捗状況

## 開発フロー（Pull Request ベース）

### ブランチ戦略

1. 作業開始時は最新の main ブランチから作業ブランチを作成
   - ブランチ名: `yuya-takeyama/feat/FEAT_NAME`
2. 機能実装・テスト追加を行う
3. ローカルでの変更を push 後、Pull Request を作成
   - 必ず `--draft` オプションで draft PR として作成
4. レビュー・テスト通過後に main ブランチにマージ

### TDD サイクル（PR 内で実施）

1. テストファイルを先に作る（`*_test.go`）
2. `go test` で RED を確認
3. 実装ファイルを作って GREEN に
4. 作業単位ごとに細かくコミット
5. PR が完成したら draft を外してレビュー依頼

## Pull Request 方針

- **1 つの PR = 1 つの機能・修正**
- 必ず draft で開始し、完成後に ready for review にする
- gh pr create コマンドを使用して作成
- コミットメッセージは Conventional Commits 準拠
- 「ポケモンのレポート」のように、こまめにコミット・push

## コーディング規約

### コメント

- **全てのコメントは英語で記述する**
- コードの意図や理由を明確に説明する
- TODO や FIXME も英語で記述
- 例: `// Parse the if statement with condition and optional else block`

## 現在の開発状況

### 完了済み Phase（詳細は @docs/ROADMAP.md 参照）

- ✅ **Phase 1: Minimal Calculator** - 四則演算の基本実装
- ✅ **Phase 2: Variables and Basic Statements** - 変数宣言・代入・スコープ
- ✅ **Phase 3: Control Flow** - if/else 文・for 文・比較演算子
- ✅ **Phase 4: Functions** - 関数定義・呼び出し・引数・戻り値
- ✅ **Phase 5: Type System** - Value システム・型チェック・型推論
- ⏭️ **Phase 6: Standard Library** - スキップ（Go 標準ライブラリ活用）
- ✅ **Phase 7: Packages and Imports** - package/import 文の基本実装
- ✅ **Phase 8: Compiler** - 真のネイティブバイナリコンパイラ完成

### 現在の Phase 9: Self-hosting 🚧

**重要な成果:**

- ✅ **真のネイティブコンパイラ完成**: Go コンパイラ非依存のネイティブバイナリ生成
- ✅ **クロスプラットフォーム対応**: macOS ARM64 と Linux x86_64 サポート
- ✅ **再帰関数の完全動作**: fibonacci 完全計算可能
- ✅ **包括的テストカバレッジ**: scanner(100%), ast(100%), parser(98%), asmgen(94.2%)

**完成機能一覧:**

- ✅ 変数再代入 (x = y)
- ✅ コメント (// と /\* \*/)
- ✅ インクリメント・デクリメント (++, --)
- ✅ 完全な for ループ (init; condition; update)
- ✅ 複合代入演算子 (+=, -=, \*=, /=)
- ✅ ARM64 & x86_64 アセンブリ生成

**利用可能コマンド:**

```bash
petitgo build file.pg    # ネイティブバイナリ生成
petitgo run file.pg      # 直接実行
petitgo ast file.pg      # AST JSON出力
petitgo asm file.pg      # アセンブリ出力
petitgo                  # REPL モード
```

### 現在のファイル構造

```
petitgo/
├── main.go              # CLI entry point
├── token/              # Token definitions
├── scanner/            # Lexical analysis
├── ast/                # AST node definitions
├── parser/             # Syntax analysis
├── eval/               # Evaluation with type system
├── asmgen/             # Native assembly generation
├── repl/               # REPL
└── examples/           # Sample programs
```

## 次のやるべきこと

### 🎯 Phase 9 の残りタスク

**1. 未実装機能の完成**

- 🚧 **switch 文の実装完了**: parser は完了済み、asmgen/eval 対応が必要
- 🚧 **構造体フィールドアクセス**: `obj.field` 構文の実装
- 🚧 **スライス操作**: 基本的なスライス操作（length, append 等）
- 🚧 **文字列比較**: 文字列リテラルとの比較演算

**2. テストカバレッジの完全化**

- 🚧 **eval パッケージ**: カバレッジ測定・改善（現在未測定）
- 🚧 **repl パッケージ**: カバレッジ測定・改善（現在未測定）
- 🏆 **asmgen 残り 5.8%**: 主に NewAsmGenerator の GOOS/GOARCH 分岐

**3. セルフホスティング準備**

- 🎯 **petitgo ソースコードの petitgo 対応**: 必要な Go 機能の洗い出し
- 🎯 **複数ファイル対応**: package/import の複数ファイルサポート
- 🎯 **ブートストラップ処理**: petitgo 自身のコンパイル検証

### 💡 優先順位の提案

1. **switch 文の実装完了** - 既に parser 完了、残りの実装も比較的簡単
2. **eval/repl パッケージのテストカバレッジ改善** - 品質向上の継続
3. **構造体フィールドアクセスの実装** - セルフホスティングに必要
4. **petitgo ソースコードの分析** - セルフホスティングに向けた次のステップ

### 📚 参考情報

- **詳細な過去の作業ログ**: 必要に応じて git log や過去のコミットメッセージを参照
- **全体のロードマップ**: @docs/ROADMAP.md
- **アーキテクチャの詳細**: 各パッケージのソースコード

### 🎯 重要なメモ (TL;DR)

- petitgo は**既に真のネイティブコンパイラとして完成**している
- Go コンパイラ非依存で ARM64/x86_64 バイナリ生成可能
- fibonacci, for 文, if 文, 関数, 変数など主要機能は全て動作
- テストカバレッジも非常に高い（scanner/ast: 100%, parser: 98%, asmgen: 94%）
- **次の目標はセルフホスティング**（petitgo で petitgo をコンパイル）
