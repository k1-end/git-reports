package reportgenerator

import (
	"testing"

	"github.com/k1-end/git-reports/src/report"
	"github.com/stretchr/testify/assert"
	"time"
)

func TestCommitsPerDevReportGenerator_LogIterationStep(t *testing.T) {
	generator := CommitsPerDevReportGenerator{
		CommitsPerDevMap: make(map[string]int),
	}

	// Test with a single commit
	commitTime1 := time.Date(2024, time.January, 15, 10, 0, 0, 0, time.UTC)
	commit1 := createMockCommit("Author A", "authora@example.com", commitTime1)
    authorA := Author{Name: "Author A", Emails: map[string]bool{"authora@example.com": true}}
	generator.LogIterationStep(commit1, authorA)
	assert.Equal(t, 1, generator.CommitsPerDevMap["Author A"], "Author A should have 1 commit")

	// Test with multiple commits by the same author
	commitTime2 := time.Date(2024, time.January, 16, 12, 0, 0, 0, time.UTC)
	commit2 := createMockCommit("Author A", "authora@example.com", commitTime2)
	generator.LogIterationStep(commit2, authorA)
	assert.Equal(t, 2, generator.CommitsPerDevMap["Author A"], "Author A should have 2 commits")

	// Test with a commit by a different author
	commitTime3 := time.Date(2024, time.January, 17, 14, 0, 0, 0, time.UTC)
	commit3 := createMockCommit("Author B", "authorb@example.com", commitTime3)
    authorB := Author{Name: "Author B", Emails: map[string]bool{"authorb@example.com": true}}
	generator.LogIterationStep(commit3, authorB)
	assert.Equal(t, 1, generator.CommitsPerDevMap["Author B"], "Author B should have 1 commit")
	assert.Equal(t, 2, generator.CommitsPerDevMap["Author A"], "Author A should still have 2 commits")
}

func TestCommitsPerDevReportGenerator_GetReport(t *testing.T) {
	generator := CommitsPerDevReportGenerator{
		CommitsPerDevMap: map[string]int{
			"Author A": 10,
			"Author B": 5,
			"Author C": 15,
		},
	}

	r := generator.GetReport()

	// Check report title
	assert.Equal(t, "Commits per developer", r.GetTitle(), "Report title should be 'Commits per developer'")

	// Check report type
	assert.Equal(t, "bar_chart", r.GetReportType(), "Report type should be 'bar_chart'")

	// Check labels (sorted by commit count)
	expectedLabels := []string{"Author C", "Author A", "Author B"}
	assert.Equal(t, expectedLabels, r.GetLabels(), "Report labels should be sorted by commit count")

	// Check data
	expectedData := []report.Data{
		{IsInt: true, IntValue: 15, StringValue: ""},
		{IsInt: true, IntValue: 10, StringValue: ""},
		{IsInt: true, IntValue: 5, StringValue: ""},
	}
	assert.Equal(t, expectedData, r.GetData(), "Report data should match commit counts in descending order")
}

func TestCommitsPerDevReportGenerator_GetReport_EmptyMap(t *testing.T) {
	generator := CommitsPerDevReportGenerator{
		CommitsPerDevMap: make(map[string]int),
	}

	report := generator.GetReport()

	// Check report title
	assert.Equal(t, "Commits per developer", report.GetTitle(), "Report title should be 'Commits per developer'")

	// Check report type
	assert.Equal(t, "bar_chart", report.GetReportType(), "Report type should be 'bar_chart'")

	// Check labels (empty)
	assert.Empty(t, report.GetLabels(), "Report labels should be empty")

	// Check data (empty)
	assert.Empty(t, report.GetData(), "Report data should be empty")
}
