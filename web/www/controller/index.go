package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/yushenli/badminton_match_table/web/lib/config"
	"github.com/yushenli/badminton_match_table/web/lib/gormmodel"
)

// RenderIndex is the controller for the root page.
func RenderIndex(ctx *gin.Context) {
	var events []gormmodel.Event
	config.DB.Order("date desc").Limit(20).Find(&events)

	for _, event := range events {
		ctx.Writer.WriteString(fmt.Sprintf("Date: %v, Key: %s, Location: %q <br>\n", event.Date, event.Key, event.Location))
	}
}
