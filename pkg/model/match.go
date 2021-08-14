package model

// Match is a relationship between two players.
type Match struct {
	Player1 *Player
	Player2 *Player
}

// MatchArrangement is a slice of matches that represents an arrangement of
// multiple simultaneous matches
type MatchArrangement []Match
