package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/k1-end/git-reports/src/reportgenerator"
	"github.com/k1-end/git-reports/src/reportprinter"
	"github.com/spf13/cobra"
)

var developerEmail string
var fromDate string
var toDate string
var path string
var printerOption string
var outputPath string
var branch string
var Version string

var authors = make(map[string]*reportgenerator.Author)

var rootCmd = &cobra.Command{
	Use:   "git-reports [options]",
	Short: "Visualize git reports",
	Long:  "Visualize git repository at path (default to current directory)",
	Args:  cobra.ExactArgs(0),

	Run: func(cmd *cobra.Command, args []string) {

		var fromTime, toTime time.Time
		var err error
		if fromDate != "" {
			fromTime, err = time.Parse("2006-01-02", fromDate)
			if err != nil {
				fmt.Println("Invalid 'from' date format. Please use YYYY-MM-DD.")
				os.Exit(1)
			}
		}

		if toDate != "" {
			toTime, err = time.Parse("2006-01-02", toDate)
			if err != nil {
				fmt.Println("Invalid 'to' date format. Please use YYYY-MM-DD.")
				os.Exit(1)
			}
		}

		// Ensure 'from' is before 'to'
		if fromDate != "" && toDate != "" && fromTime.After(toTime) {
			fmt.Println("'from' date must be before 'to' date.")
			os.Exit(1)
		}

        if outputPath != ""  {
            if !isValidFilePath(outputPath){
                fmt.Println("The given output is not a valid file path or is not writable")
                os.Exit(1)
            }
        }

        mailmapAuthors, err := ParseMailmapCommitEmailsAndName(path)
        checkIfError(err)
        for _, mailmapAuthor := range mailmapAuthors {
            for email := range mailmapAuthor.Emails {
                authors[email] = &mailmapAuthor
                
            }
        }

		commitCountDateHeatMapGenerator := reportgenerator.CommitCountDateHeatMapGenerator{CommitsMap: make(map[string]int)}
		commitsPerDevReportGenerator := reportgenerator.CommitsPerDevReportGenerator{CommitsPerDevMap: make(map[string]int)}
		commitsPerHourReportGenerator := reportgenerator.CommitsPerHourReportGenerator{CommitsPerHourMap: make([]int, 24)}
		mergeCommitsPerYearReportGenerator := reportgenerator.MergeCommitsPerYearReportGenerator{MergeCommitsPerYearMap: make(map[int]int)}
		fileTypeReportGenerator := reportgenerator.FileTypeReportGenerator{FileTypeMap: make(map[string]int)}
		generalInfoReportGenerator := reportgenerator.GeneralInfoReportGenerator{}


		r, err := git.PlainOpen(path)
		if errors.Is(err, git.ErrRepositoryNotExists) {
			fmt.Println("The provided path is not a git repository: " + path) // no model found for id
			os.Exit(1)
		}
		checkIfError(err)

		ref, err := r.Head()

        if branch != ""  {
            refName := plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch))

            ref, err = r.Reference(refName, false)
            if err != nil {
                if err == plumbing.ErrReferenceNotFound {
                    fmt.Printf("Local branch '%s' does not exist\n", branch)
                } else {
                    fmt.Printf("Error checking local branch '%s': %v\n", branch, err)
                }
                os.Exit(1)
            }
        }

		checkIfError(err)

		cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
		checkIfError(err)

		_ = cIter.ForEach(func(c *object.Commit) error {
            if _, exists := authors[c.Author.Email]; !exists {
                authors[c.Author.Email] = &reportgenerator.Author{
                    Name: c.Author.Name,
                    Emails: map[string]bool{c.Author.Email: true},
                }
            }

			if developerEmail != "_" {
                selectedAuthor, _ := authors[developerEmail];
                if selectedAuthor != authors[c.Author.Email] {
                    return nil
                }
			}

			// Filter by date range
			commitTime := c.Author.When
			if fromDate != "" && commitTime.Before(fromTime) {
				return nil
			}
			if toDate != "" && commitTime.After(toTime) {
				return nil
			}

			commitCountDateHeatMapGenerator.LogIterationStep(c, *authors[c.Author.Email])
			commitsPerDevReportGenerator.LogIterationStep(c, *authors[c.Author.Email])
			commitsPerHourReportGenerator.LogIterationStep(c, *authors[c.Author.Email])
			mergeCommitsPerYearReportGenerator.LogIterationStep(c, *authors[c.Author.Email])
            generalInfoReportGenerator.LogIterationStep(c, *authors[c.Author.Email])

			return nil
		})

		headRef, err := r.Head()
		checkIfError(err)
		commit, err := r.CommitObject(headRef.Hash())
		checkIfError(err)

        fIter, _ := commit.Files()
        fIter.ForEach(func(f *object.File) error {
            fileTypeReportGenerator.FileIterationStep(f)
            generalInfoReportGenerator.FileIterationStep(f)
            return nil
        })

		absolutePath, _ := filepath.Abs(path)
		dirName := filepath.Base(absolutePath)

		p := getPrinter(printerOption)
		p.RegisterReport(generalInfoReportGenerator.GetReport())
		p.RegisterReport(commitCountDateHeatMapGenerator.GetReport())
		p.RegisterReport(commitsPerDevReportGenerator.GetReport())
		p.RegisterReport(commitsPerHourReportGenerator.GetReport())
		p.RegisterReport(mergeCommitsPerYearReportGenerator.GetReport())
		p.RegisterReport(fileTypeReportGenerator.GetReport())
		p.SetProjectTitle(dirName)
        if outputPath != "" {
            destination, err := os.Create(outputPath)
            if err != nil {
                fmt.Println("os.Create:", err)
                return
            }
            defer destination.Close()
            p.Print(destination)
        }else{
            p.Print(os.Stdout)

        }
	},
}

func getPrinter(printerOption string) reportprinter.Printer {
	if printerOption == "" || printerOption == "console" {
		return &reportprinter.ConsolePrinter{}
	} else if printerOption == "html" {
		return &reportprinter.HtmlPrinter{}
	} else {
		fmt.Println("Invalid printer value. Valid values are `console` and `html`")
		os.Exit(1)
	}
	return nil
}

func init() {
    rootCmd.PersistentFlags().StringVarP(&path, "path", "p", ".", "Repository path (default to current directory)")
    rootCmd.PersistentFlags().StringVarP(&developerEmail, "dev", "d", "_", "choose developer by email")
    rootCmd.PersistentFlags().StringVarP(&fromDate, "from", "f", "", "Filter commits from this date (format: YYYY-MM-DD)")
    rootCmd.PersistentFlags().StringVarP(&toDate, "to", "t", "", "Filter commits up to this date (format: YYYY-MM-DD)")
    rootCmd.PersistentFlags().StringVarP(&branch, "branch", "b", "", "Set the branch to analyze")
    rootCmd.PersistentFlags().StringVar(&printerOption, "printer", "console", "Printer (default to console) (available options are console and html)")
    rootCmd.PersistentFlags().StringVar(&outputPath, "output", "", "Output path for the report")

    rootCmd.Flags().BoolP("version", "v", false, "Print the version") // Subcommands do not automatically inherit this flag
    rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
        versionFlag, err := cmd.Flags().GetBool("version")
        if err != nil {
            return err
        }

        if versionFlag {
            fmt.Println("Version:", Version)
            os.Exit(0)
        }
        return nil
    }
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, err = fmt.Fprintln(os.Stderr, err)
		checkIfError(err)
		os.Exit(1)
	}
}

// isValidFilePath checks if the given path is valid for creating a file.
// It checks if the parent directory exists and is writable.
func isValidFilePath(filePath string) bool {
        // Expand tilde
        expandedPath, err := expandTilde(filePath)
        if err != nil {
                return false
        }

        // Get the directory part of the path.
        dir := filepath.Dir(expandedPath)

        // Check if the directory exists.
        fileInfo, err := os.Stat(dir)
        if os.IsNotExist(err) {
                return false // Directory does not exist
        }

        if err != nil {
                return false // Other error during Stat
        }

        // Check if it's a directory.
        if !fileInfo.IsDir() {
                return false // Parent is not a directory
        }

        // Check write permissions for the directory.
        testFile := filepath.Join(dir, ".testwrite")
        err = os.WriteFile(testFile, []byte("test"), 0600)
        if err != nil {
                return false // Write permission denied
        }

        os.Remove(testFile)

        return true
}

// expandTilde expands a tilde (~) in a file path to the user's home directory.
func expandTilde(path string) (string, error) {
        if len(path) > 0 && path[0] == '~' {
                homeDir, err := os.UserHomeDir()
                if err != nil {
                        return "", err
                }
                return filepath.Join(homeDir, path[1:]), nil
        }
        return path, nil
}

func ParseMailmapCommitEmailsAndName(path string) ([]reportgenerator.Author, error) {
        file, err := os.Open(path + "/.mailmap")
        if err != nil {
                return nil, nil // the .mailmap file is not required
        }
        defer file.Close()

        scanner := bufio.NewScanner(file)
        authors := []reportgenerator.Author{}
        lineNum := 0

        for scanner.Scan() {
                lineNum++
                line := scanner.Text()
                line = strings.TrimSpace(line)

                if line == "" || strings.HasPrefix(line, "#") {
                        continue
                }

                author, err := parseMailmapLineCommitEmailsAndName(line)
                if err != nil {
                        return nil, errors.New(err.Error() + " on line " + string(rune(lineNum + '0')))
                }
                authors = append(authors, author)
        }

        if err := scanner.Err(); err != nil {
                return nil, err
        }

        return authors, nil
}

func parseMailmapLineCommitEmailsAndName(line string) (reportgenerator.Author, error) {
        author := reportgenerator.Author{
                Emails: make(map[string]bool),
        }

        parts := strings.Fields(line)

        if len(parts) == 0 {
                return author, errors.New("Invalid mailmap line syntax: Empty line")
        }

        emailFound := false;
        properNameParts := []string{};

        for _, part := range parts {
                if strings.HasPrefix(part, "<") && strings.HasSuffix(part, ">") {
                        author.Emails[strings.Trim(part, "<>")] = true
                        emailFound = true;
                } else{
                        properNameParts = append(properNameParts, part);
                }
        }

        author.Name = strings.Join(properNameParts, " ");

        if !emailFound {
                return author, errors.New("Invalid mailmap line syntax: No commit emails found")
        }

        return author, nil
}
