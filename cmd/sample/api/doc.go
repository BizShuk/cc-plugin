// Package main 提供 Anthropic Go SDK 與 MiniMax 相容端點的完整使用範例。
//
// 本套件採用 package main 直接編譯為可執行檔，透過子命令選擇要執行的範例。
//
// 環境設定:
//
//	export ANTHROPIC_API_KEY=<your-key>
//	# MiniMax 相容端點（可選；未設定時走 Anthropic 官方端點）
//	export ANTHROPIC_BASE_URL=https://api.minimax.io/anthropic
//
// 執行範例:
//
//	# 列出所有子命令
//	go run ./cmd/sample/api
//
//	# 執行單一範例
//	go run ./cmd/sample/api basic
//	go run ./cmd/sample/api streaming
//	go run ./cmd/sample/api tools
//
//	# 一次跑完全部
//	go run ./cmd/sample/api all
//
// 對應官方文件: https://platform.claude.com/docs/en/api/sdks/go
package main
