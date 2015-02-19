package dungeon

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

func BSPDungeon(width int, height int, minNodeWidth int, minNodeHeight int) (*Grid, error) {
	rand.Seed(time.Now().Unix())
	root := split(newNode(nil, newRect(1, 1, width-1, height-1)), minNodeWidth, minNodeHeight)
	grid := NewSolidGrid(width, height)
	err := root.delveRoom(grid)
	if err == nil {
		err = root.connectPartsWithCorridor(grid)
	}
	return grid, err
}

func split(n *node, minNodeWidth int, minNodeHeight int) *node {
	r := n.Rect
	var width, height, width2, height2 int
	var x, y int
	direction := rand.Intn(2)
	if direction == 0 {
		splitLoc := r.Width/3 + rand.Intn(r.Width/3)
		width = splitLoc
		x = r.X + width
		width2 = r.Width - width
		height = r.Height
		height2 = r.Height
		y = r.Y
	} else {
		splitLoc := r.Height/3 + rand.Intn(r.Height/3)
		width = r.Width
		width2 = r.Width
		x = r.X
		height = splitLoc
		height2 = r.Height - height
		y = r.Y + height
	}

	if height > minNodeHeight && height2 > minNodeHeight && width > minNodeWidth && width2 > minNodeWidth {
		leftRect := newRect(r.X, r.Y, width, height)
		rightRect := newRect(x, y, width2, height2)
		n.Left = split(newNode(n, leftRect), minNodeHeight, minNodeWidth)
		n.Right = split(newNode(n, rightRect), minNodeHeight, minNodeWidth)
	}

	return n
}

func (n *node) delveRoom(grid *Grid) error {
	if n.isLeaf() {
		roomWidth := 2*n.Rect.Width/3 + rand.Intn(n.Rect.Width/3)
		roomHeight := 2*n.Rect.Height/3 + rand.Intn(n.Rect.Height/3)

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
				if y == roomY || y == roomHeight-1 {
					err = grid.Set(n.Rect.X+x, n.Rect.Y+y, *NewGridCellOfType(WALL))
				} else if x == roomX || x == roomWidth-1 {
					err = grid.Set(n.Rect.X+x, n.Rect.Y+y, *NewGridCellOfType(WALL))
				} else {
					err = grid.Set(n.Rect.X+x, n.Rect.Y+y, *NewGridCellOfType(ROOM))
				}
			}
		}
		return err
	} else {
		err := n.Left.delveRoom(grid)
		if err == nil {
			err = n.Right.delveRoom(grid)
		}
		return err
	}
}

func (n *node) connectPartsWithCorridor(grid *Grid) error {
	if !n.isLeaf() {
		err := n.Left.connectPartsWithCorridor(grid)
		if err == nil {
			n.Right.connectPartsWithCorridor(grid)
		}

		startx := n.Left.Rect.X + n.Left.Rect.Width/2
		endx := n.Right.Rect.X + n.Right.Rect.Width/2
		starty := n.Left.Rect.Y + n.Left.Rect.Height/2
		endy := n.Right.Rect.Y + n.Right.Rect.Height/2

		return grid.MarkLineBresenham(startx, starty, endx, endy, CORRIDOR)
	} else {
		return nil
	}
}
