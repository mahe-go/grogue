package grid

import (
	"math/rand"
	"time"
)

type wrapper struct {
	SolidCellType  CellType
	HollowCellType CellType
	grid           *Grid
	shadow         *Grid
}

func newWrapper(width int, height int, emptySpacePercentage int, solidCellType CellType, hollowCellType CellType) *wrapper {
	return &wrapper{solidCellType, hollowCellType, NewRandomGrid(width, height, emptySpacePercentage, solidCellType, hollowCellType), NewSolidGridOfType(width, height, hollowCellType)}
}

func NewNaturalCavernGrid(width int, height int, emptySpacePercentage int, cleanUpRounds int) *Grid {
	rand.Seed(time.Now().Unix())
	wrapper := newWrapper(width, height, emptySpacePercentage, SOLID_ROCK, ROOM)

	for i := 0; i < cleanUpRounds; i++ {
		wrapper.runRoundOfCellularAutomata()
	}

	fillUnreachableCaverns(wrapper)

	wrapper.grid.buildCavernWalls()

	wrapper.grid.AddStairCases()

	return wrapper.grid
}

func (w *wrapper) runRoundOfCellularAutomata() {
	for x := 0; x < w.grid.Width; x++ {
		for y := 0; y < w.grid.Height; y++ {
			count := w.grid.CountNeighboursMatching(GridCellIsOfType(w.SolidCellType), x, y)
			if count > 3 && count <= 5 {
				cell, _ := w.grid.Get(x, y)
				w.shadow.Set(x, y, cell)
			} else if count <= 3 {
				w.shadow.ApplyToCellAtXY(GridCellTypeConverter(w.HollowCellType), x, y)
			} else {
				w.shadow.ApplyToCellAtXY(GridCellTypeConverter(w.SolidCellType), x, y)
			}
		}
	}
	tmp := w.shadow
	w.shadow = w.grid
	w.grid = tmp
}

func fillUnreachableCaverns(w *wrapper) {
	isSuitableStartingPoint := func(x int, y int) bool {
		return w.grid.TestCellAtXY(GridCellIsOfType(w.HollowCellType), x, y) &&
			w.grid.CountNeighboursMatching(GridCellIsOfType(w.HollowCellType), x, y) == 8
	}
	var x, y int
outer:
	for y = 1; y < w.grid.Height; y++ {
		for x = 1; x < w.grid.Width-1; x++ {
			if isSuitableStartingPoint(x, y) {
				break outer
			}
		}
	}

	w.grid.ApplyToConnectedCells(GridCellChecker, GridCellIsChecked.Not().And(GridCellIsOfType(w.HollowCellType)), x, y)

	w.grid.ApplyToAllMatchingCells(GridCellTypeConverter(w.SolidCellType), GridCellIsChecked.Not().And(GridCellIsOfType(w.HollowCellType)))

	w.grid.ApplyToAllCells(GridCellUnChecker)
}
