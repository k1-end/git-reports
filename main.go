package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	. "github.com/go-git/go-git/v5/_examples"
	"github.com/go-git/go-git/v5/plumbing/object"
	"golang.org/x/crypto/ssh/terminal"
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
	width, _, _ := terminal.GetSize(0)
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
				if d.CommitCount == 0 {
					fmt.Printf("\x1b[48;2;46;54;45m*\x1b[0m")
					continue
				}
				fmt.Printf("\x1b[48;2;56;232;21m*\x1b[0m")
			}
			fmt.Print("|")
		}
		fmt.Println()
	}
}

func main() {
	CheckArgs("<path>", "<autho-email>")
	path := os.Args[1]
	mainAuthorEmail := os.Args[2]

	r, err := git.PlainOpen(path)
	CheckIfError(err)


	ref, err := r.Head()
	CheckIfError(err)

	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	CheckIfError(err)

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
		iDate := strings.Split(keys[i], "-")
		jDate := strings.Split(keys[j], "-")

		iYear, err := strconv.Atoi(iDate[0])
		if err != nil {
			fmt.Println("Failed to convert string to integer", err)
			panic("Failed to convert string to integer")
		}

		jYear, err := strconv.Atoi(jDate[0])
		if err != nil {
			fmt.Println("Failed to convert string to integer", err)
			panic("Failed to convert string to integer")
		}

		if iYear < jYear {
			return true
		} else if iYear > jYear {
			return false
		}

		iMonth, err := strconv.Atoi(iDate[1])
		if err != nil {
			fmt.Println("Failed to convert string to integer", err)
			panic("Failed to convert string to integer")
		}

		jMonth, err := strconv.Atoi(jDate[1])
		if err != nil {
			fmt.Println("Failed to convert string to integer", err)
			panic("Failed to convert string to integer")
		}

		if iMonth < jMonth {
			return true
		} else if iMonth > jMonth {
			return false
		}

		iDay, err := strconv.Atoi(iDate[2])
		if err != nil {
			fmt.Println("Failed to convert string to integer", err)
			panic("Failed to convert string to integer")
		}

		jDay, err := strconv.Atoi(jDate[2])
		if err != nil {
			fmt.Println("Failed to convert string to integer", err)
			panic("Failed to convert string to integer")
		}

		if iDay < jDay {
			return true
		} else if iDay > jDay {
			return false
		}
		return false
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
