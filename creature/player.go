package creature

import (
	"errors"

	"github.com/mahe-go/grogue/grid"
)

var CANNOT_MOVE_THERE = errors.New("Not accessible")

type Player struct {
	X    int
	Y    int
	Name string
	*Species
}

func NewPlayer(name string, species *Species) *Player {
	return &Player{0, 0, name, species}
}

func (p *Player) SetLocation(x int, y int) {
	p.X = x
	p.Y = y
}

func (p *Player) MoveOne(g *grid.Grid, direction grid.Direction) error {
	tx := p.X + direction.Dx
	ty := p.Y + direction.Dy

	if g.TestCellAtXY(func(cell grid.GridCell) bool {
		return cell.Type.Traversable
	}, tx, ty) {
		p.X = tx
		p.Y = ty
		return nil
	} else {
		return CANNOT_MOVE_THERE
	}
}

func (p *Player) Move(g *grid.Grid, direction grid.Direction) error {
	var err error
	for i := 0; i < p.Movement && err == nil; i++ {
		err = p.MoveOne(g, direction)
	}
	return err
}
