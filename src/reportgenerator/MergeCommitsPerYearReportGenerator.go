package reportgenerator

import (
	"sort"
	"strconv"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/k1-end/git-reports/src/report"
)

type MergeCommitsPerYearReportGenerator struct {
    MergeCommitsPerYearMap map[int]int
}

func (r MergeCommitsPerYearReportGenerator) LogIterationStep(c *object.Commit, a Author)  {
    year, _, _ := c.Author.When.Local().Date()
    if c.NumParents() > 1 {
        if _, exists := r.MergeCommitsPerYearMap[year]; !exists {
            r.MergeCommitsPerYearMap[year] = 1
        } else {
            r.MergeCommitsPerYearMap[year]++
        }
    }
}

func (rg MergeCommitsPerYearReportGenerator) GetReport() report.Report {
    yearsKey := make([]int, 0, len(rg.MergeCommitsPerYearMap))
    for k := range rg.MergeCommitsPerYearMap {
        yearsKey = append(yearsKey, k)
    }
    sort.Ints(yearsKey)
    var data []report.Data
    var labels []string
    for y := range yearsKey{
        labels = append(labels, strconv.Itoa(yearsKey[y]))
        data = append(data, report.Data{IsInt: true, IntValue: rg.MergeCommitsPerYearMap[yearsKey[y]]})
    }
    r := report.Report{}
    r.SetTitle("Merge Commits per year")
    r.SetLabels(labels)
    r.SetData(data)
    r.SetReportType("bar_chart")
    return r
}
