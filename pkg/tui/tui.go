package tui

import (
	"fmt"
	"log"
	"text/tabwriter"

	"github.com/augmentable-dev/askgit/pkg/gitqlite"
	"github.com/jroimartin/gocui"
)

var (
	active   = 0
	query    = ""
	repoPath = ""
	usrInpt  = ""
)

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("Query", 0, 0, maxX/2-1, maxY*2/10); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Query"
		v.Editable = true
		v.Wrap = true
		fmt.Fprint(v, query)
		if _, err = SetCurrentViewOnTop(g, "Query"); err != nil {
			return err
		}

	}
	if v, err := g.SetView("Keybinds", 0, maxY*2/10+1, maxX/2-1, maxY*4/10); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Keybinds"
		w := tabwriter.NewWriter(v, 0, 0, 1, ' ', 0)

		fmt.Fprint(w, "Ctrl+C\t exit \nCtrl+E\t execute query \nCtrl+Q\t clear query box\nDefault L-click \t select a default to be displayed in the query view\n\n")

	}
	if v, err := g.SetView("Info", maxX/2, maxY*2/10+1, maxX-1, maxY*4/10); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Info"
		git, err := gitqlite.New(repoPath, &gitqlite.Options{})
		if err != nil {
			return err
		}
		err = DisplayInformation(g, git, 0)
		if err != nil {
			return err
		}

	}
	if v, err := g.SetView("Output", 0, maxY*4/10+1, maxX, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Output"
		v.Wrap = false

	}
	if v, err := g.SetView("Default", maxX/2, 0, maxX-1, maxY*2/10); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Preset Queries"
		for i := range Queries {
			fmt.Fprintf(v, "%s\n", i)
		}

	}
	return nil
}
func test(g *gocui.Gui, v *gocui.View) error {
	//for use with testing uses ctrl+t
	return nil
}
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
func RunGUI(repo string, directory string, q string) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	query = q
	repoPath = directory
	usrInpt = repo
	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen
	g.Mouse = true

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, NextView); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlQ, gocui.ModNone, ClearQuery); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.MouseLeft, gocui.ModNone, HandleClick); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlE, gocui.ModNone, RunQuery); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.MouseRelease, gocui.ModNone, HandleCursor); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.MouseWheelUp, gocui.ModNone, PreviousLine); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.MouseWheelDown, gocui.ModNone, NextLine); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, PreviousLine); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, NextLine); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("Output", gocui.KeyArrowRight, gocui.ModNone, GoRight); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("Output", gocui.KeyArrowLeft, gocui.ModNone, GoLeft); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlT, gocui.ModNone, test); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
