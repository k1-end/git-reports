package reportprinter

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"math/rand/v2"
	"os"
	"sort"
	"time"

	"github.com/k1-end/git-visualizer/src/report"
)

//go:embed templates/*
var templatesFS embed.FS

type HtmlPrinter struct {
	reports      []report.Report
	projectTitle string
}

func (p HtmlPrinter) renderDateHeatMapChart(c report.Report) string {
	keys := c.GetLabels()
	data := c.GetData()
	if len(data) == 0 {
		return ""
	}
	firstDate, _ := time.Parse("2006-1-2", keys[0])
	startDate := time.Date(firstDate.Year(), firstDate.Month(), 1, 0, 0, 0, 0, firstDate.Location())

	lastDate, _ := time.Parse("2006-1-2", keys[len(keys)-1])
	endDate := time.Date(lastDate.Year(), lastDate.Month(), 1, 0, 0, 0, 0, lastDate.Location()).AddDate(0, 1, -1)

	years := make(map[int]Tyear)
	counter := 0
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
		tDay.Date = startDate
		if counter < len(keys) && startDate.Format("2006-1-2") == keys[counter] {
			tDay.CommitCount = data[counter].IntValue
			counter += 1 // data does not contain all dates and we are iterating overall dates, so we must increment only when the date matches
		} else {
			tDay.CommitCount = 0
		}
		years[startDate.Year()].Tmonths[startDate.Month()].Tdays[startDate.Day()] = tDay
		startDate = startDate.AddDate(0, 0, 1)
	}

	yearsKey := make([]int, 0, len(years))
	for k := range years {
		yearsKey = append(yearsKey, k)
	}
	sort.Ints(yearsKey)
	tmpl, err := template.New("date-heatmap.html").ParseFS(templatesFS, "templates/date-heatmap.html")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	var anon struct {
		Years     map[int]Tyear
		FirstDate time.Time
		Range     int
	}
	anon.Years = years
	anon.FirstDate = firstDate
	anon.Range = endDate.Year() - firstDate.Year() + 1

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, anon)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return buf.String()
}

func (p *HtmlPrinter) RegisterReport(r report.Report) {
	p.reports = append(p.reports, r)
}

func (p HtmlPrinter) renderLineChart(c report.Report) string {
	tmpl, err := template.New("line-chart.html").ParseFS(templatesFS, "templates/line-chart.html")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	var anon struct {
		Title     string
		Labels    []string
		Data      []int
		ElementId int
	}
	anon.Title = c.GetTitle()
	anon.Labels = c.GetLabels()
	var data []int
	for k := range c.GetData() {
		data = append(data, c.GetData()[k].IntValue)
	}
	anon.Data = data
	anon.ElementId = rand.IntN(10000000)
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, anon)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return buf.String()
}

func (p HtmlPrinter) Print() {
	tmpl, err := template.New("main.html").ParseFS(templatesFS, "templates/main.html")

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	var renderedReports bytes.Buffer
	for k := range p.reports {
		switch p.reports[k].GetReportType() {
		case "date_heatmap":
			renderedReports.WriteString(p.renderDateHeatMapChart(p.reports[k]))
		case "line_chart":
			renderedReports.WriteString(p.renderLineChart(p.reports[k]))
		}
		renderedReports.WriteString("\n")
	}

	var anon struct {
		ProjectTitle    string
		RenderedReports template.HTML
	}
	anon.ProjectTitle = p.GetProjectTitle()
	anon.RenderedReports = template.HTML(renderedReports.String())
	err = tmpl.Execute(os.Stdout, anon)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func (p *HtmlPrinter) SetProjectTitle(s string) {
	p.projectTitle = s
}

func (p *HtmlPrinter) GetProjectTitle() string {
	return p.projectTitle
}
