package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
)

// Result は集約結果を格納するための構造体です
type Result struct {
	Path  string
	Count int
}

func main() {
	// 起動時引数のチェック
	if len(os.Args) != 3 {
		fmt.Println("使用方法: program <階層番号> <列挙数>")
		os.Exit(1)
	}

	// 階層番号と列挙数の取得
	level, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("エラー: 階層番号は数値で指定してください")
		os.Exit(1)
	}
	if level < 1 {
		fmt.Println("エラー: 階層番号は1以上の数値を指定してください")
		os.Exit(1)
	}

	enumCount, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("エラー: 列挙数は数値で指定してください")
		os.Exit(1)
	}

	// クリップボードからファイルパスの一覧を取得
	pathsString, err := clipboard.ReadAll()
	if err != nil {
		fmt.Println("エラー: クリップボードからデータを読み取れませんでした")
		os.Exit(1)
	}

	// ファイルパスを改行で分割
	// 修正ポイント: クリップボードの内容が`\r\n`で区切られる場合があるため、strings.Fields()を使用
	// または strings.Split(pathsString, "\n") のままでも動作することが多いが、より堅牢な実装として`strings.Fields`を推奨
	paths := strings.Split(pathsString, "\n")
	if len(paths) == 0 {
		return
	}

	// パスごとの集約
	aggregatedPaths := make(map[string]int)
	for _, path := range paths {
		// 先頭と末尾の空白をトリム
		path = strings.TrimSpace(path)
		if path == "" {
			continue
		}

		// ファイルパスを"\"で分割
		parts := strings.Split(path, `/`)
		if len(parts) < level {
			continue
		}

		// 指定された階層まで結合
		aggregatedPath := strings.Join(parts[:level], `/`)
		aggregatedPaths[aggregatedPath]++
	}

	// 集約結果をスライスに変換
	var results []Result
	for path, count := range aggregatedPaths {
		results = append(results, Result{Path: path, Count: count})
	}

	// 個数が大きい順にソート
	sort.Slice(results, func(i, j int) bool {
		return results[i].Count > results[j].Count
	})

	// 指定された列挙数分だけコンソールに出力
	for i, result := range results {
		if i >= enumCount {
			break
		}
		fmt.Printf("%s,%d\n", result.Path, result.Count)
	}
}
