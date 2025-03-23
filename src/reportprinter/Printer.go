package reportprinter

import "github.com/k1-end/git-visualizer/src/report"

type Printer interface {
	RegisterReport(r report.Report)
	Print()
	SetProjectTitle(s string)
}

type BasePrinter struct {
	reports      []report.Report
	projectTitle string
}

func (p *BasePrinter) RegisterReport(r report.Report) {
	p.reports = append(p.reports, r)
}

func (p *BasePrinter) SetProjectTitle(s string) {
	p.projectTitle = s
}

func (p *BasePrinter) GetProjectTitle() string {
	return p.projectTitle
}

func (p *BasePrinter) GetReports() []report.Report {
	return p.reports
}
