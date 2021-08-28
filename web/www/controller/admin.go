package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yushenli/badminton_match_table/web/lib/config"
	"github.com/yushenli/badminton_match_table/web/lib/gormmodel"
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
