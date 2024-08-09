package algorithms

import (
	"fmt"
	"math/rand"
)

type Cluster struct {
	Center Point
	Loads  Loads
}

type Clusters []Cluster

func NewClusters(loads Loads, k int) (Clusters, error) {
	var clusters []Cluster

	if len(loads) == 0 {
		return nil, fmt.Errorf("no loads to cluster")
	}

	if k <= 0 {
		return nil, fmt.Errorf("invalid number of clusters: %d", k)
	}

	for i := 0; i < k; i++ {
		p := Point{}

		p.Y = rand.Float64() * 100
		p.X = rand.Float64() * 100

		clusters = append(clusters, Cluster{
			Center: p,
		})
	}

	return clusters, nil
}

func (c Clusters) ClosestTo(load Load) int {
	var closest int
	var minTime float64

	for i := 0; i < len(c); i++ {
		travelTime := c[i].Center.TimeTo(load.Center)
		if travelTime < minTime || minTime == 0 {
			minTime = travelTime
			closest = i
		}
	}

	return closest
}

func (c Clusters) Reset() {
	for i := 0; i < len(c); i++ {
		c[i].Loads = Loads{}
	}
}

func (c *Cluster) Recenter() {
	center, err := c.Loads.Center()
	if err != nil {
		return
	}

	c.Center = center.Midpoint()
}

func (c Clusters) Recenter() {
	for i := 0; i < len(c); i++ {
		c[i].Recenter()
	}
}
