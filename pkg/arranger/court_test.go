package arranger

import (
	"reflect"
	"testing"

	"github.com/yushenli/badminton_match_table/pkg/model"
)

func TestPickPlayersForCourts(t *testing.T) {
	cases := []struct {
		title       string
		slice       model.PlayerSlice
		courtCount  int
		expected    []string
		expectedErr bool
	}{
		{
			"0Players0Court",
			[]*model.Player{},
			0,
			[]string{},
			false,
		},
		{
			"1Players1Court",
			[]*model.Player{
				{
					Name:     "Name1",
					Matches:  1.0,
					Priority: 1.0,
				},
			},
			1,
			nil,
			true,
		},
		{
			"3Players2Court",
			[]*model.Player{
				{
					Name:     "Name1",
					Matches:  1.0,
					Priority: 1.0,
				},
				{
					Name:     "Name2",
					Matches:  2.0,
					Priority: 1.0,
				},
				{
					Name:     "Name4",
					Matches:  4.0,
					Priority: 1.0,
				},
			},
			2,
			nil,
			true,
		},
		{
			"2Players1Court",
			[]*model.Player{
				{
					Name:     "Name1",
					Matches:  2.0,
					Priority: 1.0,
				},
				{
					Name:     "Name2",
					Matches:  2.0,
					Priority: 2.0,
				},
			},
			1,
			[]string{"Name2", "Name1"},
			false,
		},
		{
			"4Players1Court",
			[]*model.Player{
				{
					Name:     "Name1",
					Matches:  1.0,
					Priority: 1.0,
				},
				{
					Name:     "Name2",
					Matches:  2.0,
					Priority: 1.0,
				},
				{
					Name:     "Name4",
					Matches:  4.0,
					Priority: 1.0,
				},
				{
					Name:     "Name3",
					Matches:  3.0,
					Priority: 1.0,
				},
			},
			1,
			[]string{"Name1", "Name2", "Name3", "Name4"},
			false,
		},
		{
			"3Players1Court",
			[]*model.Player{
				{
					Name:     "Name1",
					Matches:  1.0,
					Priority: 1.0,
				},
				{
					Name:     "Name2",
					Matches:  2.0,
					Priority: 1.0,
				},
				{
					Name:     "Name3",
					Matches:  3.0,
					Priority: 1.0,
				},
			},
			1,
			[]string{"Name1", "Name2"},
			false,
		},
		{
			"6Players1Court",
			[]*model.Player{
				{
					Name:     "Name5",
					Matches:  5.0,
					Priority: 5.0,
				},
				{
					Name:     "Name1",
					Matches:  1.0,
					Priority: 1.0,
				},
				{
					Name:     "Name2",
					Matches:  2.0,
					Priority: 1.0,
				},
				{
					Name:     "Name4",
					Matches:  4.0,
					Priority: 1.0,
				},
				{
					Name:     "Name3",
					Matches:  3.0,
					Priority: 1.0,
				},
				{
					Name:     "Name6",
					Matches:  6.0,
					Priority: 5.0,
				},
			},
			1,
			[]string{"Name1", "Name2", "Name3", "Name4"},
			false,
		},
		{
			"7Players2Court",
			[]*model.Player{
				{
					Name:     "Name6b",
					Matches:  6.0,
					Priority: 2.0,
				},
				{
					Name:     "Name1",
					Matches:  1.0,
					Priority: 1.0,
				},
				{
					Name:     "Name2",
					Matches:  2.0,
					Priority: 1.0,
				},
				{
					Name:     "Name4",
					Matches:  4.0,
					Priority: 1.0,
				},
				{
					Name:     "Name3",
					Matches:  3.0,
					Priority: 1.0,
				},
				{
					Name:     "Name6a",
					Matches:  6.0,
					Priority: 1.0,
				},
				{
					Name:     "Name5",
					Matches:  5.0,
					Priority: 5.0,
				},
			},
			2,
			[]string{"Name1", "Name2", "Name3", "Name4", "Name5", "Name6b"},
			false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.title, func(t *testing.T) {
			players, err := PickPlayersForCourts(tc.slice, tc.courtCount)
			if tc.expectedErr {
				if err == nil {
					t.Errorf("Expected error but got nil.")
				}
				return
			}
			if !tc.expectedErr && err != nil {
				t.Errorf("Not expecting error but got %v.", err)
				return
			}

			playerNames := make([]string, len(players))
			for idx, player := range players {
				playerNames[idx] = player.Name
			}
			if !reflect.DeepEqual(tc.expected, playerNames) {
				t.Errorf("Unexpected player lists that can play, expected %v got %v", tc.expected, playerNames)
			}
		})
	}
}
