package arranger

import (
	"sort"

	"github.com/yushenli/badminton_match_table/pkg/model"
)

func sortPlayerSliceByScorePriority(slice model.PlayerSlice) {
	sort.Slice(slice, func(i, j int) bool {
		a := slice[i]
		b := slice[j]
		return a.Score > b.Score || (a.Score == b.Score && a.Priority > b.Priority)
	})
}
