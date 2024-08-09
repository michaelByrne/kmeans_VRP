package main

import (
	"fmt"
	"math/rand"
	"routing/algorithms"
)

func SolveVRP(loads map[int]algorithms.Load, kmeans algorithms.KMeans) ([][]int, error) {
	numLoads := len(loads)

	var allLoads []algorithms.Load
	for _, load := range loads {
		allLoads = append(allLoads, load)
	}

	var routesOut [][]int

	drivers := rand.Intn(numLoads)%30 + 1

	clusters, err := kmeans.Partition(allLoads, drivers)
	if err != nil {
		return nil, err
	}

	var totalOverflow []algorithms.Load

	// Keep reclustering overflow loads until all loads are assigned to a route
	for {
		totalOverflow = nil

		for _, cluster := range clusters {
			clusterLoads := clusterToLoadsMap(cluster)

			route, overflow := algorithms.NearestNeighbor(clusterLoads)

			if len(overflow) > 0 {
				for _, load := range overflow {
					totalOverflow = append(totalOverflow, load)
				}
			}

			routesOut = append(routesOut, route)
		}

		if len(totalOverflow) == 0 {
			break
		}

		overflowDrivers := rand.Intn(len(totalOverflow))%3 + 1

		clusters, err = kmeans.Partition(totalOverflow, overflowDrivers)
		if err != nil {
			return nil, err
		}
	}

	var count int
	for _, r := range routesOut {
		count += len(r)
	}

	if count != numLoads {
		fmt.Printf("not all loads were allocated to a route: %d/%d\n", count, numLoads)
	}

	return routesOut, nil
}

func clusterToLoadsMap(cluster algorithms.Cluster) map[int]algorithms.Load {
	loads := make(map[int]algorithms.Load)
	for _, load := range cluster.Loads {
		loads[load.LoadNumber] = load
	}
	return loads
}
