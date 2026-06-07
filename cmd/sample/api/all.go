package main

import "fmt"

// RunAll 依固定順序執行所有範例。
// 任一範例失敗會印出錯誤訊息但繼續執行下一個，方便一次看完整輸出。
func RunAll() error {
	fmt.Println("=" + repeat("=", 59))
	fmt.Println("Anthropic Go SDK - MiniMax 相容端點範例集")
	fmt.Println("=" + repeat("=", 59))

	ordered := []string{
		"basic", "messages", "streaming", "tools", "thinking",
		"metadata", "response", "models", "error", "full_request",
	}
	for _, name := range ordered {
		fmt.Printf("\n[%s]\n", name)
		if err := examples[name](); err != nil {
			fmt.Printf("  ! %v\n", err)
		}
	}
	fmt.Println("\n" + repeat("=", 60))
	fmt.Println("全部範例執行完畢")
	return nil
}

func repeat(s string, n int) string {
	out := ""
	for i := 0; i < n; i++ {
		out += s
	}
	return out
}
