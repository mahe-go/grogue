package main

import (
	"fmt"
	"github.com/mahe-go/grogue/dungeon"
)

func main() {
	grid := dungeon.NewCellularDungeon(80, 20, 50, 2)
	printGrid(grid)
}

func printGrid(grid *dungeon.Grid) {
	for y := 0; y < grid.Height; y++ {
		for x := 0; x < grid.Width; x++ {
			b, _ := grid.Get(x, y)
			fmt.Print(b)
		}
		fmt.Print("\n")
	}

}
