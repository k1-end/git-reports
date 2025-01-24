package report

import (
	"sort"
	"strconv"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/pterm/pterm"
)

type MergeCommitsPerYearReport struct {
    MergeCommitsPerYearMap map[int]int
}

func (r MergeCommitsPerYearReport) IterationStep(c *object.Commit)  {
    year, _, _ := c.Author.When.Local().Date()
    if c.NumParents() > 1 {
        if _, exists := r.MergeCommitsPerYearMap[year]; !exists {
            r.MergeCommitsPerYearMap[year] = 1
        } else {
            r.MergeCommitsPerYearMap[year]++
        }
    }
}

func (r MergeCommitsPerYearReport) Print() {
        yearsKey := make([]int, 0, len(r.MergeCommitsPerYearMap))
		for k := range r.MergeCommitsPerYearMap {
			yearsKey = append(yearsKey, k)
		}
		sort.Ints(yearsKey)
		var mergeCommitsPerYearBar []pterm.Bar
		var mergeCommit pterm.Bar
        for y := range yearsKey{
			mergeCommit.Label = strconv.Itoa(yearsKey[y])
			mergeCommit.Value = r.MergeCommitsPerYearMap[yearsKey[y]]
			mergeCommitsPerYearBar = append(mergeCommitsPerYearBar, mergeCommit)
		}
		pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgGreen)).Println("Merge Commits per year")
		_ = pterm.DefaultBarChart.WithShowValue().WithBars(mergeCommitsPerYearBar).WithHorizontal().WithWidth(100).Render()
}
