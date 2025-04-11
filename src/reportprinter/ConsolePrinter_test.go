package reportprinter

import (
	"bytes"
	"errors"

	"os"
	"testing"
	"time"

	"github.com/k1-end/git-reports/src/report"
	"github.com/stretchr/testify/assert"
)

func TestGetColor(t *testing.T) {
	// Test cases for getColor function
	testCases := []struct {
		commitCount int
		expectedColor string
		description string
	}{
		{0, "178;215;155", "No commits"},
		{3, "139;195;74", "Few commits"},
		{8, "34;139;34", "Some commits"},
		{12, "0;100;0", "More commits"},
		{18, "0;128;128", "Many commits"},
		{25, "0;64;0", "Very many commits"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			color := getColor(tc.commitCount)
			assert.Equal(t, tc.expectedColor, color, "Color for commit count %d should be %s", tc.commitCount, tc.expectedColor)
		})
	}
}

func TestYearData_getFirstMonth(t *testing.T) {
	// Test cases for yearData.getFirstMonth
	testCases := []struct {
		yearData yearData
		expectedMonth time.Month
		expectedError error
		description string
	}{
		{
			yearData: yearData{
				Year: 2024,
				Months: map[time.Month]map[int]struct {
					Date        time.Time
					CommitCount int
				}{
					time.January: {},
					time.March:   {},
				},
			},
			expectedMonth: time.January,
			expectedError: nil,
			description: "First month is January",
		},
		{
			yearData: yearData{
				Year: 2024,
				Months: map[time.Month]map[int]struct {
					Date        time.Time
					CommitCount int
				}{
					time.May:   {},
					time.August:  {},
				},
			},
			expectedMonth: time.May,
			expectedError: nil,
			description: "First month is May",
		},
		{
			yearData: yearData{
				Year: 2024,
				Months: map[time.Month]map[int]struct {
					Date        time.Time
					CommitCount int
				}{},
			},
			expectedMonth: time.Month(0),
			expectedError: errors.New("empty name"),
			description: "No months",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			month, err := tc.yearData.getFirstMonth()
			assert.Equal(t, tc.expectedMonth, month, "First month should be %s", tc.expectedMonth)
			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error(), "Error should be: %v", tc.expectedError)
			} else {
				assert.NoError(t, err, "Error should be nil")
			}

		})
	}
}

func TestConsolePrinter_PrintBarChart(t *testing.T) {
	// Test for ConsolePrinter.printBarChart
	testReport := report.Report{}
	testReport.SetTitle("Test Bar Chart")
	testReport.SetReportType("bar_chart")
	testReport.SetLabels([]string{"A", "B", "C"})
	testReport.SetData([]report.Data{
		{IntValue: 10, IsInt: true},
		{IntValue: 20, IsInt: true},
		{IntValue: 5, IsInt: true},
	})

	// Capture the output
	r, w, _ := os.Pipe()

	printer := ConsolePrinter{} // No need to initialize reports here for this test
	printer.RegisterReport(testReport)
    printer.Print(w)
	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Assert the important parts of the output.  Exact matching is difficult
	// and fragile with terminal output.
	assert.Contains(t, output, "Test Bar Chart", "Output should contain the title")
	assert.Contains(t, output, "A", "Output should contain label A")
	assert.Contains(t, output, "B", "Output should contain label B")
	assert.Contains(t, output, "C", "Output should contain label C")
	assert.Contains(t, output, "10", "Output should contain value 10")
	assert.Contains(t, output, "20", "Output should contain value 20")
	assert.Contains(t, output, "5", "Output should contain value 5")
}

func TestConsolePrinter_printTable(t *testing.T) {
	// Test for ConsolePrinter.printTable
	testReport := report.Report{}
	testReport.SetTitle("Test Table")
	testReport.SetReportType("table")
	testReport.SetLabels([]string{"Name", "Age", "City"})
	testReport.SetData([]report.Data{
		{StringValue: "Alice", IsInt: false},
		{IntValue: 30, IsInt: true},
		{StringValue: "New York", IsInt: false},
	})

	r, w, _ := os.Pipe()

	printer := ConsolePrinter{} // No need to initialize reports here for this test.
	printer.RegisterReport(testReport)
    printer.Print(w)

	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Assert the output contains the expected table elements
	assert.Contains(t, output, "Test Table", "Output should contain the title")
	assert.Contains(t, output, "Name", "Output should contain column header 'Name'")
	assert.Contains(t, output, "Age", "Output should contain column header 'Age'")
	assert.Contains(t, output, "City", "Output should contain column header 'City'")
	assert.Contains(t, output, "Alice", "Output should contain data 'Alice'")
	assert.Contains(t, output, "30", "Output should contain data '30'")
	assert.Contains(t, output, "New York", "Output should contain data 'New York'")
}

func TestConsolePrinter_printDateHeatMapChart(t *testing.T) {
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
		testReport.SetTitle("Commit Heatmap")
		testReport.SetReportType("date_heatmap")
		testReport.SetLabels(labels)
		testReport.SetData(data)
		return testReport
	}

	// Test for ConsolePrinter.printDateHeatMapChart
	startDate := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.Local)
	testReport := createHeatMapReport(startDate, 30, []int{1, 0, 5, 12, 20, 21, 0, 2, 3, 0, 0, 10, 11, 14, 16, 17, 18, 19, 22, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})

	r, w, _ := os.Pipe()

	printer := ConsolePrinter{} // Initialize printer
	printer.RegisterReport(testReport)
    printer.Print(w)

	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Assert the output.
	assert.Contains(t, output, "2024", "Output should contain the year")
	assert.Contains(t, output, "Jan", "Output should contain the month")
	assert.Contains(t, output, "Sun", "Output should contain day name")
	assert.Contains(t, output, "Mon", "Output should contain day name")
	assert.Contains(t, output, "Tue", "Output should contain day name")
	assert.Contains(t, output, "Wed", "Output should contain day name")
	assert.Contains(t, output, "Thu", "Output should contain day name")
	assert.Contains(t, output, "Fri", "Output should contain day name")
	assert.Contains(t, output, "Sat", "Output should contain day name")
	assert.Contains(t, output, "commits count guide", "Output should contain the guide")
}

func TestConsolePrinter_Print(t *testing.T) {
	// Test for ConsolePrinter.Print
	// Create a ConsolePrinter with a mix of report types.
	printer := ConsolePrinter{}
	printer.SetProjectTitle("Test Project")
    r1 := report.Report{}
    r1.SetTitle("Combined Bar Chart")
    r1.SetReportType("bar_chart")
    r1.SetLabels([]string{"X", "Y"})
    r1.SetData([]report.Data{
    {IntValue: 5, IsInt: true},
    {IntValue: 15, IsInt: true},
    })
	printer.RegisterReport(r1)

    r2 := report.Report{}
    r2.SetTitle("Combined Table")
    r2.SetReportType("table")
    r2.SetLabels([]string{"A", "B"})
    r2.SetData([]report.Data{
        {StringValue: "1", IsInt: false},
        {StringValue: "2", IsInt: false},
    })
	printer.RegisterReport(r2)

    r3 := report.Report{}
    r3.SetTitle("Combined Heatmap")
    r3.SetReportType("date_heatmap")
    r3.SetLabels([]string{"2024-01-01", "2024-01-02"})
    r3.SetData([]report.Data{
        {IntValue: 2, IsInt: true},
        {IntValue: 8, IsInt: true},
    })
	printer.RegisterReport(r3)


	// Create a temporary file for the printer to use.
	tmpFile, err := os.CreateTemp("", "test_output")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up the file after the test.

	printer.Print(tmpFile) // Use the temp file.

	var buf []byte
    buf, _ = os.ReadFile(tmpFile.Name()) // We have to reopen the tmpFile becuase the printer.Print closes the file.
	output := string(buf)

	// Assert that the output contains elements from all three report types.
	assert.Contains(t, output, "Combined Bar Chart", "Output should contain bar chart title")
	assert.Contains(t, output, "Combined Table", "Output should contain table title")
    // assert.Contains(t, output, "Combined Heatmap", "Output should contain heatmap title") //TODO: Add title to printer
	assert.Contains(t, output, "X", "Output should contain bar chart label")
	assert.Contains(t, output, "Y", "Output should contain bar chart label")
	assert.Contains(t, output, "A", "Output should contain table label")
	assert.Contains(t, output, "B", "Output should contain table label")
	assert.Contains(t, output, "2024", "Output should contain year of heatmap")
}


