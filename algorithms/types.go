package algorithms

import (
	"fmt"
	"math"
)

type Driver struct {
	Route         Route
	MinutesDriven float64
}

type Load struct {
	Origin      Point
	Destination Point
	Center      Point
	LoadNumber  int
}

type Loads []Load

func (l Loads) Center() (Load, error) {
	loadsLength := len(l)

	if loadsLength == 0 {
		return Load{}, fmt.Errorf("no loads to center")
	}

	var center Load
	for i := 0; i < loadsLength; i++ {
		center.Origin.X += l[i].Origin.X
		center.Origin.Y += l[i].Origin.Y
		center.Destination.X += l[i].Destination.X
		center.Destination.Y += l[i].Destination.Y
	}

	var mean Load
	mean.Origin.X = center.Origin.X / float64(loadsLength)
	mean.Origin.Y = center.Origin.Y / float64(loadsLength)
	mean.Destination.X = center.Destination.X / float64(loadsLength)
	mean.Destination.Y = center.Destination.Y / float64(loadsLength)

	return mean, nil
}

func (l Load) Time() float64 {
	return l.Origin.TimeTo(l.Destination)
}

func (l Load) Midpoint() Point {
	return Point{
		X: (l.Origin.X + l.Destination.X) / 2,
		Y: (l.Origin.Y + l.Destination.Y) / 2,
	}
}

type Point struct {
	X float64
	Y float64
}

func (p Point) TimeTo(other Point) float64 {
	return math.Sqrt(math.Pow(other.X-p.X, 2) + math.Pow(other.Y-p.Y, 2))
}

type Route []Load
type Unvisited []Load

func (l Loads) TotalTime(route []int) float64 {
	var total float64
	origin := Point{X: 0, Y: 0}
	last := origin

	for _, r := range route {
		load := l[r]
		begin := load.Origin
		end := load.Destination

		total += last.TimeTo(begin)
		total += begin.TimeTo(end)

		last = end
	}

	return total
}

func (u Unvisited) FindClosest(p Point) (int, float64) {
	var closest int
	var minDistance float64
	for i := 0; i < len(u); i++ {
		distance := p.TimeTo(u[i].Origin)
		if distance < minDistance || minDistance == 0 {
			minDistance = distance
			closest = i
		}
	}

	return closest, minDistance
}
