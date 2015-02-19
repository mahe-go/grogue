package dungeon

import (
	"math/rand"
	"time"
)

type cellularWrapper struct {
	EmptySpacePercentage int
	SolidCellType        CellType
	HollowCellType       CellType
	grid                 *Grid
	shadow               *Grid
}

func newCellularGrid(width int, height int, emptySpacePercentage int, solidCellType CellType, hollowCellType CellType) *cellularWrapper {
	return &cellularWrapper{emptySpacePercentage, solidCellType, hollowCellType, NewRandomGrid(width, height, emptySpacePercentage, solidCellType, hollowCellType), NewSolidGridOfType(width, height, hollowCellType)}
}

func NewCellularDungeon(width int, height int, emptySpacePercentage int, cleanUpRounds int) *Grid {
	rand.Seed(time.Now().Unix())
	wrapper := newCellularGrid(width, height, emptySpacePercentage, SOLID_ROCK, ROOM)
	wrapper.grid.fillRow(wrapper.grid.Height/2, wrapper.HollowCellType)

	for i := 0; i < cleanUpRounds; i++ {
		wrapper.runRoundOfCellularAutomata()
	}

	fillUnreachableCaverns(wrapper)

	return wrapper.grid
}

func (w *cellularWrapper) runRoundOfCellularAutomata() {
	for x := 0; x < w.grid.Width; x++ {
		for y := 0; y < w.grid.Height; y++ {
			count := w.neighbourCount(x, y, w.SolidCellType)
			if count > 3 && count <= 5 {
				cell, _ := w.grid.Get(x, y)
				w.shadow.Set(x, y, cell)
			} else if count <= 3 {
				w.shadow.Set(x, y, NewGridCellOfTypeValue(w.HollowCellType))
			} else {
				w.shadow.Set(x, y, NewGridCellOfTypeValue(w.SolidCellType))
			}
		}
	}
	tmp := w.shadow
	w.shadow = w.grid
	w.grid = tmp
}

func fillUnreachableCaverns(w *cellularWrapper) {
	var x, y int
outer:
	for x = 1; x < w.grid.Width-1; x++ {
		for y = 1; y < w.grid.Height; y++ {
			if cell, _ := w.grid.Get(x, y); cell.Type == w.HollowCellType && w.neighbourCount(x, y, w.HollowCellType) == 8 {
				break outer
			}
		}
	}

	markCellChecked := func(cell GridCell) GridCell {
		cell.Checked = true
		return cell
	}

	markCellUnChecked := func(cell GridCell) GridCell {
		cell.Checked = true
		return cell
	}

	isUnCheckedHollow := func(cell GridCell) bool {
		return !cell.Checked && cell.Type == w.HollowCellType
	}

	w.grid.ApplyToConnected(markCellChecked, isUnCheckedHollow, x, y)

	convertToSolid := func(c GridCell) GridCell { c.Type = w.SolidCellType; return c }

	w.grid.ApplyToAllMatching(convertToSolid, isUnCheckedHollow)

	w.grid.ApplyToAll(markCellUnChecked)
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
