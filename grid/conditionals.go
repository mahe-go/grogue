package grid

// Function type for function testing a grid cell at (x,y)
// GridConditional makes it possible to create conditions
// based on the environment of the cell we're testing, since we know its location
type GridConditional func(grid *Grid, x int, y int) bool

// Function type for function testing a grid cell
type GridCellConditional func(cell GridCell) bool

func (self GridCellConditional) Or(cond GridCellConditional) GridCellConditional {
	return func(g GridCell) bool {
		return self(g) || cond(g)
	}
}

func (self GridCellConditional) And(cond GridCellConditional) GridCellConditional {
	return func(g GridCell) bool {
		return self(g) && cond(g)
	}
}

func (self GridCellConditional) Not() GridCellConditional {
	return func(g GridCell) bool {
		return !self(g)
	}
}

func (self GridConditional) Or(cond GridConditional) GridConditional {
	return func(g *Grid, x int, y int) bool {
		return self(g, x, y) || cond(g, x, y)
	}
}

func (self GridConditional) And(cond GridConditional) GridConditional {
	return func(g *Grid, x int, y int) bool {
		return self(g, x, y) && cond(g, x, y)
	}
}

func (self GridConditional) Not() GridConditional {
	return func(g *Grid, x int, y int) bool {
		return !self(g, x, y)
	}
}

func GridCellIsOfType(typ CellType) GridCellConditional {
	return func(g GridCell) bool {
		return g.Type == typ
	}
}

var GridCellIsChecked GridCellConditional = func(g GridCell) bool {
	return g.Checked
}
