package gormmodel

import (
	"time"

	"gorm.io/gorm"
)

// Event represents a record in the match table.
type Event struct {
	gorm.Model
	Date     time.Time
	Key      string
	Location string
	Courts   int
}

// TableName overrides the default plural-form table name.
func (Event) TableName() string {
	return "event"
}
