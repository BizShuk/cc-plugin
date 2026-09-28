package model

type Observation struct {
	Source    string `json:"source"`
	SourceID  string `json:"source_id"`
	Timestamp int64  `json:"timestamp"`
	Text      string `json:"text"`
}
