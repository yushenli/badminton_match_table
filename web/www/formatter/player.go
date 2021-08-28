package formatter

import (
	"strings"

	"github.com/yushenli/badminton_match_table/web/lib/gormmodel"
)

// CommaSeparatedPlayerNames returns a comman separated string of all the passed in players
func CommaSeparatedPlayerNames(players []*gormmodel.Player) string {
	var sb strings.Builder
	for idx := range players {
		if idx > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(players[idx].Name)
	}
	return sb.String()
}
