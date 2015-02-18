package dungeon

import (
	"math/rand"
	"time"
)

type cellularWrapper struct {
	EmptySpacePercentage int
	LiveCellType         CellType
	DeadCellType         CellType
	grid                 *Grid
	shadow               *Grid
}

func newCellularGrid(width int, height int, emptySpacePercentage int, liveCellType CellType, deadCellType CellType) *cellularWrapper {
	return &cellularWrapper{emptySpacePercentage, liveCellType, deadCellType, NewRandomGrid(width, height, emptySpacePercentage, liveCellType, deadCellType), NewSolidGridOfType(width, height, deadCellType)}
}

func NewCellularDungeon(width int, height int, emptySpacePercentage int, cleanUpRounds int) *Grid {
	rand.Seed(time.Now().Unix())
	wrapper := newCellularGrid(width, height, emptySpacePercentage, SOLID_ROCK, ROOM)
	wrapper.grid.fillRow(wrapper.grid.Height/2, wrapper.LiveCellType)

	for i := 0; i < cleanUpRounds; i++ {
		wrapper.runRoundOfCellularAutomata()
	}

	return wrapper.grid
}

func (w *cellularWrapper) runRoundOfCellularAutomata() {
	for x := 0; x < w.grid.Width; x++ {
		for y := 0; y < w.grid.Height; y++ {
			count := w.neighbourCount(x, y, w.LiveCellType)
			if count > 3 && count <= 5 {
				cell, _ := w.grid.Get(x, y)
				w.shadow.Set(x, y, cell)
			} else if count <= 3 {
				w.shadow.Set(x, y, *NewGridCell(w.DeadCellType))
			} else {
				w.shadow.Set(x, y, *NewGridCell(w.LiveCellType))
			}
		}
	}
	tmp := w.shadow
	w.shadow = w.grid
	w.grid = tmp
}

func (w *cellularWrapper) neighbourCount(x int, y int, celltype CellType) int {
	count := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			} else {
				cell, err := w.grid.Get(x+i, y+j)
				if err == GRID_OVERFLOW {
					count++
				} else if cell.Type == celltype {
					count++
				}
			}
		}
	}
	return count
}
