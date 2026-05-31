package model

type Candidate struct {
	Text             string     `json:"text"`
	Entities         []string   `json:"entities"`
	Kind             string     `json:"kind"` // "fact" | "experience" | "preference" | "inference"
	FirstPerson      bool       `json:"first_person"`
	ConfirmedByHuman bool       `json:"confirmed_by_human"`
	SourceRefs       [][]string `json:"source_refs"` // [[source, source_id], ...]
}
