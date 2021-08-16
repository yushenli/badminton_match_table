package arranger

import (
	"log"
	"math"
	"sort"

	"github.com/pbnjay/clustering"
	"github.com/yushenli/badminton_match_table/pkg/model"
)

func clusterByScores(players model.PlayerSlice) []float32 {
	if len(players) == 0 {
		return []float32{}
	}

	minScore := math.MaxFloat32
	maxScore := -math.MaxFloat32
	distanceMap := make(clustering.DistanceMap)
	for _, player := range players {
		if float64(player.Score) < minScore {
			minScore = float64(player.Score)
		}
		if float64(player.Score) > maxScore {
			maxScore = float64(player.Score)
		}
		distances := make(map[clustering.ClusterItem]float64)
		for _, playerj := range players {
			if player == playerj {
				continue
			}
			distances[playerj] = math.Abs(float64(player.Score) - float64(playerj.Score))
		}
		distanceMap[player] = distances
	}

	// By default we consider scores to be integers, so we make the default max distance within a cluster to be 0.5,
	// so that different scores will be in different clusters at the begginning of a game.
	// As the game goes on, when the difference between maxScore and minScores becomes much larger, we can have
	// more than one integer in one cluster.
	// When minScore and maxScore is far enough apart, we by default consider
	// there are 6 bands between the best players and the most begginners.
	// However, when there are even more (more than 14) users we consider more bands.
	maxDistance := math.Max(0.5, (maxScore-minScore)/math.Max(6, float64(len(players))/2))

	clusters := clustering.NewDistanceMapClusterSet(distanceMap)
	clustering.Cluster(clusters, clustering.Threshold(maxDistance), clustering.CompleteLinkage())

	lowerBonds := make([]float32, clusters.Count())
	clusters.EachCluster(-1, func(cluster int) {
		min := float32(math.MaxFloat32)
		clusters.EachItem(cluster, func(x clustering.ClusterItem) {
			if x.(*model.Player).Score < min {
				min = x.(*model.Player).Score
			}
		})
		lowerBonds[cluster] = min
	})
	sort.Slice(lowerBonds, func(i, j int) bool { return lowerBonds[i] > lowerBonds[j] })

	return lowerBonds
}

type separateRange struct {
	left     int
	right    int
	endFixed bool
}

func findSeparateRanges(players model.PlayerSlice, clusters []float32) []separateRange {
	var ranges []separateRange
	if len(players) == 0 {
		return ranges
	}

	// Find the cluster that the first player's score belongs to.
	clusterIdx := 0
	for ; clusters[clusterIdx] > players[0].Score; clusterIdx++ {
	}

	i := 0
	j := 1
	for j <= len(players) {
		if j < len(players) && players[j].Score >= clusters[clusterIdx] {
			// This means the scores from player i to j are still in the same band.
			j++
			continue
		}

		if (j-1)/2 > i/2 {
			// Player j-1 is the last player in the same band as player i
			// Here we check if player i and j-1 belong to the same 2-sized pair,
			// if so, there is no need to rearrange them.

			// For the ranges to be rearranged, always start with an even number-th player
			// and end with an odd number-th player. If the end of the range is on an even
			// indexed player, do the separation process including the very next element in the
			// next band, but set endFixed so that that player won't be moved.
			var r separateRange
			if i%2 == 0 {
				r.left = i
			} else {
				r.left = i - 1
			}
			if (j-1)%2 == 1 {
				r.right = j - 1
				r.endFixed = false
			} else {
				r.right = j
				r.endFixed = true
			}
			ranges = append(ranges, r)
		}

		i = j
		j++
		clusterIdx++
	}

	return ranges
}

// SeparateCompetedPlayersWithinBands scans a given sorted player list and break them into bands.
// For each band, it will call SeparateCompetedPlayers to rearrange (in place) the players within the band,
// so that people with similar scores will be shuffled in a way that everyone will play with an opponent
// with whom they have played as few times as possible.
//
// For a given list of players, the bands are determined using hierarchical clustering. The max distance
// within clusters is the difference between the highest score and lowest score divided by
// Max(6, Player Count / 2), or 0.5, whichever is bigger.
func SeparateCompetedPlayersWithinBands(allPlayers, playingPlayers model.PlayerSlice) error {
	clusters := clusterByScores(allPlayers)
	ranges := findSeparateRanges(playingPlayers, clusters)
	for _, r := range ranges {
		err := SeparateCompetedPlayers(playingPlayers, r.left, r.right, r.endFixed)
		if err != nil {
			log.Printf("Failed to separate competed players within bands among allPlayers %v and playingPlayers %v : %v", allPlayers, playingPlayers, err)
		}
	}

	return nil
}
