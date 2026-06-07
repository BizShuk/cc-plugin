package main

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
)

// Messages 示範各種訊息內容區塊的組裝方式。
// 對應 SDK 文件: Content blocks（文字、圖片 URL、圖片 base64）。
func Messages() error {
	if err := multiTextBlock(); err != nil {
		return fmt.Errorf("multiTextBlock: %w", err)
	}
	if err := imageURLBlock(); err != nil {
		return fmt.Errorf("imageURLBlock: %w", err)
	}
	if err := imageBase64Block(); err != nil {
		return fmt.Errorf("imageBase64Block: %w", err)
	}
	return nil
}

// multiTextBlock — 同一訊息內放入多個文字區塊
func multiTextBlock() error {
	fmt.Println("[messages.multiTextBlock]")
	client := getClient()
	message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		Model:     defaultModel,
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			{
				Role: anthropic.MessageParamRoleUser,
				Content: []anthropic.ContentBlockParamUnion{
					{OfText: &anthropic.TextBlockParam{Text: "What is the capital of France?"}},
					{OfText: &anthropic.TextBlockParam{Text: "What about Germany?"}},
				},
			},
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

// imageURLBlock — 透過 URL 引用遠端圖片
// 支援格式: JPEG / PNG / GIF / WEBP；單檔上限 10 MB
func imageURLBlock() error {
	fmt.Println("\n[messages.imageURLBlock]")
	client := getClient()
	message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		Model:     defaultModel,
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			{
				Role: anthropic.MessageParamRoleUser,
				Content: []anthropic.ContentBlockParamUnion{
					{OfText: &anthropic.TextBlockParam{Text: "What is in this image?"}},
					{OfImage: &anthropic.ImageBlockParam{
						Source: anthropic.ImageBlockParamSourceUnion{
							OfURL: &anthropic.URLImageSourceParam{
								URL: "https://upload.wikimedia.org/wikipedia/commons/thumb/0/0c/GoldenGateBridge-001.jpg/1200px-GoldenGateBridge-001.jpg",
							},
						},
					}},
				},
			},
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

// imageBase64Block — 直接以 base64 編碼內嵌小型圖片
// 範例使用 1x1 透明 PNG；實務上會用 os.ReadFile 讀入真實檔案
func imageBase64Block() error {
	fmt.Println("\n[messages.imageBase64Block]")
	client := getClient()
	// 1x1 透明 PNG
	const png1x1 = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg=="
	data, _ := base64.StdEncoding.DecodeString(png1x1)

	message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		Model:     defaultModel,
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			{
				Role: anthropic.MessageParamRoleUser,
				Content: []anthropic.ContentBlockParamUnion{
					{OfText: &anthropic.TextBlockParam{Text: "Describe this image."}},
					{OfImage: &anthropic.ImageBlockParam{
						Source: anthropic.ImageBlockParamSourceUnion{
							OfBase64: &anthropic.Base64ImageSourceParam{
								MediaType: anthropic.Base64ImageSourceMediaTypeImagePNG,
								Data:      string(data),
							},
						},
					}},
				},
			},
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
