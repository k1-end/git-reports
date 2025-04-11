package reportgenerator

import (
	"testing"
	"time"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/k1-end/git-reports/src/report"
	"github.com/stretchr/testify/assert"
)

func TestGeneralInfoReportGenerator_LogIterationStep(t *testing.T) {
	generator := GeneralInfoReportGenerator{}

	// Test with the first commit
	commitTime1 := time.Date(2024, time.January, 15, 10, 0, 0, 0, time.UTC)
	commit1 := createMockCommit("Author A", "authora@example.com", commitTime1)
	authorA := Author{Name: "Author A", Emails: map[string]bool{"authora@example.com": true}}
	generator.LogIterationStep(commit1, authorA)
	assert.Equal(t, 1, generator.ContributorsNo, "ContributorsNo should be 1")
	assert.Equal(t, 1, generator.CommitsNo, "CommitsNo should be 1")
	assert.Contains(t, generator.contributors, "Author A", "Author A should be in contributors")

	// Test with another commit from the same author
	commitTime2 := time.Date(2024, time.January, 16, 12, 0, 0, 0, time.UTC)
	commit2 := createMockCommit("Author A", "authora@example.com", commitTime2)
	generator.LogIterationStep(commit2, authorA)
	assert.Equal(t, 1, generator.ContributorsNo, "ContributorsNo should still be 1")
	assert.Equal(t, 2, generator.CommitsNo, "CommitsNo should be 2")
	assert.Contains(t, generator.contributors, "Author A", "Author A should still be in contributors")

	// Test with a commit from a different author
	commitTime3 := time.Date(2024, time.January, 17, 14, 0, 0, 0, time.UTC)
	commit3 := createMockCommit("Author B", "authorb@example.com", commitTime3)
	authorB := Author{Name: "Author B", Emails: map[string]bool{"authorb@example.com": true}}
	generator.LogIterationStep(commit3, authorB)
	assert.Equal(t, 2, generator.ContributorsNo, "ContributorsNo should be 2")
	assert.Equal(t, 3, generator.CommitsNo, "CommitsNo should be 3")
	assert.Contains(t, generator.contributors, "Author B", "Author B should be in contributors")
}

func TestGeneralInfoReportGenerator_FileIterationStep(t *testing.T) {
	generator := GeneralInfoReportGenerator{}

	// Test with a file
	file1 := &object.File{
		Name: "test.txt",
		Mode: 0644, // Example file mode
		Blob: object.Blob{ // Use Blob struct
			Size: 1024,
			Hash: plumbing.ComputeHash(plumbing.BlobObject, []byte("dummy content")), //Need to provide a hash.
		},
	}
	generator.FileIterationStep(file1)
	assert.Equal(t, 1, generator.FilesNo, "FilesNo should be 1")
	assert.Equal(t, uint64(1024), generator.ProjectSize, "ProjectSize should be 1024")

	// Test with another file
	file2 := &object.File{
		Name: "image.jpg",
		Mode: 0644,
		Blob: object.Blob{
			Size: 512000,
			Hash: plumbing.ComputeHash(plumbing.BlobObject, []byte("dummy content")), //Need to provide a hash.
		},
	}
	generator.FileIterationStep(file2)
	assert.Equal(t, 2, generator.FilesNo, "FilesNo should be 2")
	assert.Equal(t, uint64(513024), generator.ProjectSize, "ProjectSize should be 513024")
}

func TestGeneralInfoReportGenerator_GetReport(t *testing.T) {
	generator := GeneralInfoReportGenerator{
		ContributorsNo: 3,
		CommitsNo:      150,
		ProjectSize:    2048500,
		FilesNo:        50,
	}

	r := generator.GetReport()

	// Check report title
	assert.Equal(t, "General Info", r.GetTitle(), "Report title should be 'General Info'")

	// Check report type
	assert.Equal(t, "table", r.GetReportType(), "Report type should be 'table'")

	// Check labels
	expectedLabels := []string{"Number of contributors", "Number of commits", "Project size", "Number of files"}
	assert.Equal(t, expectedLabels, r.GetLabels(), "Report labels should be correct")

	// Check data
	expectedData := []report.Data{
		{IsInt: false, StringValue: "3"},
		{IsInt: false, StringValue: "150"},
		{IsInt: false, StringValue: "2,048 KB"},
		{IsInt: false, StringValue: "50"},
	}
	assert.Equal(t, expectedData, r.GetData(), "Report data should be correct")
}

func TestGeneralInfoReportGenerator_GetReport_ZeroValues(t *testing.T) {
	generator := GeneralInfoReportGenerator{}

	r := generator.GetReport()

	// Check report title
	assert.Equal(t, "General Info", r.GetTitle(), "Report title should be 'General Info'")

	// Check report type
	assert.Equal(t, "table", r.GetReportType(), "Report type should be 'table'")

	// Check labels
	expectedLabels := []string{"Number of contributors", "Number of commits", "Project size", "Number of files"}
	assert.Equal(t, expectedLabels, r.GetLabels(), "Report labels should be correct")

	// Check data
	expectedData := []report.Data{
		{IsInt: false, StringValue: "0"},
		{IsInt: false, StringValue: "0"},
		{IsInt: false, StringValue: "0 KB"},
		{IsInt: false, StringValue: "0"},
	}
	assert.Equal(t, expectedData, r.GetData(), "Report data should be all zeros")
}
