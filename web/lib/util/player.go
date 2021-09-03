package util

import (
	"log"

	"github.com/yushenli/badminton_match_table/pkg/model"
	"github.com/yushenli/badminton_match_table/web/lib/config"
	"github.com/yushenli/badminton_match_table/web/lib/gormmodel"
	"gorm.io/gorm"
)

// PlayerWithCounter is a gormmodel.Player embedded with games counters and a score.
type PlayerWithCounter struct {
	gormmodel.Player
	Games int
	Win   int
	Loss  int
	Score float32
}

// PopulatePlayers fetches all sides under an event and put them in a slice as well as a unique-key based map
func PopulatePlayers(eid int) ([]PlayerWithCounter, map[int]*PlayerWithCounter, error) {
	var players []PlayerWithCounter
	playerMap := make(map[int]*PlayerWithCounter)
	ret := config.DB.Where("eid = ?", eid).Find(&players)
	if ret.Error != nil {
		log.Printf("Failed to list players under event %d: %v", eid, ret.Error)
		return nil, nil, ret.Error
	}
	for idx := range players {
		playerMap[int(players[idx].ID)] = &(players[idx])
		players[idx].Score = players[idx].InitialScore
	}

	return players, playerMap, nil
}

// PopulateSides fetches all sides under an event and put them in a slice as well as a unique-key based map
// The player pointers inside the side objects will point to the players.
func PopulateSides(eid int, playerMap map[int]*PlayerWithCounter, execludeRound *int) ([]gormmodel.Side, map[int]*gormmodel.Side, error) {
	var sides []gormmodel.Side
	sideMap := make(map[int]*gormmodel.Side)
	var ret *gorm.DB
	if execludeRound == nil {
		ret = config.DB.Where("eid = ?", eid).Find(&sides)
	} else {
		ret = config.DB.Joins("JOIN `match` ON `match`.sid1 = side.ID OR `match`.sid2 = side.ID").
			Where("side.eid = ?", eid).Where("round != ?", *execludeRound).Find(&sides)
	}
	if ret.Error != nil {
		log.Printf("Failed to list sides under event %d: %v", eid, ret.Error)
		return nil, nil, ret.Error
	}
	for idx := range sides {
		sideMap[int(sides[idx].ID)] = &(sides[idx])
		sides[idx].Player1 = &playerMap[sides[idx].Pid1].Player
		if sides[idx].Pid2 == nil {
			continue
		}
		if _, ok := playerMap[*sides[idx].Pid2]; ok {
			sides[idx].Player2 = &playerMap[*sides[idx].Pid2].Player
		}
	}

	return sides, sideMap, nil
}

// PopulateMatches fetches all matches under an event and put them in a slice as well as a unique-key based map
// The side pointers inside the match objects will point to the sides.
func PopulateMatches(eid, currentRound int, sideMap map[int]*gormmodel.Side) ([]gormmodel.Match, [][]*gormmodel.Match, error) {
	var matches []gormmodel.Match
	matchesByRound := make([][]*gormmodel.Match, currentRound)

	ret := config.DB.Where("eid = ?", eid).Order("round").Order("court").Find(&matches)
	if ret.Error != nil {
		log.Printf("Failed to list matches under event %d: %v", eid, ret.Error)
		return nil, nil, ret.Error
	}

	for idx := range matches {
		roundIdx := matches[idx].Round - 1
		matchesByRound[roundIdx] = append(matchesByRound[roundIdx], &matches[idx])

		matches[idx].Side1 = sideMap[matches[idx].Sid1]
		matches[idx].Side2 = sideMap[matches[idx].Sid2]
	}

	return matches, matchesByRound, nil
}

// FillPlayerCounter calculates the game counter and scores for given players using sides data.
func FillPlayerCounter(playerMap map[int]*PlayerWithCounter, sides []gormmodel.Side) {
	for _, side := range sides {
		player, ok := playerMap[side.Pid1]
		if ok {
			player.Score += side.Score
			if side.Score > 0 {
				player.Games++
				player.Win++
			}
			if side.Score < 0 {
				player.Games++
				player.Loss++
			}
		}

		if side.Pid2 == nil {
			// This means the side object is associated with a single match
			continue
		}

		player, ok = playerMap[*side.Pid2]
		if ok {
			player.Score += side.Score
			if side.Score > 0 {
				player.Games++
				player.Win++
			}
			if side.Score < 0 {
				player.Games++
				player.Loss++
			}
		}
	}
}

// FilterActivePlayers returns a slice of pointers to only the ones not in break in the given players slice.
func FilterActivePlayers(players []PlayerWithCounter) []*PlayerWithCounter {
	var keptPlayers []*PlayerWithCounter
	for idx := range players {
		if players[idx].InBreak {
			continue
		}
		keptPlayers = append(keptPlayers, &players[idx])
	}
	return keptPlayers
}

// ToArrangerPlayers converts an array to PlayerWithCounter to a PlayerSlice usable
// by the arranger library.
func ToArrangerPlayers(players []*PlayerWithCounter) model.PlayerSlice {
	var ret model.PlayerSlice
	for _, player := range players {
		ret = append(ret, &model.Player{
			Name:      player.Name,
			ID:        int(player.ID),
			Priority:  player.Priority,
			Score:     player.Score,
			Matches:   float32(player.Games),
			Opponents: make(map[*model.Player]int),
		})
	}
	return ret
}

// ToArrangerPlayersP converts an array to PlayerWithCounter to a PlayerSlice usable
// by the arranger library.
func ToArrangerPlayersP(players []PlayerWithCounter) model.PlayerSlice {
	var ret model.PlayerSlice
	for _, player := range players {
		ret = append(ret, &model.Player{
			Name:      player.Name,
			ID:        int(player.ID),
			Priority:  player.Priority,
			Score:     player.Score,
			Matches:   float32(player.Games),
			Opponents: make(map[*model.Player]int),
		})
	}
	return ret
}

// FillArrangerPlayersOpponents populates the Opponents field for all model.Player objects
// using their match history in the sides data.
func FillArrangerPlayersOpponents(players model.PlayerSlice, sides []gormmodel.Side) {
	playerMap := make(map[int]*model.Player)
	for idx := range players {
		playerMap[players[idx].ID] = players[idx]
	}

	for _, side := range sides {
		if side.Pid2 == nil {
			continue
		}

		// Opponents may be in their breaks and won't exist in the above map.
		// For separation purpose, there is no need to include those players in break.
		if _, ok := playerMap[side.Pid1]; !ok {
			continue
		}
		if _, ok := playerMap[*side.Pid2]; !ok {
			continue
		}

		playerMap[side.Pid1].Opponents[playerMap[*side.Pid2]] =
			playerMap[side.Pid1].Opponents[playerMap[*side.Pid2]] + 1
		playerMap[*side.Pid2].Opponents[playerMap[side.Pid1]] =
			playerMap[*side.Pid2].Opponents[playerMap[side.Pid1]] + 1
	}
}
