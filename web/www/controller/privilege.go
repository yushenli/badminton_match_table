package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/yushenli/badminton_match_table/web/lib/util"
)

// SetAdminCookie sets the admin cookie using the string passed in the URL.
// It also shows the current key set in the cookie before changing, or just show the current
// cookie if the cookie field is not given.
func SetAdminCookie(ctx *gin.Context) {
	cookieStr := ctx.Query("cookie")

	currentCookie := util.GetAdminCookie(ctx)

	if cookieStr != "" {
		ctx.SetCookie(util.AdminCookieKey, cookieStr, 86400*7, "/", ctx.Request.Host, false, false)
	}

	if currentCookie == "" {
		ctx.Writer.WriteString("Currently no admin key is set.\n")
	} else {
		ctx.Writer.WriteString(fmt.Sprintf("Current admin key: %s\n", currentCookie))
	}

	if cookieStr != "" {
		ctx.Writer.WriteString(fmt.Sprintf("Set admin key to %s\n", cookieStr))
	}
}
