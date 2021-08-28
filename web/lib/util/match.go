package util

import (
	"github.com/yushenli/badminton_match_table/pkg/model"
	"github.com/yushenli/badminton_match_table/web/lib/gormmodel"
)

// FromArrangerMatchArrangement converts a MatchArrangement provided by the arrenger
// into Match objects under gormmodel. The Sides in Match objects and Players in Side
// objects will be filled.
func FromArrangerMatchArrangement(arrangement model.MatchArrangement, event gormmodel.Event) []gormmodel.Match {
	matches := make([]gormmodel.Match, len(arrangement))

	for idx, arrangerMatch := range arrangement {
		matches[idx] = gormmodel.Match{
			Eid:    int(event.ID),
			Round:  event.CurrentRound,
			Court:  idx + 1,
			Status: gormmodel.PLAYING,
			Side1: &gormmodel.Side{
				Eid:  int(event.ID),
				Pid1: arrangerMatch.Side1.Player1.ID,
			},
			Side2: &gormmodel.Side{
				Eid:  int(event.ID),
				Pid1: arrangerMatch.Side2.Player1.ID,
			},
		}

		if arrangerMatch.Side1.Player2 != nil {
			matches[idx].Side1.Pid2 = &arrangerMatch.Side1.Player2.ID
			matches[idx].Side2.Pid2 = &arrangerMatch.Side2.Player2.ID
		}
	}

	return matches
}
