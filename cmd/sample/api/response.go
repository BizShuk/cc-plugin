package main

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
)

// Response 示範如何讀取回應物件上的常見屬性：id、model、role、stop_reason、usage、container 等。
// 對應 SDK 文件: Response objects 章節。
func Response() error {
	fmt.Println("[response]")
	client := getClient()
	message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		Model:     defaultModel,
		MaxTokens: 256,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("Say hello in three languages.")),
		},
	})
	if err != nil {
		return err
	}

	fmt.Printf("  id:          %s\n", message.ID)
	fmt.Printf("  type:        %s\n", message.Type)
	fmt.Printf("  role:        %s\n", message.Role)
	fmt.Printf("  model:       %s\n", message.Model)
	fmt.Printf("  stop_reason: %s\n", message.StopReason)
	fmt.Printf("  stop_seq:    %q\n", message.StopSequence)
	fmt.Printf("  usage:       input=%d output=%d cache_read=%d cache_write=%d\n",
		message.Usage.InputTokens,
		message.Usage.OutputTokens,
		message.Usage.CacheReadInputTokens,
		message.Usage.CacheCreationInputTokens,
	)

	// 取出第一個文字區塊
	for _, b := range message.Content {
		if t, ok := b.AsAny().(anthropic.TextBlock); ok {
			fmt.Printf("  text:        %s\n", t.Text)
			break
		}
	}
	return nil
}
