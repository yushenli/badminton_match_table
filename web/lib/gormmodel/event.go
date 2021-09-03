package gormmodel

import (
	"time"

	"gorm.io/gorm"
)

// Event represents a record in the event table.
type Event struct {
	gorm.Model
	Date         time.Time
	Key          string
	Location     string
	Courts       int
	CurrentRound int
	AdminKey     string
}

// TableName overrides the default plural-form table name.
func (Event) TableName() string {
	return "event"
}
