package main

import (
	"fmt"
	"crystallization/fastrand"
	"flag"
	"os"
	"encoding/binary"
)

const WIDTH = 1000
const HEIGHT = 1000
var AREA = WIDTH * HEIGHT
//var FLUX = 1000
//var CRITICAL = 8
var FLUX int
var CRITICAL int
var iterations int
var imageName string
var dataName string
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
		coord := fastrand.Random.Int() % AREA
		g.grid[coord]++
	}
}

func brownian() {
	// moves the monomers around unless they are part of a crystal
	grid := g.grid
	grid2 := g2.grid
	var k int32
	for i := 0; i < AREA; i++ {
		cell := grid[i]
		if cell < int32(CRITICAL) {
			for k = 0; k < cell; k++ {
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
	monomers := 0
	for i := 0; i < AREA; i++ {
		if g.grid[i] >= int32(CRITICAL) {
			x, y := g.indexToXY(i)
			crystal := [3]int{x,y, int(g.grid[i])}
			cl = append(cl, crystal)
			count++
		} else {
			monomers += int(g.grid[i])
		}
		if int(g.grid[i]) >= max {
			max = int(g.grid[i])
		}
	}
	//fmt.Print(cl)
	fmt.Printf("Max # monomers: %d ", max)
	fmt.Printf("# Crystals: %d ", count)
	fmt.Printf("Avg monomers: %f",float64(monomers)/float64(g.area))
}

func initialize() {
	flag.IntVar(&FLUX, "flux", 1000, "number of monomers added / step")
	flag.IntVar(&CRITICAL, "crit", 5, "critical number for crystallization")
	flag.IntVar(&iterations, "iter", 3000, "number of iterations")
	flag.StringVar(&imageName, "image", "crystals.png", "name for file containing image")
	flag.StringVar(&dataName, "file", "data", "name for data file")
	flag.Parse()
}

func writeGrid(filename string) {
	f, _ := os.Create(filename)
	binary.Write(f, binary.LittleEndian, uint32(0x21464245))//magic number: "EBF!"
	binary.Write(f, binary.LittleEndian, uint32(g.W))
	binary.Write(f, binary.LittleEndian, uint32(g.H))
	binary.Write(f, binary.LittleEndian, uint32(FLUX))
	binary.Write(f, binary.LittleEndian, uint32(CRITICAL))
	binary.Write(f, binary.LittleEndian, uint32(iterations))
	err := binary.Write(f, binary.LittleEndian, g.grid)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	f.Close()
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
	g.createImage(imageName)
	writeGrid(dataName)
	
}
