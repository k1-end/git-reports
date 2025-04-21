package reportprinter

import (
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/k1-end/git-reports/src/report"
	"github.com/stretchr/testify/assert"
)

func TestHtmlPrinter_renderTable(t *testing.T) {
	// Test for HtmlPrinter.renderTable
	testReport := report.Report{}
	testReport.SetTitle("Test Table")
	testReport.SetReportType("table")
	testReport.SetLabels([]string{"Name", "Age", "City"})
	testReport.SetData([]report.Data{
		{StringValue: "Alice", IsInt: false},
		{IntValue: 30, IsInt: true},
		{StringValue: "New York", IsInt: false},
	})

	printer := HtmlPrinter{}
	result := printer.renderTable(testReport, 1)

	// Assert the output contains the expected HTML elements.
    assert.Regexp(t, regexp.MustCompile(`<table.*id="elem-1"`), result, "Output should contain the table element with the correct ID")
	assert.Contains(t, result, "Test Table", "Output should contain the title")
	assert.Contains(t, result, "Name", "Output should contain the label 'Name'")
	assert.Contains(t, result, "Age", "Output should contain the label 'Age'")
	assert.Contains(t, result, "City", "Output should contain the label 'City'")
	assert.Contains(t, result, "Alice", "Output should contain the value 'Alice'")
	assert.Contains(t, result, "30", "Output should contain the value '30'")
	assert.Contains(t, result, "New York", "Output should contain the value 'New York'")
}

func TestHtmlPrinter_renderDateHeatMapChart(t *testing.T) {
	// Helper function to create a report with sample data
	createHeatMapReport := func(startDate time.Time, days int, commits []int) report.Report {
		labels := make([]string, 0, days)
		data := make([]report.Data, 0, days)
		for i := 0; i < days; i++ {
			date := startDate.AddDate(0, 0, i)
			labels = append(labels, date.Format("2006-1-2"))
			commitCount := 0
			if i < len(commits) {
				commitCount = commits[i]
			}
			data = append(data, report.Data{IntValue: commitCount, IsInt: true})
		}
		testReport := report.Report{}
		testReport.SetReportType("date_heatmap")
		testReport.SetLabels(labels)
		testReport.SetData(data)
		return testReport
	}

	// Test for HtmlPrinter.renderDateHeatMapChart
	startDate := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	testReport := createHeatMapReport(startDate, 30, []int{1, 0, 5, 12, 20, 21, 0, 2, 3, 0, 0, 10, 11, 14, 16, 17, 18, 19, 22, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})

	printer := HtmlPrinter{}
	result := printer.renderDateHeatMapChart(testReport, 2)

	// Assert the output contains the expected HTML and data.
	assert.Regexp(t, regexp.MustCompile(`<div.*id="elem-[0-9]*"`), result, "Output should contain the div element with the correct ID")
	assert.Contains(t, result, "2024", "Output should contain the year")
	assert.Contains(t, result, "Sun", "Output should contain day name")
	assert.Contains(t, result, "Tue", "Output should contain day name")
	assert.Contains(t, result, "Fri", "Output should contain day name")
}

func TestHtmlPrinter_renderBartChart(t *testing.T) {
	// Test for HtmlPrinter.renderBartChart
	testReport := report.Report{}
	testReport.SetTitle("Bar Chart Example")
	testReport.SetReportType("bar_chart")
	testReport.SetLabels([]string{"A", "B", "C", "D"})
	testReport.SetData([]report.Data{
		{IntValue: 10, IsInt: true},
		{IntValue: 25, IsInt: true},
		{IntValue: 15, IsInt: true},
		{IntValue: 5, IsInt: true},
	})

	printer := HtmlPrinter{}
	result := printer.renderBartChart(testReport, 3)

	// Assert the output contains the expected HTML and data.
	assert.Regexp(t, regexp.MustCompile(`<div.*id="elem-3"`), result, "Output should contain the div element with the correct ID")
	assert.Contains(t, result, "Bar Chart Example", "Output should contain the title")
	assert.Contains(t, result, `"A"`, "Output should contain label A")
	assert.Contains(t, result, `"B"`, "Output should contain label B")
	assert.Contains(t, result, `"C"`, "Output should contain label C")
	assert.Contains(t, result, `"D"`, "Output should contain label D")
	assert.Contains(t, result, `10`, "Output should contain value 10")
	assert.Contains(t, result, `25`, "Output should contain value 25")
	assert.Contains(t, result, `15`, "Output should contain value 15")
	assert.Contains(t, result, `5`, "Output should contain value 5")
}

func TestHtmlPrinter_Print(t *testing.T) {
	// Test for HtmlPrinter.Print
	printer := HtmlPrinter{}
	printer.SetProjectTitle("My Project")

	tableReport := report.Report{}
	tableReport.SetTitle("Example Table")
	tableReport.SetReportType("table")
	tableReport.SetLabels([]string{"Header 1", "Header 2"})
	tableReport.SetData([]report.Data{
		{StringValue: "Data 1", IsInt: false},
		{StringValue: "Data 2", IsInt: false},
	})

	heatmapReport := report.Report{}
	heatmapReport.SetTitle("Example Heatmap")
	heatmapReport.SetReportType("date_heatmap")
	heatmapReport.SetLabels([]string{"2024-01-01", "2024-01-02"})
	heatmapReport.SetData([]report.Data{
		{IntValue: 5, IsInt: true},
		{IntValue: 10, IsInt: true},
	})

	barChartReport := report.Report{}
	barChartReport.SetTitle("Example Bar Chart")
	barChartReport.SetReportType("bar_chart")
	barChartReport.SetLabels([]string{"Label A", "Label B"})
	barChartReport.SetData([]report.Data{
		{IntValue: 20, IsInt: true},
		{IntValue: 30, IsInt: true},
	})

	printer.RegisterReport(tableReport)
	printer.RegisterReport(heatmapReport)
	printer.RegisterReport(barChartReport)

	// Create a temporary file to write the output to.
	tmpFile, err := os.CreateTemp("", "test_output.html")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up the file after the test.

	// Call the Print method.
	printer.Print(tmpFile)
	err = tmpFile.Close()
	if err != nil{
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Read the content of the temporary file.
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read temporary file: %v", err)
	}
	output := string(content)

	// Assert the output contains the expected HTML structure and report data.
	assert.Contains(t, output, "My Project", "Output should contain the project title")
	assert.Regexp(t, regexp.MustCompile(`<table.*id="elem-0"`), output, "Output should contain the table")
	assert.Regexp(t, regexp.MustCompile(`<div.*id="elem-1"`), output, "Output should contain the heatmap")
	assert.Regexp(t, regexp.MustCompile(`<div.*id="elem-2"`), output, "Output should contain the bar chart")
	assert.Contains(t, output, "Example Table", "Output should contain the table title")
	assert.Contains(t, output, "Example Heatmap", "Output should contain the heatmap title")
	assert.Contains(t, output, "Example Bar Chart", "Output should contain the bar chart title")
	assert.Contains(t, output, "Header 1", "Output should contain table header")
	assert.Contains(t, output, "Header 2", "Output should contain table header")
	assert.Contains(t, output, "Data 1", "Output should contain table data")
	assert.Contains(t, output, "Data 2", "Output should contain table data")
	assert.Contains(t, output, "2024", "Output should contain heatmap year")
	assert.Contains(t, output, "Label A", "Output should contain bar chart label")
	assert.Contains(t, output, "Label B", "Output should contain bar chart label")
	assert.Contains(t, output, `20`, "Output should contain bar chart data")
	assert.Contains(t, output, `30`, "Output should contain bar chart data")
}

