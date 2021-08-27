package controller

import (
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/yushenli/badminton_match_table/web/lib/config"
	"github.com/yushenli/badminton_match_table/web/lib/gormmodel"
)

type playerWithCounter struct {
	gormmodel.Player
	Games int
	Win   int
	Loss  int
	Score float32
}

func populatePlayers(eid int) ([]playerWithCounter, map[int]*playerWithCounter, error) {
	var players []playerWithCounter
	playerMap := make(map[int]*playerWithCounter)
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

func populateSides(eid int, playerMap map[int]*playerWithCounter) ([]gormmodel.Side, map[int]*gormmodel.Side, error) {
	var sides []gormmodel.Side
	sideMap := make(map[int]*gormmodel.Side)
	ret := config.DB.Where("eid = ?", eid).Find(&sides)
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

func fillPlayerCounter(playerMap map[int]*playerWithCounter, sides []gormmodel.Side) {
	for _, side := range sides {
		player, ok := playerMap[side.Pid1]
		if ok {
			player.Score += side.Score
			player.Games++
			if side.Score > 0 {
				player.Win++
			}
			if side.Score < 0 {
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
			player.Games++
			if side.Score > 0 {
				player.Win++
			}
			if side.Score < 0 {
				player.Loss++
			}
		}
	}
}

func sortPlayerSlice(players []playerWithCounter) {
	sort.Slice(players, func(i, j int) bool {
		p1 := players[i]
		p2 := players[j]

		if p1.InBreak && !p2.InBreak {
			return false
		}
		if !p1.InBreak && p2.InBreak {
			return true
		}

		if p1.Score != p2.Score {
			return p1.Score > p2.Score
		}
		return p1.Priority > p2.Priority
	})
}

// RenderEvent is the controller for the event page.
func RenderEvent(ctx *gin.Context) {
	if config.DB == nil {
		RenderError(ctx, http.StatusInternalServerError, "Unable to connect to database. Please contact the admin.")
		return
	}

	eventKey := ctx.Param("key")
	var event gormmodel.Event
	ret := config.DB.Where("`key` = ?", eventKey).First(&event)
	if ret.Error != nil {
		RenderError(ctx, http.StatusNotFound,
			fmt.Sprintf("Unable to find an event with key %q", eventKey))
		return
	}

	// Fill the Players section
	players, playerMap, err := populatePlayers(int(event.ID))
	if err != nil {
		RenderError(ctx, http.StatusInternalServerError,
			fmt.Sprintf("Failed to list players under event %d", event.ID))
		return
	}

	sides, _, err := populateSides(int(event.ID), playerMap)
	if ret.Error != nil {
		RenderError(ctx, http.StatusInternalServerError,
			fmt.Sprintf("Failed to list sides under event %d", event.ID))
		return
	}

	fillPlayerCounter(playerMap, sides)
	sortPlayerSlice(players)

	ctx.HTML(http.StatusOK, "event.html", gin.H{
		"eventKey": eventKey,
		"players":  players,
	})
}