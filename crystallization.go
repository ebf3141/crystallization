package main

import (
	"fmt"
	"math/rand"
	"flag"
)

const WIDTH = 1000
const HEIGHT = 1000
var AREA = WIDTH * HEIGHT
//var FLUX = 1000
//var CRITICAL = 8
var FLUX int
var CRITICAL int
var iterations int
/*
var grid [WIDTH][HEIGHT]int
var grid2 [WIDTH][HEIGHT]int
var g = &grid
var g2 = &grid2

var wg sync.WaitGroup
*/

var g = makeGrid(WIDTH, HEIGHT)
var g2 = makeGrid(WIDTH, HEIGHT)

func addMonomers(f int) {
	// adds f monomers
	for i := 0; i < f; i++ {
		coord := rand.Intn(AREA)
		g.grid[coord]++
	}
}

func brownian() {
	// moves the monomers around unless they are part of a crystal
	grid := g.grid
	grid2 := g2.grid

	for i := 0; i < AREA; i++ {
		cell := grid[i]
		if cell < CRITICAL {
			for k := 0; k < cell; k++ {
				newcoord := g.getRandomNeighbor(i)
				grid2[newcoord]++
			}
		} else {
			grid2[i]=grid[i]
		}
		grid[i] = 0
	}
	temp := g
	g = g2
	g2 = temp
}
/*
func cBrownian(startRow int, endRow int){
	// moves the monomers around unless they are part of a crystal
	for i := startRow; i < endRow; i++ {
		for j := 0; j < HEIGHT; j++ {
			cell := g[i][j]
			if cell < CRITICAL {
				for k := 0; k < cell; k++ {
					newX := i + rand.Intn(3) - 1
					newY := j + rand.Intn(3) - 1
					if newX == 1000 {
						newX = 0
					}
					if newX == -1 {
						newX = 999
					}
					if newY == 1000 {
						newY = 0
					}
					if newY == -1 {
						newY = 999
					}
					g2[newX][newY]++
				}
			} else {
				g2[i][j] = g[i][j]
			}
			g[i][j] = 0
		}
	}
	wg.Done()
}

func syncBrownian() {
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go cBrownian(HEIGHT/4*i, HEIGHT/4*(i+1))
	}
	wg.Wait()
	temp := g
	g = g2
	g2 = temp
}
*/

func findCrystals() {
	var crystalList [1][3]int
	cl := crystalList[:]
	max := 0
	count := 0
	for i := 0; i < AREA; i++ {
		if g.grid[i] >= CRITICAL {
			x, y := g.indexToXY(i)
			crystal := [3]int{x,y, g.grid[i]}
			cl = append(cl, crystal)
			count++
		}
		if g.grid[i] >= max {
			max = g.grid[i]
		}
	}
	//fmt.Print(cl)
	fmt.Printf("Max # monomers: %d ", max)
	fmt.Printf("# Crystals: %d ", count)
}

func initialize() {
	flag.IntVar(&FLUX, "flux", 1000, "number of monomers added / step")
	flag.IntVar(&CRITICAL, "crit", 5, "critical number for crystallization")
	flag.IntVar(&iterations, "iter", 3000, "number of iterations")
	flag.Parse()
}

func main() {
	initialize()
	for i := 0; i < iterations; i++ {
		addMonomers(FLUX)
		brownian()
		if i%200 == 0 {
			fmt.Printf("\n Step: %d ", i)
			findCrystals()

		}
	}
	g.createImage("test.png")
}
