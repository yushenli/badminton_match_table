package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RenderRules is the controller for the rules page.
func RenderRules(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "rules.html", gin.H{})
}
