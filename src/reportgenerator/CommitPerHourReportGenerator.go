package reportgenerator

import (
	"strconv"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/k1-end/git-visualizer/src/report"
)

type CommitsPerHourReportGenerator struct {
    CommitsPerHourMap []int
}

func (r CommitsPerHourReportGenerator) LogIterationStep(c *object.Commit)  {
	r.CommitsPerHourMap[c.Author.When.Local().Hour()]++
}

func (rg CommitsPerHourReportGenerator) GetReport() report.Report {
    var data []report.Data
    var labels []string
    for i := 1; i < 24; i++ {
        labels = append(labels, strconv.Itoa(i))
        data = append(data, report.Data{IsInt: true, IntValue: rg.CommitsPerHourMap[i]})
    }
    r := report.Report{}
    r.SetData(data)
    r.SetLabels(labels)
    r.SetTitle("Commits per hour of day")
    r.SetReportType("bar_chart")
    return r

}
