package dungeon

import (
	"errors"
	"github.com/mahe-go/grogue/util"
)

type CellType uint

const CELLTYPE_CHARS = " #.*"

const SOLID_ROCK CellType = 0
const WALL CellType = 1
const ROOM CellType = 2
const CORRIDOR CellType = 3

type GridCell struct {
	Typ CellType
}

func NewGridCell(typ CellType) *GridCell {
	return &GridCell{typ}
}

func (g GridCell) String() string {
	return string([]rune(CELLTYPE_CHARS)[g.Typ])
}

var GRID_OVERFLOW error = errors.New("Grid overflow")

type CorridorFunc func(grid *Grid) error

type RoomFunc func(grid *Grid) error

type Grid struct {
	cells         []GridCell
	Width         int
	Height        int
	DelveRoom     RoomFunc
	DelveCorridor CorridorFunc
}

func NewGrid(width int, height int, rf RoomFunc, cf CorridorFunc) *Grid {
	return &Grid{make([]GridCell, width*height), width, height, rf, cf}
}

func (g *Grid) Set(x int, y int, value GridCell) error {
	size := g.Width * g.Height
	index := g.Width*y + x
	if index >= size || index < 0 {
		return GRID_OVERFLOW
	} else {
		g.cells[index] = value
		return nil
	}
}

func (g *Grid) Get(x int, y int) (*GridCell, error) {
	size := g.Width * g.Height
	index := g.Width*y + x
	if index > size || index < 0 {
		return nil, GRID_OVERFLOW
	} else {
		return &(g.cells[index]), nil
	}
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
		} else if old != nil {
			if old.Typ == WALL || old.Typ == SOLID_ROCK {
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
