package model

type Fact struct {
	Fingerprint string     `json:"fingerprint"`
	Text        string     `json:"text"`
	Entities    []string   `json:"entities"`
	Evidence    [][]string `json:"evidence"`
	CreatedAt   int64      `json:"created_at"`
}
