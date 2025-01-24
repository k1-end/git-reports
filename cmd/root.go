package cmd

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
	"github.com/k1-end/git-visualizer/src"
)


var developerEmail string

var rootCmd = &cobra.Command{
	Use:   "git-reports <path>",
	Short: "Visualize git reports",
	Long:  "Visualize git reports",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]

        heatMapReport := report.HeatMapReport{CommitsMap: make(map[string]int) }
        commitsPerDevReport := report.CommitsPerDevReport{CommitsPerDevMap: make(map[string]int)}
        commitsPerHourReport := report.CommitsPerHourReport{CommitsPerHourMap: make([]int, 24)}
        mergeCommitsPerYearReport := report.MergeCommitsPerYearReport{MergeCommitsPerYearMap: make(map[int]int)}

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

            heatMapReport.IterationStep(c)
            commitsPerDevReport.IterationStep(c)
            commitsPerHourReport.IterationStep(c)
            mergeCommitsPerYearReport.IterationStep(c)

			return nil
		})

        heatMapReport.Print()
		fmt.Println()
        commitsPerDevReport.Print()
		fmt.Println()
        commitsPerHourReport.Print()
		fmt.Println()
        mergeCommitsPerYearReport.Print()

		headRef, err := r.Head()
		commit, err := r.CommitObject(headRef.Hash())

        fileTypeReport := report.FileTypeReport{FileTypeMap: make(map[string]int)}
        fileTypeReport.Iterate(commit)
        fileTypeReport.Print()
	},
}

func Execute() {
	rootCmd.PersistentFlags().StringVar(&developerEmail, "dev", "_", "choose developer by email")

	if err := rootCmd.Execute(); err != nil {
		_, err = fmt.Fprintln(os.Stderr, err)
		checkIfError(err)
		os.Exit(1)
	}
}
