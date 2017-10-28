package grid

import (
	"errors"
	"math/rand"

	"github.com/mahe-go/grogue/util"
)

type Direction struct {
	Dx int
	Dy int
}

var North = Direction{0, -1}
var NorthEast = Direction{1, -1}
var East = Direction{1, 0}
var SouthEast = Direction{1, 1}
var South = Direction{0, 1}
var SouthWest = Direction{-1, 1}
var West = Direction{-1, 0}
var NorthWest = Direction{-1, -1}

var GRID_OVERFLOW error = errors.New("Grid overflow")

type Grid struct {
	cells  []GridCell
	Width  int
	Height int
}

type shadowWrapper struct {
	grid   *Grid
	shadow *Grid
}

//Return a new grid of default cells
func NewSolidGrid(width int, height int) *Grid {
	return &Grid{make([]GridCell, width*height), width, height}
}

//Return a new grid with all cells of type cellType
func NewSolidGridOfType(width int, height int, cellType CellType) *Grid {
	grid := &Grid{make([]GridCell, width*height), width, height}
	for i, _ := range grid.cells {
		grid.cells[i].Type = cellType
	}
	return grid
}

//Return a grid with cells of type emptyCellType and solidCellPercentage% cells of type solidCellType at random locations
func NewRandomGrid(width int, height int, solidCellPercentage int, solidCellType CellType, emptyCellType CellType) *Grid {
	grid := &Grid{make([]GridCell, width*height), width, height}

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			probability := rand.Intn(100)
			if solidCellPercentage > probability {
				grid.Set(x, y, NewGridCellOfTypeValue(solidCellType))
			} else {
				grid.Set(x, y, NewGridCellOfTypeValue(emptyCellType))
			}
		}
	}
	return grid
}

//Set cell at (x,y)
func (g *Grid) Set(x int, y int, value GridCell) error {
	index, err := g.cellIndex(x, y)
	if err == nil {
		g.cells[index] = value
	}
	return err
}

//Return cell at (x,y)
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

//Test cell at (x,y) for condition
func (g *Grid) TestCellAtXY(condition CellPredicate, x int, y int) bool {
	index, err := g.cellIndex(x, y)
	if err == nil {
		return condition(g.cells[index])
	} else {
		return false
	}
}

//Test a condition at (x,y)
func (g *Grid) TestAtXY(condition LocationPredicate, x int, y int) bool {
	return condition(g, x, y)
}

//Apply modification to cell at (x,y)
func (g *Grid) ApplyToCellAtXY(mod CellModification, x int, y int) error {
	index, err := g.cellIndex(x, y)
	if err == nil {
		g.cells[index] = mod(g.cells[index])
	}
	return err
}

//Apply modification at (or starting from) (x,y)
func (g *Grid) ApplyaAtXY(mod Modification, x int, y int) error {
	return mod(g, x, y)
}

//Apply modification to all cells
func (g *Grid) ApplyToAllCells(mod CellModification) error {
	for i, _ := range g.cells {
		g.cells[i] = mod(g.cells[i])
	}
	return nil
}

//Apply modification everywhere
func (g *Grid) ApplyEverywhere(mod Modification) error {
	var err error
	for x := 0; x < g.Width && err == nil; x++ {
		for y := 0; y < g.Height && err == nil; y++ {
			err = mod(g, x, y)
		}
	}
	return err
}

//Apply modification to all cells matching condition
func (g *Grid) ApplyToAllMatchingCells(mod CellModification, condition CellPredicate) error {
	for i := range g.cells {
		if condition(g.cells[i]) {
			g.cells[i] = mod(g.cells[i])
		}
	}
	return nil
}

// Apply modification everywhere the condtion matches
func (g *Grid) ApplyEverywhereMatching(mod Modification, condition LocationPredicate) error {
	var err error
	for x := 0; x < g.Width && err == nil; x++ {
		for y := 0; y < g.Height && err == nil; y++ {
			if condition(g, x, y) {
				err = mod(g, x, y)
			}
		}
	}
	return err
}

// Apply modification to a cell if condition matches at its location
func (g *Grid) ApplyToCellEverywhereMatching(mod CellModification, condition LocationPredicate) error {
	var err error
	for x := 0; x < g.Width && err == nil; x++ {
		for y := 0; y < g.Height && err == nil; y++ {
			if condition(g, x, y) {
				err = g.ApplyToCellAtXY(mod, x, y)
			}
		}
	}
	return err
}

//Apply modification to cell in (x,y) if it matches condition
func (g *Grid) ApplyToCellAtXYMatching(mod CellModification, condition CellPredicate, x int, y int) error {
	cell, err := g.Get(x, y)
	if err == nil && condition(cell) {
		err = g.ApplyToCellAtXY(mod, x, y)
	}
	return err
}

// Apply modification at(x,y) if condition matches at(x,y)
func (g *Grid) ApplyAtXYMAtching(mod Modification, condition LocationPredicate, x int, y int) error {
	if condition(g, x, y) {
		return mod(g, x, y)
	} else {
		return nil
	}

}

//Apply modification to all connected cells matching a condition starting from (x,y), using 4-way floodFill algorithm
func (g *Grid) ApplyToConnectedCells(mod CellModification, selectCondition CellPredicate, x int, y int) {
	left := x - 1
	right := x + 1
	up := y - 1
	down := y + 1

	origo, err := g.Get(x, y)
	if err == nil && selectCondition(origo) {
		g.ApplyToCellAtXY(mod, x, y)
		g.ApplyToConnectedCells(mod, selectCondition, left, y)
		g.ApplyToConnectedCells(mod, selectCondition, right, y)
		g.ApplyToConnectedCells(mod, selectCondition, x, up)
		g.ApplyToConnectedCells(mod, selectCondition, x, down)
	} else {
		return
	}
}

// Surround all empty space with walls.
func (grid *Grid) buildCavernWalls() {
	hasOnlyWallsAroundAtXY := func(g *Grid, x int, y int) bool {
		return grid.CountNeighboursMatching(GridCellIsOfType(WALL), x, y) == 8
	}

	grid.ApplyToAllMatchingCells(GridCellTypeConverter(WALL), GridCellIsOfType(SOLID_ROCK))

	// check all cells that have only wall neighbours
	grid.ApplyToCellEverywhereMatching(GridCellChecker, hasOnlyWallsAroundAtXY)

	// convert checked cells back to solid rock and uncheck
	grid.ApplyToAllMatchingCells(GridCellTypeConverter(SOLID_ROCK).And(GridCellUnChecker), GridCellIsChecked)
}

// Apply modification to cell matching condition on a straight line from (startX, startY) to (endX, endY).
// Line is calculated with Bresenham algorithm.
func (grid *Grid) ApplyOnLine(mod CellModification, cond CellPredicate, startx int, starty int, endx int, endy int) {

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
			return
		}
		err := grid.ApplyToCellAtXYMatching(mod, cond, cx, cy)
		if err != nil {
			return
		}
		if (cx == endx) && (cy == endy) {
			return
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

func (g *Grid) CountNeighboursMatching(condition CellPredicate, x int, y int) int {
	count := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			} else {
				cell, err := g.Get(x+i, y+j)
				if err == GRID_OVERFLOW {
					count++
				} else if condition(cell) {
					count++
				}
			}
		}
	}
	return count
}

func (g Grid) String() string {
	var buffer string
	for y := 0; y < g.Height; y++ {
		line := make([]rune, g.Width)
		for x := 0; x < g.Width; x++ {
			var cell GridCell
			cell, err := g.Get(x, y)
			if err != nil {
				line[x] = '?'
			} else {
				line[x] = cell.Type.Rune
			}
		}
		buffer += string(line) + "\n"
	}
	return buffer
}

func (g *Grid) AddStairCases() {
	var x, y int
	for x, y = rand.Intn(g.Width), rand.Intn(g.Height); g.TestCellAtXY(CellIsTraversable.Not(), x, y); x, y = rand.Intn(g.Width), rand.Intn(g.Height) {
	}
	g.ApplyToCellAtXY(GridCellTypeConverter(STAIRCASE_UP), x, y)

	for x, y = rand.Intn(g.Width), rand.Intn(g.Height); g.TestCellAtXY(CellIsTraversable.Not(), x, y); x, y = rand.Intn(g.Width), rand.Intn(g.Height) {
	}
	g.ApplyToCellAtXY(GridCellTypeConverter(STAIRCASE_DOWN), x, y)
}
