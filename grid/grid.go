package grid

import (
	"errors"
	"github.com/mahe-go/grogue/util"
	"math/rand"
)

type CellType uint

const CELLTYPE_CHARS = " #.."
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

type GridCell struct {
	Type    CellType
	Checked bool
}

// Function type for function testing a grid cell
type GridCellConditional func(cell GridCell) bool

// Function type for function modifying a grid cell. Returns the modified cell.
type GridCellModification func(cell GridCell) GridCell

// Function type for function testing a grid cell at (x,y)
// GridConditional makes it possible to create conditions
// based on the environment of the cell we're testing, since we know its location
type GridConditional func(grid *Grid, x int, y int) bool

// Function type for function modifying a grid cell at (x,y) (or region of grid starting from (x,y) as well)
type GridModification func(grid *Grid, x int, y int) error

func NewGridCellOfType(typ CellType) *GridCell {
	return &GridCell{typ, false}
}

func NewGridCellOfTypeValue(typ CellType) GridCell {
	return GridCell{typ, false}
}

func (g GridCell) String() string {
	return string([]rune(CELLTYPE_CHARS)[g.Type])
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
func (g *Grid) TestCellXY(condition GridCellConditional, x int, y int) bool {
	index, err := g.cellIndex(x, y)
	if err == nil {
		return condition(g.cells[index])
	} else {
		return false
	}
}

//Test a condition at (x,y)
func (g *Grid) TestXY(condition GridConditional, x int, y int) bool {
	return condition(g, x, y)
}

//Apply modification to cell at (x,y)
func (g *Grid) ApplyToCellAt(mod GridCellModification, x int, y int) error {
	index, err := g.cellIndex(x, y)
	if err == nil {
		g.cells[index] = mod(g.cells[index])
	}
	return err
}

//Apply modification grid at (or starting from) (x,y)
func (g *Grid) ApplyaAtXY(mod GridModification, x int, y int) error {
	return mod(g, x, y)
}

//Apply modification to all cells
func (g *Grid) ApplyToAllCells(mod GridCellModification) error {
	for i, _ := range g.cells {
		g.cells[i] = mod(g.cells[i])
	}
	return nil
}

//Apply modification everywhere
func (g *Grid) ApplyEverywhere(mod GridModification) error {
	var err error
	for x := 0; x < g.Width && err == nil; x++ {
		for y := 0; y < g.Height && err == nil; y++ {
			err = mod(g, x, y)
		}
	}
	return err
}

//Apply modification to all cells matching condition
func (g *Grid) ApplyToAllMatchingCells(mod GridCellModification, condition GridCellConditional) error {
	for i := range g.cells {
		if condition(g.cells[i]) {
			g.cells[i] = mod(g.cells[i])
		}
	}
	return nil
}

// Apply modification everywhere
func (g *Grid) ApplyEverywhereMatching(mod GridModification, condition GridConditional) error {
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

func (grid *Grid) buildCavernWalls() {
	grid.ApplyToAllMatchingCells(func(cell GridCell) GridCell {
		cell.Type = WALL
		return cell
	}, func(cell GridCell) bool {
		return cell.Type == SOLID_ROCK
	})

	shadow := NewSolidGridOfType(grid.Width, grid.Height, WALL)

	shadow.ApplyEverywhereMatching(func(g *Grid, x int, y int) error {
		return g.Set(x, y, NewGridCellOfTypeValue(SOLID_ROCK))
	}, func(g *Grid, x int, y int) bool {
		return grid.neighbourCount(x, y, WALL) == 8
	})

	grid.ApplyEverywhereMatching(func(g *Grid, x int, y int) error {
		return g.Set(x, y, NewGridCellOfTypeValue(SOLID_ROCK))
	}, func(g *Grid, x int, y int) bool {
		cell, err := shadow.Get(x, y)
		return err == nil && cell.Type == SOLID_ROCK
	})
}

//Apply modification to cell in (x,y) if it matches condition
func (g *Grid) ApplyCellAtXYMatching(mod GridCellModification, condition GridCellConditional, x int, y int) error {
	cell, err := g.Get(x, y)
	if err == nil && condition(cell) {
		err = g.ApplyToCellAt(mod, x, y)
	}
	return err
}

// Apply modification at(x,y) if condition matches at(x,y)
func (g *Grid) ApplyAtXYMAtching(mod GridModification, condition GridConditional, x int, y int) error {
	if condition(g, x, y) {
		return mod(g, x, y)
	} else {
		return nil
	}

}

//Apply modification to all connected cells matching a condition starting from (x,y), using 4-way floodFill algorithm
func (g *Grid) ApplyToConnected(mod GridCellModification, selectCondition GridCellConditional, x int, y int) {
	left := x - 1
	right := x + 1
	up := y - 1
	down := y + 1

	origo, err := g.Get(x, y)
	if err == nil && selectCondition(origo) {
		g.ApplyToCellAt(mod, x, y)
		g.ApplyToConnected(mod, selectCondition, left, y)
		g.ApplyToConnected(mod, selectCondition, right, y)
		g.ApplyToConnected(mod, selectCondition, x, up)
		g.ApplyToConnected(mod, selectCondition, x, down)
	} else {
		return
	}
}

// Change the type of cells on a straight line from (startX, startY) to (endX, endY). Line is calculated with Bresenham algorithm.
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

func (g *Grid) neighbourCount(x int, y int, celltype CellType) int {
	count := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			} else {
				cell, err := g.Get(x+i, y+j)
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
