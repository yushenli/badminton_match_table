package arranger

import (
	"fmt"
	"math/rand"
	"sort"

	"github.com/yushenli/badminton_match_table/pkg/model"
)

// canPlayCount returns, for given number of courts and players, how many player can play
// the next round, and how many courts will be used for single and double matches.
// If there are too few players, an error will be thrown.
func canPlayCount(courtCount, playerCount int) (canPlayCount, singleCount, doubleCount int, err error) {
	if playerCount < courtCount*2 {
		return 0, 0, 0, fmt.Errorf("%d players are not enough for %d courts even for singles", playerCount, courtCount)
	}
	singleCount = courtCount
	doubleCount = 0
	playerCount -= courtCount * 2

	doubleCount = playerCount / 2
	if doubleCount > courtCount {
		doubleCount = courtCount
	}
	singleCount -= doubleCount

	return singleCount*2 + doubleCount*4, singleCount, doubleCount, nil
}

// PickPlayersForCourts picks players to play based on the number of courts available.
// PickPlayersForCourts will schedule double matches on all courts at first, if there are not enough players
// to play doubles on all courts, it will schedule some courts to host singles.
// If there are not even enough players to fill single matches on all courts, an error will be returned.
// Players will be picked from those who have played the least times. For ties, players with higher prioritiy
// will be picked first.
// The passed in player slice will not be affected. A new slice of players will be returned.
func PickPlayersForCourts(players model.PlayerSlice, courtCount int) (model.PlayerSlice, error) {
	canPlay, _, _, err := canPlayCount(courtCount, len(players))
	if err != nil {
		return nil, err
	}

	ret := make(model.PlayerSlice, len(players))
	copy(ret, players)

	sort.Slice(ret, func(i, j int) bool {
		a := ret[i]
		b := ret[j]
		return a.Matches < b.Matches || (a.Matches == b.Matches && a.Priority > b.Priority)
	})
	ret = ret[:canPlay]

	return ret, nil
}

// MakeMatchArrangements put a processed slice of players into courtCount matches.
// When not all matches are doubles, which matches will be single and doubles are
// randomized using the provided seed.
func MakeMatchArrangements(players model.PlayerSlice, courtCount int, seed int) (model.MatchArrangement, error) {
	_, singles, doubles, err := canPlayCount(courtCount, len(players))
	if err != nil {
		return nil, err
	}

	isDoubles := make([]bool, singles+doubles)
	for i := 0; i < doubles; i++ {
		isDoubles[i] = true
	}

	rand.Seed(int64(seed))
	rand.Shuffle(singles+doubles, func(i, j int) {
		isDoubles[i], isDoubles[j] = isDoubles[j], isDoubles[i]
	})

	matches := make(model.MatchArrangement, singles+doubles)
	idx := 0
	for i := 0; i < singles+doubles; i++ {
		if isDoubles[i] {
			SeparateCompetedPlayers(players, idx, idx+3, false)
			matches[i] = model.Match{
				Side1: model.Side{
					Player1: players[idx],
					Player2: players[idx+1],
				},
				Side2: model.Side{
					Player1: players[idx+2],
					Player2: players[idx+3],
				},
			}
			idx += 4
		} else {
			matches[i] = model.Match{
				Side1: model.Side{
					Player1: players[idx],
					Player2: nil,
				},
				Side2: model.Side{
					Player1: players[idx+1],
					Player2: nil,
				},
			}
			idx += 2
		}
	}

	return matches, nil
}
