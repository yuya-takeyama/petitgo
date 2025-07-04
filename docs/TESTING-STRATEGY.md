# petitgo Testing Strategy

## 概要

petitgo のテストは段階的にアプローチを変えていきます：
1. **初期**: 既存の Go でユニットテスト
2. **中期**: 統合テストツールで petitgo の実行結果を検証
3. **後期**: petitgo 自身でテストを記述・実行

## petitgo の出力検証方法

### 1. Go ライブラリとして実装（Phase 1-4）

**実現性**: ⭐⭐⭐⭐⭐（最も現実的）

**実装方法**:
```go
// eval.go
package main

func Eval(input string) (int, error) {
    lexer := NewLexer(input)
    parser := NewParser(lexer)
    ast := parser.Parse()
    return EvalAST(ast)
}
```

**テスト方法**:
```go
// eval_test.go
func TestEval(t *testing.T) {
    tests := []struct {
        input    string
        expected int
    }{
        {"1 + 2", 3},
        {"2 * 3 + 4", 10},
        {"(2 + 3) * 4", 20},
    }
    
    for _, tt := range tests {
        result, err := Eval(tt.input)
        if err != nil {
            t.Fatalf("eval error: %v", err)
        }
        if result != tt.expected {
            t.Errorf("input: %s, got: %d, want: %d", 
                tt.input, result, tt.expected)
        }
    }
}
```

**メリット**:
- Go の標準的なテスト手法が使える
- デバッグが簡単
- カバレッジ測定も容易
- CI/CD との統合も簡単

### 2. petitgo run コマンドの実装（Phase 5-6）

**実現性**: ⭐⭐⭐⭐（実装は簡単だが、テストがやや複雑）

**実装方法**:
```go
// main.go
func main() {
    if len(os.Args) > 1 && os.Args[1] == "run" {
        // ファイルを読み込んで実行
        content, _ := os.ReadFile(os.Args[2])
        result, _ := Eval(string(content))
        fmt.Println(result)
    }
}
```

**テスト方法**:
```go
// integration_test.go
func TestPetitgoRun(t *testing.T) {
    // テスト用の .pg ファイルを作成
    testFile := "test.pg"
    os.WriteFile(testFile, []byte("1 + 2"), 0644)
    defer os.Remove(testFile)
    
    // petitgo run を実行
    cmd := exec.Command("go", "run", ".", "run", testFile)
    output, err := cmd.Output()
    
    if string(output) != "3\n" {
        t.Errorf("unexpected output: %s", output)
    }
}
```

**別案: テストヘルパー関数**:
```go
func runPetitgo(input string) (string, error) {
    r, w, _ := os.Pipe()
    
    // 標準出力を一時的にパイプにリダイレクト
    oldStdout := os.Stdout
    os.Stdout = w
    
    // petitgo を実行
    Eval(input)
    
    // 標準出力を元に戻す
    w.Close()
    os.Stdout = oldStdout
    
    // 出力を読み取る
    output, _ := io.ReadAll(r)
    return string(output), nil
}
```

### 3. petitgo build コマンドの実装（Phase 8-9）

**実現性**: ⭐⭐（最も難しい）

**実装の前提条件**:
- コンパイラの実装が必要
- 実行可能バイナリの生成が必要
- ターゲットアーキテクチャの考慮

**テスト方法**:
```go
func TestPetitgoBuild(t *testing.T) {
    // ソースファイルを作成
    os.WriteFile("test.pg", []byte("print(1 + 2)"), 0644)
    
    // petitgo build を実行
    cmd := exec.Command("./petitgo", "build", "test.pg", "-o", "test.out")
    cmd.Run()
    
    // 生成されたバイナリを実行
    output, _ := exec.Command("./test.out").Output()
    
    if string(output) != "3\n" {
        t.Errorf("unexpected output: %s", output)
    }
}
```

**課題**:
- バイナリ生成は複雑（LLVM IR、アセンブリ、または Go のソース生成）
- クロスプラットフォーム対応
- 実行環境の違いの考慮

## テスト手法の提案

### 1. Go の標準テストフレームワーク（初期〜中期）

**メリット**: 
- すぐに始められる
- Go エコシステムとの統合が簡単
- カバレッジ計測も可能

**実装例**:
```go
// lexer_test.go
func TestLexer(t *testing.T) {
    input := "1 + 2"
    lexer := NewLexer(input)
    
    tests := []struct{
        expectedType TokenType
        expectedLiteral string
    }{
        {NUMBER, "1"},
        {PLUS, "+"},
        {NUMBER, "2"},
        {EOF, ""},
    }
    
    for i, tt := range tests {
        tok := lexer.NextToken()
        // アサーション
    }
}
```

### 2. スナップショットテスト（実行結果の検証）

**アプローチ A: シェルスクリプト + diff**
```bash
#!/bin/bash
# test_runner.sh

echo "Testing: 1 + 2"
echo "1 + 2" | ./petitgo > output.txt
echo "3" > expected.txt
diff output.txt expected.txt || exit 1
```

**アプローチ B: Go のテストヘルパー**
```go
// e2e_test.go
func TestPetitgoExecution(t *testing.T) {
    tests := []struct{
        input    string
        expected string
    }{
        {"1 + 2", "3"},
        {"2 * 3 + 4", "10"},
        {"(2 + 3) * 4", "20"},
    }
    
    for _, tt := range tests {
        output := runPetitgo(tt.input)
        if output != tt.expected {
            t.Errorf("input: %s, got: %s, want: %s", 
                tt.input, output, tt.expected)
        }
    }
}
```

### 3. YAML/JSON ベースのテストケース

**test_cases.yaml**:
```yaml
tests:
  - name: "basic addition"
    input: "1 + 2"
    expected: "3"
  - name: "operator precedence"
    input: "2 + 3 * 4"
    expected: "14"
  - name: "parentheses"
    input: "(2 + 3) * 4"
    expected: "20"
```

**実行ツール（Go/Python/Ruby など）**:
```go
// test_runner.go
type TestCase struct {
    Name     string `yaml:"name"`
    Input    string `yaml:"input"`
    Expected string `yaml:"expected"`
}

func runTests() {
    // YAML を読み込んで各テストケースを実行
}
```

### 4. 既存の言語処理系テストフレームワークの活用

**Test262 風アプローチ**:
- JavaScript の Test262 のようなテストスイート形式
- 各機能ごとにテストファイルを配置
- メタデータでテストの意図を記述

**例**:
```
// test/arithmetic/addition.pg.test
/*---
description: Basic addition operation
expected: 3
---*/
1 + 2
```

### 5. Property-based Testing（中期〜）

**QuickCheck スタイル**:
- ランダムな入力を生成してプロパティを検証
- 例: `(a + b) + c == a + (b + c)` （結合法則）

### 6. Differential Testing（後期）

- petitgo と Go の実行結果を比較
- 同じコードで同じ結果が出ることを確認

## 推奨アプローチと実現性の評価

### Phase 1-4: Go ライブラリとしてテスト（推奨）

**なぜこれが最適か**:
1. すぐに始められる
2. 既存の Go テストツールが全て使える
3. デバッグが容易
4. リファクタリングも安全

**実装順序**:
1. Lexer を実装 → `lexer_test.go` でテスト
2. Parser を実装 → `parser_test.go` でテスト
3. Evaluator を実装 → `eval_test.go` でテスト
4. 統合テスト → `integration_test.go` で全体をテスト

### Phase 5-6: CLI 実装後の選択肢

**Option A: exec.Command でサブプロセス実行**
- メリット: 実際の使用に近い
- デメリット: やや遅い、デバッグしづらい

**Option B: 標準出力のキャプチャ**
- メリット: 高速、プロセス内で完結
- デメリット: やや複雑な実装

### Phase 8-9: コンパイラ実装後

**現実的な選択肢**:
1. Go ソースコード生成 → `go build` でコンパイル
2. WebAssembly 出力 → wasmtime で実行
3. 独自 VM のバイトコード → VM で実行

## 結論

**Phase 1 では Go ライブラリとしてのテストに集中すべき**

理由:
- 最も簡単で確実
- 開発速度が最速
- バグの早期発見が可能
- TDD が実践しやすい

```go
// 最初はこれで十分！
func TestEval(t *testing.T) {
    result, _ := Eval("1 + 2")
    if result != 3 {
        t.Errorf("got %d, want 3", result)
    }
}
```

## 将来の展望

### Self-testing への道筋

1. **Phase 5-6**: petitgo で簡単なアサーション関数を実装
   ```go
   func assert(condition bool, message string) {
       if !condition {
           print("FAIL: " + message)
           exit(1)
       }
   }
   ```

2. **Phase 7-8**: テストランナーを petitgo で実装
   - ファイルグロブ
   - テスト関数の自動検出
   - 結果レポート

3. **Phase 9**: 完全なセルフテスト
   - petitgo のテストを petitgo で実行
   - ブートストラップの検証

## Next Steps

1. Go の標準テストで lexer のテストを書く
2. YAML ベースのテストケースファイルを作成
3. 統合テストランナーを実装（Go or シェルスクリプト）