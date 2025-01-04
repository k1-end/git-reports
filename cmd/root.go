package cmd

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"golang.org/x/term"
	"os"
	"sort"
	"strconv"
	"time"
)

type Tday struct {
	CommitCount int
	Date        time.Time
}

type Tmonth struct {
	Month time.Month
	Tdays map[int]Tday
}

type Tyear struct {
	Tmonths map[time.Month]Tmonth
	Year    int
}

func (y Tyear) getFirstMonth() (time.Month, error) {
	for i := 1; i < 13; i++ {
		_, ok := y.Tmonths[time.Month(i)]
		if ok {
			return time.Month(i), nil
		}
	}
	return time.Month(0), errors.New("empty name")
}

var shortDayNames = []string{
	"Sun",
	"Mon",
	"Tue",
	"Wed",
	"Thu",
	"Fri",
	"Sat",
}

var commitCountRange = []string{
	"0    ",
	"1-5  ",
	"6-10 ",
	"11-15",
	"16-20",
	"20<  ",
}

func commitCountGuide() {
	commitCount := 0
	pterm.DefaultBasicText.Println(pterm.Blue("commits count guide:"))
	for i := 0; i < len(commitCountRange); i++ {
		pterm.DefaultBasicText.Print(pterm.Blue(commitCountRange[i]))
		pterm.DefaultBasicText.Print(pterm.Yellow("=> "))
		commitCount = i * 5
		color := getColor(commitCount)
		var char string
		if commitCount == 0 {
			char = " . "
		} else {
			char = " * "
		}
		fmt.Printf("\x1b[48;2;%sm%s\x1b[0m", color, char)
		fmt.Println()
	}
	fmt.Println()
}

func (y Tyear) p() {
	fmt.Println()
	newHeader := pterm.HeaderPrinter{
		TextStyle:       pterm.NewStyle(pterm.FgBlack),
		BackgroundStyle: pterm.NewStyle(pterm.BgLightGreen),
		// Margin:          20,
	}

	newHeader.WithFullWidth().Println(y.Year)
	commitCountGuide()
	width, _, _ := term.GetSize(int(os.Stdout.Fd()))
	// border := strings.Repeat("-", width)
	// fmt.Println(border)
	monthW := 6
	offset := 1 + 5
	for {
		if (width+1-offset)%(monthW+1) == 0 {
			break
		}
		offset++
	}

	monthPerLine := (width + 1 - offset) / (monthW + 1)
	firstMonth, err := y.getFirstMonth()
	if err != nil {
		return
	}
	monthIndex := int(firstMonth)
	lineIndex := 1
	for monthIndex < 13 {
		fmt.Print("     ")
		for monthIndex-int(firstMonth)+1 <= lineIndex*monthPerLine && monthIndex < 13 {
			value, ok := y.Tmonths[time.Month(monthIndex)]
			if !ok {
				monthIndex = monthIndex + 1
				continue
			}
			fmt.Print("  ")
			pterm.DefaultBasicText.Print(pterm.Green(value.Month.String()[0:3]))
			pterm.DefaultBasicText.Print(pterm.Yellow(" |"))
			monthIndex = monthIndex + 1
		}
		fmt.Println()
		for i := 0; i < 7; i++ {
			monthIndex = monthPerLine*(lineIndex-1) + 1
			pterm.DefaultBasicText.Print(pterm.Blue(shortDayNames[i]))
			pterm.DefaultBasicText.Print(pterm.Yellow(": "))
			for monthIndex-int(firstMonth)+1 <= lineIndex*monthPerLine && monthIndex < 13 {
				value, ok := y.Tmonths[time.Month(monthIndex)]
				if !ok {
					monthIndex = monthIndex + 1
					continue
				}
				firstWeekDay := value.Tdays[1].Date.Weekday()
				for j := 1; j < 7; j++ {
					dayIndex := 7*j - int(firstWeekDay) - 6 + i
					d, ok := value.Tdays[dayIndex]
					if !ok {
						fmt.Print(" ")
						continue
					}
					color := getColor(d.CommitCount)
					var char string
					if d.CommitCount == 0 {
						char = "."
					} else {
						char = "*"
					}
					fmt.Printf("\x1b[48;2;%sm%s\x1b[0m", color, char)
				}
				pterm.DefaultBasicText.Print(pterm.Yellow("|"))
				monthIndex = monthIndex + 1
			}
			fmt.Println()
		}
		lineIndex = lineIndex + 1
	}
}

var developerEmail string

var rootCmd = &cobra.Command{
	Use:   "git-reports <path>",
	Short: "Visualize git reports",
	Long:  "Visualize git reports",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]

		r, err := git.PlainOpen(path)
		checkIfError(err)

		ref, err := r.Head()
		checkIfError(err)

		cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
		checkIfError(err)

		commits := make(map[string]int)
		commitsPerDev := make(map[string]int)
		var commitsPerHour [24]int

		_ = cIter.ForEach(func(c *object.Commit) error {
			if developerEmail != "_" {
				if c.Author.Email != developerEmail {
					return nil
				}
			}
			year, month, date := c.Author.When.Local().Date()
			key := fmt.Sprintf("%d-%d-%d", year, month, date)
			_, exists := commits[key]
			if !exists {
				commits[key] = 1
			} else {
				commits[key]++
			}

			_, exists = commitsPerDev[c.Author.Name]
			if !exists {
				commitsPerDev[c.Author.Name] = 1
			} else {
				commitsPerDev[c.Author.Name]++
			}

			commitsPerHour[c.Author.When.Local().Hour()]++

			return nil
		})

		if len(commits) == 0 {
			fmt.Println("No commits where found!")
			os.Exit(1)
		}

		keys := make([]string, 0, len(commits))

		for k := range commits {
			keys = append(keys, k)
		}

		sort.Slice(keys, func(i, j int) bool {
			timeI, err := time.Parse("2006-1-2", keys[i])
			checkIfError(err)
			timeJ, err := time.Parse("2006-1-2", keys[j])
			checkIfError(err)

			if timeI.Before(timeJ) {
				return true
			} else {
				return false
			}
		})

		firstDate, _ := time.Parse("2006-1-2", keys[0])
		startDate := time.Date(firstDate.Year(), firstDate.Month(), 1, 0, 0, 0, 0, firstDate.Location())

		lastDate, _ := time.Parse("2006-1-2", keys[len(keys)-1])
		endDate := time.Date(lastDate.Year(), lastDate.Month(), 1, 0, 0, 0, 0, lastDate.Location()).AddDate(0, 1, -1)

		years := make(map[int]Tyear)
		for startDate.Before(endDate) {
			_, exists := years[startDate.Year()]
			if !exists {
				years[startDate.Year()] = Tyear{Tmonths: make(map[time.Month]Tmonth), Year: startDate.Year()}
			}

			_, exists = years[startDate.Year()].Tmonths[startDate.Month()]
			if !exists {
				years[startDate.Year()].Tmonths[startDate.Month()] = Tmonth{Tdays: make(map[int]Tday), Month: startDate.Month()}
			}

			tDay, exists := years[startDate.Year()].Tmonths[startDate.Month()].Tdays[startDate.Day()]
			if !exists {
				years[startDate.Year()].Tmonths[startDate.Month()].Tdays[startDate.Day()] = Tday{}
			}
			tDay.CommitCount = commits[startDate.Format("2006-1-2")]
			tDay.Date = startDate

			years[startDate.Year()].Tmonths[startDate.Month()].Tdays[startDate.Day()] = tDay

			startDate = startDate.AddDate(0, 0, 1)
		}

		yearsKey := make([]int, 0, len(years))
		for k := range years {
			yearsKey = append(yearsKey, k)
		}
		sort.Ints(yearsKey)
		for _, k := range yearsKey {
			years[k].p()
		}

		authorNames := make([]string, 0, len(commitsPerDev))

		for k := range commitsPerDev {
			authorNames = append(authorNames, k)
		}

		sort.SliceStable(authorNames, func(i, j int) bool {
			return commitsPerDev[authorNames[i]] > commitsPerDev[authorNames[j]]
		})

		var barData []pterm.Bar
		var bar pterm.Bar
		for _, authorName := range authorNames {
			bar.Label = authorName
			bar.Value = commitsPerDev[authorName]
			barData = append(barData, bar)
		}

		fmt.Println()
		pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgGreen)).Println("Commit count per developer")
		err = pterm.DefaultBarChart.WithBars(barData).WithHorizontal().WithWidth(90).WithShowValue().Render()
		checkIfError(err)

		var hourData []pterm.Bar
		var hour pterm.Bar
		for i := 1; i < 24; i++ {
			hour.Label = strconv.Itoa(i)
			hour.Value = commitsPerHour[i]
			hourData = append(hourData, hour)
		}
		pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgGreen)).Println("Commits per hour of day (local)")
		err = pterm.DefaultBarChart.WithShowValue().WithBars(hourData).WithHorizontal().WithWidth(100).Render()
		checkIfError(err)
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
