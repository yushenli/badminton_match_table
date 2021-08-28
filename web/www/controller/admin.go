package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yushenli/badminton_match_table/pkg/arranger"
	"github.com/yushenli/badminton_match_table/web/lib/config"
	"github.com/yushenli/badminton_match_table/web/lib/gormmodel"
	"github.com/yushenli/badminton_match_table/web/lib/util"
	"gorm.io/gorm"
)

// ChangeMatchStatus handles the reuqest to change the status of a match record.
func ChangeMatchStatus(ctx *gin.Context) {
	midStr := ctx.Query("mid")
	sideStr := ctx.Query("side")
	if midStr == "" || sideStr == "" {
		RenderError(ctx, http.StatusBadRequest,
			fmt.Sprintf("You must provide mid and side parameters: %s", ctx.Request.URL.String()))
		return
	}

	mid, err := strconv.Atoi(midStr)
	if err != nil {
		RenderError(ctx, http.StatusBadRequest,
			fmt.Sprintf("Invalid mid provided: %q", midStr))
		return
	}
	if sideStr != "1" && sideStr != "0" && sideStr != "2" {
		RenderError(ctx, http.StatusBadRequest,
			fmt.Sprintf("Invalid side provided: %q", sideStr))
		return
	}

	log.Println(fmt.Sprintf("ChangeMatchStatus called: mid=%d, side=%s", mid, sideStr))

	if config.DB == nil {
		RenderError(ctx, http.StatusInternalServerError, "Unable to connect to database. Please contact the admin.")
		return
	}

	var match gormmodel.Match
	ret := config.DB.First(&match, mid)
	if ret.Error != nil {
		RenderError(ctx, http.StatusBadRequest,
			fmt.Sprintf("Failed to locate the match by mid %d: %v", mid, ret.Error))
	}

	switch sideStr {
	case "1":
		match.Status = gormmodel.SIDE1WON
	case "0":
		match.Status = gormmodel.PLAYING
	case "2":
		match.Status = gormmodel.SIDE2WON
	}
	ret = config.DB.Save(&match)
	if ret.Error != nil {
		RenderError(ctx, http.StatusBadRequest,
			fmt.Sprintf("Failed to update the match by mid %d to status %s: %v", mid, sideStr, ret.Error))
	}
}

// ChangeBreakStatus handles the reuqest to change whether a player is in a break.
func ChangeBreakStatus(ctx *gin.Context) {
	pidStr := ctx.Query("pid")
	inBreakStr := ctx.Query("in_break")
	if pidStr == "" || inBreakStr == "" {
		RenderError(ctx, http.StatusBadRequest,
			fmt.Sprintf("You must provide pid and in_break parameters: %s", ctx.Request.URL.String()))
		return
	}

	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		RenderError(ctx, http.StatusBadRequest,
			fmt.Sprintf("Invalid pid provided: %q", pidStr))
		return
	}
	if inBreakStr != "1" && inBreakStr != "0" {
		RenderError(ctx, http.StatusBadRequest,
			fmt.Sprintf("Invalid in_break provided: %q", inBreakStr))
		return
	}

	log.Println(fmt.Sprintf("ChangeBreakStatus called: pid=%d, in_break=%s", pid, inBreakStr))

	if config.DB == nil {
		RenderError(ctx, http.StatusInternalServerError, "Unable to connect to database. Please contact the admin.")
		return
	}

	var player gormmodel.Player
	ret := config.DB.First(&player, pid)
	if ret.Error != nil {
		RenderError(ctx, http.StatusBadRequest,
			fmt.Sprintf("Failed to locate the player by pid %d: %v", pid, ret.Error))
	}

	if inBreakStr == "1" {
		player.InBreak = true
	} else {
		player.InBreak = false
	}
	ret = config.DB.Save(&player)
	if ret.Error != nil {
		RenderError(ctx, http.StatusBadRequest,
			fmt.Sprintf("Failed to update the player by pid %d: %v", pid, ret.Error))
	}
}

// CompleteRound submits the scroes for each side involved in matches in the given round,
// based on the status of each match record.
// In the end it increment the current round field in the event record by one
// if the current round == the given round.
func CompleteRound(ctx *gin.Context) {
	eidStr := ctx.Query("eid")
	roundStr := ctx.Query("round")
	if eidStr == "" || roundStr == "" {
		RenderError(ctx, http.StatusBadRequest,
			fmt.Sprintf("You must provide eid and round parameters: %s", ctx.Request.URL.String()))
		return
	}

	eid, err := strconv.Atoi(eidStr)
	if err != nil {
		RenderError(ctx, http.StatusBadRequest,
			fmt.Sprintf("Invalid eid provided: %q", eidStr))
		return
	}
	round, err := strconv.Atoi(roundStr)
	if err != nil {
		RenderError(ctx, http.StatusBadRequest,
			fmt.Sprintf("Invalid round provided: %q", roundStr))
		return
	}

	if config.DB == nil {
		RenderError(ctx, http.StatusInternalServerError, "Unable to connect to database. Please contact the admin.")
		return
	}

	var event gormmodel.Event
	ret := config.DB.First(&event, eid)
	if ret.Error != nil {
		RenderError(ctx, http.StatusInternalServerError,
			fmt.Sprintf("Failed to locate the event by eid %d: %v", eid, ret.Error))
	}

	var matches []gormmodel.Match
	ret = config.DB.Where("eid = ?", eid).Where("round = ?", round).Find(&matches)
	if ret.Error != nil {
		log.Printf("Failed to list matches under event %d in round %d: %v", eid, round, ret.Error)
		RenderError(ctx, http.StatusInternalServerError, "Failed to list matches under event")
		return
	}

	for _, match := range matches {
		if match.Status == gormmodel.PLAYING {
			RenderError(ctx, http.StatusBadRequest,
				fmt.Sprintf("There are still match(es) with PLAYING status: %+v", match))
			return
		}
	}

	var output strings.Builder
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		for _, match := range matches {
			sidWon := match.Sid1
			sidLost := match.Sid2
			if match.Status == gormmodel.SIDE2WON {
				sidWon = match.Sid2
				sidLost = match.Sid1
			}

			var sideWon gormmodel.Side
			var sideLost gormmodel.Side
			ret = tx.First(&sideWon, sidWon)
			if ret.Error != nil {
				return ret.Error
			}
			ret = tx.First(&sideLost, sidLost)
			if ret.Error != nil {
				return ret.Error
			}

			sideWon.Score = 1.0
			output.WriteString(fmt.Sprintf("Setting side %d score to 1.0<br>\n", sideWon.ID))
			sideLost.Score = -1.0
			output.WriteString(fmt.Sprintf("Setting side %d score to -1.0<br>\n", sideLost.ID))

			ret = tx.Save(&sideWon)
			if ret.Error != nil {
				return ret.Error
			}
			ret = tx.Save(&sideLost)
			if ret.Error != nil {
				return ret.Error
			}
		}

		if event.CurrentRound == round {
			event.CurrentRound++
			output.WriteString(fmt.Sprintf("Increased event %d current round to %d<br>\n", event.ID, event.CurrentRound))
			ret = tx.Save(&event)
			if ret.Error != nil {
				return ret.Error
			}
		}

		return nil
	})
	if err != nil {
		log.Printf("Failed when modifying sides and/or event %d in round %d: %v", eid, round, ret.Error)
		RenderError(ctx, http.StatusInternalServerError, "Failed when modifying sides and/or event")
		return
	}

	// The output must wait until no error will be thrown, since errors are thrown
	// under different HTTP status codes.
	ctx.Writer.WriteString(output.String())
}

// ScheduleCurrentRound generates a new match table for the current round.
// Existing match tables for the same round will be deleted without asking.
func ScheduleCurrentRound(ctx *gin.Context) {
	eidStr := ctx.Query("eid")
	if eidStr == "" {
		RenderError(ctx, http.StatusBadRequest,
			fmt.Sprintf("You must provide the eid parameter: %s", ctx.Request.URL.String()))
		return
	}

	eid, err := strconv.Atoi(eidStr)
	if err != nil {
		RenderError(ctx, http.StatusBadRequest,
			fmt.Sprintf("Invalid eid provided: %q", eidStr))
		return
	}

	if config.DB == nil {
		RenderError(ctx, http.StatusInternalServerError, "Unable to connect to database. Please contact the admin.")
		return
	}

	var event gormmodel.Event
	ret := config.DB.First(&event, eid)
	if ret.Error != nil {
		RenderError(ctx, http.StatusInternalServerError,
			fmt.Sprintf("Failed to locate the event by eid %d: %v", eid, ret.Error))
	}

	players, playerMap, err := util.PopulatePlayers(eid)
	if err != nil {
		RenderError(ctx, http.StatusInternalServerError,
			fmt.Sprintf("Failed to list players under event %d", eid))
		return
	}

	sides, _, err := util.PopulateSides(int(event.ID), playerMap)
	if ret.Error != nil {
		RenderError(ctx, http.StatusInternalServerError,
			fmt.Sprintf("Failed to list sides under event %d", eid))
		return
	}

	util.FillPlayerCounter(playerMap, sides)

	activePlayers := util.FilterActivePlayers(players)
	allArrangerPlayers := util.ToArrangerPlayersP(players)
	activeArrangerPlayers := util.ToArrangerPlayers(activePlayers)
	util.FillArrangerPlayersOpponents(activeArrangerPlayers, sides)

	ctx.Writer.WriteString("Active players with opponents filled:\n")
	for idx := range activeArrangerPlayers {
		ctx.Writer.WriteString(fmt.Sprintf("%+v", activeArrangerPlayers[idx]))
		ctx.Writer.WriteString("\n")
	}

	playingPlayers, err := arranger.PickPlayersForCourts(activeArrangerPlayers, event.Courts)
	if err != nil {
		RenderError(ctx, http.StatusInternalServerError,
			fmt.Sprintf("Error when picking players based on number of courts %d: %v", event.Courts, err))
		return
	}

	arranger.SortPlayerSliceByScorePriority(playingPlayers)
	err = arranger.SeparateCompetedPlayersWithinBands(allArrangerPlayers, playingPlayers)
	if err != nil {
		RenderError(ctx, http.StatusInternalServerError,
			fmt.Sprintf("Error when separating competed players within bands %v", err))
		return
	}

	ctx.Writer.WriteString("\n\nPlayers clustered by score and separated between competed:\n")
	for idx := range playingPlayers {
		ctx.Writer.WriteString(fmt.Sprintf("%+v", playingPlayers[idx]))
		ctx.Writer.WriteString("\n")
	}

	arrangerMatches, err := arranger.MakeMatchArrangements(playingPlayers, event.Courts, event.CurrentRound)
	if err != nil {
		RenderError(ctx, http.StatusInternalServerError,
			fmt.Sprintf("Error when making match arrangement based on playing players: %v", err))
		return
	}
	matches := util.FromArrangerMatchArrangement(arrangerMatches, event)
	ctx.Writer.WriteString("\n\nMatch arrangement:\n")
	for idx := range matches {
		ctx.Writer.WriteString(fmt.Sprintf("%+v\n", matches[idx]))
		ctx.Writer.WriteString(fmt.Sprintf("    %+v\n", matches[idx].Side1))
		ctx.Writer.WriteString(fmt.Sprintf("    %+v", matches[idx].Side2))
		ctx.Writer.WriteString("\n")
	}

	if ctx.Query("proceed") != "1" {
		return
	}

	// Clean up existing arrangements
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		var oldMatches []gormmodel.Match
		ret := tx.Where("eid = ?", event.ID).Where("round = ?", event.CurrentRound).Find(&oldMatches)
		if ret.Error != nil {
			return ret.Error
		}

		for idx := range oldMatches {
			ret = tx.Delete(&gormmodel.Side{}, oldMatches[idx].Sid1)
			if ret.Error != nil {
				return ret.Error
			}
			ret = tx.Delete(&gormmodel.Side{}, oldMatches[idx].Sid2)
			if ret.Error != nil {
				return ret.Error
			}
			ret = tx.Delete(&oldMatches[idx])
			if ret.Error != nil {
				return ret.Error
			}
		}

		return nil
	})
	if err != nil {
		log.Printf("Failed when cleaning up old matches/sides %d in round %d: %v", eid, event.CurrentRound, err)
		RenderError(ctx, http.StatusInternalServerError, "Failed when creating new matches/sides")
		return
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		for idx := range matches {
			ret = tx.Create(&matches[idx])
			if ret.Error != nil {
				return ret.Error
			}

			matches[idx].Side1.Mid = int(matches[idx].ID)
			ret = tx.Save(matches[idx].Side1)
			if ret.Error != nil {
				return ret.Error
			}
			matches[idx].Side2.Mid = int(matches[idx].ID)
			ret = tx.Save(matches[idx].Side2)
			if ret.Error != nil {
				return ret.Error
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("Failed when creating new matches/sides %d in round %d: %v", eid, event.CurrentRound, err)
		RenderError(ctx, http.StatusInternalServerError, "Failed when creating new matches/sides")
		return
	}

	ctx.Writer.WriteString("\n\nMatch arrangement persisted:\n")
	for idx := range matches {
		ctx.Writer.WriteString(fmt.Sprintf("%+v\n", matches[idx]))
		ctx.Writer.WriteString(fmt.Sprintf("    %+v\n", matches[idx].Side1))
		ctx.Writer.WriteString(fmt.Sprintf("    %+v", matches[idx].Side2))
		ctx.Writer.WriteString("\n")
	}
}
