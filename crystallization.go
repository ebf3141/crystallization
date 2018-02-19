package main

import (
	"fmt"
	"crystallization/fastrand"
	"flag"
	"os"
	"encoding/binary"
)

//Width, height, area of grid
const WIDTH = 1000
const HEIGHT = 1000
var AREA = WIDTH * HEIGHT

//Number of monomers added each step
var FLUX int

//Number of monomers necessary for a crystal to form
var CRITICAL int

//Number of iterations the program runs for
var iterations int

//Names for output files
var imageName string
var dataName string

//Number of times to run brownian motion in each iteration
var diffusion int

// How often to output a data file
var outputs int

//A separate boolean array for storing obstruction locations
var obstructions [WIDTH * HEIGHT]bool

var x_size int
var y_size int
var x_spacing int
var y_spacing int

//Fills in obstructions array
func setObstructions(xSize, xSpacing, ySize, ySpacing int) {
	for x := 0; x < WIDTH; x++ {
		for y := 0; y<HEIGHT; y++ {
			if ((x%(xSize+xSpacing) >= xSpacing) && (y%(ySize+ySpacing) >= ySpacing)) {
				obstructions[x + y * WIDTH] = true
			}
		}
	}
}

//Grid is main data storage
//One-dimensional array of cells, each storing a number of monomers
//Separate file for grid functions
var g = makeGrid(WIDTH, HEIGHT, &obstructions)

//List of free  monomers so I don't have to iterate through the whole grid
var list = make([]int, 500000)
var numMonomers int //number of free monomers



//Lookup tables for more efficient monomer movement
var lookup0, lookup1, lookup2, lookup3, lookup4, lookup5, lookup6, lookup7, lookup8 [WIDTH*HEIGHT]int

//An array of lookup tables
var lookups = [9][WIDTH*HEIGHT]int{lookup0,lookup1,lookup2,lookup3,lookup4,lookup5,lookup6,lookup7,lookup8}



//Creates lookup tables
func createLookups() {
	nearestNeighborOffsets := [9]int {-WIDTH - 1, -WIDTH, -WIDTH + 1, -1, 0, 1, WIDTH-1, WIDTH, WIDTH+1}
	for i := 0; i < 9; i++ {
		for j := 0; j < AREA; j++ {
			newIndex := j + nearestNeighborOffsets[i]
			if newIndex < 0 {
				newIndex = newIndex + AREA
			}
			if newIndex >= AREA {
				newIndex = newIndex - AREA
			}
			if obstructions[newIndex] {
				newIndex = -1
			}
			lookups[i][j] = newIndex
		}
	}
}
		

func addMonomers(f int) {
	// adds f monomers
	for i := 0; i < f; i++ {
		//Random coordinate
		coord := fastrand.Random.Int() % AREA
		//Only adds to unobstructed sites
		if !obstructions[coord] {
			//Adds monomer to grid and list
			g.grid[coord]++
			list[numMonomers] = coord
			numMonomers++
		}
	}
}

func brownian() {
	// moves the monomers around unless they are part of a crystal
	grid := g.grid
	//nextOpenSpace is index in list to write to
	//to prevent empty spaces in list, separate read/write locations
	//monomers all moved to next open space
	//when monomers join crystals, they are overwritten
	nextOpenSpace := 0
	//loops through monomer-containing cells
	for i := 0; i < numMonomers; i++ {
		cell := grid[list[i]]
		if cell < int32(CRITICAL) { 
			//generates a new coordinate for the monomer to move to
			newcoord := -1
			for newcoord < 0 {
				//Uses customized rng for more efficiency
				newcoord = lookups[fastrand.Rand9()][list[i]]
			}
			//moves monomer in grid
			grid[newcoord]++
			grid[list[i]]--
			//moves monomer in list and increments write location
			list[nextOpenSpace] = newcoord
			nextOpenSpace++
		}
		//if monomer part of crystal, nextOpenSpace not incremented and crystal overwritten
		
	}
	//adjusts number of free monomers
	numMonomers = nextOpenSpace
}

//finds all crystals
func findCrystals() {
	var crystalList [1][3]int //stores coordinates, num monomers
	cl := crystalList[:] 
	max := 0 //max monomers in one cell
	count := 0 //count of crystals
	monomers := 0 //total monomers not part of crystals
	for i := 0; i < AREA; i++ { //loops through grid
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
	//prints results
	fmt.Printf("Max # monomers: %d ", max)
	fmt.Printf("# Crystals: %d ", count)
	fmt.Printf("Avg monomers: %f",float64(monomers)/float64(g.area))
}

//takes command-line arguments, sets variables
func initialize() {
	flag.IntVar(&FLUX, "flux", 1000, "number of monomers added / step")
	flag.IntVar(&CRITICAL, "crit", 5, "critical number for crystallization")
	flag.IntVar(&iterations, "iter", 3000, "number of iterations")
	flag.IntVar(&diffusion, "diff", 1, "diffusion constant")
	flag.IntVar(&outputs, "out", 5000, "how many iterations to wait between outputting file")
	flag.StringVar(&imageName, "image", "crystals.png", "name for file containing image")
	flag.StringVar(&dataName, "file", "data", "name for data file")
	flag.IntVar(&x_size, "xSize", 80, "width of obstructions")
	flag.IntVar(&x_spacing, "xSpacing", 80, "spacing between obstructions (edge-to-edge")
	flag.IntVar(&y_size, "ySize", 80, "height of obstructions")
	flag.IntVar(&y_spacing, "ySpacing", 80, "spacing between obstructions (edge-to-edge")
	flag.Parse()
}

//writes contents of grid to file
func writeGrid(filename string) {
	f, _ := os.Create(filename)
	//file header has magic # plus constants
	binary.Write(f, binary.LittleEndian, uint32(0x21464245))//magic number: "EBF!" for id
	binary.Write(f, binary.LittleEndian, uint32(g.W)) 
	binary.Write(f, binary.LittleEndian, uint32(g.H)) 
	binary.Write(f, binary.LittleEndian, uint32(FLUX))
	binary.Write(f, binary.LittleEndian, uint32(CRITICAL))
	binary.Write(f, binary.LittleEndian, uint32(iterations))
	//then writes grid itself
	err := binary.Write(f, binary.LittleEndian, g.grid)
	if err != nil {
		fmt.Println("binary.Write failed:", err) //error handler
	}
	f.Close()
}
	
//main function

func main() {
	//sets stuff up
	initialize()
	numMonomers = 0
	setObstructions(80,160,80,160)
	createLookups()
	//main loop of program
	for i := 0; i < iterations; i++ {
		//each step: add monomer, do brownian motion
		addMonomers(FLUX)
		for j := 0; j < diffusion; j++ {
			brownian()
		}
		//occasionally: run findCrystals and create files
		if i%200 == 0 {
			fmt.Printf("\n Step: %d ", i)
			findCrystals()
		}
		if i%outputs == 0{
			writeGrid(dataName + "-" + fmt.Sprint(i))
		}
	}
	//at end: create image and one last file
	g.createImage(imageName)
	writeGrid(dataName)
	
}

/*
func main() {
	initialize()
	setObstructions(40,60,40,60)
	addMonomers(100000)
	for i:=0; i<100; i++ {
		brownian()
	}
	g.createImage("imagetest.png")
}

*/
