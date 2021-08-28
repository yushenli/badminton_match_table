package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yushenli/badminton_match_table/web/lib/config"
	"github.com/yushenli/badminton_match_table/web/lib/gormmodel"
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
