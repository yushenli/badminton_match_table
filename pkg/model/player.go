package model

// Player represents a player, with it's necessary information to be arranged.
type Player struct {
	Priority  float32
	Name      string
	Score     float32
	Matches   float32
	Opponents map[*Player]int
}

// PlayerSlice is an alias for a slice of Player pointers.
type PlayerSlice []*Player
