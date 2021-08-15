package arranger

import (
	"fmt"
	"math"
	"sort"

	"github.com/yushenli/badminton_match_table/pkg/model"
)

func sortPlayerSliceByScorePriority(slice model.PlayerSlice) {
	sort.Slice(slice, func(i, j int) bool {
		a := slice[i]
		b := slice[j]
		return a.Score > b.Score || (a.Score == b.Score && a.Priority > b.Priority)
	})
}

// SeparateCompetedPlayers will, for a given subslice of 2N players, find the best combination
// of N pairs where the total times that the two players within each pair have played with each other
// if minimized. When multiple solution exists with the same total times, SeparateCompetedPlayers
// will try to pair the players with the closest scores.
// In the end, the provided PlayerSlice will be rearranged. Every two consecutive players will
// be considered a pair.
//
// The input PlayerSlice is expected to have been sorted by sortPlayerSliceByScorePriority() already.
//
// start: the index of the first player in the subslice within the PlayerSlice
// end: the index of the last player (inclusive) in the subslice within the PlayerSlice
// endFix: whether the last player can be rearranged in the slice.
// Note: you might want to fix the last player in the slice because they may
// be at different levels than the other players to be rearranged. In that case, the last player
// will always stay at the last in the subslice, and it will be only used to calculate the
// time of matches against the second last player in the arrangement.
//
// Returns error if there are odd number of players between start and end (inclusive).
func SeparateCompetedPlayers(players model.PlayerSlice, start, end int, endFixed bool) error {
	if start >= end {
		return fmt.Errorf("start has to be smaller than end, you provided start(%d) and end(%d)", start, end)
	}
	if start < 0 || end >= len(players) {
		return fmt.Errorf("start and/or out of range, PlayerSlice length %d, you provided start(%d) and end(%d)",
			len(players), start, end)
	}
	if (end-start+1)%2 != 0 {
		return fmt.Errorf("Sub slice to be rearranged must have even number of elements, you provided start(%d) and end(%d)", start, end)
	}

	bestPairs := make([]model.Side, (end-start+1)/2)
	currentPairs := make([]model.Side, (end-start+1)/2)
	used := make([]bool, len(players))
	var minTotalMatches int
	minTotalMatches = math.MaxInt32
	var minReversedScores float32
	minReversedScores = math.MaxFloat32

	tryArrange(players, &bestPairs, currentPairs, 0, start, end, endFixed, used, &minTotalMatches, &minReversedScores)

	for i := 0; i < len(bestPairs); i++ {
		players[start+i*2] = bestPairs[i].Player1
		players[start+i*2+1] = bestPairs[i].Player2
	}

	return nil
}

func tryArrange(players model.PlayerSlice, bestPairs *[]model.Side, pairs []model.Side, level, start, end int, endFixed bool, used []bool, minTotalMatches *int, minReversedScores *float32) {
	i := start
	for used[i] {
		i++
	}
	used[i] = true

	maxj := end
	if endFixed && level < len(*bestPairs)-1 {
		// If the last player is supposed to be fixed in position, only the last level of the
		// recursion can (and will) use the last player.
		maxj = end - 1
	}

	for j := i + 1; j <= maxj; j++ {
		if used[j] {
			continue
		}

		pairs[level] = model.Side{
			Player1: players[i],
			Player2: players[j],
		}
		used[j] = true

		if level == len(*bestPairs)-1 {
			checkBestArrangement(bestPairs, pairs, minTotalMatches, minReversedScores)
		} else {
			tryArrange(players, bestPairs, pairs, level+1, i+1, end, endFixed, used, minTotalMatches, minReversedScores)
		}

		used[j] = false
	}
	used[i] = false
}

func checkBestArrangement(bestPairs *[]model.Side, pairs []model.Side, minTotalMatches *int, minReversedScores *float32) {
	var totalMatches int
	var reversedScores float32
	var lastPlayer2 *model.Player

	for _, pair := range pairs {
		if matches, ok := pair.Player1.Opponents[pair.Player2]; ok {
			totalMatches += matches
		}
		if totalMatches > *minTotalMatches {
			return
		}

		if lastPlayer2 != nil && pair.Player1.Score > lastPlayer2.Score {
			reversedScores += pair.Player1.Score - lastPlayer2.Score
		}

		lastPlayer2 = pair.Player2
	}

	if totalMatches < *minTotalMatches || reversedScores < *minReversedScores {
		*bestPairs = make([]model.Side, len(pairs))
		copy(*bestPairs, pairs)
		*minTotalMatches = totalMatches
		*minReversedScores = reversedScores
	}
}
