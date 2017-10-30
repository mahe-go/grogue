package gui

import (
	"math/rand"

	"github.com/jroimartin/gocui"
	"github.com/mahe-go/grogue/creature"
	"github.com/mahe-go/grogue/grid"
)

func PlayerMovementHandler(g *grid.Grid, player *creature.Player, direction grid.Direction) gocui.KeybindingHandler {
	return func(gcui *gocui.Gui, v *gocui.View) error {
		player.Move(g, direction)
		Layout(g, player, gcui)
		return nil
	}
}

func StaircaseUpHandler(g *grid.Grid, player *creature.Player) gocui.KeybindingHandler {
	return func(gcui *gocui.Gui, v *gocui.View) error {
		if g.TestCellAtXY(grid.GridCellIsOfType(grid.STAIRCASE_UP), player.X, player.Y) {
			if rand.Intn(2) == 0 {
				*g = *grid.NewNaturalCavernGrid(g.Width, g.Height, 45, 2)
			} else {
				*g = *grid.NewRectangularCavernGrid(g.Width, g.Height, 7, 7)
			}

			creature.PlacePlayerToGridAtMatching(player, g, grid.GridCellIsOfType(grid.STAIRCASE_DOWN))
			Layout(g, player, gcui)
			return nil
		} else {
			return nil
		}
	}
}

func StaircaseDownHandler(g *grid.Grid, player *creature.Player) gocui.KeybindingHandler {
	return func(gcui *gocui.Gui, v *gocui.View) error {
		if g.TestCellAtXY(grid.GridCellIsOfType(grid.STAIRCASE_DOWN), player.X, player.Y) {
			if rand.Intn(2) == 0 {
				*g = *grid.NewNaturalCavernGrid(g.Width, g.Height, 45, 2)
			} else {
				*g = *grid.NewRectangularCavernGrid(g.Width, g.Height, 7, 7)
			}

			creature.PlacePlayerToGridAtMatching(player, g, grid.GridCellIsOfType(grid.STAIRCASE_UP))
			Layout(g, player, gcui)
			return nil
		} else {
			return nil
		}
	}
}
