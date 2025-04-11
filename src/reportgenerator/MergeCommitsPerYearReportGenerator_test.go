package reportgenerator

import (
	"testing"
	"time"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/k1-end/git-reports/src/report"
	"github.com/stretchr/testify/assert"
)

func TestMergeCommitsPerYearReportGenerator_LogIterationStep(t *testing.T) {
	generator := MergeCommitsPerYearReportGenerator{
		MergeCommitsPerYearMap: make(map[int]int),
	}

	// Test with a merge commit
	commitTime1 := time.Date(2024, time.January, 15, 10, 0, 0, 0, time.UTC)
	commit1 := createMockCommit("Author A", "authora@example.com", commitTime1)
	commit1.ParentHashes = []plumbing.Hash{plumbing.NewHash("parent1"), plumbing.NewHash("parent2")} // Set ParentHashes directly
	authorA := Author{Name: "Author A", Emails: map[string]bool{"authora@example.com": true}}
	generator.LogIterationStep(commit1, authorA)
	assert.Equal(t, 1, generator.MergeCommitsPerYearMap[2024], "2024 should have 1 merge commit")

	// Test with a non-merge commit
	commitTime2 := time.Date(2024, time.February, 20, 12, 0, 0, 0, time.UTC)
	commit2 := createMockCommit("Author B", "authorb@example.com", commitTime2)
	authorB := Author{Name: "Author B", Emails: map[string]bool{"authorb@example.com": true}}
	generator.LogIterationStep(commit2, authorB)
	assert.Equal(t, 1, generator.MergeCommitsPerYearMap[2024], "2024 should still have 1 merge commit")

	// Test with another merge commit in the same year
	commitTime3 := time.Date(2024, time.March, 10, 14, 0, 0, 0, time.UTC)
	commit3 := createMockCommit("Author C", "authorc@example.com", commitTime3)
	commit3.ParentHashes = []plumbing.Hash{plumbing.NewHash("parent3"), plumbing.NewHash("parent4")} // Set ParentHashes directly
	authorC := Author{Name: "Author C", Emails: map[string]bool{"authorc@example.com": true}}
	generator.LogIterationStep(commit3, authorC)
	assert.Equal(t, 2, generator.MergeCommitsPerYearMap[2024], "2024 should have 2 merge commits")

	// Test with a merge commit in a different year
	commitTime4 := time.Date(2025, time.April, 5, 16, 0, 0, 0, time.UTC)
	commit4 := createMockCommit("Author D", "authord@example.com", commitTime4)
	commit4.ParentHashes = []plumbing.Hash{plumbing.NewHash("parent5"), plumbing.NewHash("parent6")} // Set ParentHashes directly
	authorD := Author{Name: "Author D", Emails: map[string]bool{"authord@example.com": true}}
	generator.LogIterationStep(commit4, authorD)
	assert.Equal(t, 1, generator.MergeCommitsPerYearMap[2025], "2025 should have 1 merge commit")
	assert.Equal(t, 2, generator.MergeCommitsPerYearMap[2024], "2024 should still have 2 merge commits")
}

func TestMergeCommitsPerYearReportGenerator_GetReport(t *testing.T) {
	generator := MergeCommitsPerYearReportGenerator{
		MergeCommitsPerYearMap: map[int]int{
			2022: 5,
			2023: 12,
			2024: 8,
			2025: 0,
		},
	}

	r := generator.GetReport()

	// Check report title
	assert.Equal(t, "Merge Commits per year", r.GetTitle(), "Report title should be 'Merge Commits per year'")

	// Check report type
	assert.Equal(t, "bar_chart", r.GetReportType(), "Report type should be 'bar_chart'")

	// Check labels (sorted years)
	expectedLabels := []string{"2022", "2023", "2024", "2025"}
	assert.Equal(t, expectedLabels, r.GetLabels(), "Report labels should be sorted years")

	// Check data
	expectedData := []report.Data{
		{IsInt: true, IntValue: 5, StringValue: ""},
		{IsInt: true, IntValue: 12, StringValue: ""},
		{IsInt: true, IntValue: 8, StringValue: ""},
		{IsInt: true, IntValue: 0, StringValue: ""},
	}
	assert.Equal(t, expectedData, r.GetData(), "Report data should match merge commit counts")
}

func TestMergeCommitsPerYearReportGenerator_GetReport_EmptyMap(t *testing.T) {
	generator := MergeCommitsPerYearReportGenerator{
		MergeCommitsPerYearMap: make(map[int]int),
	}

	r := generator.GetReport()

	// Check report title
	assert.Equal(t, "Merge Commits per year", r.GetTitle(), "Report title should be 'Merge Commits per year'")

	// Check report type
	assert.Equal(t, "bar_chart", r.GetReportType(), "Report type should be 'bar_chart'")

	// Check labels (empty)
	assert.Empty(t, r.GetLabels(), "Report labels should be empty")

	// Check data (empty)
	assert.Empty(t, r.GetData(), "Report data should be empty")
}
