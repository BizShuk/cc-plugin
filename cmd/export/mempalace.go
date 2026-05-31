package export

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DrawerRow represents a drawer's data retrieved from the Chroma database.
type DrawerRow struct {
	EmbeddingID string `gorm:"column:embedding_id"`
	Wing        string `gorm:"column:wing"`
	Room        string `gorm:"column:room"`
	Document    string `gorm:"column:document"`
	SourceFile  string `gorm:"column:source_file"`
	FiledAt     string `gorm:"column:filed_at"`
	AddedBy     string `gorm:"column:added_by"`
}

// ReadPalaceData connects to the SQLite database and queries all drawer records.
func ReadPalaceData(dbPath string) ([]DrawerRow, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	sqlDB, err := db.DB()
	if err == nil {
		defer sqlDB.Close()
	}

	var rows []DrawerRow
	query := `
SELECT 
    e.embedding_id,
    MAX(CASE WHEN m.key = 'wing' THEN m.string_value END) AS wing,
    MAX(CASE WHEN m.key = 'room' THEN m.string_value END) AS room,
    MAX(CASE WHEN m.key = 'chroma:document' THEN m.string_value END) AS document,
    MAX(CASE WHEN m.key = 'source_file' THEN m.string_value END) AS source_file,
    MAX(CASE WHEN m.key = 'filed_at' THEN m.string_value END) AS filed_at,
    MAX(CASE WHEN m.key = 'added_by' THEN m.string_value END) AS added_by
FROM embeddings e
JOIN embedding_metadata m ON e.id = m.id
WHERE e.embedding_id LIKE 'drawer_%'
GROUP BY e.id
ORDER BY wing, room, e.embedding_id;
`
	if err := db.Raw(query).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("failed to query drawer rows: %w", err)
	}

	return rows, nil
}

// ExportCategories exports categories (wing, room, drawer_id) into a CSV format.
func ExportCategories(rows []DrawerRow, outputFolder string) error {
	var out *os.File
	var err error

	if outputFolder != "" {
		if err := os.MkdirAll(outputFolder, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", outputFolder, err)
		}
		filePath := filepath.Join(outputFolder, "mempalace_categories.csv")
		out, err = os.Create(filePath)
		if err != nil {
			return fmt.Errorf("failed to create CSV file: %w", err)
		}
		defer out.Close()
		fmt.Printf("Exporting categories to %s...\n", filePath)
	} else {
		out = os.Stdout
	}

	writer := csv.NewWriter(out)
	defer writer.Flush()

	if err := writer.Write([]string{"wing", "room", "drawer_id"}); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	for _, r := range rows {
		wing := r.Wing
		if wing == "" {
			wing = "unknown"
		}
		room := r.Room
		if room == "" {
			room = "general"
		}
		if err := writer.Write([]string{wing, room, r.EmbeddingID}); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	return nil
}

// ExportData exports drawer contents to Markdown files organized by wing/room and creates an index.md.
func ExportData(rows []DrawerRow, outputFolder string) error {
	if outputFolder == "" {
		outputFolder = "mempalace_export"
	}

	if err := os.MkdirAll(outputFolder, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", outputFolder, err)
	}

	type WingRoomKey struct {
		Wing string
		Room string
	}

	grouped := make(map[WingRoomKey][]DrawerRow)
	wingStats := make(map[string]map[string]int) // wing -> room -> count
	var wingsList []string
	roomsInWing := make(map[string][]string)

	for _, r := range rows {
		wing := r.Wing
		if wing == "" {
			wing = "unknown"
		}
		room := r.Room
		if room == "" {
			room = "general"
		}
		key := WingRoomKey{Wing: wing, Room: room}
		grouped[key] = append(grouped[key], r)

		if _, ok := wingStats[wing]; !ok {
			wingStats[wing] = make(map[string]int)
			wingsList = append(wingsList, wing)
		}
		if wingStats[wing][room] == 0 {
			roomsInWing[wing] = append(roomsInWing[wing], room)
		}
		wingStats[wing][room]++
	}

	sort.Strings(wingsList)
	for w := range roomsInWing {
		sort.Strings(roomsInWing[w])
	}

	fmt.Printf("Streaming %d drawers...\n", len(rows))

	for _, wing := range wingsList {
		safeWing := sanitizePathComponent(wing)
		wingDir := filepath.Join(outputFolder, safeWing)
		if err := os.MkdirAll(wingDir, 0755); err != nil {
			return fmt.Errorf("failed to create wing directory %s: %w", wingDir, err)
		}

		for _, room := range roomsInWing[wing] {
			safeRoom := sanitizePathComponent(room)
			roomPath := filepath.Join(wingDir, fmt.Sprintf("%s.md", safeRoom))

			file, err := os.Create(roomPath)
			if err != nil {
				return fmt.Errorf("failed to create room file: %w", err)
			}

			// Write Header
			fmt.Fprintf(file, "# %s / %s\n\n", wing, room)

			key := WingRoomKey{Wing: wing, Room: room}
			drawers := grouped[key]
			for _, d := range drawers {
				source := d.SourceFile
				if source == "" {
					source = "unknown"
				}
				filed := d.FiledAt
				if filed == "" {
					filed = "unknown"
				}
				addedBy := d.AddedBy
				if addedBy == "" {
					addedBy = "unknown"
				}

				quotedDoc := quoteContent(d.Document)

				fmt.Fprintf(file, "## %s\n\n> %s\n\n| Field | Value |\n|-------|-------|\n| Source | %s |\n| Filed | %s |\n| Added by | %s |\n\n---\n\n",
					d.EmbeddingID,
					quotedDoc,
					source,
					filed,
					addedBy,
				)
			}
			file.Close()
		}
	}

	// Write index.md
	today := time.Now().Format("2006-01-02")
	indexPath := filepath.Join(outputFolder, "index.md")
	indexFile, err := os.Create(indexPath)
	if err != nil {
		return fmt.Errorf("failed to create index file: %w", err)
	}
	defer indexFile.Close()

	fmt.Fprintf(indexFile, "# Palace Export — %s\n\n", today)
	fmt.Fprintln(indexFile, "| Wing | Rooms | Drawers |")
	fmt.Fprintln(indexFile, "|------|-------|---------|")

	var totalDrawers int
	for _, wing := range wingsList {
		rooms := wingStats[wing]
		var wingDrawers int
		for _, count := range rooms {
			wingDrawers += count
		}
		totalDrawers += wingDrawers
		fmt.Fprintf(indexFile, "| [%s](%s/) | %d | %d |\n", wing, sanitizePathComponent(wing), len(rooms), wingDrawers)
		fmt.Printf("  %s: %d rooms, %d drawers\n", wing, len(rooms), wingDrawers)
	}

	fmt.Printf("\n  Exported %d drawers across %d wings, %d rooms\n", totalDrawers, len(wingsList), len(grouped))
	fmt.Printf("  Output: %s\n", outputFolder)

	return nil
}

// MempalaceCmd returns the mempalace subcommand for export.
func MempalaceCmd() *cobra.Command {
	var dataFlag bool
	var outputFolder string

	cmd := &cobra.Command{
		Use:   "mempalace",
		Short: "Export MemPalace drawers data",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Find palace path
			configDir := viper.GetString("stores.mempalace.config_dir")
			configDirExpanded, err := homedir.Expand(configDir)
			if err != nil {
				configDirExpanded = configDir
			}
			palacePath := filepath.Join(configDirExpanded, "palace")
			configPath := filepath.Join(configDirExpanded, "config.json")

			if _, err := os.Stat(configPath); err == nil {
				data, err := os.ReadFile(configPath)
				if err == nil {
					var config struct {
						PalacePath string `json:"palace_path"`
					}
					if err := json.Unmarshal(data, &config); err == nil && config.PalacePath != "" {
						palacePath = config.PalacePath
					}
				}
			}

			palacePathExpanded, err := homedir.Expand(palacePath)
			if err != nil {
				palacePathExpanded = palacePath
			}
			dbPath := filepath.Join(palacePathExpanded, "chroma.sqlite3")

			if _, err := os.Stat(dbPath); os.IsNotExist(err) {
				return fmt.Errorf("palace database not found at: %s", dbPath)
			}

			fmt.Printf("Connecting to MemPalace at: %s\n", palacePathExpanded)

			rows, err := ReadPalaceData(dbPath)
			if err != nil {
				return err
			}

			if dataFlag {
				return ExportData(rows, outputFolder)
			}

			return ExportCategories(rows, outputFolder)
		},
	}

	cmd.Flags().BoolVar(&dataFlag, "data", false, "Export all details and document content to markdown files")
	cmd.Flags().StringVarP(&outputFolder, "output", "o", "", "Output folder name")

	return cmd
}

func sanitizePathComponent(name string) string {
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range invalidChars {
		name = strings.ReplaceAll(name, char, "_")
	}
	name = strings.Trim(name, ". ")
	if name == "" {
		return "unknown"
	}
	return name
}

func quoteContent(text string) string {
	text = strings.TrimRight(text, "\n")
	lines := strings.Split(text, "\n")
	return strings.Join(lines, "\n> ")
}
