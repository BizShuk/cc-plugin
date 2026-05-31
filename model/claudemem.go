package model

type ClaudeMemObservation struct {
	ID             string `gorm:"column:id;primaryKey"`
	CreatedAtEpoch int64  `gorm:"column:created_at_epoch"`
	Text           string `gorm:"column:text"`
}

func (ClaudeMemObservation) TableName() string {
	return "observations"
}
