package main

import(
	"math/rand"
	"image"
	"image/color"
	"image/png"
	"os"
)

type Grid struct {
	W int
	H int
	grid []int
}

func makeGrid(width, height int) *Grid {
	return &Grid{width, height, make([]int, width*height)}
}
	

func (g *Grid) xyToIndex(x, y int) int {
	return x + y * g.W
}

func (g *Grid) indexToXY(c int) (x, y int) {
	x = c % g.W
	y = c / g.W
	return
}

func (g *Grid) indexDelta(index, dx, dy int) int {
	newX, newY := g.indexToXY(index)
	newX += dx
	newY += dy
	if newX < 0 {
		newX += 1000
	}
	if newY < 0 {
		newY += 1000
	}
	if newX > 999 {
		newX -= 1000
	}
	if newY > 999 {
		newY -= 1000
	}
	return g.xyToIndex(newX, newY)
}

func (g *Grid) getRandomNeighbor(c int) int {
	deltaX := rand.Intn(3) - 1
	deltaY := rand.Intn(3) - 1
	return g.indexDelta(c, deltaX, deltaY)
}

func (g *Grid) ColorModel() color.Model {
	return color.GrayModel
}

func (g *Grid) Bounds() image.Rectangle {
	return image.Rectangle{image.Point{0,0}, image.Point{g.W, g.H}}
}

func (g *Grid) At(x, y int) color.Color {
	index := g.xyToIndex(x,y)
	return color.Gray{uint8(g.grid[index])}
}

func (g *Grid) createImage(name string) {
	f, _ := os.Create(name)
	png.Encode(f,g)
	f.Close()
}

