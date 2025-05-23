package reportgenerator

import (
	"sort"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/k1-end/git-reports/src/report"
)

type CommitsPerDevReportGenerator struct {
    CommitsPerDevMap map[string]int
}

func (r CommitsPerDevReportGenerator) LogIterationStep(c *object.Commit, a Author)  {
    _, exists := r.CommitsPerDevMap[a.Name]
    if !exists {
        r.CommitsPerDevMap[a.Name] = 1
    } else {
        r.CommitsPerDevMap[a.Name]++
    }
}

func (rg CommitsPerDevReportGenerator) GetReport() report.Report {
    keys := make([]string, 0, len(rg.CommitsPerDevMap))

    for k := range rg.CommitsPerDevMap {
        keys = append(keys, k)
    }

    sort.SliceStable(keys, func(i, j int) bool {
        return rg.CommitsPerDevMap[keys[i]] > rg.CommitsPerDevMap[keys[j]]
    })
    var data []report.Data
    for k := range keys {
        data = append(data, report.Data{IsInt: true, IntValue: rg.CommitsPerDevMap[keys[k]], StringValue: ""})
    }

    r := report.Report{}
    r.SetLabels(keys)
    r.SetData(data)
    r.SetTitle("Commits per developer")
    r.SetReportType("bar_chart")
    return r
}
