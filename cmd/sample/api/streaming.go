package main

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
)

// Streaming 示範使用 NewStreaming 取得事件流並以 Accumulate 組合最終訊息。
// 對應 SDK 文件: Streaming 章節。
func Streaming() error {
	fmt.Println("[streaming]")
	client := getClient()

	content := "What is a quaternion?"
	stream := client.Messages.NewStreaming(context.TODO(), anthropic.MessageNewParams{
		Model:     defaultModel,
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(content)),
		},
	})

	// message 會在迴圈中透過 Accumulate 累積成完整的最終訊息
	message := anthropic.Message{}
	fmt.Print("  ")
	for stream.Next() {
		event := stream.Current()
		if err := message.Accumulate(event); err != nil {
			return fmt.Errorf("accumulate: %w", err)
		}

		// 取出文字增量即時列印
		if deltaEvent, ok := event.AsAny().(anthropic.ContentBlockDeltaEvent); ok {
			if delta, ok := deltaEvent.Delta.AsAny().(anthropic.TextDelta); ok {
				fmt.Print(delta.Text)
			}
		}
	}
	fmt.Println()

	if stream.Err() != nil {
		return fmt.Errorf("stream err: %w", stream.Err())
	}
	fmt.Printf("  usage: input=%d output=%d\n", message.Usage.InputTokens, message.Usage.OutputTokens)
	return nil
}
