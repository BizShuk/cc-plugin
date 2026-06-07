package main

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
)

// Thinking 示範 extended thinking 的三種模式：自適應、停用、預算上限。
// 對應 SDK 文件: Extended thinking 章節。
func Thinking() error {
	if err := adaptive(); err != nil {
		return fmt.Errorf("adaptive: %w", err)
	}
	if err := disabled(); err != nil {
		return fmt.Errorf("disabled: %w", err)
	}
	if err := budgeted(); err != nil {
		return fmt.Errorf("budgeted: %w", err)
	}
	return nil
}

// adaptive — 啟用 adaptive thinking，由模型自行決定思考預算（不指定 budget_tokens）
func adaptive() error {
	fmt.Println("[thinking.adaptive]")
	client := getClient()
	message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		Model:     defaultModel,
		MaxTokens: 4096,
		Thinking: anthropic.ThinkingConfigParamUnion{
			OfAdaptive: &anthropic.ThinkingConfigAdaptiveParam{},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("What is the meaning of life?")),
		},
	})
	if err != nil {
		return err
	}
	for _, b := range message.Content {
		switch v := b.AsAny().(type) {
		case anthropic.ThinkingBlock:
			fmt.Printf("  [thinking] %d chars\n", len(v.Thinking))
		case anthropic.TextBlock:
			fmt.Printf("  [answer] %s\n", v.Text)
		}
	}
	return nil
}

// disabled — 明確關閉 extended thinking
func disabled() error {
	fmt.Println("\n[thinking.disabled]")
	client := getClient()
	disabled := anthropic.NewThinkingConfigDisabledParam()
	message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		Model:     defaultModel,
		MaxTokens: 1024,
		Thinking: anthropic.ThinkingConfigParamUnion{
			OfDisabled: &disabled,
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("What is 2+2?")),
		},
	})
	if err != nil {
		return err
	}
	for _, b := range message.Content {
		if t, ok := b.AsAny().(anthropic.TextBlock); ok {
			fmt.Printf("  %s\n", t.Text)
		}
	}
	return nil
}

// budgeted — 給定 thinking budget（token 數，必須 ≥1024 且 < max_tokens），模型會在預算內思考
func budgeted() error {
	fmt.Println("\n[thinking.budgeted]")
	client := getClient()
	message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		Model:     defaultModel,
		MaxTokens: 4096,
		Thinking: anthropic.ThinkingConfigParamUnion{
			OfEnabled: &anthropic.ThinkingConfigEnabledParam{
				BudgetTokens: 1024,
			},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("Solve: 17 * 23. Show your reasoning.")),
		},
	})
	if err != nil {
		return err
	}
	for _, b := range message.Content {
		switch v := b.AsAny().(type) {
		case anthropic.ThinkingBlock:
			fmt.Printf("  [thinking] %d chars\n", len(v.Thinking))
		case anthropic.TextBlock:
			fmt.Printf("  [answer] %s\n", v.Text)
		}
	}
	return nil
}
