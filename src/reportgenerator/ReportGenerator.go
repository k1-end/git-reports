package reportgenerator

import (
	"github.com/k1-end/git-visualizer/src/report"
)

type ReportGenerator interface {
    GetReport() report.Report
}
