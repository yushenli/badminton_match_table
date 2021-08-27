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

func fillPlayerCounter(eid int, playerMap map[int]*playerWithCounter) error {
	var sides []gormmodel.Side
	ret := config.DB.Where("eid = ?", eid).Find(&sides)
	if ret.Error != nil {
		return ret.Error
	}

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

	return nil
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

	var players []playerWithCounter
	playerMap := make(map[int]*playerWithCounter)
	ret = config.DB.Where("eid = ?", event.ID).Find(&players)
	if ret.Error != nil {
		log.Printf("Failed to list players under event %d: %v", event.ID, ret.Error)
		RenderError(ctx, http.StatusInternalServerError,
			fmt.Sprintf("Failed to list players under event %d", event.ID))
		return
	}
	for idx, _ := range players {
		playerMap[int(players[idx].ID)] = &(players[idx])
		players[idx].Score = players[idx].InitialScore
	}

	err := fillPlayerCounter(int(event.ID), playerMap)
	if err != nil {
		log.Printf("Failed to fill player score and counters for event %d: %v", event.ID, ret.Error)
		RenderError(ctx, http.StatusInternalServerError,
			fmt.Sprintf("Failed to fill player score and counters for event %d", event.ID))
		return
	}

	sortPlayerSlice(players)

	ctx.HTML(http.StatusOK, "event.html", gin.H{
		"eventsKey": eventKey,
		"players":   players,
	})
}
