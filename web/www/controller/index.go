package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yushenli/badminton_match_table/web/lib/config"
	"github.com/yushenli/badminton_match_table/web/lib/gormmodel"
)

// RenderIndex is the controller for the root page.
func RenderIndex(ctx *gin.Context) {
	var events []gormmodel.Event
	if config.DB != nil {
		config.DB.Order("date desc").Limit(20).Find(&events)
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"events": events,
	})
}
