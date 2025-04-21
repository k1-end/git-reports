package reportprinter

import (
	"testing"

	"github.com/k1-end/git-reports/src/report"
	"github.com/stretchr/testify/assert"
)

func TestBasePrinter_RegisterReport(t *testing.T) {
	// Test for BasePrinter.RegisterReport
	printer := BasePrinter{}
	testReport := report.Report{}

	// Register the report
	printer.RegisterReport(testReport)

	// Assert that the report is added to the reports slice
	assert.Len(t, printer.reports, 1, "Report should be added to the reports slice")
	assert.Equal(t, testReport, printer.reports[0], "The registered report should be the same as the added report")

	// Register another report
	anotherReport := report.Report{}
	printer.RegisterReport(anotherReport)

	// Assert that the new report is also added
	assert.Len(t, printer.reports, 2, "Another report should be added to the slice")
	assert.Equal(t, anotherReport, printer.reports[1], "The second registered report should be in the slice")
}

func TestBasePrinter_SetProjectTitle(t *testing.T) {
	// Test for BasePrinter.SetProjectTitle
	printer := BasePrinter{}
	title := "My Project"

	// Set the project title
	printer.SetProjectTitle(title)

	// Assert that the project title is set correctly
	assert.Equal(t, title, printer.projectTitle, "Project title should be set correctly")

	// Set a different title
	newTitle := "Another Project"
	printer.SetProjectTitle(newTitle)

	// Assert that the project title is updated
	assert.Equal(t, newTitle, printer.projectTitle, "Project title should be updated")
}

func TestBasePrinter_GetProjectTitle(t *testing.T) {
	// Test for BasePrinter.GetProjectTitle
	printer := BasePrinter{}
	title := "Sample Project"
	printer.projectTitle = title // Directly set for this test

	// Get the project title
	returnedTitle := printer.GetProjectTitle()

	// Assert that the returned title is correct
	assert.Equal(t, title, returnedTitle, "GetProjectTitle should return the correct title")

	// Test with an empty title
	printer.projectTitle = ""
	returnedTitle = printer.GetProjectTitle()
	assert.Empty(t, returnedTitle, "GetProjectTitle should return an empty string for an empty title")
}

func TestBasePrinter_GetReports(t *testing.T) {
	// Test for BasePrinter.GetReports
	printer := BasePrinter{}
	report1 := report.Report{}
    report1.SetTitle("Report 1")
	report2 := report.Report{}
    report2.SetTitle("Report 2")
	printer.reports = []report.Report{report1, report2} // Directly set for this test

	// Get the reports
	reports := printer.GetReports()

	// Assert that the returned reports are correct
	assert.Len(t, reports, 2, "GetReports should return the correct number of reports")
	assert.Equal(t, report1, reports[0], "GetReports should return the correct first report")
	assert.Equal(t, report2, reports[1], "GetReports should return the correct second report")

	// Test with no reports
	printer.reports = []report.Report{}
	reports = printer.GetReports()
	assert.Empty(t, reports, "GetReports should return an empty slice when there are no reports")
}
