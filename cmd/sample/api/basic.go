package main

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
)

// Basic 示範最常見的基礎請求：最小請求、會話延續、系統提示詞、溫度、Top-P。
// 對應 SDK 文件: Usage、Conversations、System prompts 章節。
func Basic() error {
	if err := minimal(); err != nil {
		return fmt.Errorf("minimal: %w", err)
	}
	if err := conversation(); err != nil {
		return fmt.Errorf("conversation: %w", err)
	}
	if err := withSystem(); err != nil {
		return fmt.Errorf("withSystem: %w", err)
	}
	if err := withTemperature(); err != nil {
		return fmt.Errorf("withTemperature: %w", err)
	}
	if err := withTopP(); err != nil {
		return fmt.Errorf("withTopP: %w", err)
	}
	return nil
}

// minimal — 最小請求：只需要 model、messages、max_tokens
func minimal() error {
	fmt.Println("[basic.minimal]")
	client := getClient()
	message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		Model:     defaultModel,
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("What is a quaternion?")),
		},
	})
	if err != nil {
		return err
	}
	for _, block := range message.Content {
		if t, ok := block.AsAny().(anthropic.TextBlock); ok {
			fmt.Println("  ", t.Text)
		}
	}
	return nil
}

// conversation — 多輪會話：用 message.ToParam() 將回應回填到訊息列表
func conversation() error {
	fmt.Println("\n[basic.conversation]")
	client := getClient()
	messages := []anthropic.MessageParam{
		anthropic.NewUserMessage(anthropic.NewTextBlock("What is my first name?")),
	}
	message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		Model:     defaultModel,
		Messages:  messages,
		MaxTokens: 1024,
	})
	if err != nil {
		return err
	}
	for _, b := range message.Content {
		if t, ok := b.AsAny().(anthropic.TextBlock); ok {
			fmt.Println("  turn 1:", t.Text)
		}
	}

	messages = append(messages, message.ToParam())
	messages = append(messages, anthropic.NewUserMessage(
		anthropic.NewTextBlock("My full name is John Doe"),
	))
	message, err = client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		Model:     defaultModel,
		Messages:  messages,
		MaxTokens: 1024,
	})
	if err != nil {
		return err
	}
	for _, b := range message.Content {
		if t, ok := b.AsAny().(anthropic.TextBlock); ok {
			fmt.Println("  turn 2:", t.Text)
		}
	}
	return nil
}

// withSystem — 系統提示詞：透過 System 欄位傳遞角色設定
func withSystem() error {
	fmt.Println("\n[basic.withSystem]")
	client := getClient()
	message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		Model:     defaultModel,
		MaxTokens: 1024,
		System: []anthropic.TextBlockParam{
			{Text: "Be very serious at all times."},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("Hello!")),
		},
	})
	if err != nil {
		return err
	}
	for _, b := range message.Content {
		if t, ok := b.AsAny().(anthropic.TextBlock); ok {
			fmt.Println("  ", t.Text)
		}
	}
	return nil
}

// withTemperature — 溫度參數 (0.0 - 1.0)，控制取樣隨機性
func withTemperature() error {
	fmt.Println("\n[basic.withTemperature]")
	client := getClient()
	temp := 0.7
	message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		Model:       defaultModel,
		MaxTokens:   1024,
		Temperature: anthropic.Float(temp),
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("Write a short poem about the sea.")),
		},
	})
	if err != nil {
		return err
	}
	for _, b := range message.Content {
		if t, ok := b.AsAny().(anthropic.TextBlock); ok {
			fmt.Println("  ", t.Text)
		}
	}
	return nil
}

// withTopP — Top-P 核取樣 (0.0 - 1.0)
func withTopP() error {
	fmt.Println("\n[basic.withTopP]")
	client := getClient()
	topP := 0.95
	message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		Model:     defaultModel,
		MaxTokens: 1024,
		TopP:      anthropic.Float(topP),
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("Explain quantum physics in simple terms.")),
		},
	})
	if err != nil {
		return err
	}
	for _, b := range message.Content {
		if t, ok := b.AsAny().(anthropic.TextBlock); ok {
			fmt.Println("  ", t.Text)
		}
	}
	return nil
}
