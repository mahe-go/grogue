package dungeon

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
				grid.Set(x, y, *NewGridCell(solidCellType))
			} else {
				grid.Set(x, y, *NewGridCell(emptyCellType))
			}
		}
	}
	return grid
}

type GridCell struct {
	Type CellType
}

type GridCellConditional func(cell GridCell) bool
type GridCellModification func(cell GridCell) (GridCell, error)

func NewGridCell(typ CellType) *GridCell {
	return &GridCell{typ}
}

func (g GridCell) String() string {
	return string([]rune(CELLTYPE_CHARS)[g.Type])
}

func (g *Grid) cellIndex(x int, y int) (int, error) {
	size := g.Width * g.Height
	index := g.Width*y + x
	if index >= size || index < 0 {
		return -1, GRID_OVERFLOW
	} else {
		return index, nil
	}
}

func (g *Grid) Set(x int, y int, value GridCell) (err error) {
	index, err := g.cellIndex(x, y)
	if err == nil {
		g.cells[index] = value
	}
	return
}

func (g *Grid) Apply(mod GridCellModification, x int, y int) (err error) {
	index, err := g.cellIndex(x, y)
	if err == nil {
		g.cells[index], err = mod(g.cells[index])
	}
	return
}

func (g *Grid) Get(x int, y int) (cell GridCell, err error) {
	index, err := g.cellIndex(x, y)
	if err == nil {
		cell = g.cells[index]
	}
	return
}

func (g *Grid) fillRow(y int, cellType CellType) (err error) {
	if y < 0 || y >= g.Height {
		err = GRID_OVERFLOW
	} else {
		for i := 0; i < g.Width && err == nil; i++ {
			err = g.Set(i, y, GridCell{cellType})
		}
	}
	return
}

func (g *Grid) fillColumn(x int, cellType CellType) (err error) {
	if x < 0 || x >= g.Width {
		err = GRID_OVERFLOW
	} else {
		for i := 0; i < g.Height && err == nil; i++ {
			err = g.Set(x, i, GridCell{cellType})
		}
	}
	return
}

func (g *Grid) convertAll(from CellType, to CellType) {
	for i, _ := range g.cells {
		if g.cells[i].Type == from {
			g.cells[i] = GridCell{to}
		}
	}
}

func (g *Grid) convertAllMatching(condition GridCellConditional, to CellType) (err error) {
	for x := 0; x < g.Width && err == nil; x++ {
		for y := 0; y < g.Height && err == nil; y++ {
			err = g.convertSingleMatching(condition, x, y, to)
		}
	}
	return
}

func (g *Grid) convertSingleMatching(condition GridCellConditional, x int, y int, to CellType) (err error) {
	var cell GridCell
	cell, err = g.Get(x, y)
	if condition(cell) {
		err = g.Set(x, y, GridCell{to})
	}
	return
}

func (g *Grid) floodFill(x int, y int, cellType CellType) {
	left := -1
	right := 1
	up := -1
	down := 1
	cell, err := g.Get(x, y)
	if err != nil {
		return
	}

	origType := cell.Type
	if cell, err = g.Get(x+left, y); err == nil && cell.Type == origType {
		g.floodFill(x+left, y, cellType)
	}
	if cell, err = g.Get(x+right, y); err == nil && cell.Type == origType {
		g.floodFill(x+right, y, cellType)
	}
	if cell, err = g.Get(x, y+up); err == nil && cell.Type == origType {
		g.floodFill(x, y+up, cellType)
	}
	if cell, err = g.Get(x, y+down); err == nil && cell.Type == origType {
		g.floodFill(x, y+down, cellType)
	}
	_ = g.Set(x, y, GridCell{cellType})
	return
}

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
				err = grid.Set(cx, cy, *NewGridCell(cellType))
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
