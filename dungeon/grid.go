package dungeon

import (
	"errors"
	"github.com/mahe-go/grogue/util"
	"math/rand"
)

type CellType uint

const CELLTYPE_CHARS = " #.*"
const SOLID_ROCK CellType = 0
const WALL CellType = 1
const ROOM CellType = 2
const CORRIDOR CellType = 3

var GRID_OVERFLOW error = errors.New("Grid overflow")

type Grid struct {
	cells  []GridCell
	Width  int
	Height int
}

func NewSolidGrid(width int, height int) *Grid {
	return &Grid{make([]GridCell, width*height), width, height}
}

func NewSolidGridOfType(width int, height int, cellType CellType) *Grid {
	grid := &Grid{make([]GridCell, width*height), width, height}
	for i, _ := range grid.cells {
		grid.cells[i].Type = cellType
	}
	return grid
}

func NewRandomGrid(width int, height int, solidCellPercentage int, solidCellType CellType, emptyCellType CellType) *Grid {
	grid := &Grid{make([]GridCell, width*height), width, height}

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			probability := rand.Intn(100)
			if solidCellPercentage > probability {
				grid.Set(x, y, *NewGridCellOfType(solidCellType))
			} else {
				grid.Set(x, y, *NewGridCellOfType(emptyCellType))
			}
		}
	}
	return grid
}

type GridCell struct {
	Type    CellType
	Checked bool
}

type GridCellConditional func(cell GridCell) bool
type GridCellModification func(cell GridCell) GridCell

func NewGridCellOfType(typ CellType) *GridCell {
	return &GridCell{typ, false}
}

func NewGridCellOfTypeValue(typ CellType) GridCell {
	return GridCell{typ, false}
}

func (g GridCell) String() string {
	return string([]rune(CELLTYPE_CHARS)[g.Type])
}

/*
Set cell at (x,y)
*/
func (g *Grid) Set(x int, y int, value GridCell) error {
	index, err := g.cellIndex(x, y)
	if err == nil {
		g.cells[index] = value
	}
	return err
}

/*
Return cell at (x,y)
*/
func (g *Grid) Get(x int, y int) (GridCell, error) {
	index, err := g.cellIndex(x, y)
	if err == nil {
		return g.cells[index], nil
	} else {
		return GridCell{}, err
	}
}

func (g *Grid) cellIndex(x int, y int) (int, error) {
	if x < 0 || x >= g.Width || y < 0 || y >= g.Height {
		return -1, GRID_OVERFLOW
	} else {
		return g.Width*y + x, nil
	}
}

func (g *Grid) fillRow(y int, cellType CellType) (err error) {
	if y < 0 || y >= g.Height {
		err = GRID_OVERFLOW
	} else {
		for i := 0; i < g.Width && err == nil; i++ {
			err = g.Set(i, y, NewGridCellOfTypeValue(cellType))
		}
	}
	return
}

func (g *Grid) fillColumn(x int, cellType CellType) (err error) {
	if x < 0 || x >= g.Width {
		err = GRID_OVERFLOW
	} else {
		for i := 0; i < g.Height && err == nil; i++ {
			err = g.Set(x, i, NewGridCellOfTypeValue(cellType))
		}
	}
	return
}

/*
Test cell at (x,y) for condition
*/
func (g *Grid) Test(condition GridCellConditional, x int, y int) bool {
	index, err := g.cellIndex(x, y)
	if err == nil {
		return condition(g.cells[index])
	} else {
		return false
	}
}

/*
Apply modification to cell at (x,y)
*/
func (g *Grid) Apply(mod GridCellModification, x int, y int) error {
	index, err := g.cellIndex(x, y)
	if err == nil {
		g.cells[index] = mod(g.cells[index])
	}
	return err
}

/*
Apply modification to all cells
*/
func (g *Grid) ApplyToAll(mod GridCellModification) error {
	for i, _ := range g.cells {
		g.cells[i] = mod(g.cells[i])
	}
	return nil
}

/*
Apply modification to all cells matching condition
*/
func (g *Grid) ApplyToAllMatching(mod GridCellModification, condition GridCellConditional) error {
	for i := range g.cells {
		if condition(g.cells[i]) {
			g.cells[i] = mod(g.cells[i])
		}
	}
	return nil
}

/*
Apply modification to cell in (x,y) if it matches condition
*/
func (g *Grid) ApplyToSingleMatching(mod GridCellModification, condition GridCellConditional, x int, y int) error {
	cell, err := g.Get(x, y)
	if err == nil && condition(cell) {
		err = g.Apply(mod, x, y)
	}
	return err
}

/*
Apply modification to all connected cells matching a condition starting from (x,y), using 4-way floodFill algorithm
*/
func (g *Grid) ApplyToConnected(mod GridCellModification, selectCondition GridCellConditional, x int, y int) {
	left := x - 1
	right := x + 1
	up := y - 1
	down := y + 1

	origo, err := g.Get(x, y)
	if err == nil && selectCondition(origo) {
		g.Apply(mod, x, y)
		g.ApplyToConnected(mod, selectCondition, left, y)
		g.ApplyToConnected(mod, selectCondition, right, y)
		g.ApplyToConnected(mod, selectCondition, x, up)
		g.ApplyToConnected(mod, selectCondition, x, down)
	} else {
		return
	}
}

/*
Change type of cells on a straight line from (startX, startY) to (endX, endY). Line is calculated with Bresenham algorithm.
*/
func (grid *Grid) MarkLineBresenham(startx int, starty int, endx int, endy int, cellType CellType) error {

	// Bresenham's line drawing algorithm
	var cx int = startx
	var cy int = starty

	var dx int = util.Abs(endx - cx)
	var dy int = util.Abs(endy - cy)

	var sx int
	var sy int
	if cx < endx {
		sx = 1
	} else {
		sx = -1
	}
	if cy < endy {
		sy = 1
	} else {
		sy = -1
	}
	var e int = dx - dy

	for {
		if cy >= grid.Height || cy < 0 || cx >= grid.Width || cx < 0 {
			return nil
		}
		old, err := grid.Get(cx, cy)
		if err != nil {
			return err
		} else {
			if old.Type == WALL || old.Type == SOLID_ROCK {
				err = grid.Set(cx, cy, *NewGridCellOfType(cellType))
			}
		}

		if err != nil {
			return err
		}
		if (cx == endx) && (cy == endy) {
			return nil
		}
		var e2 int = 2 * e
		if e2 > (0 - dy) {
			e = e - dy
			cx = cx + sx
		}
		if e2 < dx {
			e = e + dx
			cy = cy + sy
		}
	}
}
