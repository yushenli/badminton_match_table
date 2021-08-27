package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// RenderError renders error pages including 404 and 500 series.
func RenderError(ctx *gin.Context, errorCode int, message string) {
	ctx.AbortWithStatus(errorCode)
	ctx.Writer.WriteString(fmt.Sprintf("%d %s", errorCode, message))
}
