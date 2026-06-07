package main

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
)

// Models 列出目前可用的模型 ID 與顯示名稱。
// 對應 SDK 文件: Models 章節與可用模型清單。
func Models() error {
	fmt.Println("[models]")
	client := getClient()
	page, err := client.Models.List(context.TODO(), anthropic.ModelListParams{
		Limit: anthropic.Int(20),
	})
	if err != nil {
		return err
	}
	for page != nil {
		for _, m := range page.Data {
			fmt.Printf("  %-32s %s\n", m.ID, m.DisplayName)
		}
		page, err = page.GetNextPage()
		if err != nil {
			return err
		}
	}
	return nil
}
