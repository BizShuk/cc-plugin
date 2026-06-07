package main

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
)

// Metadata 示範如何在請求中夾帶 metadata（user_id）以利後續追蹤與濫用偵測，
// 以及如何透過 CountTokens 在不發送請求的情況下預估 token 用量。
// 對應 SDK 文件: Metadata、Count tokens 章節。
func Metadata() error {
	if err := withUserMetadata(); err != nil {
		return fmt.Errorf("withUserMetadata: %w", err)
	}
	if err := countTokens(); err != nil {
		return fmt.Errorf("countTokens: %w", err)
	}
	return nil
}

// withUserMetadata — 傳遞不透明的使用者 ID（uuid / hash），不可含個資
func withUserMetadata() error {
	fmt.Println("[metadata.withUserMetadata]")
	client := getClient()
	message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		Model:     defaultModel,
		MaxTokens: 256,
		Metadata: anthropic.MetadataParam{
			UserID: anthropic.String("user_abc123"),
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("Hello!")),
		},
	})
	if err != nil {
		return err
	}
	fmt.Printf("  usage: input=%d output=%d\n", message.Usage.InputTokens, message.Usage.OutputTokens)
	return nil
}

// countTokens — 呼叫 CountTokens 預估 input tokens，不會實際生成回應
func countTokens() error {
	fmt.Println("\n[metadata.countTokens]")
	client := getClient()
	res, err := client.Messages.CountTokens(context.TODO(), anthropic.MessageCountTokensParams{
		Model: defaultModel,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("What is the meaning of life?")),
		},
	})
	if err != nil {
		return err
	}
	fmt.Printf("  input tokens: %d\n", res.InputTokens)
	return nil
}
