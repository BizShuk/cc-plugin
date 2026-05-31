package cmd

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bizshuk/cc-plugin/model"
	"github.com/spf13/viper"
)

func TestOllamaExtract(t *testing.T) {
	mockResponse := map[string]interface{}{
		"message": map[string]string{
			"content": `{"candidates": [{"text": "likes tea", "entities": ["alice"], "kind": "preference", "first_person": false, "confirmed_by_human": true}]}`,
		},
	}
	responseBytes, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseBytes)
	}))
	defer server.Close()

	viper.Set("llm.host", server.URL)
	viper.Set("llm.model", "test-model")

	svc := NewOllamaService()
	obs := []model.Observation{{Source: "test", SourceID: "1", Text: "alice likes tea"}}
	cands, err := svc.Extract(context.Background(), obs)
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	if len(cands) != 1 {
		t.Fatalf("expected 1 candidate, got %d", len(cands))
	}
	if cands[0].Text != "likes tea" || cands[0].Kind != "preference" {
		t.Errorf("unexpected candidate values: %+v", cands[0])
	}
	if len(cands[0].SourceRefs) != 1 || cands[0].SourceRefs[0][0] != "test" || cands[0].SourceRefs[0][1] != "1" {
		t.Errorf("unexpected source refs: %v", cands[0].SourceRefs)
	}
}
