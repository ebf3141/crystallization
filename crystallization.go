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

var grid [WIDTH][HEIGHT]int
var grid2 [WIDTH][HEIGHT]int
var g = &grid
var g2 = &grid2

func add_monomers(f int) {
	// adds f monomers
	for i := 0; i < f; i++ {
		xcoord := rand.Intn(WIDTH)
		ycoord := rand.Intn(HEIGHT)
		g[xcoord][ycoord]++
	}
}

func brownian() {
	// moves the monomers around unless they are part of a crystal
	for i := 0; i < WIDTH; i++ {
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
	temp := g
	g = g2
	g2 = temp
}

func find_crystals() {
	var crystal_list [1][3]int
	cl := crystal_list[:]
	max := 0
	count := 0
	for i := 0; i < WIDTH; i++ {
		for j := 0; j < HEIGHT; j++ {
			if g[i][j] >= CRITICAL {
				crystal := [3]int{i, j, g[i][j]}
				cl = append(cl, crystal)
				count++
			}
			if g[i][j] >= max {
				max = g[i][j]
			}
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
		add_monomers(FLUX)
		brownian()
		if i%200 == 0 {
			fmt.Printf("\n Step: %d ", i)
			find_crystals()

		}
	}
}
