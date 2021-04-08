package ui

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/pkg/browser"
	"google-tui/search"
)

type Results struct {
	List    *widgets.List
	Search  *widgets.Paragraph
	Results []*search.Result
	focused bool
}

func NewResults(list *widgets.List, search *widgets.Paragraph) *Results {
	return &Results{
		List:    list,
		Search:  search,
		focused: false,
	}
}

func (results *Results) AppendResult(result *search.Result) {
	results.Results = append(results.Results, result)
}

func (results *Results) GetStyle(d bool) ui.Style {
	if d {
		return ui.NewStyle(ui.ColorWhite)
	}
	return ui.NewStyle(ui.ColorCyan)
}

func (results *Results) SetFocused(focused bool) {
	results.focused = focused

	results.List.SelectedRowStyle = results.GetStyle(results.focused)
	results.Search.TextStyle = results.GetStyle(!results.focused)
}

func (results *Results) Up() {
	if results.List.SelectedRow - 1 < 0 && len(results.List.Rows) == 0 {
		return
	}

	if results.List.SelectedRow == 0 {
		results.SetFocused(false)
		return
	}

	results.List.ScrollUp()
}

func (results *Results) Down() {
	if results.List.SelectedRow + 1 > len(results.List.Rows) - 1 {
		return
	}

	if !results.focused {
		results.SetFocused(true)
		return
	}

	results.List.ScrollDown()
}

func (results *Results) OpenResult() {
	url := results.Results[results.List.SelectedRow].URL

	_ = browser.OpenURL(url)
}