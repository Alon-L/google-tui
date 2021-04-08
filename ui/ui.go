package ui

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"google-tui/search"
	"log"
)

func submitSearch(grid *ui.Grid, results *Results, text string) {
	results.List.Rows = nil
	results.Results = nil

	searcher := search.NewSearcher(text, &results.List.Rows)
	err, r := searcher.Search()

	if err != nil {
		log.Fatal(err)
	}

	results.SetFocused(true)

	for result := range r {
		results.AppendResult(result)
		searcher.AppendResult(result)
	}

	ui.Clear()
	ui.Render(grid)
}

func InitUI() {
	if err := ui.Init(); err != nil {
		log.Fatalf("Failed to initialize UI: %v", err)
	}
	defer ui.Close()

	p := widgets.NewParagraph()
	p.Title = "Search"
	p.TextStyle = ui.NewStyle(ui.ColorCyan)

	l := widgets.NewList()
	l.Rows = []string{}
	l.TextStyle = ui.NewStyle(ui.ColorCyan)

	results := NewResults(l, p)

	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	grid.Set(
		ui.NewRow(1.0/8, p),
		ui.NewRow(1.0-1.0/8, l),
	)

	ui.Render(grid)

	typeEvent := func(char string) bool {
		defer ui.Clear()
		defer ui.Render(grid)

		if char == "<C-c>" {
			return false
		}

		pLen := len(p.Text)

		if char == "<Up>" {
			results.Up()
		}

		if char == "<Down>" {
			results.Down()
		}

		if !results.focused {
			if char == "<Delete>" || char == "<Backspace>" {
				if pLen > 0 {
					p.Text = p.Text[:pLen-1]
				}
			}

			if char == "<Space>" {
				p.Text += " "
			}

			if char == "<Enter>" {
				go submitSearch(grid, results, p.Text)
			}

			if len(char) <= 1 {
				p.Text += char
			}
		} else {
			if char == "<Enter>" {
				go results.OpenResult()
			}
		}

		return true
	}

	uiEvents := ui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			switch e.Type {
			case ui.KeyboardEvent:
				char := e.ID
				if !typeEvent(char) {
					return
				}
			case ui.ResizeEvent:
				payload := e.Payload.(ui.Resize)
				grid.SetRect(0, 0, payload.Width, payload.Height)
				ui.Clear()
				ui.Render(grid)
			}
		}
	}
}
