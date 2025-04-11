package reportgenerator

import (
	"testing"
	"time"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/k1-end/git-reports/src/report"
	"github.com/stretchr/testify/assert"
)

// Helper function to create a mock commit
func createMockCommit(authorName string, authorEmail string, commitTime time.Time) *object.Commit {
	return &object.Commit{
		Author: object.Signature{
			Name:  authorName,
			Email: authorEmail,
			When:  commitTime,
		},
		Committer: object.Signature{
			Name:  authorName,
			Email: authorEmail,
			When:  commitTime,
		},
		Message: "Test commit message",
	}
}

func TestCommitCountDateHeatMapGenerator_LogIterationStep(t *testing.T) {
	generator := CommitCountDateHeatMapGenerator{
		CommitsMap: make(map[string]int),
	}

	// Test with a single commit
	commitTime1 := time.Date(2024, time.January, 15, 10, 0, 0, 0, time.UTC)
	commit1 := createMockCommit("Test Author", "test@example.com", commitTime1)
    author1 := Author{Name: "Test Author", Emails: map[string]bool{"test@example.com": true}}
	generator.LogIterationStep(commit1, author1)
	assert.Equal(t, 1, generator.CommitsMap["2024-1-15"], "Commit count for 2024-01-15 should be 1")

	// Test with multiple commits on the same day
	commitTime2 := time.Date(2024, time.January, 15, 12, 0, 0, 0, time.UTC)
	commit2 := createMockCommit("Test Author", "test@example.com", commitTime2)
	generator.LogIterationStep(commit2, author1)
	assert.Equal(t, 2, generator.CommitsMap["2024-1-15"], "Commit count for 2024-01-15 should be 2")

	// Test with a commit on a different day
	commitTime3 := time.Date(2024, time.January, 20, 10, 0, 0, 0, time.UTC)
	commit3 := createMockCommit("Test Author", "test@example.com", commitTime3)
	generator.LogIterationStep(commit3, author1)
	assert.Equal(t, 1, generator.CommitsMap["2024-1-20"], "Commit count for 2024-01-20 should be 1")
	assert.Equal(t, 2, generator.CommitsMap["2024-1-15"], "Commit count for 2024-01-15 should still be 2")
}

func TestCommitCountDateHeatMapGenerator_GetReport(t *testing.T) {
	generator := CommitCountDateHeatMapGenerator{
		CommitsMap: map[string]int{
			"2024-01-10": 5,
			"2024-01-15": 10,
			"2024-01-05": 2,
		},
	}

	r := generator.GetReport()

	// Check report title
	assert.Equal(t, "Commit count heat map", r.GetTitle(), "Report title should be 'Commit count heat map'")

	// Check report type
	assert.Equal(t, "date_heatmap", r.GetReportType(), "Report type should be 'date_heatmap'")

	// Check labels (sorted dates)
	expectedLabels := []string{"2024-01-05", "2024-01-10", "2024-01-15"}
	assert.Equal(t, expectedLabels, r.GetLabels(), "Report labels should be sorted dates")

    // Check data
    expectedData := []report.Data{
        {IsInt: true, IntValue: 2, StringValue: ""},
        {IsInt: true, IntValue: 5, StringValue: ""},
        {IsInt: true, IntValue: 10, StringValue: ""},
    }
    assert.Equal(t, expectedData, r.GetData(), "Report data should match commit counts")
}

func TestCommitCountDateHeatMapGenerator_GetReport_EmptyMap(t *testing.T) {
	generator := CommitCountDateHeatMapGenerator{
		CommitsMap: make(map[string]int),
	}

	report := generator.GetReport()

	// Check report title
	assert.Equal(t, "Commit count heat map", report.GetTitle(), "Report title should be 'Commit count heat map'")

	// Check report type
	assert.Equal(t, "date_heatmap", report.GetReportType(), "Report type should be 'date_heatmap'")

	// Check labels (empty)
	assert.Empty(t, report.GetLabels(), "Report labels should be empty")

	// Check data (empty)
	assert.Empty(t, report.GetData(), "Report data should be empty")
}

