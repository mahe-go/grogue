package main

import (
	"log"

	"github.com/jroimartin/gocui"
	"github.com/mahe-go/grogue/creature"
	"github.com/mahe-go/grogue/grid"
	"github.com/mahe-go/grogue/gui"
)

func main() {
	gcui := gocui.NewGui()
	if err := gcui.Init(); err != nil {
		log.Panicln(err)
	}
	defer gcui.Close()

	currentGrid, player := naturalGrid(80, 20)

	gui.Layout(currentGrid, player, gcui)

	if err := gcui.SetKeybinding("", rune('q'), 0, quit); err != nil {
		log.Panicln(err)
	}

	if err := gcui.SetKeybinding("Map", rune('s'), 0, gui.PlayerMovementHandler(currentGrid, player, grid.South)); err != nil {
		log.Panicln(err)
	}
	if err := gcui.SetKeybinding("Map", rune('w'), 0, gui.PlayerMovementHandler(currentGrid, player, grid.North)); err != nil {
		log.Panicln(err)
	}
	if err := gcui.SetKeybinding("Map", rune('d'), 0, gui.PlayerMovementHandler(currentGrid, player, grid.East)); err != nil {
		log.Panicln(err)
	}
	if err := gcui.SetKeybinding("Map", rune('a'), 0, gui.PlayerMovementHandler(currentGrid, player, grid.West)); err != nil {
		log.Panicln(err)
	}

	if err := gcui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func placePlayerToGrid(player *creature.Player, g *grid.Grid) {
	x := g.Width / 2
	y := g.Height / 2
	for !g.TestCellAtXY(grid.CellIsTraversable, x, y) {
		x++
		y++
	}
	player.SetLocation(x, y)
}

func naturalGrid(width int, height int) (*grid.Grid, *creature.Player) {
	g := grid.NewNaturalCavernGrid(width, height, 45, 2)
	p := creature.NewPlayer("Mahe", creature.NewSpecies(1, '@'))
	placePlayerToGrid(p, g)
	return g, p
}
