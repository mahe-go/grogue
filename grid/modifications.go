package grid

// Function type for function modifying a grid cell. Returns the modified cell.
type GridCellModification func(cell GridCell) GridCell

// Function type for function modifying a grid cell at (x,y) (or region of grid starting from (x,y) as well)
type GridModification func(grid *Grid, x int, y int) error

func (self GridCellModification) And(mod GridCellModification) GridCellModification {
	return func(g GridCell) GridCell {
		return mod(self(g))
	}
}

func GridCellTypeConverter(typ CellType) GridCellModification {
	return func(g GridCell) GridCell {
		g.Type = typ
		return g
	}
}

var GridCellChecker GridCellModification = func(g GridCell) GridCell {
	g.Checked = true
	return g
}

var GridCellUnChecker GridCellModification = func(g GridCell) GridCell {
	g.Checked = false
	return g
}
