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

func TestSeparateCompetedPlayers(t *testing.T) {
	cases := []struct {
		title       string
		slice       model.PlayerSlice
		matches     []map[int]int
		start       int
		end         int
		endFixed    bool
		expectedErr bool
		expected    []string
	}{
		{
			"4PlayersFirstNonePlayed",
			[]*model.Player{
				{
					Name:  "Name1",
					Score: 4.0,
				},
				{
					Name:  "Name2",
					Score: 3.0,
				},
				{
					Name:  "Name3",
					Score: 2.0,
				},
				{
					Name:  "Name4",
					Score: 1.0,
				},
			},
			[]map[int]int{
				{},
				{},
				{},
				{},
			},
			0,
			3,
			false,
			false,
			[]string{"Name1", "Name2", "Name3", "Name4"},
		},
		{
			"4PlayersFirst2Played",
			[]*model.Player{
				{
					Name:  "Name1",
					Score: 4.0,
				},
				{
					Name:  "Name2",
					Score: 3.0,
				},
				{
					Name:  "Name3",
					Score: 2.0,
				},
				{
					Name:  "Name4",
					Score: 1.0,
				},
			},
			[]map[int]int{
				{1: 1},
				{0: 1},
				{},
				{},
			},
			0,
			3,
			false,
			false,
			[]string{"Name1", "Name3", "Name2", "Name4"},
		},
		{
			"4Players0Vs1And2Played",
			[]*model.Player{
				{
					Name:  "Name1",
					Score: 4.0,
				},
				{
					Name:  "Name2",
					Score: 3.0,
				},
				{
					Name:  "Name3",
					Score: 2.0,
				},
				{
					Name:  "Name4",
					Score: 1.0,
				},
			},
			[]map[int]int{
				{1: 1, 2: 1},
				{0: 1},
				{0: 1},
				{},
			},
			0,
			3,
			false,
			false,
			[]string{"Name1", "Name4", "Name2", "Name3"},
		},
		{
			"4Players0Vs1And2PlayedLastFixed",
			[]*model.Player{
				{
					Name:  "Name1",
					Score: 4.0,
				},
				{
					Name:  "Name2",
					Score: 3.0,
				},
				{
					Name:  "Name3",
					Score: 2.0,
				},
				{
					Name:  "Name4",
					Score: 1.0,
				},
			},
			[]map[int]int{
				{1: 1, 2: 1},
				{0: 1},
				{0: 1},
				{},
			},
			0,
			3,
			true,
			false,
			[]string{"Name1", "Name2", "Name3", "Name4"},
		},
		{
			"4Players0Vs1MoreThan2PlayedLastFixed",
			[]*model.Player{
				{
					Name:  "Name1",
					Score: 4.0,
				},
				{
					Name:  "Name2",
					Score: 3.0,
				},
				{
					Name:  "Name3",
					Score: 2.0,
				},
				{
					Name:  "Name4",
					Score: 1.0,
				},
			},
			[]map[int]int{
				{1: 2, 2: 1},
				{0: 2},
				{0: 1},
				{},
			},
			0,
			3,
			true,
			false,
			[]string{"Name1", "Name3", "Name2", "Name4"},
		},
		{
			"Mid4OutOf6PlayersFirst2Played",
			[]*model.Player{
				{
					Name:  "Name1",
					Score: 6.0,
				},
				{
					Name:  "Name2",
					Score: 5.0,
				},
				{
					Name:  "Name3",
					Score: 4.0,
				},
				{
					Name:  "Name4",
					Score: 3.0,
				},
				{
					Name:  "Name5",
					Score: 2.0,
				},
				{
					Name:  "Name6",
					Score: 1.0,
				},
			},
			[]map[int]int{
				{},
				{2: 1},
				{1: 1},
				{},
				{},
				{},
			},
			1,
			4,
			false,
			false,
			[]string{"Name1", "Name2", "Name4", "Name3", "Name5", "Name6"},
		},
		{
			"Last4OutOf6PlayersFirst2Played",
			[]*model.Player{
				{
					Name:  "Name1",
					Score: 6.0,
				},
				{
					Name:  "Name2",
					Score: 5.0,
				},
				{
					Name:  "Name3",
					Score: 4.0,
				},
				{
					Name:  "Name4",
					Score: 3.0,
				},
				{
					Name:  "Name5",
					Score: 2.0,
				},
				{
					Name:  "Name6",
					Score: 1.0,
				},
			},
			[]map[int]int{
				{},
				{},
				{3: 1},
				{2: 1},
				{},
				{},
			},
			2,
			5,
			false,
			false,
			[]string{"Name1", "Name2", "Name3", "Name5", "Name4", "Name6"},
		},
		{
			"OdderNumberSlice",
			[]*model.Player{
				{
					Name:  "Name1",
					Score: 4.0,
				},
				{
					Name:  "Name2",
					Score: 3.0,
				},
				{
					Name:  "Name3",
					Score: 2.0,
				},
				{
					Name:  "Name4",
					Score: 1.0,
				},
			},
			[]map[int]int{
				{},
				{},
				{},
				{},
			},
			0,
			2,
			false,
			true,
			[]string{},
		},
		{
			"RangeTooBig",
			[]*model.Player{
				{
					Name:  "Name1",
					Score: 4.0,
				},
				{
					Name:  "Name2",
					Score: 3.0,
				},
				{
					Name:  "Name3",
					Score: 2.0,
				},
				{
					Name:  "Name4",
					Score: 1.0,
				},
			},
			[]map[int]int{
				{},
				{},
				{},
				{},
			},
			0,
			4,
			false,
			true,
			[]string{},
		},
	}

	for _, tc := range cases {
		for i := 0; i < len(tc.slice); i++ {
			tc.slice[i].Opponents = make(map[*model.Player]int)
			for key, val := range tc.matches[i] {
				tc.slice[i].Opponents[tc.slice[key]] = val
			}
		}

		t.Run(tc.title, func(t *testing.T) {
			err := SeparateCompetedPlayers(tc.slice, tc.start, tc.end, tc.endFixed)
			if tc.expectedErr {
				if err == nil {
					t.Errorf("Expected error but no error was returned")
				}
				return
			}

			if err != nil {
				t.Errorf("Not expected error but got: %v", err)
				return
			}

			if len(tc.slice) != len(tc.expected) {
				t.Errorf("Different length of PlayerSlice after arranged, expected %d got %d", len(tc.expected), len(tc.slice))
				return
			}

			gotNames := []string{}
			for _, player := range tc.slice {
				gotNames = append(gotNames, player.Name)
			}

			for i := 0; i < len(tc.slice); i++ {
				if gotNames[i] != tc.expected[i] {
					t.Errorf("Different order of players returned, expected %+v got +%v", tc.expected, gotNames)
					return
				}
			}
		})
	}
}
