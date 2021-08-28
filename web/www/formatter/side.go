package formatter

import (
	"html/template"
	"strings"

	"github.com/yushenli/badminton_match_table/web/lib/gormmodel"
)

// SideInMatchTable returns the names of a given side in mutiple lines.
func SideInMatchTable(side *gormmodel.Side) template.HTML {
	var sb strings.Builder
	sb.WriteString(side.Player1.Name)
	if side.Player2 != nil {
		sb.WriteString("<br><br>")
		sb.WriteString(side.Player2.Name)
	}
	// You can't just return the string here, otherwise the HTML tags will be escaped
	return template.HTML(sb.String())
}

// SideInResults returns the names of a given side in a single line.
func SideInResults(side *gormmodel.Side) template.HTML {
	var sb strings.Builder
	sb.WriteString(side.Player1.Name)
	if side.Player2 != nil {
		sb.WriteString(" / ")
		sb.WriteString(side.Player2.Name)
	}
	// You can't just return the string here, otherwise the HTML tags will be escaped
	return template.HTML(sb.String())
}

// SideResult returns "LOST" or "WON" for a given side.
func SideResult(side *gormmodel.Side) string {
	if side.Score > 0 {
		return "WON"
	}
	if side.Score < 0 {
		return "LOST"
	}
	return "N/A"
}
