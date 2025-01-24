package report

import (
	"sort"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/pterm/pterm"
)

type CommitsPerDevReport struct {
    CommitsPerDevMap map[string]int
}

func (r CommitsPerDevReport) IterationStep(c *object.Commit)  {
    _, exists := r.CommitsPerDevMap[c.Author.Name]
    if !exists {
        r.CommitsPerDevMap[c.Author.Name] = 1
    } else {
        r.CommitsPerDevMap[c.Author.Name]++
    }
}

func (r CommitsPerDevReport) Print() {
    authorNames := make([]string, 0, len(r.CommitsPerDevMap))

    for k := range r.CommitsPerDevMap {
        authorNames = append(authorNames, k)
    }

    sort.SliceStable(authorNames, func(i, j int) bool {
        return r.CommitsPerDevMap[authorNames[i]] > r.CommitsPerDevMap[authorNames[j]]
    })

    var barData []pterm.Bar
    var bar pterm.Bar
    for _, authorName := range authorNames {
        bar.Label = authorName
        bar.Value = r.CommitsPerDevMap[authorName]
        barData = append(barData, bar)
    }

    pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgGreen)).Println("Commit count per developer")
    _ = pterm.DefaultBarChart.WithBars(barData).WithHorizontal().WithWidth(90).WithShowValue().Render()
}
