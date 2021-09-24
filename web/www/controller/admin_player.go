package controller

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yushenli/badminton_match_table/web/lib/config"
	"github.com/yushenli/badminton_match_table/web/lib/gormmodel"
	"github.com/yushenli/badminton_match_table/web/lib/util"
	"gorm.io/gorm"
)

// PlayersForm returns the form for adding/updating players for a given event.
// The text form will be pre-filled with existing players in csv format, if any exists.
func PlayersForm(ctx *gin.Context) {
	if config.DB == nil {
		RenderError(ctx, http.StatusInternalServerError, "Unable to connect to database. Please contact the admin.")
		return
	}

	eidStr := ctx.Param("eid")
	eid, err := strconv.Atoi(eidStr)
	if err != nil {
		RenderError(ctx, http.StatusBadRequest,
			fmt.Sprintf("Invalid eid provided: %q", eidStr))
		return
	}

	var event gormmodel.Event
	ret := config.DB.First(&event, eid)
	if ret.Error != nil {
		RenderError(ctx, http.StatusInternalServerError,
			fmt.Sprintf("Failed to locate the event by eid %d: %v", eid, ret.Error))
		return
	}

	if !util.HasAdminPrivilege(ctx, event) {
		RenderError(ctx, http.StatusForbidden,
			fmt.Sprintf("You do not have admin privilege to event %d", eid))
		return
	}

	var players []gormmodel.Player
	ret = config.DB.Where("eid = ?", eid).Find(&players)
	if ret.Error != nil {
		RenderError(ctx, http.StatusInternalServerError,
			fmt.Sprintf("Failed to list players by eid %d: %v", eid, ret.Error))
		return
	}

	var playerEntries [][]string
	for _, player := range players {
		playerEntries = append(playerEntries, []string{
			player.Name,
			fmt.Sprintf("%0.1f", player.Priority),
			fmt.Sprintf("%0.1f", player.InitialScore),
		})
	}

	buffer := new(bytes.Buffer)
	writer := csv.NewWriter(buffer)
	writer.WriteAll(playerEntries)

	ctx.Writer.WriteString(`
<html>
<head>
	<style>
		body {
			font-family: Courier New;
			font-weight: bold;
		}
	</style>
</head>
<body>
	<form id="playersform" method="post">
		<textarea rows="24" cols="50" name="players">` +
		buffer.String() +
		`</textarea>
		<p>
		<input type="submit">
	</form>
</body>
</html>
	`)
}

// PlayersSubmit takes the form for adding/updating players for a given event.
// The input is expected to be a CSV where each line contains the player's name, priority and initial score.
// Whether is player is an existing one or to be added is determined by their name.
func PlayersSubmit(ctx *gin.Context) {
	if config.DB == nil {
		RenderError(ctx, http.StatusInternalServerError, "Unable to connect to database. Please contact the admin.")
		return
	}

	eidStr := ctx.Param("eid")
	eid, err := strconv.Atoi(eidStr)
	if err != nil {
		RenderError(ctx, http.StatusBadRequest,
			fmt.Sprintf("Invalid eid provided: %q", eidStr))
		return
	}

	var event gormmodel.Event
	ret := config.DB.First(&event, eid)
	if ret.Error != nil {
		RenderError(ctx, http.StatusInternalServerError,
			fmt.Sprintf("Failed to locate the event by eid %d: %v", eid, ret.Error))
		return
	}

	if !util.HasAdminPrivilege(ctx, event) {
		RenderError(ctx, http.StatusForbidden,
			fmt.Sprintf("You do not have admin privilege to event %d", eid))
		return
	}

	var players []gormmodel.Player
	playerMap := make(map[string]*gormmodel.Player)
	ret = config.DB.Where("eid = ?", eid).Find(&players)
	if ret.Error != nil {
		RenderError(ctx, http.StatusInternalServerError,
			fmt.Sprintf("Failed to list players by eid %d: %v", eid, ret.Error))
		return
	}
	for idx := range players {
		playerMap[players[idx].Name] = &players[idx]
	}

	reader := csv.NewReader(strings.NewReader(ctx.PostForm("players")))
	playerEntries, err := reader.ReadAll()
	if err != nil {
		RenderError(ctx, http.StatusBadRequest,
			fmt.Sprintf("Unable to parse the players in CSV: %v", err))
		return
	}

	var playersToUpdate []*gormmodel.Player
	var playersToCreate []gormmodel.Player
	for idx, entry := range playerEntries {
		if len(entry) == 0 {
			continue
		}
		if len(entry) != 3 {
			RenderError(ctx, http.StatusBadRequest,
				fmt.Sprintf("Invalid entry on row %d , 3 fields are expected: %+v", idx+1, entry))
			return
		}

		name := strings.TrimSpace(entry[0])
		if name == "" {
			RenderError(ctx, http.StatusBadRequest,
				fmt.Sprintf("Invalid entry on row %d , name cannot be empty: %+v", idx+1, entry))
			return
		}

		priority, err := strconv.ParseFloat(strings.TrimSpace(entry[1]), 32)
		if err != nil {
			RenderError(ctx, http.StatusBadRequest,
				fmt.Sprintf("Invalid entry on row %d , priority needs to be a valid float: %+v", idx+1, entry))
			return
		}

		initialScore, err := strconv.ParseFloat(strings.TrimSpace(entry[2]), 32)
		if err != nil {
			RenderError(ctx, http.StatusBadRequest,
				fmt.Sprintf("Invalid entry on row %d , initial score needs to be a valid float: %+v", idx+1, entry))
			return
		}

		player, ok := playerMap[name]
		if ok {
			player.Priority = float32(priority)
			player.InitialScore = float32(initialScore)
			playersToUpdate = append(playersToUpdate, player)
		} else {
			playersToCreate = append(playersToCreate, gormmodel.Player{
				Name:         name,
				Eid:          eid,
				Priority:     float32(priority),
				InitialScore: float32(initialScore),
				InBreak:      true,
			})
		}

		if len(playersToUpdate) == 0 && len(playersToCreate) == 0 {
			RenderError(ctx, http.StatusBadRequest,
				fmt.Sprintf("No player entries is provided."))
			return
		}
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		for _, player := range playersToUpdate {
			ctx.Writer.WriteString(fmt.Sprintf("To be updated: %+v\n", player))
		}
		if len(playersToUpdate) > 0 {
			ret := tx.Save(&playersToUpdate)
			if ret.Error != nil {
				return ret.Error
			}
		}

		for idx := range playersToCreate {
			ctx.Writer.WriteString(fmt.Sprintf("To be created: %+v\n", playersToCreate[idx]))
		}
		if len(playersToCreate) > 0 {
			tx.Create(&playersToCreate)
			if ret.Error != nil {
				return ret.Error
			}
		}

		return nil
	})
	if err != nil {
		RenderError(ctx, http.StatusInternalServerError,
			fmt.Sprintf("Failed to create/update players: %v", err))
		return
	}
}
