package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// ErrorExample 示範 SDK 的錯誤處理模式：errors.As 取出 *anthropic.Error 以讀取
// Request ID、序列化請求/回應；以及用 option.WithMaxRetries 控制重試次數。
// 對應 SDK 文件: Error handling、Retries 章節。
func ErrorExample() error {
	if err := inspectAPIError(); err != nil {
		return fmt.Errorf("inspectAPIError: %w", err)
	}
	if err := configureRetries(); err != nil {
		return fmt.Errorf("configureRetries: %w", err)
	}
	return nil
}

// inspectAPIError — 故意使用過大的 max_tokens 觸發 400 錯誤，展示錯誤物件的取用
func inspectAPIError() error {
	fmt.Println("[error.inspectAPIError]")
	client := getClient()
	_, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		Model:     defaultModel,
		MaxTokens: 999_999, // 超出模型上限 -> 400
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("What is a quaternion?")),
		},
	})
	if err == nil {
		return fmt.Errorf("expected an error, got nil")
	}
	var apiErr *anthropic.Error
	if errors.As(err, &apiErr) {
		fmt.Printf("  status:    %d\n", apiErr.StatusCode)
		fmt.Printf("  requestID: %s\n", apiErr.RequestID)
	} else {
		fmt.Printf("  non-API error: %v\n", err)
	}
	return nil
}

// configureRetries — 展示 WithMaxRetries 與 WithRequestTimeout 的設定方式
func configureRetries() error {
	fmt.Println("\n[error.configureRetries]")
	client := anthropic.NewClient(
		option.WithAPIKey("dummy"),
		option.WithMaxRetries(5),
	)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := client.Messages.New(
		ctx,
		anthropic.MessageNewParams{
			Model:     defaultModel,
			MaxTokens: 64,
			Messages: []anthropic.MessageParam{
				anthropic.NewUserMessage(anthropic.NewTextBlock("hi")),
			},
		},
		option.WithRequestTimeout(10*time.Second),
		option.WithMaxRetries(0),
	)
	if err != nil {
		// 此處使用 dummy key 預期會失敗，目的在展示 option 設定本身能通過型別檢查
		fmt.Printf("  expected failure with dummy key: %v\n", err)
	} else {
		fmt.Println("  unexpected success")
	}
	return nil
}
