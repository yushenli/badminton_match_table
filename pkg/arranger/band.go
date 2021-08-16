package arranger

import (
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

	upperBonds := make([]float32, clusters.Count())
	clusters.EachCluster(-1, func(cluster int) {
		max := float32(-math.MaxFloat32)
		clusters.EachItem(cluster, func(x clustering.ClusterItem) {
			if x.(*model.Player).Score > max {
				max = x.(*model.Player).Score
			}
			//log.Println(cluster, x)
		})
		upperBonds[cluster] = max
	})
	sort.Slice(upperBonds, func(i, j int) bool { return upperBonds[i] < upperBonds[j] })

	return upperBonds
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

	return nil
}
