package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/k1-end/git-visualizer/src/reportgenerator"
	"github.com/k1-end/git-visualizer/src/reportprinter"
	"github.com/spf13/cobra"
)


var developerEmail string
var fromDate string
var toDate string
var printerOption string

var rootCmd = &cobra.Command{
	Use:   "git-reports [path]",
	Short: "Visualize git reports",
	Long:  "Visualize git repository at path (default to current directory)",
	Args:  cobra.RangeArgs(0, 1),

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


        commitCountDateHeatMapGenerator := reportgenerator.CommitCountDateHeatMapGenerator{CommitsMap: make(map[string]int) }
        commitsPerDevReportGenerator := reportgenerator.CommitsPerDevReportGenerator{CommitsPerDevMap: make(map[string]int)}
        commitsPerHourReportGenerator := reportgenerator.CommitsPerHourReportGenerator{CommitsPerHourMap: make([]int, 24)}
        mergeCommitsPerYearReportGenerator := reportgenerator.MergeCommitsPerYearReportGenerator{MergeCommitsPerYearMap: make(map[int]int)}

		path := "."
        if len(args) == 1 {
		    path = args[0]
        }

        r, err := git.PlainOpen(path)
        if errors.Is(err, git.ErrRepositoryNotExists) {
            fmt.Println("The provided path is not a git repository: " + path) // no model found for id
            os.Exit(1)
        }
		checkIfError(err)

		ref, err := r.Head()
		checkIfError(err)

		cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
		checkIfError(err)


		_ = cIter.ForEach(func(c *object.Commit) error {
			if developerEmail != "_" {
				if c.Author.Email != developerEmail {
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

            commitCountDateHeatMapGenerator.IterationStep(c)
            commitsPerDevReportGenerator.IterationStep(c)
            commitsPerHourReportGenerator.IterationStep(c)
            mergeCommitsPerYearReportGenerator.IterationStep(c)

			return nil
		})
        
		headRef, err := r.Head()
		commit, err := r.CommitObject(headRef.Hash())

        fileTypeReportGenerator := reportgenerator.FileTypeReportGenerator{FileTypeMap: make(map[string]int)}
        fileTypeReportGenerator.Iterate(commit)


        absolutePath, _ := filepath.Abs(path)
        dirName := filepath.Base(absolutePath)

        if printerOption == "" || printerOption == "console" {
            p := reportprinter.ConsolePrinter{}
            p.RegisterReport(commitCountDateHeatMapGenerator.GetReport())
            p.RegisterReport(commitsPerDevReportGenerator.GetReport())
            p.RegisterReport(commitsPerHourReportGenerator.GetReport())
            p.RegisterReport(mergeCommitsPerYearReportGenerator.GetReport())
            p.RegisterReport(fileTypeReportGenerator.GetReport())
            p.SetProjectTitle(dirName)
            p.Print()
        }else if printerOption == "html" {
            p := reportprinter.HtmlPrinter{}
            p.RegisterReport(commitCountDateHeatMapGenerator.GetReport())
            p.RegisterReport(commitsPerDevReportGenerator.GetReport())
            p.RegisterReport(commitsPerHourReportGenerator.GetReport())
            p.RegisterReport(mergeCommitsPerYearReportGenerator.GetReport())
            p.RegisterReport(fileTypeReportGenerator.GetReport())
            p.SetProjectTitle(dirName)
            p.Print()
        }else {
            fmt.Println("Invalid printer value. Valid values are `console` and `html`")
            os.Exit(1)
        }
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&developerEmail, "dev", "d", "_", "choose developer by email")
    rootCmd.PersistentFlags().StringVarP(&fromDate, "from", "f", "", "Filter commits from this date (format: YYYY-MM-DD)")
	rootCmd.PersistentFlags().StringVarP(&toDate, "to", "t", "", "Filter commits up to this date (format: YYYY-MM-DD)")
	rootCmd.PersistentFlags().StringVarP(&printerOption, "printer", "p", "console", "Printer (default to console) (available options are console and html")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, err = fmt.Fprintln(os.Stderr, err)
		checkIfError(err)
		os.Exit(1)
	}
}
