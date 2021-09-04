package controller

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yushenli/badminton_match_table/web/lib/config"
	"github.com/yushenli/badminton_match_table/web/lib/gormmodel"
	"github.com/yushenli/badminton_match_table/web/lib/util"
)

// matchTableColStyle returns the css style for the match table columns so that
// the width adapts the number of courts.
func matchTableColStyle(courts int) string {
	if courts == 1 {
		return "w-col-12" // width: 100%
	}
	if courts == 2 {
		return "w-col-6" // width: 50%
	}
	if courts == 3 {
		return "w-col-4" // width: 33%
	}
	return "w-col-3" // width: 25%
}

func sortPlayerSlice(players []*util.PlayerWithCounter) {
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

func findUnscheduledPlayers(matches []*gormmodel.Match, players []util.PlayerWithCounter) []*gormmodel.Player {
	scheduled := make(map[int]bool)
	for _, match := range matches {
		scheduled[match.Side1.Pid1] = true
		if match.Side1.Pid2 != nil {
			scheduled[*match.Side1.Pid2] = true
		}
		scheduled[match.Side2.Pid1] = true
		if match.Side2.Pid2 != nil {
			scheduled[*match.Side2.Pid2] = true
		}
	}

	var unscheduled []*gormmodel.Player
	for idx, player := range players {
		if player.InBreak {
			continue
		}
		if _, ok := scheduled[int(player.ID)]; ok {
			continue
		}
		unscheduled = append(unscheduled, &players[idx].Player)
	}

	return unscheduled
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
	players, playerMap, err := util.PopulatePlayers(int(event.ID))
	if err != nil {
		RenderError(ctx, http.StatusInternalServerError,
			fmt.Sprintf("Failed to list players under event %d", event.ID))
		return
	}

	sides, sideMap, err := util.PopulateSides(int(event.ID), playerMap, nil)
	if ret.Error != nil {
		RenderError(ctx, http.StatusInternalServerError,
			fmt.Sprintf("Failed to list sides under event %d", event.ID))
		return
	}

	util.FillPlayerCounter(playerMap, sides)
	sortedPlayers := make([]*util.PlayerWithCounter, len(players))
	for idx := range players {
		sortedPlayers[idx] = &players[idx]
	}
	sortPlayerSlice(sortedPlayers)

	// Fill the current round match table and match results
	_, matchesByRound, err := util.PopulateMatches(int(event.ID), event.CurrentRound, sideMap)

	round := event.CurrentRound
	if ctx.Query("round") != "" {
		round1, err := strconv.Atoi(ctx.Query("round"))
		if err == nil {
			round = round1
		}
	}

	unscheduledPlayers := findUnscheduledPlayers(matchesByRound[round-1], players)

	ctx.HTML(http.StatusOK, "event.html", gin.H{
		"event":              event,
		"players":            sortedPlayers,
		"displayRound":       round,
		"currentMatches":     matchesByRound[round-1],
		"matchesByRound":     matchesByRound,
		"matchTableColStyle": matchTableColStyle(len(matchesByRound[round-1])),
		"unscheduledPlayers": unscheduledPlayers,
		"hasAdminPrivilege":  util.HasAdminPrivilege(ctx, event),
	})
}

// RedirctToToday sends HTTP relocation header to the event who has a key of today in the
// form of YYYYMMDD. If such an event does not exist, it will redirect to the match list
// on the home page.
func RedirctToToday(ctx *gin.Context) {
	if config.DB == nil {
		RenderError(ctx, http.StatusInternalServerError, "Unable to connect to database. Please contact the admin.")
		return
	}

	eventKey := time.Now().Format("20060102")
	url := fmt.Sprintf("/event/%s", eventKey)

	var event *gormmodel.Event
	ret := config.DB.Where("`key` = ?", eventKey).First(&event)
	if ret.Error != nil || event == nil {
		url = "/#matches"
	}

	ctx.Redirect(http.StatusTemporaryRedirect, url)
}
