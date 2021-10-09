package util

import (
	"github.com/gin-gonic/gin"
	"github.com/yushenli/badminton_match_table/web/lib/gormmodel"
)

const (
	// AdminCookieKey is the name of the cookie field that stores the admin key
	AdminCookieKey = "admin"
)

// HasAdminPrivilege returns if the current visitor has admin privilege to the given event.
func HasAdminPrivilege(ctx *gin.Context, event gormmodel.Event) bool {
	currentCookie, err := ctx.Request.Cookie(AdminCookieKey)
	if currentCookie == nil || err != nil {
		// No admin cookie set, no admin privilege.
		return false
	}

	return event.AdminKey != "" && currentCookie.Value == event.AdminKey
}

// GetAdminCookie returns the value of the admin cookie
func GetAdminCookie(ctx *gin.Context) string {
	currentCookie, err := ctx.Request.Cookie(AdminCookieKey)
	adminCookie := ""
	if currentCookie != nil && err == nil {
		adminCookie = currentCookie.Value
	}

	return adminCookie
}
