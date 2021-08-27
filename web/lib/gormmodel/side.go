package gormmodel

import (
	"gorm.io/gorm"
)

// Side represents a record in the side table.
type Side struct {
	gorm.Model
	Eid   int
	Mid   int
	Pid1  int
	Pid2  *int // Pid2 is nullable
	Score float32
}

// TableName overrides the default plural-form table name.
func (Side) TableName() string {
	return "side"
}
