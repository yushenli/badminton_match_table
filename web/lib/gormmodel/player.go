package gormmodel

import (
	"gorm.io/gorm"
)

// Player represents a record in the player table.
type Player struct {
	gorm.Model
	Eid          int
	Name         string
	Priority     float32
	InitialScore float32
	InBreak      bool
}

// TableName overrides the default plural-form table name.
func (Player) TableName() string {
	return "player"
}
