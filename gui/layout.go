package gui

import (
	"fmt"

	"github.com/jroimartin/gocui"
	"github.com/mahe-go/grogue/creature"
	"github.com/mahe-go/grogue/grid"
)

func Layout(g *grid.Grid, player *creature.Player, gui *gocui.Gui) {
	gui.SetLayout(func(gui *gocui.Gui) error {
		if mapView, err := gui.SetView("Map", 0, 0, g.Width+1, g.Height+1); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			mapView.Clear()
			_, err = fmt.Fprint(mapView, *g)
			if err != nil {
				return err
			}
			mapView.Overwrite = true
			mapView.SetCursor(player.X, player.Y)
			mapView.EditWrite(player.Rune)
			mapView.Overwrite = false
		}
		return gui.SetCurrentView("Map")
	})
}
