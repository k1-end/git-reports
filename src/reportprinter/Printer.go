package reportprinter

import "github.com/k1-end/git-visualizer/src/report"

type Printer interface {
	RegisterReport(r report.Report)
    PrintAllReports()
}
