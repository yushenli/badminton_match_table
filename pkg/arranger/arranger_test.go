package arranger

import (
	"testing"

	"github.com/yushenli/badminton_match_table/pkg/model"
)

func TestSortPlayerSliceByScorePriority(t *testing.T) {
	cases := []struct {
		title    string
		slice    model.PlayerSlice
		expected []string
	}{
		{
			"DifferentScorePriority",
			[]*model.Player{
				{
					Name:     "Name1",
					Score:    2.0,
					Priority: 1.0,
				},
				{
					Name:     "Name2",
					Score:    2.0,
					Priority: 2.0,
				},
				{
					Name:     "Name3",
					Score:    3.0,
					Priority: 0.5,
				},
				{
					Name:     "Name4",
					Score:    1.0,
					Priority: 4.0,
				},
			},
			[]string{"Name3", "Name2", "Name1", "Name4"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.title, func(t *testing.T) {
			sortPlayerSliceByScorePriority(tc.slice)
			if len(tc.expected) != len(tc.slice) {
				t.Errorf("Unpected length of the sorted PlayerSlice, expected %v, got %v",
					len(tc.expected),
					len(tc.slice))
			}
			for i, expectedName := range tc.expected {
				if tc.slice[i].Name != expectedName {
					t.Errorf("Unpected %d-th player in the sorted PlayerSlice, expected %q, got %q",
						i,
						expectedName,
						tc.slice[i].Name)
				}
			}
		})
	}
}
