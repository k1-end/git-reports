package reportgenerator

import (
	"path/filepath"
	"sort"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/k1-end/git-visualizer/src/report"
)

type FileTypeReportGenerator struct {
    FileTypeMap  map[string]int
}

func (r FileTypeReportGenerator) Iterate(c *object.Commit)  {
    fIter, _ := c.Files()
    fIter.ForEach(func(f *object.File) error {
        mtype := filepath.Ext(f.Name)
        if _, exists := r.FileTypeMap[mtype]; !exists {
            r.FileTypeMap[mtype] = int(f.Size)
        } else {
            r.FileTypeMap[mtype] = r.FileTypeMap[mtype] + int(f.Size)
        }
        return nil
    })
}

func (rg FileTypeReportGenerator) GetReport() report.Report {

    mimeTypes := make([]string, 0, len(rg.FileTypeMap))

    for k := range rg.FileTypeMap {
        mimeTypes = append(mimeTypes, k)
    }

    sort.SliceStable(mimeTypes, func(i, j int) bool {
        return rg.FileTypeMap[mimeTypes[i]] > rg.FileTypeMap[mimeTypes[j]]
    })

    var data []report.Data
    var labels []string
    for k := range mimeTypes {
        v := int(rg.FileTypeMap[mimeTypes[k]] / 1000)
        if v == 0 {
            continue
        }
        labels = append(labels, mimeTypes[k])
        data = append(data, report.Data{IsInt: true, IntValue: v})
    }
    r := report.Report{}
    r.SetTitle("File Types (KB)")
    r.SetData(data)
    r.SetLabels(labels)
    return r
}
