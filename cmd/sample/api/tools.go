package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/invopop/jsonschema"
)

// GetCoordinatesInput 是 get_coordinates 工具的輸入參數結構。
// 透過 jsonschema tag 與 jsonschema_description 自動產生 JSON Schema。
type GetCoordinatesInput struct {
	Location string `json:"location" jsonschema_description:"The location to look up."`
}

// getCoordinatesSchema 在 init 時一次性反射產生，執行階段重複使用。
var getCoordinatesSchema = generateSchema[GetCoordinatesInput]()

// GetCoordinates 回傳寫死的假座標；實務上會串接真實的地理 API。
func GetCoordinates(location string) Coordinate {
	fmt.Printf("  [mock] lookup %q -> SF\n", location)
	return Coordinate{Long: -122.4194, Lat: 37.7749}
}

type Coordinate struct {
	Long float64 `json:"long"`
	Lat  float64 `json:"lat"`
}

// Tools 示範完整的工具呼叫迴圈：模型請求呼叫 → 執行函式 → 回傳結果 → 讓模型產出最終回答。
// 對應 SDK 文件: Tool calling 章節。
func Tools() error {
	fmt.Println("[tools]")
	client := getClient()

	toolParams := []anthropic.ToolParam{
		{
			Name:        "get_coordinates",
			Description: anthropic.String("Accepts a place as an address, then returns the latitude and longitude coordinates."),
			InputSchema: getCoordinatesSchema,
		},
	}
	tools := make([]anthropic.ToolUnionParam, len(toolParams))
	for i, tp := range toolParams {
		tools[i] = anthropic.ToolUnionParam{OfTool: &tp}
	}

	messages := []anthropic.MessageParam{
		anthropic.NewUserMessage(anthropic.NewTextBlock("Where is San Francisco?")),
	}
	fmt.Println("  [user]: Where is San Francisco?")

	for {
		message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
			Model:     defaultModel,
			MaxTokens: 1024,
			Messages:  messages,
			Tools:     tools,
		})
		if err != nil {
			return err
		}

		fmt.Print("  [assistant]: ")
		toolResults := []anthropic.ContentBlockParamUnion{}
		for _, block := range message.Content {
			switch v := block.AsAny().(type) {
			case anthropic.TextBlock:
				fmt.Println(v.Text)
			case anthropic.ToolUseBlock:
				raw, _ := json.Marshal(v.Input)
				fmt.Printf("\n  [tool_use] %s: %s\n", v.Name, string(raw))

				var input GetCoordinatesInput
				if err := json.Unmarshal(v.Input, &input); err != nil {
					return fmt.Errorf("unmarshal tool input: %w", err)
				}
				resp := GetCoordinates(input.Location)
				respJSON, _ := json.Marshal(resp)
				fmt.Printf("  [tool_result]: %s\n", string(respJSON))

				toolResults = append(toolResults, anthropic.NewToolResultBlock(v.ID, string(respJSON), false))
			}
		}

		// 把助手訊息與工具回填都加入對話歷史
		messages = append(messages, message.ToParam())
		if len(toolResults) == 0 {
			break // 模型已給出最終文字回答
		}
		messages = append(messages, anthropic.NewUserMessage(toolResults...))
	}
	return nil
}

// generateSchema 透過 invopop/jsonschema 反射產生 ToolInputSchemaParam。
// DoNotReference 避免使用 $ref，讓子屬性直接展開，與官方範例一致。
func generateSchema[T any]() anthropic.ToolInputSchemaParam {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	return anthropic.ToolInputSchemaParam{
		Properties: schema.Properties,
		Required:   schema.Required,
	}
}
