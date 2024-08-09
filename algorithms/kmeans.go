package algorithms

import (
	"fmt"
	"math/rand"
)

type KMeans struct {
	iterations int
}

func NewKMeans(iterations int) KMeans {
	return KMeans{
		iterations: iterations,
	}
}

func (m KMeans) Partition(loads []Load, k int) (Clusters, error) {
	if k > len(loads) {
		return nil, fmt.Errorf("cannot partition %d loads into %d clusters", len(loads), k)
	}

	clusters, err := NewClusters(loads, k)
	if err != nil {
		return nil, err
	}

	loadMidpoints := make([]int, len(loads))
	changes := 1

	for i := 0; changes > 0; i++ {
		changes = 0
		clusters.Reset()

		for idx, load := range loads {
			nearest := clusters.ClosestTo(load)
			clusters[nearest].Loads = append(clusters[nearest].Loads, load)

			if loadMidpoints[idx] != nearest {
				loadMidpoints[idx] = nearest
				changes++
			}
		}

		// Ensure no clusters are empty
		for j := 0; j < len(clusters); j++ {
			if len(clusters[j].Loads) == 0 {
				var ri int
				for {
					ri = rand.Intn(len(loads))
					if len(clusters[loadMidpoints[ri]].Loads) > 1 {
						break
					}
				}
				clusters[j].Loads = append(clusters[j].Loads, loads[ri])
				clusters[loadMidpoints[ri]].Loads = removeLoad(clusters[loadMidpoints[ri]].Loads, loads[ri])
				loadMidpoints[ri] = j
				changes = len(loads)
			}
		}

		if changes > 0 {
			clusters.Recenter()
		}

		if i == m.iterations {
			break
		}
	}

	return clusters, nil
}

func removeLoad(loads []Load, loadToRemove Load) []Load {
	for i, load := range loads {
		if load == loadToRemove {
			return append(loads[:i], loads[i+1:]...)
		}
	}
	return loads
}
