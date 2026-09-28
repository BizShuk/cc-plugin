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
