package model

type Cursor struct {
	Source string `gorm:"primaryKey;column:source"`
	LastTs int64  `gorm:"column:last_ts;not null"`
}

func (Cursor) TableName() string {
	return "cursor"
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
