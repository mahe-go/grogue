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
	if err != nil && err != gocui.Quit {
		fmt.Errorf("%s", err)
		return
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.Quit
}

type LayoutFunc func(gui *gocui.Gui) error

type KeyHandlerFunc func(gui *gocui.Gui, view *gocui.View) error

func changeGrid() gocui.KeybindingHandler {
	return func(gui *gocui.Gui, view *gocui.View) error {
		gui.SetLayout(gridGui(gridFuncs[currentGridFunc](80, 20)))
		return gui.Flush()
	}
}

func changeGridType() gocui.KeybindingHandler {
	return func(gui *gocui.Gui, view *gocui.View) error {
		currentGridFunc = (currentGridFunc + 1) % len(gridFuncs)
		return changeGrid()(gui, view)
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
		if mapView, err := gui.SetView("Map", 0, 0, g.Width+1, g.Height+1); err != nil {
			if err != gocui.ErrorUnkView {
				return err
			}
			mapView.Clear()
			err = printGrid(mapView, g)
			if err != nil {
				return err
			}
		}
		return gui.SetCurrentView("Map")
	}
}

func printGrid(w io.Writer, g *grid.Grid) (err error) {
	for y := 0; y < g.Height; y++ {
		line := make([]rune, g.Width)
		line[0] = '?'
		for x := 0; x < g.Width; x++ {
			var cell grid.GridCell
			cell, err = g.Get(x, y)
			if err != nil {
				return
			}
			line[x] = cell.Type.Rune
		}
		_, err = fmt.Fprintln(w, string(line))
		if err != nil {
			return
		}
	}
	return
}
