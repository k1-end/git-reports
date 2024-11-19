package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"golang.org/x/term"
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

func (y Tyear) p() {
	fmt.Println(y.Year)
	width, _, _ := term.GetSize(0)
	border := strings.Repeat("-", width)
	fmt.Println(border)
	monthW := 7
	offset := 0
	for true {
		if (width+1-offset)%(monthW+1) == 0 {
			break
		}
		offset++
	}
	for m := time.January; m <= time.December; m++ {
		value, ok := y.Tmonths[m]
		if !ok {
			continue
		}
		fmt.Print("  ")
		fmt.Print(value.Month.String()[0:3])
		fmt.Print("  |")
	}

	fmt.Println()

	for i := 0; i < 5; i++ {
		for m := time.January; m <= time.December; m++ {
			value, ok := y.Tmonths[m]
			if !ok {
				continue
			}
			for j := i*7 + 1; j < (i+1)*7+1; j++ {
				d, ok := value.Tdays[j]
				if !ok {
					fmt.Print(" ")
					continue
				}
                color := getColor(d.CommitCount)
				if d.CommitCount == 0 {
					fmt.Printf("\x1b[48;2;%sm.\x1b[0m", color)
					continue
				}
				fmt.Printf("\x1b[48;2;%sm*\x1b[0m", color)
			}
			fmt.Print("|")
		}
		fmt.Println()
	}
}



func main() {
	if len(os.Args) < 3 {
	    fmt.Printf("\x1b[36;1m%s\x1b[0m\n", fmt.Sprintf("Usage: %s <path> <author-email>", os.Args[0]))
		os.Exit(1)
	}
	path := os.Args[1]
	mainAuthorEmail := os.Args[2]

	r, err := git.PlainOpen(path)
	checkIfError(err)


	ref, err := r.Head()
	checkIfError(err)

	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	checkIfError(err)

	commits := make(map[string]int)

	err = cIter.ForEach(func(c *object.Commit) error {
		if c.Author.Email == mainAuthorEmail {
			year, month, date := c.Author.When.Local().Date()
			key := fmt.Sprintf("%d-%d-%d", year, month, date)
			_, exists := commits[key]
			if !exists {
				commits[key] = 1
			} else {
				commits[key]++
			}
		}
		return nil
	})

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

	firstDate, err := time.Parse("2006-1-2", keys[0])
	startDate := time.Date(firstDate.Year(), firstDate.Month(), 1, 0, 0, 0, 0, firstDate.Location())

	lastDate, err := time.Parse("2006-1-2", keys[len(keys)-1])
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

	for _, year := range years {
		year.p()
	}

}
