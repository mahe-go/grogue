package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/mahe-go/grogue/grid"
	"io"
)

func main() {
	g := gridFuncs[currentGridFunc](80, 20)

	var err error
	gui := gocui.NewGui()
	if err := gui.Init(); err != nil {
		fmt.Errorf("Error initializing gui: %s", err)
		return
	}
	defer gui.Close()
	gui.SetLayout(gridGui(g))
	if err := gui.SetKeybinding("", rune('q'), 0, quit); err != nil {
		fmt.Errorf("%s", err)
		return
	}
	if err := gui.SetKeybinding("", rune('n'), 0, changeGrid()); err != nil {
		fmt.Errorf("%s", err)
		return
	}
	if err := gui.SetKeybinding("", rune('c'), 0, changeGridType()); err != nil {
		fmt.Errorf("%s", err)
		return
	}
	err = gui.MainLoop()
	if err != nil && err != gocui.ErrorQuit {
		fmt.Errorf("%s", err)
		return
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrorQuit
}

type LayoutFunc func(gui *gocui.Gui) error

type KeyHandlerFunc func(gui *gocui.Gui, view *gocui.View) error

func changeGrid() gocui.KeybindingHandler {
	return func(gui *gocui.Gui, view *gocui.View) error {
		gui.SetLayout(gridGui(gridFuncs[currentGridFunc](80, 20)))
		return nil
	}
}

func changeGridType() gocui.KeybindingHandler {
	return func(gui *gocui.Gui, view *gocui.View) error {
		currentGridFunc = (currentGridFunc + 1) % 2
		return nil
	}
}

type GridFunc func(width int, height int) *grid.Grid

func naturalGrid(width int, height int) *grid.Grid { return grid.NewNaturalCavernGrid(80, 20, 45, 2) }
func rectangularGrid(width int, height int) *grid.Grid {
	return grid.NewRectangularCavernGrid(width, height, 7, 7)
}

var gridFuncs [2]GridFunc = [2]GridFunc{naturalGrid, rectangularGrid}

var currentGridFunc int = 0

func gridGui(g *grid.Grid) LayoutFunc {
	return func(gui *gocui.Gui) error {
		maxX, _ := gui.Size()
		if mapView, err := gui.SetView("Map", maxX-g.Width-2, 0, maxX-1, g.Height+1); err != nil {
			if err != gocui.ErrorUnkView {
				return err
			}
			printGrid(mapView, g)
		}
		return gui.SetCurrentView("Map")
	}
}

func printGrid(w io.Writer, g *grid.Grid) {
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			b, _ := g.Get(x, y)
			fmt.Fprint(w, b)
		}
		fmt.Fprint(w, "\n")
	}

}
