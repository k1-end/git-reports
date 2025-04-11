package reportgenerator

import (
	"testing"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/k1-end/git-reports/src/report"
	"github.com/stretchr/testify/assert"
)

func TestFileTypeReportGenerator_FileIterationStep(t *testing.T) {
	generator := FileTypeReportGenerator{
		FileTypeMap: make(map[string]int),
	}

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
	assert.Equal(t, 1024, generator.FileTypeMap[".txt"], "File size for .txt should be 1024")

	// Test with another file of the same type
	file2 := &object.File{
		Name: "another.txt",
		Mode: 0644,
		Blob: object.Blob{
			Size: 2048,
			Hash: plumbing.ComputeHash(plumbing.BlobObject, []byte("another dummy content")),
		},
	}
	generator.FileIterationStep(file2)
	assert.Equal(t, 3072, generator.FileTypeMap[".txt"], "File size for .txt should be 3072")

	// Test with a file of a different type
	file3 := &object.File{
		Name: "image.jpg",
		Mode: 0644,
		Blob: object.Blob{
			Size: 512000,
			Hash: plumbing.ComputeHash(plumbing.BlobObject, []byte("image dummy content")),
		},
	}
	generator.FileIterationStep(file3)
	assert.Equal(t, 512000, generator.FileTypeMap[".jpg"], "File size for .jpg should be 512000")
	assert.Equal(t, 3072, generator.FileTypeMap[".txt"], "File size for .txt should still be 3072")
}

func TestFileTypeReportGenerator_GetReport(t *testing.T) {
	generator := FileTypeReportGenerator{
		FileTypeMap: map[string]int{
			".txt":  3000,
			".jpg":  2500000,
			".html": 10000,
			".css":  500, // This will be filtered out in GetReport
			".go":   1000,
		},
	}

	r := generator.GetReport()

	// Check report title
	assert.Equal(t, "File Types (KB)", r.GetTitle(), "Report title should be 'File Types (KB)'")

	// Check report type
	assert.Equal(t, "bar_chart", r.GetReportType(), "Report type should be 'bar_chart'")

	// Check labels (sorted by size, filtered for > 1000)
	expectedLabels := []string{".jpg", ".html", ".txt", ".go"} // Corrected order
	assert.Equal(t, expectedLabels, r.GetLabels(), "Report labels should be sorted by size and filtered")

	// Check data (sizes in KB, filtered for > 1000)
	expectedData := []report.Data{
		{IsInt: true, IntValue: 2500, StringValue: ""}, // .jpg: 2500000 / 1000
		{IsInt: true, IntValue: 10, StringValue: ""},    // .html: 10000 / 1000
		{IsInt: true, IntValue: 3, StringValue: ""},     // .txt: 3000 / 1000
		{IsInt: true, IntValue: 1, StringValue: ""},     // .go: 1000 / 1000
	}
	assert.Equal(t, expectedData, r.GetData(), "Report data should match file sizes in KB and be filtered")
}

func TestFileTypeReportGenerator_GetReport_EmptyMap(t *testing.T) {
	generator := FileTypeReportGenerator{
		FileTypeMap: make(map[string]int),
	}

	r := generator.GetReport()

	// Check report title
	assert.Equal(t, "File Types (KB)", r.GetTitle(), "Report title should be 'File Types (KB)'")

	// Check report type
	assert.Equal(t, "bar_chart", r.GetReportType(), "Report type should be 'bar_chart'")

	// Check labels (empty)
	assert.Empty(t, r.GetLabels(), "Report labels should be empty")

	// Check data (empty)
	assert.Empty(t, r.GetData(), "Report data should be empty")
}

