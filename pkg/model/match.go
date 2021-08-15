package model

// Side represents one side of a match, which can be one player for singles
// or two players for doubles.
type Side struct {
	Player1 *Player
	Player2 *Player
}

// Match is a relationship between two sides.
type Match struct {
	Side1 Side
	Side2 Side
}

// MatchArrangement is a slice of matches that represents an arrangement of
// multiple simultaneous matches.
type MatchArrangement []Match
