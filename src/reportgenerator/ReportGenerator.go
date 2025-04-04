package reportgenerator

import (
	"github.com/k1-end/git-reports/src/report"
)

type ReportGenerator interface {
    GetReport() report.Report
}

type Author struct {
    Name  string
    Emails map[string]bool
}

