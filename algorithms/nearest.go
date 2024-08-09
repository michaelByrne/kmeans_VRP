package algorithms

import (
	"fmt"
	"math"
)

func NearestNeighbor(loads map[int]Load) ([]int, map[int]Load) {
	var route []int
	var unassigned []int
	overflow := make(map[int]Load)

	for loadNumber := range loads {
		unassigned = append(unassigned, loadNumber)
	}

	origin := Point{X: 0, Y: 0}
	totalTime := 0.0

	for len(unassigned) > 0 {
		var lastLoadIndex int
		if len(route) == 0 {
			// First load, start from origin
			lastLoadIndex = -1
		} else {
			lastLoadIndex = route[len(route)-1]
		}

		closest := -1
		closestTime := math.MaxFloat64
		for _, loadIndex := range unassigned {
			var distance float64

			if lastLoadIndex == -1 {
				// Distance from origin to load's origin
				distance = origin.TimeTo(loads[loadIndex].Origin)
			} else {
				// Distance from last destination to next load's origin
				distance = loads[lastLoadIndex].Destination.TimeTo(loads[loadIndex].Origin)
			}

			// Calculate the total time if this load is added
			newTotalTime := totalTime + distance + loads[loadIndex].Origin.TimeTo(loads[loadIndex].Destination) + loads[loadIndex].Destination.TimeTo(origin)

			if newTotalTime <= 720 && distance < closestTime {
				closest = loadIndex
				closestTime = distance
			}
		}

		if closest == -1 {
			// No valid load found, break the loop
			break
		}

		route = append(route, closest)
		totalTime += closestTime + loads[closest].Origin.TimeTo(loads[closest].Destination)
		unassigned = removeIndex(unassigned, closest)
	}

	// Any remaining unassigned loads are overflow
	for _, loadIndex := range unassigned {
		overflow[loadIndex] = loads[loadIndex]
	}

	return route, overflow
}

func removeIndex(unassigned []int, closest int) []int {
	for i := 0; i < len(unassigned); i++ {
		if unassigned[i] == closest {
			return append(unassigned[:i], unassigned[i+1:]...)
		}
	}
	return unassigned
}

func TotalTime(loads map[int]Load, routes [][]int) (float64, error) {
	totalTime := 0.0

	for _, route := range routes {
		routeTime := TotalRouteTime(loads, route)

		if routeTime > float64(720) {
			return 0, fmt.Errorf("route time exceeds 12 hours: %f", routeTime)
		}

		totalTime = totalTime + routeTime + 500
	}

	return totalTime, nil
}

func TotalRouteTime(loads map[int]Load, route []int) float64 {
	totalTime := 0.0
	origin := Point{X: 0, Y: 0}
	last := origin

	for _, r := range route {
		load := loads[r]
		begin := load.Origin
		end := load.Destination

		totalTime += last.TimeTo(begin)
		totalTime += begin.TimeTo(end)

		last = end
	}

	return totalTime
}
