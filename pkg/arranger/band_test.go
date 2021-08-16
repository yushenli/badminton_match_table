package arranger

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/yushenli/badminton_match_table/pkg/model"
)

func TestClusterByScores(t *testing.T) {
	cases := []struct {
		title    string
		slice    model.PlayerSlice
		expected []float32
	}{
		{
			"EmptySlice",
			[]*model.Player{},
			[]float32{},
		},
		{
			"SinglePlayer",
			[]*model.Player{
				{
					Name:  "Name1",
					Score: 1.23,
				},
			},
			[]float32{1.23},
		},
		{
			"4Players1FarOff",
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
					Score: -8.0,
				},
			},
			[]float32{2.0, -8.0},
		},
		{
			"4PlayersSomeSame",
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
					Score: 2.0,
				},
			},
			[]float32{4.0, 3.0, 2.0},
		},
		{
			"8PlayersBandSizeGreaterThan1",
			[]*model.Player{
				{
					Name:  "Name1",
					Score: 1.0,
				},
				{
					Name:  "Name2a",
					Score: 2.0,
				},
				{
					Name:  "Name2b",
					Score: 2.1,
				},
				{
					Name:  "Name2c",
					Score: 2.0,
				},
				{
					Name:  "Name5",
					Score: 5.0,
				},
				{
					Name:  "Name6a",
					Score: 6.0,
				},
				{
					Name:  "Name6b",
					Score: 6.2,
				},
				{
					Name:  "Name8",
					Score: 8,
				},
			},
			[]float32{8.0, 6.0, 5.0, 1.0},
		},
	}

	for _, tc := range cases {
		t.Run(tc.title, func(t *testing.T) {
			clusters := clusterByScores(tc.slice)
			if !reflect.DeepEqual(tc.expected, clusters) {
				t.Errorf("Unexpected clustering results: expected %v got %v", tc.expected, clusters)
			}
		})
	}
}

func TestFindSeparateRanges(t *testing.T) {
	cases := []struct {
		title    string
		scores   []float32
		expected []separateRange
	}{
		{
			"EmptySlice",
			[]float32{},
			[]separateRange{},
		},
		{
			"NoOneInSameBand",
			[]float32{6.0, 5.0, 4.0, 3.0, 2.0, 1.0},
			[]separateRange{},
		},
		{
			"TwoSameBandInSamePair",
			[]float32{6.0, 5.0, 4.0, 4.0, 2.0, 1.0},
			[]separateRange{},
		},
		{
			"TwoSameBandInAdjacentPair",
			[]float32{6.0, 4.0, 4.0, 3.0, 2.0, 1.0},
			[]separateRange{
				{0, 3, true},
			},
		},
		{
			"SameBandAtHead",
			[]float32{6.0, 6.0, 6.0, 3.0, 2.0, 1.0},
			[]separateRange{
				{0, 3, true},
			},
		},
		{
			"SameBandAtTail",
			[]float32{6.0, 5.0, 4.0, 2.0, 2.0, 2.0},
			[]separateRange{
				{2, 5, false},
			},
		},
		{
			"4In1Band",
			[]float32{4.0, 3.0, 2.0, 2.0, 2.0, 2.0, 1.0, 0},
			[]separateRange{
				{2, 5, false},
			},
		},
		{
			"StartingOnOddIndex",
			[]float32{4.0, 2.0, 2.0, 2.0, 2.0, 2.0, 1.0, 0},
			[]separateRange{
				{0, 5, false},
			},
		},
		{
			"EndingOnEvenIndex",
			[]float32{4.0, 3.0, 2.0, 2.0, 2.0, 1.0, 0, -1.0},
			[]separateRange{
				{2, 5, true},
			},
		},
		{
			"TwoBands",
			// Note (max-min)/6 is greater than 0.5 in this case
			[]float32{4.0, 3.7, 3.4, 1.5, 0.0, -3.0, -3.5, -4.0},
			[]separateRange{
				{0, 3, true},
				{4, 7, false},
			},
		},
		{
			"TwoBandsWithOverlaps",
			// Note (max-min)/6 is greater than 0.5 in this case
			[]float32{4.0, 3.7, 3.4, -1.0, -1.2, -1.4, -1.5, -2.0, -4.0, -4.0},
			[]separateRange{
				{0, 3, true},
				{2, 7, false},
			},
		},
	}

	for _, tc := range cases {
		var players model.PlayerSlice
		for i, score := range tc.scores {
			players = append(players, &model.Player{
				Name:  fmt.Sprintf("Name%d", i),
				Score: score,
			})
		}
		t.Run(tc.title, func(t *testing.T) {
			ranges := findSeparateRanges(players, clusterByScores(players))
			if len(tc.expected) == 0 && len(ranges) == 0 {
				return
			}
			if !reflect.DeepEqual(tc.expected, ranges) {
				t.Errorf("Unexpected separation ranges: expected %v got %v", tc.expected, ranges)
			}
		})
	}
}
