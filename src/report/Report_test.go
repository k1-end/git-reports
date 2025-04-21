package report

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestData(t *testing.T) {
	// Test for Data struct
	t.Run("Test Data Int", func(t *testing.T) {
		data := Data{
			IntValue:    10,
			IsInt:       true,
			StringValue: "invalid", // should be ignored
		}
		assert.Equal(t, 10, data.IntValue, "IntValue should be set")
		assert.True(t, data.IsInt, "IsInt should be true")
	})

	t.Run("Test Data String", func(t *testing.T) {
		data := Data{
			StringValue: "hello",
			IsInt:       false,
			IntValue:    100, // should be ignored
		}
		assert.Equal(t, "hello", data.StringValue, "StringValue should be set")
		assert.False(t, data.IsInt, "IsInt should be false")
	})

	t.Run("Test Data Zero Values", func(t *testing.T) {
		data := Data{}
		assert.Equal(t, 0, data.IntValue, "IntValue should default to 0")
		assert.False(t, data.IsInt, "IsInt should default to false")
		assert.Equal(t, "", data.StringValue, "StringValue should default to empty string")
	})
}

func TestReport_SetTitle(t *testing.T) {
	// Test for Report.SetTitle
	report := Report{}
	title := "My Report"

	report.SetTitle(title)
	assert.Equal(t, title, report.title, "Title should be set correctly")

	report.SetTitle("Another Report")
	assert.Equal(t, "Another Report", report.title, "Title should be updated")
}

func TestReport_SetData(t *testing.T) {
	// Test for Report.SetData
	report := Report{}
	data := []Data{{IntValue: 1}, {IntValue: 2}}

	report.SetData(data)
	assert.Equal(t, data, report.data, "Data should be set correctly")

	newData := []Data{{StringValue: "a"}, {StringValue: "b"}}
	report.SetData(newData)
	assert.Equal(t, newData, report.data, "Data should be updated")
}

func TestReport_SetLabels(t *testing.T) {
	// Test for Report.SetLabels
	report := Report{}
	labels := []string{"A", "B"}

	report.SetLabels(labels)
	assert.Equal(t, labels, report.labels, "Labels should be set correctly")

	newLabels := []string{"C", "D"}
	report.SetLabels(newLabels)
	assert.Equal(t, newLabels, report.labels, "Labels should be updated")
}

func TestReport_SetReportType(t *testing.T) {
	// Test for Report.SetReportType
	report := Report{}
	reportType := "table"

	report.SetReportType(reportType)
	assert.Equal(t, reportType, report.reportType, "ReportType should be set correctly")

	newReportType := "chart"
	report.SetReportType(newReportType)
	assert.Equal(t, newReportType, report.reportType, "ReportType should be updated")
}

func TestReport_GetTitle(t *testing.T) {
	// Test for Report.GetTitle
	report := Report{title: "My Report"} // Initialize directly for test
	title := report.GetTitle()
	assert.Equal(t, "My Report", title, "GetTitle should return the correct title")
}

func TestReport_GetData(t *testing.T) {
	// Test for Report.GetData
	data := []Data{{IntValue: 1}, {IntValue: 2}}
	report := Report{data: data} // Initialize directly for test
	retrievedData := report.GetData()
	assert.Equal(t, data, retrievedData, "GetData should return the correct data")
}

func TestReport_GetLabels(t *testing.T) {
	// Test for Report.GetLabels
	labels := []string{"A", "B"}
	report := Report{labels: labels} // Initialize directly for test
	retrievedLabels := report.GetLabels()
	assert.Equal(t, labels, retrievedLabels, "GetLabels should return the correct labels")
}

func TestReport_GetReportType(t *testing.T) {
	// Test for Report.GetReportType
	reportType := "table"
	report := Report{reportType: reportType} // Initialize directly for test
	retrievedReportType := report.GetReportType()
	assert.Equal(t, reportType, retrievedReportType, "GetReportType should return the correct report type")
}
