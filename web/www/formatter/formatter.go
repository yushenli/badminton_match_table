package formatter

import (
	"fmt"
	"html/template"

	"github.com/gin-gonic/gin"
)

// RegisterFormatters registers all the defined formatters into the given router.
func RegisterFormatters(router *gin.Engine) {
	router.SetFuncMap(template.FuncMap{
		"asDashedDate":              AsDashedDate,
		"sideInMatchTable":          SideInMatchTable,
		"sideInResults":             SideInResults,
		"sideResult":                SideResult,
		"commaSeparatedPlayerNames": CommaSeparatedPlayerNames,
		"add1": func(n int) string {
			return fmt.Sprintf("%d", n+1)
		},
	})
}
