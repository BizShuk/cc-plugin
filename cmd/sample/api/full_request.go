package main

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
)

// FullRequest 示範一次塞滿所有常用欄位的完整請求：含 system、tools、metadata、
// temperature、top_p，並由 stop_reason 判斷是否進入工具迴圈。
// 對應 SDK 文件: Request fields 章節（omitzero / opt 語意示範）。
func FullRequest() error {
	fmt.Println("[full_request]")
	client := getClient()

	schema := generateSchema[GetCoordinatesInput]()
	toolParam := anthropic.ToolParam{
		Name:        "get_coordinates",
		Description: anthropic.String("Accepts a place, returns lat/long."),
		InputSchema: schema,
	}

	message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		Model:       defaultModel,
		MaxTokens:   1024,
		Temperature: anthropic.Float(0.5),
		TopP:        anthropic.Float(0.9),
		Metadata: anthropic.MetadataParam{
			UserID: anthropic.String("user_42"),
		},
		System: []anthropic.TextBlockParam{
			{Text: "You are a concise geography assistant. Always use the get_coordinates tool when relevant."},
		},
		Tools: []anthropic.ToolUnionParam{
			{OfTool: &toolParam},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("Where is Tokyo?")),
		},
	})
	if err != nil {
		return err
	}

	fmt.Printf("  stop_reason: %s\n", message.StopReason)
	fmt.Printf("  usage:       input=%d output=%d\n",
		message.Usage.InputTokens, message.Usage.OutputTokens)

	for _, b := range message.Content {
		switch v := b.AsAny().(type) {
		case anthropic.TextBlock:
			fmt.Printf("  text: %s\n", v.Text)
		case anthropic.ToolUseBlock:
			fmt.Printf("  tool_use: %s id=%s\n", v.Name, v.ID)
		}
	}
	return nil
}
