package grid

type CellType struct {
	Traversable bool
	Rune        rune
	Description string
}

var SOLID_ROCK = CellType{false, ' ', "solid rock"}
var WALL = CellType{false, '#', "wall"}
var ROOM = CellType{true, '.', "thin air"}
var CORRIDOR = CellType{true, '.', "corridor"}
var STAIRCASE_UP = CellType{true, '<', "staircase up"}
var STAIRCASE_DOWN = CellType{true, '>', "staircase down"}

type GridCell struct {
	Type    CellType
	Checked bool
}

func NewGridCellOfType(typ CellType) *GridCell {
	return &GridCell{typ, false}
}

func NewGridCellOfTypeValue(typ CellType) GridCell {
	return GridCell{typ, false}
}
