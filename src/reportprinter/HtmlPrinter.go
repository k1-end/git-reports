package reportprinter

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"math/rand/v2"
	"os"
	"time"

	"github.com/k1-end/git-visualizer/src/report"
)

//go:embed templates/*
var templatesFS embed.FS

type HtmlPrinter struct {
	BasePrinter
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

	type dateData struct {
		Date        time.Time
		CommitCount int
	}
	yearData := make(map[int]map[time.Month]map[int]dateData)

	counter := 0
	for startDate.Before(endDate) {
		year := startDate.Year()
		month := startDate.Month()
		day := startDate.Day()

		if _, exists := yearData[year]; !exists {
			yearData[year] = make(map[time.Month]map[int]dateData)
		}
		if _, exists := yearData[year][month]; !exists {
			yearData[year][month] = make(map[int]dateData)
		}

		commitCount := 0
		if counter < len(keys) && startDate.Format("2006-1-2") == keys[counter] {
			commitCount = data[counter].IntValue
			counter++
		}

		yearData[year][month][day] = dateData{
			Date:        startDate,
			CommitCount: commitCount,
		}

		startDate = startDate.AddDate(0, 0, 1)
	}

	tmpl, err := template.New("date-heatmap.html").ParseFS(templatesFS, "templates/date-heatmap.html")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	var anon struct {
		YearData  map[int]map[time.Month]map[int]dateData
		FirstDate time.Time
		Range     int
	}
	anon.YearData = yearData
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
