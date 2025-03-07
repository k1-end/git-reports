package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/k1-end/git-visualizer/cmd/serve"
	"github.com/k1-end/git-visualizer/src/reportgenerator"
	"github.com/k1-end/git-visualizer/src/reportprinter"
	"github.com/spf13/cobra"
)


var developerEmail string
var fromDate string
var toDate string

var rootCmd = &cobra.Command{
	Use:   "git-reports <path>",
	Short: "Visualize git reports",
	Long:  "Visualize git reports",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]

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

		r, err := git.PlainOpen(path)
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
        p := reportprinter.ConsolePrinter {}
        commitCountDateHeatMap:= commitCountDateHeatMapGenerator.GetReport()
        p.PrintDateHeatMapChart(commitCountDateHeatMap)
        commitsPerDevReport := commitsPerDevReportGenerator.GetReport()
        p.PrintLineChart(commitsPerDevReport)
        commitsPerHourReport := commitsPerHourReportGenerator.GetReport()
        p.PrintLineChart(commitsPerHourReport)
        mergeCommitsPerYearReport := mergeCommitsPerYearReportGenerator.GetReport()
        p.PrintLineChart(mergeCommitsPerYearReport)

		headRef, err := r.Head()
		commit, err := r.CommitObject(headRef.Hash())

        fileTypeReportGenerator := reportgenerator.FileTypeReportGenerator{FileTypeMap: make(map[string]int)}
        fileTypeReportGenerator.Iterate(commit)
        fileTypeReport := fileTypeReportGenerator.GetReport()
        p.PrintLineChart(fileTypeReport)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&developerEmail, "dev", "_", "choose developer by email")
    rootCmd.PersistentFlags().StringVar(&fromDate, "from", "", "Filter commits from this date (format: YYYY-MM-DD)")
	rootCmd.PersistentFlags().StringVar(&toDate, "to", "", "Filter commits up to this date (format: YYYY-MM-DD)")

    rootCmd.AddCommand(serve.ServeCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, err = fmt.Fprintln(os.Stderr, err)
		checkIfError(err)
		os.Exit(1)
	}
}
