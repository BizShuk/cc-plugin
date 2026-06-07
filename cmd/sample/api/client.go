package main

import (
	"fmt"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// MiniMax-M3 是 MiniMax Anthropic 相容端點使用的模型名稱。
// 官方端點可改用 anthropic.ModelClaudeOpus4_8 等常數。
const defaultModel = "MiniMax-M3"

// getClient 依環境變數建立 Anthropic SDK Client。
//   - ANTHROPIC_API_KEY 必要，未設定時直接退出並提示
//   - ANTHROPIC_BASE_URL 選填，設定時切換到自訂端點（MiniMax 相容模式）
func getClient() anthropic.Client {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "請先設定環境變數 ANTHROPIC_API_KEY")
		os.Exit(1)
	}

	opts := []option.RequestOption{option.WithAPIKey(apiKey)}
	if base := os.Getenv("ANTHROPIC_BASE_URL"); base != "" {
		opts = append(opts, option.WithBaseURL(base))
		fmt.Fprintf(os.Stderr, "[client] using base URL: %s\n", base)
	}
	return anthropic.NewClient(opts...)
}
