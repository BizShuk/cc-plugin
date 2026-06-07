package main

import (
	"fmt"
	"os"
	"sort"
)

// examples 註冊所有可執行的子命令。
// 每個鍵為子命令名稱，值為對應的執行函式。
var examples = map[string]func() error{
	"basic":        Basic,
	"messages":     Messages,
	"streaming":    Streaming,
	"tools":        Tools,
	"thinking":     Thinking,
	"metadata":     Metadata,
	"response":     Response,
	"models":       Models,
	"error":        ErrorExample,
	"full_request": FullRequest,
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	if cmd == "all" {
		if err := RunAll(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	fn, ok := examples[cmd]
	if !ok {
		fmt.Fprintf(os.Stderr, "未知子命令: %q\n\n", cmd)
		printUsage()
		os.Exit(1)
	}

	if err := fn(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "用法: go run ./cmd/sample/api <subcommand>")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "可用子命令:")
	names := make([]string, 0, len(examples)+1)
	names = append(names, "all")
	for k := range examples {
		names = append(names, k)
	}
	sort.Strings(names)
	for _, n := range names {
		fmt.Fprintf(os.Stderr, "  %s\n", n)
	}
}
