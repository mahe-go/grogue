package grid

// Function type for function modifying a grid cell. Returns the modified cell.
type CellModification func(cell GridCell) GridCell

// Function type for function modifying a grid cell at (x,y) (or region of grid starting from (x,y) as well)
type Modification func(grid *Grid, x int, y int) error

func (self CellModification) And(mod CellModification) CellModification {
	return func(g GridCell) GridCell {
		return mod(self(g))
	}
}

func GridCellTypeConverter(typ CellType) CellModification {
	return func(g GridCell) GridCell {
		g.Type = typ
		return g
	}
}

var GridCellChecker CellModification = func(g GridCell) GridCell {
	g.Checked = true
	return g
}

var GridCellUnChecker CellModification = func(g GridCell) GridCell {
	g.Checked = false
	return g
}
