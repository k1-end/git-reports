package report

import (
	"strconv"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/pterm/pterm"
)

type CommitsPerHourReport struct {
    CommitsPerHourMap []int
}

func (r CommitsPerHourReport) IterationStep(c *object.Commit)  {
	r.CommitsPerHourMap[c.Author.When.Local().Hour()]++
}

func (r CommitsPerHourReport) Print() {
    var hourData []pterm.Bar
    var hour pterm.Bar
    for i := 1; i < 24; i++ {
        hour.Label = strconv.Itoa(i)
        hour.Value = r.CommitsPerHourMap[i]
        hourData = append(hourData, hour)
    }
    pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgGreen)).Println("Commits per hour of day (local)")
    _ = pterm.DefaultBarChart.WithShowValue().WithBars(hourData).WithHorizontal().WithWidth(100).Render()
}
