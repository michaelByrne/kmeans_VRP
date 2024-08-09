package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"routing/algorithms"
	"strconv"
	"strings"
)

func main() {
	fileName := os.Args[1]
	if fileName == "" {
		fmt.Println("Please provide a filename")
		os.Exit(1)
	}

	loads, err := parseLoadFile(fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	kmeans := algorithms.NewKMeans(100)

	bestTime := math.MaxFloat64
	var bestRoutes [][]int

	for i := 0; i < len(loads)/2; i++ {
		routes, err := SolveVRP(loads, kmeans)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		totalTime, err := algorithms.TotalTime(loads, routes)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if totalTime < bestTime {
			bestTime = totalTime
			bestRoutes = routes
		}

	}

	for _, route := range bestRoutes {
		fmt.Print("[")
		for j, loadIndex := range route {
			if loadIndex == -1 {
				continue
			}

			load := loads[loadIndex]
			fmt.Printf("%d", load.LoadNumber)
			if j < len(route)-1 {
				fmt.Print(", ")
			}
		}

		fmt.Println("]")
	}
}

func parseLoadFile(filename string) (map[int]algorithms.Load, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	loads := make(map[int]algorithms.Load)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "loadNumber") {
			continue
		}

		load, err := parseLoadLine(line)
		if err != nil {
			return nil, err
		}

		loads[load.LoadNumber] = load
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return loads, nil
}

func parseLoadLine(line string) (algorithms.Load, error) {
	fields := strings.Fields(line)
	if len(fields) != 3 {
		return algorithms.Load{}, fmt.Errorf("invalid line format: %s", line)
	}

	loadNumber, err := strconv.Atoi(fields[0])
	if err != nil {
		return algorithms.Load{}, err
	}

	origin, err := parsePoint(fields[1])
	if err != nil {
		return algorithms.Load{}, err
	}

	destination, err := parsePoint(fields[2])
	if err != nil {
		return algorithms.Load{}, err
	}

	loadOut := algorithms.Load{
		LoadNumber:  loadNumber,
		Origin:      origin,
		Destination: destination,
	}

	loadOut.Center = loadOut.Midpoint()

	return loadOut, nil
}

func parsePoint(pointStr string) (algorithms.Point, error) {
	pointStr = strings.Trim(pointStr, "()")
	coords := strings.Split(pointStr, ",")
	if len(coords) != 2 {
		return algorithms.Point{}, fmt.Errorf("invalid point format: %s", pointStr)
	}

	x, err := strconv.ParseFloat(coords[0], 64)
	if err != nil {
		return algorithms.Point{}, err
	}

	y, err := strconv.ParseFloat(coords[1], 64)
	if err != nil {
		return algorithms.Point{}, err
	}

	return algorithms.Point{X: x, Y: y}, nil
}
