package reportprinter

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/k1-end/git-visualizer/src/report"
	"github.com/pterm/pterm"
)

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
	"  .0.  ",
	" *1-5* ",
	"*06-10*",
	"*11-15*",
	"*16-20*",
	" *20<* ",
}

func getColor(commitCount int) string {
	if commitCount == 0 {
		return "178;215;155" // #B2D79B (Light Green)
	} else if commitCount <= 5 {
		return "139;195;74" // #8BC34A (Medium Green)
	} else if commitCount <= 10 {
		return "34;139;34" // #228B22 (Forest Green)
	} else if commitCount <= 15 {
		return "0;100;0" // #006400 (Dark Green)
	} else if commitCount <= 20 {
		return "0;128;128" // #008080 (Emerald Green)
	} else {
		return "0;64;0" // #004000 (Darker Green)
	}
}

func commitCountGuide() {
	fmt.Println()
	commitCount := 0
	pterm.DefaultBasicText.Print(pterm.Blue("commits count guide:"))
	for i, char := range commitCountRange {
		commitCount = i * 5
		color := getColor(commitCount)
		pterm.DefaultBasicText.Printf(" \x1b[48;2;%sm%s\x1b[0m ", color, pterm.Red(char))
	}
	fmt.Println()
}

type ConsolePrinter struct {
	BasePrinter
}

func (y yearData) getFirstMonth() (time.Month, error) {
	for i := 1; i < 13; i++ {
		if _, ok := y.Months[time.Month(i)]; ok {
			return time.Month(i), nil
		}
	}
	return time.Month(0), errors.New("empty name")
}

func (y yearData) print() {
	fmt.Println()
	newHeader := pterm.HeaderPrinter{
		TextStyle:       pterm.NewStyle(pterm.FgBlack),
		BackgroundStyle: pterm.NewStyle(pterm.BgLightGreen),
	}

	newHeader.WithFullWidth().Println(y.Year)
	width := pterm.GetTerminalWidth()
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
			if _, ok := y.Months[time.Month(monthIndex)]; ok {
				fmt.Print("  ")
				pterm.DefaultBasicText.Print(pterm.Green(time.Month(monthIndex).String()[0:3]))
				pterm.DefaultBasicText.Print(pterm.Yellow(" |"))
			}
			monthIndex++
		}
		fmt.Println()

		for i := 0; i < 7; i++ {
			monthIndex = monthPerLine*(lineIndex-1) + 1
			pterm.DefaultBasicText.Print(pterm.Blue(shortDayNames[i]))
			pterm.DefaultBasicText.Print(pterm.Yellow(": "))

			for monthIndex-int(firstMonth)+1 <= lineIndex*monthPerLine && monthIndex < 13 {
				if monthData, ok := y.Months[time.Month(monthIndex)]; ok {
					firstWeekDay := monthData[1].Date.Weekday()
					for j := 1; j < 7; j++ {
						dayIndex := 7*j - int(firstWeekDay) - 6 + i
						if dayData, ok := monthData[dayIndex]; ok {
							color := getColor(dayData.CommitCount)
							char := "."
							if dayData.CommitCount > 0 {
								char = "*"
							}
							fmt.Printf("\x1b[48;2;%sm%s\x1b[0m", color, char)
						} else {
							fmt.Print(" ")
						}
					}
					pterm.DefaultBasicText.Print(pterm.Yellow("|"))
				}
				monthIndex++
			}
			fmt.Println()
		}
		lineIndex++
	}
}

func (p ConsolePrinter) printBarChart(c report.Report) {
	var barData []pterm.Bar
	var bar pterm.Bar
	labels := c.GetLabels()
	data := c.GetData()
	for i := 0; i < len(labels); i++ {
		bar.Label = labels[i]
		bar.Value = data[i].IntValue
		barData = append(barData, bar)
	}

	pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgGreen)).Println(c.GetTitle())
	_ = pterm.DefaultBarChart.WithBars(barData).WithHorizontal().WithWidth(90).WithShowValue().Render()
}

func (p ConsolePrinter) printTable(r report.Report) {
    tableData := pterm.TableData{}
    tableData = append(tableData, []string{"", ""})
    labels := r.GetLabels()
    for index, data := range r.GetData() {
        label := labels[index]
        var value string
        switch data.IsInt {
        case true:
            value = strconv.Itoa(data.IntValue)
        case false:
            value = data.StringValue
        }
        tableData = append(tableData, []string{label, value})
    }
    pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgGreen)).Println(r.GetTitle())
    pterm.DefaultTable.WithHasHeader().WithBoxed().WithData(tableData).Render()
}

func (p ConsolePrinter) printDateHeatMapChart(c report.Report) {
	keys := c.GetLabels()
	data := c.GetData()
	if len(data) == 0 {
		fmt.Println("No commits where found!")
		return
	}

	// Parse the first commit date from the input data
	firstDate, _ := time.Parse("2006-1-2", keys[0])
	// Set startDate to beginning of the month containing the first commit
	startDate := time.Date(firstDate.Year(), firstDate.Month(), 1, 0, 0, 0, 0, firstDate.Location())

	// Parse the last commit date from the input data
	lastDate, _ := time.Parse("2006-1-2", keys[len(keys)-1])
	// Set endDate to end of the month containing the last commit
	endDate := time.Date(lastDate.Year(), lastDate.Month(), 1, 0, 0, 0, 0, lastDate.Location()).AddDate(0, 1, -1)

	years := make(map[int]yearData)
	counter := 0

	for startDate.Before(endDate) {
		year := startDate.Year()
		month := startDate.Month()
		day := startDate.Day()

		if _, exists := years[year]; !exists {
			years[year] = yearData{
				Year: year,
				Months: make(map[time.Month]map[int]struct {
					Date        time.Time
					CommitCount int
				}),
			}
		}

		if _, exists := years[year].Months[month]; !exists {
			years[year].Months[month] = make(map[int]struct {
				Date        time.Time
				CommitCount int
			})
		}

		commitCount := 0
		if counter < len(keys) && startDate.Format("2006-1-2") == keys[counter] {
			commitCount = data[counter].IntValue
			counter++
		}

		years[year].Months[month][day] = struct {
			Date        time.Time
			CommitCount int
		}{
			Date:        startDate,
			CommitCount: commitCount,
		}

		startDate = startDate.AddDate(0, 0, 1)
	}

	yearsKey := make([]int, 0, len(years))
	for k := range years {
		yearsKey = append(yearsKey, k)
	}
	sort.Ints(yearsKey)

	for _, k := range yearsKey {
		years[k].print()
	}
	commitCountGuide()
}

func (p *ConsolePrinter) Print(s *os.File) {
    pterm.FallbackTerminalWidth = 100
	for k := range p.reports {
		switch p.reports[k].GetReportType() {
		case "bar_chart":
			p.printBarChart(p.reports[k])
		case "date_heatmap":
			p.printDateHeatMapChart(p.reports[k])
        case "table":
            p.printTable(p.reports[k])
		}
	}
}
