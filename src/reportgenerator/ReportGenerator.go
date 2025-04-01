package reportgenerator

import (
	"github.com/k1-end/git-visualizer/src/report"
)

type ReportGenerator interface {
    GetReport() report.Report
}

type Author struct {
    Name  string
    Emails map[string]bool
}

