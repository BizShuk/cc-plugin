package model

type Cursor struct {
	Source string `gorm:"primaryKey;column:source"`
	LastTs int64  `gorm:"column:last_ts;not null"`
	LastID int64  `gorm:"column:last_id;not null;default:0"`
}

func (Cursor) TableName() string {
	return "cursor"
}

// CursorPosition identifies the last exported record by timestamp and source ID.
type CursorPosition struct {
	LastTS int64
	LastID int64
}

type Seen struct {
	Fingerprint string `gorm:"primaryKey;column:fingerprint"`
	Source      string `gorm:"primaryKey;column:source"`
	FirstSeen   int64  `gorm:"column:first_seen;not null"`
}

func (Seen) TableName() string {
	return "seen"
}

type Distilled struct {
	Source      string `gorm:"primaryKey;column:source"`
	SourceID    string `gorm:"primaryKey;column:source_id"`
	DistilledAt int64  `gorm:"column:distilled_at;not null"`
}

func (Distilled) TableName() string {
	return "distilled"
}
