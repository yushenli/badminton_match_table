package gormmodel

import (
	"gorm.io/gorm"
)

// Match represents a record in the match table.
type Match struct {
	gorm.Model
	Eid    int
	Round  int
	Sid1   int
	Sid2   int
	Court  int
	Status string
	Side1  *Side `gorm:"foreignKey:sid1"`
	Side2  *Side `gorm:"foreignKey:sid2"`
}

// TableName overrides the default plural-form table name.
func (Match) TableName() string {
	return "match"
}
