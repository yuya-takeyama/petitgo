package main

import (
	"os"

	"github.com/yuya-takeyama/petitgo/repl"
)

func main() {
	if len(os.Args) > 1 {
		// TODO: 将来的にはサブコマンドやファイル実行を実装
		// 今はREPLのみ
	}

	// REPLを起動
	repl.StartREPL()
}
