package reportprinter

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"os"
	"strconv"
	"time"

	"github.com/k1-end/git-visualizer/src/report"
)

//go:embed templates/*
var templatesFS embed.FS

type HtmlPrinter struct {
	BasePrinter
}

func (p HtmlPrinter) renderTable(r report.Report, elementId int) string {
	var anon struct {
        Rows map[string]string
        ElementId int
        Title string
	}
    anon.ElementId = elementId
    anon.Rows = make(map[string]string)
    anon.Title = r.GetTitle()
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
        anon.Rows[label] = value
    }

	tmpl, err := template.New("table.html").ParseFS(templatesFS, "templates/table.html")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, anon)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return buf.String()
}

func (p HtmlPrinter) renderDateHeatMapChart(c report.Report, elementId int) string {
	keys := c.GetLabels()
	data := c.GetData()
	if len(data) == 0 {
		return ""
	}
	firstDate, _ := time.Parse("2006-1-2", keys[0])
	startDate := time.Date(firstDate.Year(), firstDate.Month(), 1, 0, 0, 0, 0, firstDate.Location())

	lastDate, _ := time.Parse("2006-1-2", keys[len(keys)-1])
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

	tmpl, err := template.New("date-heatmap.html").ParseFS(templatesFS, "templates/date-heatmap.html")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	var anon struct {
		Years     map[int]yearData
		FirstDate time.Time
		Range     int
        ElementId int
	}
	anon.Years = years
	anon.FirstDate = firstDate
	anon.Range = endDate.Year() - firstDate.Year() + 1
    anon.ElementId = elementId

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, anon)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return buf.String()
}

func (p HtmlPrinter) renderBartChart(c report.Report, elmentId int) string {
	tmpl, err := template.New("bar-chart.html").ParseFS(templatesFS, "templates/bar-chart.html")
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
	anon.ElementId = elmentId
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, anon)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return buf.String()
}

func (p HtmlPrinter) Print(s *os.File) {
	tmpl, err := template.New("main.html").ParseFS(templatesFS, "templates/main.html")

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	var anon struct {
		ProjectTitle    string
		RenderedReports template.HTML
        Reports []struct {
            Title string
            Id string
        }
	}

	var renderedReports bytes.Buffer
	for k := range p.reports {
		switch p.reports[k].GetReportType() {
		case "date_heatmap":
			renderedReports.WriteString(p.renderDateHeatMapChart(p.reports[k], k))
		case "bar_chart":
			renderedReports.WriteString(p.renderBartChart(p.reports[k], k))
		case "table":
			renderedReports.WriteString(p.renderTable(p.reports[k], k))
		}
		renderedReports.WriteString("\n")
        anon.Reports = append(anon.Reports, struct {
            Title string
            Id string
        }{
                Title: p.reports[k].GetTitle(), 
                Id: "elem-" + strconv.Itoa(k),
            })
	}

	anon.ProjectTitle = p.GetProjectTitle()
	anon.RenderedReports = template.HTML(renderedReports.String())
	err = tmpl.Execute(s, anon)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
