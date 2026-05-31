package model

type Memory struct {
	Fingerprint string   `json:"fingerprint"`
	Text        string   `json:"text"`
	Entities    []string `json:"entities"`
	Kind        string   `json:"kind"`
	CreatedAt   int64    `json:"created_at"`
}
