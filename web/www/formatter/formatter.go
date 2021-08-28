package formatter

import (
	"html/template"

	"github.com/gin-gonic/gin"
)

// RegisterFormatters registers all the defined formatters into the given router.
func RegisterFormatters(router *gin.Engine) {
	router.SetFuncMap(template.FuncMap{
		"asDashedDate":              AsDashedDate,
		"sideInMatchTable":          SideInMatchTable,
		"commaSeparatedPlayerNames": CommaSeparatedPlayerNames,
	})
}
