package report

import (
	"path/filepath"
	"sort"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/pterm/pterm"
)

type FileTypeReport struct {
    FileTypeMap  map[string]int
}

func (r FileTypeReport) Iterate(c *object.Commit)  {
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

func (r FileTypeReport) Print() {

    mimeTypes := make([]string, 0, len(r.FileTypeMap))

    for k := range r.FileTypeMap {
        mimeTypes = append(mimeTypes, k)
    }

    sort.SliceStable(mimeTypes, func(i, j int) bool {
        return r.FileTypeMap[mimeTypes[i]] > r.FileTypeMap[mimeTypes[j]]
    })

    var fileTypeData []pterm.Bar
    var fileType pterm.Bar
    for k := range mimeTypes {
        fileType.Label = mimeTypes[k]
        v := int(r.FileTypeMap[mimeTypes[k]] / 1000)
        if v == 0 {
            continue
        }
        fileType.Value = v
        fileTypeData = append(fileTypeData, fileType)
    }
    pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgGreen)).Println("File Types (KB)")
    _ = pterm.DefaultBarChart.WithShowValue().WithBars(fileTypeData).WithHorizontal().WithWidth(100).Render()
}
