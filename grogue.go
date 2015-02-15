package main

import (
	"fmt"
	"github.com/mahe-go/grogue/dungeon"
)

func main() {
	grid, err := dungeon.BSPDungeon(80, 40, 7, 7)
	if err != nil {
		fmt.Errorf("Error %s", err)
	}
	for y := 0; y < grid.Height; y++ {
		for x := 0; x < grid.Width; x++ {
			b, _ := grid.Get(x, y)
			fmt.Print(b)
		}
		fmt.Print("\n")
	}
}
