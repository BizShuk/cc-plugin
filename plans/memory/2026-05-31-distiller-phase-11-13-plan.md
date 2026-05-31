# Go Distiller Phase 11-13 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 實作 `distiller` 的 Ollama 記憶抽取服務與 `extract` 指令，將設定載入重構為使用 `gosdk` 的 `settings.json` 並落腳於 `~/.config/cc-plugin/settings.json`，整合並實作 `distiller run` 完整蒸餾管線，最後產出手動設定的 `crontab.txt`。

**Architecture:**

1. 建立 `OllamaService` 結構體呼叫本地 Ollama 服務 API，並提供 `distiller extract` 指令支援直接的命令列語義提取。
2. 整合 `gosdk/config` 在 Root CLI 初始化中載入預設 `settings.json`（放置於 `~/.config/cc-plugin/settings.json`）。
3. 實作 `distiller run` 指令作為編排器，整合：讀取來源 -> LLM 提取 -> 寫入長期與真實記憶 -> 更新游標，並在結尾調用 Retention 進行清理。
4. 建立 `crontab.txt` 供手動設定排程。

**Tech Stack:** Go 1.25, `github.com/bizshuk/gosdk`, `spf13/cobra`, `spf13/viper`, SQLite3

---

## User Review Required

> [!IMPORTANT]
>
> - 設定檔會遷移至 `~/.config/cc-plugin/settings.json`（由 `gosdk/config` 自動接管與合併）。
> - 整合測試將會在本地模擬 Ollama HTTP 伺服器，因此不依賴本地 Ollama 服務的實際運作。
> - 我們只會建立 `crontab.txt` 檔案，需要您在發佈後手動設定 `crontab`。

---

## Proposed Changes

### [Component Name] Go Distiller Codebase

#### [MODIFY] [root.go](file:///Users/shuk/projects/cc-plugin/cmd/root.go)

#### [NEW] [ollama.go](file:///Users/shuk/projects/cc-plugin/cmd/ollama.go)

#### [NEW] [run.go](file:///Users/shuk/projects/cc-plugin/cmd/run.go)

#### [NEW] [crontab.txt](file:///Users/shuk/projects/cc-plugin/crontab.txt)

#### [MODIFY] [state.go](file:///Users/shuk/projects/cc-plugin/cmd/state.go)

#### [MODIFY] [read_gbrain.go](file:///Users/shuk/projects/cc-plugin/cmd/read_gbrain.go)

#### [MODIFY] [read_claudemem.go](file:///Users/shuk/projects/cc-plugin/cmd/read_claudemem.go)

#### [MODIFY] [write_agentmemory.go](file:///Users/shuk/projects/cc-plugin/cmd/write_agentmemory.go)

#### [MODIFY] [write_mempalace.go](file:///Users/shuk/projects/cc-plugin/cmd/write_mempalace.go)

#### [MODIFY] [retain.go](file:///Users/shuk/projects/cc-plugin/cmd/retain.go)

---

## Bite-Sized Tasks

### Task 1: Phase 11 — OllamaLLM Service & Extract Command

**Files:**

- Create: `cmd/ollama.go`
- Modify: `cmd/state.go`
- Test: `cmd/ollama_test.go`

- [x] **Step 1: 於 `cmd/state.go` 中新增 Candidate 與輔助函數**

    ```go
    type Candidate struct {
     Text             string     `json:"text"`
     Entities         []string   `json:"entities"`
     Kind             string     `json:"kind"` // "fact" | "experience" | "preference" | "inference"
     FirstPerson      bool       `json:"first_person"`
     ConfirmedByHuman bool       `json:"confirmed_by_human"`
     SourceRefs       [][]string `json:"source_refs"` // [[source, source_id], ...]
    }
    ```

- [x] **Step 2: 撰寫 `cmd/ollama.go` 中的 `OllamaService` 與 `extract` 指令**

    ```go
    package cmd

    import (
     "bytes"
     "encoding/json"
     "fmt"
     "io"
     "net/http"
     "os"
     "strings"
     "time"

     "github.com/spf13/cobra"
     "github.com/spf13/viper"
    )

    type OllamaService struct {
     Model   string
     Host    string
     Timeout time.Duration
    }

    func NewOllamaService() *OllamaService {
     host := viper.GetString("llm.host")
     if host == "" {
      host = "http://localhost:11434"
     }
     model := viper.GetString("llm.model")
     if model == "" {
      model = "qwen2.5"
     }
     return &OllamaService{
      Model:   model,
      Host:    strings.TrimSuffix(host, "/"),
      Timeout: 120 * time.Second,
     }
    }

    const ExtractSystemPrompt = `You extract durable, reusable memories from agent/chat observations. ` +
     `Return ONLY a JSON object: {"candidates": [...]}. Each candidate has: ` +
     `text (verbatim statement worth remembering), ` +
     `entities (list of canonical names: people/projects/topics), ` +
     `kind (one of "fact","experience","preference","inference"), ` +
     `first_person (true if it is the human's own first-person life fact/experience), ` +
     `confirmed_by_human (true only if the human explicitly confirmed it). ` +
     `Omit chit-chat and transient operational noise.`

    func (s *OllamaService) Extract(observations []Observation) ([]Candidate, error) {
     if len(observations) == 0 {
      return nil, nil
     }

     var parts []string
     for _, o := range observations {
      parts = append(parts, fmt.Sprintf("[%s:%s] %s", o.Source, o.SourceID, o.Text))
     }
     joined := strings.Join(parts, "\n\n")

     payload := map[string]interface{}{
      "model": s.Model,
      "messages": []map[string]string{
       {"role": "system", "content": ExtractSystemPrompt},
       {"role": "user", "content": joined},
      },
      "format": "json",
      "stream": false,
     }

     payloadBytes, err := json.Marshal(payload)
     if err != nil {
      return nil, err
     }

     url := s.Host + "/api/chat"
     req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
     if err != nil {
      return nil, err
     }
     req.Header.Set("Content-Type", "application/json")

     client := &http.Client{Timeout: s.Timeout}
     resp, err := client.Do(req)
     if err != nil {
      return nil, fmt.Errorf("failed to call Ollama API: %w", err)
     }
     defer resp.Body.Close()

     if resp.StatusCode != http.StatusOK {
      body, _ := io.ReadAll(resp.Body)
      return nil, fmt.Errorf("Ollama returned status %d: %s", resp.StatusCode, string(body))
     }

     var reply struct {
      Message struct {
       Content string `json:"content"`
      } `json:"message"`
     }
     if err := json.NewDecoder(resp.Body).Decode(&reply); err != nil {
      return nil, fmt.Errorf("failed to decode Ollama reply: %w", err)
     }

     var responseObj struct {
      Candidates []Candidate `json:"candidates"`
     }
     if err := json.Unmarshal([]byte(reply.Message.Content), &responseObj); err != nil {
      return nil, fmt.Errorf("failed to parse structured JSON from Ollama content: %w. raw content: %s", err, reply.Message.Content)
     }

     var refs [][]string
     for _, o := range observations {
      refs = append(refs, []string{o.Source, o.SourceID})
     }

     for i := range responseObj.Candidates {
      responseObj.Candidates[i].SourceRefs = refs
     }

     return responseObj.Candidates, nil
    }

    func ExtractCmd() *cobra.Command {
     cmd := &cobra.Command{
      Use:   "extract",
      Short: "Directly extract memories from JSON observations on stdin using Ollama",
      RunE: func(cmd *cobra.Command, args []string) error {
       var observations []Observation
       if err := json.NewDecoder(os.Stdin).Decode(&observations); err != nil {
        return fmt.Errorf("failed to parse observations from stdin: %w", err)
       }

       svc := NewOllamaService()
       candidates, err := svc.Extract(observations)
       if err != nil {
        return err
       }

       output, err := json.MarshalIndent(candidates, "", "  ")
       if err != nil {
        return err
       }
       fmt.Println(string(output))
       return nil
      },
     }
     return cmd
    }
    ```

- [x] **Step 3: 撰寫 `cmd/ollama_test.go` 單元測試**

    ```go
    package cmd

    import (
     "encoding/json"
     "net/http"
     "net/http/httptest"
     "testing"

     "github.com/spf13/viper"
    )

    func TestOllamaExtract(t *testing.T) {
     mockResponse := map[string]interface{}{
      "message": map[string]string{
       "content": `{"candidates": [{"text": "likes tea", "entities": ["alice"], "kind": "preference", "first_person": false, "confirmed_by_human": true}]}`,
      },
     }
     responseBytes, _ := json.Marshal(mockResponse)

     server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Header) {
      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(http.StatusOK)
      w.Write(responseBytes)
     }))
     defer server.Close()

     viper.Set("llm.host", server.URL)
     viper.Set("llm.model", "test-model")

     svc := NewOllamaService()
     obs := []Observation{{Source: "test", SourceID: "1", Text: "alice likes tea"}}
     cands, err := svc.Extract(obs)
     if err != nil {
      t.Fatalf("Extract failed: %v", err)
     }

     if len(cands) != 1 {
      t.Fatalf("expected 1 candidate, got %d", len(cands))
     }
     if cands[0].Text != "likes tea" || cands[0].Kind != "preference" {
      t.Errorf("unexpected candidate values: %+v", cands[0])
     }
    }
    ```

- [x] **Step 4: 執行並確認測試通過**

    Run: `go test -v -run TestOllamaExtract ./cmd`
    Expected: PASS

- [x] **Step 5: Commit**

    ```bash
    git add cmd/ollama.go cmd/state.go cmd/ollama_test.go
    git commit -m "feat(distiller): add OllamaLLM service and extract command"
    ```

---

### Task 2: Phase 12 — settings.json config & pipeline run command

**Files:**

- Modify: `cmd/root.go`
- Modify: `cmd/state.go`
- Modify: `cmd/read_gbrain.go`
- Modify: `cmd/read_claudemem.go`
- Modify: `cmd/write_agentmemory.go`
- Modify: `cmd/write_mempalace.go`
- Modify: `cmd/retain.go`
- Create: `cmd/run.go`
- Test: `cmd/run_test.go`

- [x] **Step 1: 在 `cmd/state.go` 中新增 `Fingerprint` 輔助函數與 `expandPath` 函數**

    ```go
    // In cmd/state.go:
    import (
     "crypto/sha256"
     "encoding/hex"
     "sort"
     "strings"
    )

    func Fingerprint(text string, entities []string) string {
     normalizedText := strings.ToLower(strings.Join(strings.Fields(text), " "))
     sortedEntities := make([]string, len(entities))
     copy(sortedEntities, entities)
     sort.Strings(sortedEntities)

     h := sha256.New()
     h.Write([]byte(normalizedText))
     h.Write([]byte("|"))
     h.Write([]byte(strings.Join(sortedEntities, "|")))
     return hex.EncodeToString(h.Sum(nil))
    }

    func expandPath(p string) string {
     if strings.HasPrefix(p, "~") {
      home, _ := os.UserHomeDir()
      return filepath.Join(home, p[1:])
     }
     return p
    }
    ```

- [x] **Step 2: 修改 `cmd/root.go` 引入 `gosdk/config` 與預設 JSON 載入**

    ```go
    package cmd

    import (
     "fmt"
     "os"

     "github.com/bizshuk/gosdk/config"
     "github.com/spf13/cobra"
    )

    var RootCmd = &cobra.Command{
     Use:   "distiller",
     Short: "Distiller CLI manages memory systems cross laptop and server",
    }

    const defaultSettings = `{
      "state": {
        "db_path": "~/.distiller/state.db"
      },
      "retention": {
        "max_age_days": 30
      },
      "llm": {
        "model": "qwen2.5",
        "host": "http://localhost:11434"
      },
      "sources": {
        "claude_mem": {
          "db_path": "~/.claude-mem/claude-mem.db",
          "table": "observations",
          "id_col": "id",
          "ts_col": "created_at_epoch",
          "text_col": "text"
        },
        "gbrain_working": {
          "root": "~/brain/working"
        }
      },
      "stores": {
        "agentmemory": {
          "url": "http://localhost:3111/agentmemory/remember"
        },
        "mempalace": {
          "wing": "main",
          "temp_dir": "/tmp/mempalace-temp"
        }
      }
    }`

    func init() {
     config.Default(
      config.WithAppName("cc-plugin"),
      config.WithDefaultValue(defaultSettings),
     )

     RootCmd.AddCommand(RetainCmd())
     RootCmd.AddCommand(ReadGbrainCmd())
     RootCmd.AddCommand(ReadClaudeMemCmd())
     RootCmd.AddCommand(WriteAgentMemoryCmd())
     RootCmd.AddCommand(WriteMempalaceCmd())
     RootCmd.AddCommand(ExtractCmd())
     RootCmd.AddCommand(RunCmd())
    }

    func Execute() {
     if err := RootCmd.Execute(); err != nil {
      fmt.Fprintf(os.Stderr, "Error: %v\n", err)
      os.Exit(1)
     }
    }
    ```

- [x] **Step 3: 重構各指令檔案，使其邏輯功能為導出函數並自 Viper 讀取設定值**

    重構 `read_gbrain.go`、`read_claudemem.go`、`write_agentmemory.go`、`write_mempalace.go`、`retain.go`。
    例如，在 `read_gbrain.go`：

    ```go
    func readGbrain(store *StateStore, workingDir string) ([]Observation, int64, error) { ... }
    ```

    在各 CLI Cmd 的 `RunE` 呼叫這些導出函數，並且當 flags 沒有傳入時，優先讀取 `viper.GetString` 設定（路徑透過 `expandPath` 展開）。

- [x] **Step 4: 實作 `cmd/run.go` 主管線整合指令**

    `cmd/run.go` 將執行以下動作：
    1. 開啟 `StateStore`。
    2. 讀取 `gbrain-working` 與 `claude-mem` 新增的 observations。
    3. 合併 observations 並呼叫 `OllamaService.Extract` 進行記憶抽取。
    4. 對於每一個 `Candidate`：
        - 用 `Fingerprint` 算雜湊值。
        - 用 `RecordSeen` 佐證。
        - 寫入 `agentmemory` API。
        - 用 `QualifiesForTruth` 驗證政策判斷是否可升級寫入 `mempalace`。
    5. 標記 observations 為 distilled 且更新 cursors。
    6. 在結尾調用 `retain` 清理 (除非傳入 `--no-retain`)。

    ```go
    package cmd

    import (
     "fmt"
     "os"
     "path/filepath"
     "time"

     "github.com/spf13/cobra"
     "github.com/spf13/viper"
    )

    func QualifiesForTruth(c Candidate, corroboration int) bool {
     if c.Kind == "inference" {
      return false
     }
     if c.ConfirmedByHuman {
      return true
     }
     if c.FirstPerson && (c.Kind == "fact" || c.Kind == "experience") {
      return true
     }
     if corroboration >= 2 {
      return true
     }
     return false
    }

    func RunCmd() *cobra.Command {
     var noRetain bool

     cmd := &cobra.Command{
      Use:   "run",
      Short: "Run the full memory distillation pipeline",
      RunE: func(cmd *cobra.Command, args []string) error {
       statePath := expandPath(viper.GetString("state.db_path"))
       store, err := NewStateStore(statePath)
       if err != nil {
        return err
       }
       defer store.Close()

       // 1. Read gbrain
       gbrainRoot := expandPath(viper.GetString("sources.gbrain_working.root"))
       gbrainObs, gbrainMaxTS, err := readGbrainLogic(store, gbrainRoot)
       if err != nil {
        return err
       }

       // 2. Read claude-mem
       cmDB := expandPath(viper.GetString("sources.claude_mem.db_path"))
       cmTable := viper.GetString("sources.claude_mem.table")
       cmIdCol := viper.GetString("sources.claude_mem.id_col")
       cmTsCol := viper.GetString("sources.claude_mem.ts_col")
       cmTextCol := viper.GetString("sources.claude_mem.text_col")
       cmObs, cmMaxTS, err := readClaudeMemLogic(store, cmDB, cmTable, cmIdCol, cmTsCol, cmTextCol)
       if err != nil {
        return err
       }

       allObs := append(gbrainObs, cmObs...)
       if len(allObs) == 0 {
        fmt.Println("[distiller] No new observations found.")
        return nil
       }

       // 3. Extract via Ollama
       llm := NewOllamaService()
       candidates, err := llm.Extract(allObs)
       if err != nil {
        return err
       }

       // 4. Process candidates
       var memories []Memory
       var facts []Fact
       now := time.Now().Unix()

       for _, c := range candidates {
        fp := Fingerprint(c.Text, c.Entities)

        // Record seen source count
        var corroboration int
        for _, ref := range c.SourceRefs {
         if len(ref) > 0 {
          count, err := store.RecordSeen(fp, ref[0])
          if err != nil {
           return err
          }
          corroboration = count
         }
        }

        memories = append(memories, Memory{
         Fingerprint: fp,
         Text:        c.Text,
         Entities:    c.Entities,
         Kind:        c.Kind,
         CreatedAt:   now,
        })

        if QualifiesForTruth(c, corroboration) {
         facts = append(facts, Fact{
          Fingerprint: fp,
          Text:        c.Text,
          Entities:    c.Entities,
          Evidence:    c.SourceRefs,
          CreatedAt:   now,
         })
        }
       }

       // Write Memories to agentmemory
       if len(memories) > 0 {
        url := viper.GetString("stores.agentmemory.url")
        if err := writeAgentMemoryLogic(memories, url); err != nil {
         return err
        }
       }

       // Write Facts to mempalace
       if len(facts) > 0 {
        tempDir := expandPath(viper.GetString("stores.mempalace.temp_dir"))
        wing := viper.GetString("stores.mempalace.wing")
        if err := writeMempalaceLogic(facts, tempDir, wing); err != nil {
         return err
        }
       }

       // Mark distilled and update cursors
       for _, o := range allObs {
        if err := store.MarkDistilled(o.Source, o.SourceID, now); err != nil {
         return err
        }
       }

       if gbrainMaxTS > 0 {
        store.SetCursor("gbrain-working", gbrainMaxTS)
       }
       if cmMaxTS > 0 {
        store.SetCursor("claude-mem", cmMaxTS)
       }

       fmt.Printf("[distiller] Pipeline ran successfully. Sources read=%d, Memories written=%d, Facts written=%d\n", len(allObs), len(memories), len(facts))

       // 5. Sweep retention
       if !noRetain {
        maxAgeDays := viper.GetInt("retention.max_age_days")
        if maxAgeDays == 0 {
         maxAgeDays = 30
        }
        pruneGbrainDir := expandPath(viper.GetString("sources.gbrain_working.root"))
        if err := retainLogic(store, maxAgeDays, pruneGbrainDir); err != nil {
         return err
        }
       }

       return nil
      },
     }

     cmd.Flags().BoolVar(&noRetain, "no-retain", false, "Disable retention sweep after run")
     return cmd
    }
    ```

- [x] **Step 5: 撰寫 `cmd/run_test.go` 單元測試**

    ```go
    package cmd

    import (
     "testing"
    )

    func TestQualifiesForTruth(t *testing.T) {
     c1 := Candidate{Kind: "inference", ConfirmedByHuman: true}
     if QualifiesForTruth(c1, 5) {
      t.Error("inference should never qualify")
     }

     c2 := Candidate{Kind: "preference", ConfirmedByHuman: true}
     if !QualifiesForTruth(c2, 1) {
      t.Error("human confirmed should qualify")
     }

     c3 := Candidate{Kind: "fact", FirstPerson: true}
     if !QualifiesForTruth(c3, 1) {
      t.Error("first person fact should qualify")
     }

     c4 := Candidate{Kind: "fact", FirstPerson: false}
     if QualifiesForTruth(c4, 1) {
      t.Error("unconfirmed single source fact should not qualify")
     }
     if !QualifiesForTruth(c4, 2) {
      t.Error("corroborated fact should qualify")
     }
    }
    ```

- [x] **Step 6: 執行測試並確認綠燈**

    Run: `go test -v ./cmd`
    Expected: PASS

- [x] **Step 7: Commit**

    ```bash
    git add cmd/run.go cmd/run_test.go cmd/root.go
    git commit -m "feat(distiller): implement distiller run pipeline and settings loading"
    ```

---

### Task 3: Phase 13 — crontab.txt creation

**Files:**

- Create: `crontab.txt`

- [x] **Step 1: 建立 `crontab.txt` 檔案於專案根目錄**

    ```text
    # 每晨 03:00 執行 distiller 記憶蒸餾整合管道
    0 3 * * * cd /Users/shuk/projects/cc-plugin && ./distiller run >> $HOME/.distiller/logs/run.log 2>&1
    ```

- [x] **Step 2: Commit**

    ```bash
    git add crontab.txt
    git commit -m "chore(distiller): add crontab.txt template for manual crontab registration"
    ```

---

## Verification Plan

### Automated Tests

- 在工作區根目錄下執行 `go test -v ./...` 以確認單元測試通過。
- 再次執行整合測試腳本 `verify_distiller.sh` 以確保新架構不影響原有功能。
