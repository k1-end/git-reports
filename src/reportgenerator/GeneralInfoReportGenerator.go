package reportgenerator

import (

	"github.com/go-git/go-git/v5/plumbing/object"
    "github.com/k1-end/git-visualizer/src/report"
)

type GeneralInfoReportGenerator struct {
    ContributorsNo int
    CommitsNo int
    ProjectSize uint64
    FilesNo int

    contributors map[string]bool // email => true
}

func (r *GeneralInfoReportGenerator) LogIterationStep(c *object.Commit)  {
    if r.contributors == nil { // Check if the map is nil
        r.contributors = make(map[string]bool) // Initialize the map
    }
	if _, exists := r.contributors[c.Author.Email]; !exists {
        r.contributors[c.Author.Email] = true
        r.ContributorsNo += 1
    }
    r.CommitsNo += 1
}

func (r *GeneralInfoReportGenerator) FileIterationStep(f *object.File)  {
    r.FilesNo += 1
    r.ProjectSize += uint64(f.Size)
}

func (rg GeneralInfoReportGenerator) GetReport() report.Report {
    keys := []string{"Number of contributors", "Number of commits", "Project size", "Number of files"}


    var data []report.Data
    data = append(data, report.Data{IsInt: true, IntValue: rg.ContributorsNo})
    data = append(data, report.Data{IsInt: true, IntValue: rg.CommitsNo})
    data = append(data, report.Data{IsInt: true, IntValue: int(rg.ProjectSize)})
    data = append(data, report.Data{IsInt: true, IntValue: rg.FilesNo})

    r := report.Report{}
    r.SetLabels(keys)
    r.SetData(data)
    r.SetTitle("General Info")
    r.SetReportType("table")
    return r
}
