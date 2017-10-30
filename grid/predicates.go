package grid

// Function type for function testing a grid cell at (x,y)
// LocationPredicate makes it possible to create conditions
// based on the environment of the cell we're testing, since we know its location
type LocationPredicate func(grid *Grid, x int, y int) bool

// Function type for function testing a grid cell
type CellPredicate func(cell GridCell) bool

func (self CellPredicate) Or(cond CellPredicate) CellPredicate {
	return func(g GridCell) bool {
		return self(g) || cond(g)
	}
}

func (self CellPredicate) And(cond CellPredicate) CellPredicate {
	return func(g GridCell) bool {
		return self(g) && cond(g)
	}
}

func (self CellPredicate) Not() CellPredicate {
	return func(g GridCell) bool {
		return !self(g)
	}
}

func (self LocationPredicate) Or(cond LocationPredicate) LocationPredicate {
	return func(g *Grid, x int, y int) bool {
		return self(g, x, y) || cond(g, x, y)
	}
}

func (self LocationPredicate) And(cond LocationPredicate) LocationPredicate {
	return func(g *Grid, x int, y int) bool {
		return self(g, x, y) && cond(g, x, y)
	}
}

func (self LocationPredicate) Not() LocationPredicate {
	return func(g *Grid, x int, y int) bool {
		return !self(g, x, y)
	}
}

func GridCellIsOfType(typ CellType) CellPredicate {
	return func(g GridCell) bool {
		return g.Type == typ
	}
}

var GridCellIsChecked CellPredicate = func(g GridCell) bool {
	return g.Checked
}

var CellIsTraversable CellPredicate = func(c GridCell) bool {
	return c.Type.Traversable
}