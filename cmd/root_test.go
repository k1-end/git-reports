package cmd

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/k1-end/git-reports/src/reportgenerator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to create a temporary directory and file for testing.
func createTempDirAndFile(t *testing.T, dirName, fileName, fileContent string) (string, string, func()) {
	tempDir, err := os.MkdirTemp("", dirName)
	require.NoError(t, err)

	filePath := filepath.Join(tempDir, fileName)
	err = os.WriteFile(filePath, []byte(fileContent), 0644)
	require.NoError(t, err)

	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return tempDir, filePath, cleanup
}

func TestIsValidFilePath(t *testing.T) {
	// Test cases for isValidFilePath
	t.Run("Valid file path", func(t *testing.T) {
		_, tempFile, cleanup := createTempDirAndFile(t, "test_valid", "testfile.txt", "test content")
		defer cleanup()
		assert.True(t, isValidFilePath(tempFile), "Should return true for a valid file path")
	})

	t.Run("Invalid file path - directory does not exist", func(t *testing.T) {
		assert.False(t, isValidFilePath("/nonexistent/file.txt"), "Should return false for a nonexistent directory")
	})

	t.Run("Invalid file path - not a directory", func(t *testing.T) {
		tempDir, _, cleanup := createTempDirAndFile(t, "test_not_dir", "testfile.txt", "test content")
		defer cleanup()
		assert.False(t, isValidFilePath(tempDir), "Should return false for a path that is not a directory")
	})

	t.Run("Invalid file path - no write permissions", func(t *testing.T) {
		if runtime.GOOS != "windows" { // Skip on Windows, permissions are handled differently
			tempDir, err := os.MkdirTemp("", "test_no_write")
			require.NoError(t, err)
			defer os.RemoveAll(tempDir)

			err = os.Chmod(tempDir, 0500) // Remove write permissions
			require.NoError(t, err)
			defer os.Chmod(tempDir, 0755) //restore permissions

			filePath := filepath.Join(tempDir, "testfile.txt")

			assert.False(t, isValidFilePath(filePath), "Should return false for a directory with no write permissions")
		}
	})
	t.Run("Valid file path with tilde", func(t *testing.T) {
		if runtime.GOOS != "windows" { //Tilde expansion not reliable on windows in tests
			homeDir, err := os.UserHomeDir()
			require.NoError(t, err)
			tempDir, err := os.MkdirTemp("", "test_tilde")
			require.NoError(t, err)
			defer os.RemoveAll(tempDir)
			relPath := strings.Replace(tempDir, homeDir, "~", 1)

			filePath := filepath.Join(relPath, "testfile.txt")
			assert.True(t, isValidFilePath(filePath), "Should return true for a valid file path with tilde")
		}
	})
}

func TestExpandTilde(t *testing.T) {
	// Test cases for expandTilde
	t.Run("No tilde", func(t *testing.T) {
		path := "/path/to/file"
		expandedPath, err := expandTilde(path)
		require.NoError(t, err)
		assert.Equal(t, path, expandedPath, "Should return the original path if no tilde")
	})

	t.Run("Tilde at the beginning", func(t *testing.T) {
		if runtime.GOOS != "windows" { // Tilde expansion in tests on Windows is not reliable.
			homeDir, err := os.UserHomeDir()
			require.NoError(t, err)
			path := "~/path/to/file"
			expandedPath, err := expandTilde(path)
			require.NoError(t, err)
			assert.Equal(t, filepath.Join(homeDir, "path/to/file"), expandedPath, "Should expand tilde to home directory")
		}
	})

	t.Run("Tilde not at the beginning", func(t *testing.T) {
		path := "/path/to/~file"
		expandedPath, err := expandTilde(path)
		require.NoError(t, err)
		assert.Equal(t, path, expandedPath, "Should not expand tilde if not at the beginning")
	})
}

func TestParseMailmapCommitEmailsAndName(t *testing.T) {
	// Test for ParseMailmapCommitEmailsAndName
	t.Run("Valid mailmap file", func(t *testing.T) {
		mailmapContent := `Proper Name <commitemail@example.com> <otheremail@example.com>
# A comment
Other Name <other@example.org>`
		tempDir, _, cleanup := createTempDirAndFile(t, "test_mailmap", ".mailmap", mailmapContent)
		defer cleanup()

		authors, err := ParseMailmapCommitEmailsAndName(tempDir)
		require.NoError(t, err)
		require.Len(t, authors, 2, "Should parse two authors")

		expectedAuthor1 := reportgenerator.Author{
			Name:   "Proper Name",
			Emails: map[string]bool{"commitemail@example.com": true, "otheremail@example.com": true},
		}
		expectedAuthor2 := reportgenerator.Author{
			Name:   "Other Name",
			Emails: map[string]bool{"other@example.org": true},
		}

		assert.Equal(t, expectedAuthor1, authors[0], "Should parse the first author correctly")
		assert.Equal(t, expectedAuthor2, authors[1], "Should parse the second author correctly")
	})

	t.Run("Empty mailmap file", func(t *testing.T) {
		tempDir, _, cleanup := createTempDirAndFile(t, "test_mailmap", ".mailmap", "")
		defer cleanup()

		authors, err := ParseMailmapCommitEmailsAndName(tempDir)
		require.NoError(t, err)
		assert.Empty(t, authors, "Should return an empty slice for an empty file")
	})

	t.Run("Mailmap file with comments and empty lines", func(t *testing.T) {
		mailmapContent := `# This is a comment
 
Name <email@example.com>
# Another comment`
		tempDir, _, cleanup := createTempDirAndFile(t, "test_mailmap", ".mailmap", mailmapContent)
		defer cleanup()

		authors, err := ParseMailmapCommitEmailsAndName(tempDir)
		require.NoError(t, err)
		require.Len(t, authors, 1, "Should parse one author, ignoring comments and empty lines")
		expectedAuthor := reportgenerator.Author{
			Name:   "Name",
			Emails: map[string]bool{"email@example.com": true},
		}

		assert.Equal(t, expectedAuthor, authors[0], "Should parse the author correctly")
	})

	t.Run("Invalid mailmap file", func(t *testing.T) {
		mailmapContent := `Invalid line` // missing email
		tempDir, _, cleanup := createTempDirAndFile(t, "test_mailmap", ".mailmap", mailmapContent)
		defer cleanup()

		_, err := ParseMailmapCommitEmailsAndName(tempDir)
		assert.Error(t, err, "Should return an error for an invalid mailmap line")
		assert.Contains(t, err.Error(), "Invalid mailmap line syntax", "Error should contain invalid syntax message")
	})
}

func TestParseMailmapLineCommitEmailsAndName(t *testing.T) {
	// Test for parseMailmapLineCommitEmailsAndName
	t.Run("Valid line with one email", func(t *testing.T) {
		line := "Proper Name <commitemail@example.com>"
		author, err := parseMailmapLineCommitEmailsAndName(line)
		require.NoError(t, err)
		expectedAuthor := reportgenerator.Author{
			Name:   "Proper Name",
			Emails: map[string]bool{"commitemail@example.com": true},
		}
		assert.Equal(t, expectedAuthor, author, "Should parse the line correctly")
	})

	t.Run("Valid line with multiple emails", func(t *testing.T) {
		line := "Proper Name <commitemail@example.com> <otheremail@example.com>"
		author, err := parseMailmapLineCommitEmailsAndName(line)
		require.NoError(t, err)
		expectedAuthor := reportgenerator.Author{
			Name:   "Proper Name",
			Emails: map[string]bool{"commitemail@example.com": true, "otheremail@example.com": true},
		}
		assert.Equal(t, expectedAuthor, author, "Should parse the line with multiple emails correctly")
	})

	t.Run("Line with extra spaces", func(t *testing.T) {
		line := "  Proper  Name  <commitemail@example.com>  "
		author, err := parseMailmapLineCommitEmailsAndName(line)
		require.NoError(t, err)
		expectedAuthor := reportgenerator.Author{
			Name:   "Proper Name",
			Emails: map[string]bool{"commitemail@example.com": true},
		}
		assert.Equal(t, expectedAuthor, author, "Should handle extra spaces correctly")
	})

	t.Run("Line with no email", func(t *testing.T) {
		line := "Invalid line"
		_, err := parseMailmapLineCommitEmailsAndName(line)
		assert.Error(t, err, "Should return an error for a line with no email")
		assert.Contains(t, err.Error(), "No commit emails found", "Error should contain no emails message")
	})

	t.Run("Empty line", func(t *testing.T) {
		line := ""
		_, err := parseMailmapLineCommitEmailsAndName(line)
		assert.Error(t, err, "Should return an error for an empty line")
		assert.Contains(t, err.Error(), "Empty line", "Error should contain empty line message")
	})
}
