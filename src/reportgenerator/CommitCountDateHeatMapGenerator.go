package reportgenerator

import (
	"fmt"
	"sort"
	"time"

	"github.com/go-git/go-git/v5/plumbing/object"
    "github.com/k1-end/git-visualizer/src/report"
)

type CommitCountDateHeatMapGenerator struct {
    CommitsMap map[string]int
}

func (r CommitCountDateHeatMapGenerator) IterationStep(c *object.Commit)  {
    year, month, date := c.Author.When.Local().Date()
    key := fmt.Sprintf("%d-%d-%d", year, month, date)
    _, exists := r.CommitsMap[key]
    if !exists {
        r.CommitsMap[key] = 1
    } else {
        r.CommitsMap[key]++
    }
}

func (rg CommitCountDateHeatMapGenerator) GetReport() report.Report {
    labels := make([]string, 0, len(rg.CommitsMap))

    for k := range rg.CommitsMap {
        labels = append(labels, k)
    }

    sort.Slice(labels, func(i, j int) bool {
        timeI, _ := time.Parse("2006-1-2", labels[i])
        timeJ, _ := time.Parse("2006-1-2", labels[j])

        if timeI.Before(timeJ) {
            return true
        } else {
            return false
        }
    })
    var data []report.Data
    for k := range labels {
        data = append(data, report.Data{IsInt: true, IntValue: rg.CommitsMap[labels[k]], StringValue: ""})
    }
    r := report.Report{}
    r.SetData(data)
    r.SetTitle("Commit count heat map")
    r.SetLabels(labels)
    r.SetReportType("date_heatmap")
    return r
}


