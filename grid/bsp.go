package grid

import (
	"math/rand"
	"time"
)

type rect struct {
	X      int
	Y      int
	Width  int
	Height int
}

func newRect(x int, y int, width int, height int) *rect {
	return &rect{x, y, width, height}
}

type node struct {
	Parent *node
	Rect   *rect
	Left   *node
	Right  *node
}

func newNode(parent *node, r *rect) *node {
	return &node{parent, r, nil, nil}
}

func (n *node) isLeaf() bool {
	return n.Left == nil && n.Right == nil
}

func NewRectangularCavernGrid(width int, height int, minNodeWidth int, minNodeHeight int) *Grid {
	rand.Seed(time.Now().Unix())
	root := split(newNode(nil, newRect(1, 1, width-1, height-1)), minNodeWidth, minNodeHeight)
	grid := NewSolidGridOfType(width, height, SOLID_ROCK)
	root.delveRoom(grid)
	root.connectPartsWithCorridor(grid)

	grid.buildCavernWalls()

	return grid
}

func split(n *node, minNodeWidth int, minNodeHeight int) *node {
	r := n.Rect
	var width, height, width2, height2 int
	var x, y int

	verticalSplitPossible := r.Width > minNodeWidth*2
	horizontalSplitPossible := r.Height > minNodeHeight*2

	if !verticalSplitPossible && !horizontalSplitPossible {
		return n
	}

	direction := rand.Intn(2)
	if verticalSplitPossible && !horizontalSplitPossible {
		direction = 0
	} else if !verticalSplitPossible && horizontalSplitPossible {
		direction = 1
	}

	if direction == 0 {
		splitLoc := minNodeWidth + rand.Intn(r.Width-2*minNodeWidth)
		width = splitLoc
		x = r.X + width
		width2 = r.Width - width
		height = r.Height
		height2 = r.Height
		y = r.Y
	} else {
		splitLoc := minNodeHeight + rand.Intn(r.Height-2*minNodeHeight)
		width = r.Width
		width2 = r.Width
		x = r.X
		height = splitLoc
		height2 = r.Height - height
		y = r.Y + height
	}

	leftRect := newRect(r.X, r.Y, width, height)
	rightRect := newRect(x, y, width2, height2)

	n.Left = split(newNode(n, leftRect), minNodeHeight, minNodeWidth)
	n.Right = split(newNode(n, rightRect), minNodeHeight, minNodeWidth)

	return n
}

func (n *node) delveRoom(grid *Grid) {
	if n.isLeaf() {
		roomWidth := n.Rect.Width/2 + rand.Intn(n.Rect.Width/2)
		roomHeight := n.Rect.Height/2 + rand.Intn(n.Rect.Height/2)

		roomX := 0
		if roomWidth < n.Rect.Width {
			roomX = rand.Intn(n.Rect.Width - roomWidth)
		}

		roomY := 0
		if roomWidth < n.Rect.Height {
			roomY = rand.Intn(n.Rect.Height - roomHeight)
		}

		var err error
		for x := roomX; x < roomWidth && err == nil; x++ {
			for y := roomY; y < roomHeight && err == nil; y++ {
				err = grid.ApplyToCellAtXY(GridCellTypeConverter(ROOM), n.Rect.X+x, n.Rect.Y+y)
			}
		}
		return
	} else {
		n.Left.delveRoom(grid)
		n.Right.delveRoom(grid)
		return
	}
}

func (n *node) connectPartsWithCorridor(grid *Grid) {
	if !n.isLeaf() {
		n.Left.connectPartsWithCorridor(grid)
		n.Right.connectPartsWithCorridor(grid)

		startx := n.Left.Rect.X + n.Left.Rect.Width/2
		endx := n.Right.Rect.X + n.Right.Rect.Width/2
		starty := n.Left.Rect.Y + n.Left.Rect.Height/2
		endy := n.Right.Rect.Y + n.Right.Rect.Height/2

		grid.ApplyOnLine(GridCellTypeConverter(CORRIDOR), GridCellIsOfType(SOLID_ROCK), startx, starty, endx, endy)
	}
}
