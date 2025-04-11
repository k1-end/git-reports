package reportgenerator

import (
	"testing"
	"time"

	"github.com/k1-end/git-reports/src/report"
	"github.com/stretchr/testify/assert"
)

func TestCommitsPerHourReportGenerator_LogIterationStep(t *testing.T) {
	// Initialize a generator with a non-nil CommitsPerHourMap
	generator := CommitsPerHourReportGenerator{
		CommitsPerHourMap: make([]int, 24),
	}

	// Test with a commit at 00:00
	commitTime1 := time.Date(2024, time.January, 15, 0, 0, 0, 0, time.Local)
	commit1 := createMockCommit("Author A", "test@example.com", commitTime1)
	author1 := Author{Name: "Author A", Emails: map[string]bool{"test@example.com": true}}
	generator.LogIterationStep(commit1, author1)
	assert.Equal(t, 1, generator.CommitsPerHourMap[0], "Commit count for hour 0 should be 1")

	// Test with a commit at 13:00
	commitTime2 := time.Date(2024, time.January, 15, 13, 0, 0, 0, time.Local)
	commit2 := createMockCommit("Author B", "test2@example.com", commitTime2)
	author2 := Author{Name: "Author B", Emails: map[string]bool{"test2@example.com": true}}
	generator.LogIterationStep(commit2, author2)
	assert.Equal(t, 1, generator.CommitsPerHourMap[13], "Commit count for hour 13 should be 1")

	// Test with multiple commits at the same hour
	commitTime3 := time.Date(2024, time.January, 15, 13, 30, 0, 0, time.Local)
	commit3 := createMockCommit("Author C", "test3@example.com", commitTime3)
	author3 := Author{Name: "Author C", Emails: map[string]bool{"test3@example.com": true}}
	generator.LogIterationStep(commit3, author3)
	assert.Equal(t, 2, generator.CommitsPerHourMap[13], "Commit count for hour 13 should be 2")
}

func TestCommitsPerHourReportGenerator_GetReport(t *testing.T) {
	// Initialize CommitsPerHourMap with some data
	generator := CommitsPerHourReportGenerator{
		CommitsPerHourMap: []int{
			0, 2, 5, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0,
		},
	}

	r := generator.GetReport()

	// Check report title
	assert.Equal(t, "Commits per hour of day", r.GetTitle(), "Report title should be 'Commits per hour of day'")

	// Check report type
	assert.Equal(t, "bar_chart", r.GetReportType(), "Report type should be 'bar_chart'")

	// Check labels
	expectedLabels := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23"}
	assert.Equal(t, expectedLabels, r.GetLabels(), "Report labels should be hours 1-23")

	// Check data
	expectedData := []report.Data{
		{IsInt: true, IntValue: 2, StringValue: ""},  // Hour 1
		{IsInt: true, IntValue: 5, StringValue: ""},  // Hour 2
		{IsInt: true, IntValue: 1, StringValue: ""},  // Hour 3
		{IsInt: true, IntValue: 0, StringValue: ""},  // Hour 4
		{IsInt: true, IntValue: 0, StringValue: ""},  // Hour 5
		{IsInt: true, IntValue: 0, StringValue: ""},  // Hour 6
		{IsInt: true, IntValue: 0, StringValue: ""},  // Hour 7
		{IsInt: true, IntValue: 0, StringValue: ""},  // Hour 8
		{IsInt: true, IntValue: 0, StringValue: ""},  // Hour 9
		{IsInt: true, IntValue: 0, StringValue: ""},  // Hour 10
		{IsInt: true, IntValue: 0, StringValue: ""},  // Hour 11
		{IsInt: true, IntValue: 0, StringValue: ""},  // Hour 12
		{IsInt: true, IntValue: 3, StringValue: ""},  // Hour 13
		{IsInt: true, IntValue: 0, StringValue: ""},  // Hour 14
		{IsInt: true, IntValue: 0, StringValue: ""},  // Hour 15
		{IsInt: true, IntValue: 0, StringValue: ""},  // Hour 16
		{IsInt: true, IntValue: 1, StringValue: ""},  // Hour 17
		{IsInt: true, IntValue: 0, StringValue: ""},  // Hour 18
		{IsInt: true, IntValue: 0, StringValue: ""},  // Hour 19
		{IsInt: true, IntValue: 0, StringValue: ""},  // Hour 20
		{IsInt: true, IntValue: 0, StringValue: ""},  // Hour 21
		{IsInt: true, IntValue: 0, StringValue: ""},  // Hour 22
		{IsInt: true, IntValue: 0, StringValue: ""},  // Hour 23
	}
	assert.Equal(t, expectedData, r.GetData(), "Report data should match the provided commit counts")
}

func TestCommitsPerHourReportGenerator_GetReport_EmptyMap(t *testing.T) {
	// Initialize CommitsPerHourMap as empty
	generator := CommitsPerHourReportGenerator{
		CommitsPerHourMap: make([]int, 24),
	}

	r := generator.GetReport()

	// Check report title
	assert.Equal(t, "Commits per hour of day", r.GetTitle(), "Report title should be 'Commits per hour of day'")

	// Check report type
	assert.Equal(t, "bar_chart", r.GetReportType(), "Report type should be 'bar_chart'")

	// Check labels
	expectedLabels := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23"}
	assert.Equal(t, expectedLabels, r.GetLabels(), "Report labels should be hours 1-23")

	// Check data
	expectedData := []report.Data{
		{IsInt: true, IntValue: 0, StringValue: ""},
		{IsInt: true, IntValue: 0, StringValue: ""},
		{IsInt: true, IntValue: 0, StringValue: ""},
		{IsInt: true, IntValue: 0, StringValue: ""},
		{IsInt: true, IntValue: 0, StringValue: ""},
		{IsInt: true, IntValue: 0, StringValue: ""},
		{IsInt: true, IntValue: 0, StringValue: ""},
		{IsInt: true, IntValue: 0, StringValue: ""},
		{IsInt: true, IntValue: 0, StringValue: ""},
		{IsInt: true, IntValue: 0, StringValue: ""},
		{IsInt: true, IntValue: 0, StringValue: ""},
		{IsInt: true, IntValue: 0, StringValue: ""},
		{IsInt: true, IntValue: 0, StringValue: ""},
		{IsInt: true, IntValue: 0, StringValue: ""},
		{IsInt: true, IntValue: 0, StringValue: ""},
		{IsInt: true, IntValue: 0, StringValue: ""},
		{IsInt: true, IntValue: 0, StringValue: ""},
		{IsInt: true, IntValue: 0, StringValue: ""},
		{IsInt: true, IntValue: 0, StringValue: ""},
		{IsInt: true, IntValue: 0, StringValue: ""},
		{IsInt: true, IntValue: 0, StringValue: ""},
		{IsInt: true, IntValue: 0, StringValue: ""},
		{IsInt: true, IntValue: 0, StringValue: ""},
	}
	assert.Equal(t, expectedData, r.GetData(), "Report data should be all zeros")
}

