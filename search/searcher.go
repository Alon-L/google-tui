package search

import (
	"fmt"
	"github.com/gocolly/colly"
	"net/url"
	"strconv"
)

const (
	ResultsNum = 25
	GoogleURL = "https://www.google.com/search"
)

type Result struct {
	Title string
	URL   string
}

type Searcher struct {
	Text    string
	PageNum int
	Results *[]string
}

func NewSearcher(text string, results *[]string) *Searcher {
	return &Searcher{
		Text:    text,
		PageNum: 1,
		Results: results,
	}
}

func (searcher *Searcher) GetURL() string {
	params := url.Values{}

	params.Add("q", searcher.Text)
	params.Add("num", strconv.Itoa(ResultsNum))
	params.Add("start", strconv.Itoa((searcher.PageNum-1)*ResultsNum))

	return fmt.Sprintf("%s?%s", GoogleURL, params.Encode())
}

func (searcher *Searcher) AppendResult(result *Result) {
	*searcher.Results = append(*searcher.Results, searcher.DecorateResult(result))
}

func (searcher *Searcher) DecorateResult(result *Result) string {
	return fmt.Sprintf("%d. %s", len(*searcher.Results) + 1, result.Title)
}

func (searcher *Searcher) Search() (error, <-chan *Result) {
	results := make(chan *Result, ResultsNum)

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:65.0) Gecko/20100101 Firefox/65.0"))

	c.OnHTML("[class=\"g\"]", func(e *colly.HTMLElement) {
		result := &Result{}

		result.Title = e.ChildText("h3")
		result.URL = e.ChildAttr("a", "href")

		results <- result
	})

	c.OnScraped(func(response *colly.Response) {
		close(results)
	})

	err := c.Visit(searcher.GetURL())
	if err != nil {
		return err, nil
	}

	return nil, results
}

func (searcher *Searcher) IncPage() {
	searcher.PageNum++
}

func (searcher *Searcher) DecPage() {
	searcher.PageNum--
}

// http://go-colly.org/docs/introduction/start/
