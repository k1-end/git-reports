package reportprinter

import (
	"bytes"
	"fmt"
	"html/template"
	"math/rand/v2"
	"os"
	"sort"
	"time"

	"github.com/k1-end/git-visualizer/src/report"
)


type HtmlPrinter struct {
    reports []report.Report
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
        }else{
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
    tmpl, err := template.New("HeatMap").Parse(`
        <script>
        var data = [
        {{range .Years}}
            {{range .Tmonths}}
                {{range .Tdays}}
                    { date: '{{.Date.Format "2006-01-02"}}', value: {{.CommitCount}} },
                {{end}}
            {{end}}
        {{end}}
        ];

        var cal = new CalHeatmap();
        cal.paint({
        data: { source: data, x: 'date', y: 'value' },
        domain: { type: 'year'},
        subDomain: { type: 'day', width: 13, height: 13},
        scale: { color: { type: 'linear', domain: [0, 20], range: ['white', 'green'], interpolate: 'hsl',}, },
        verticalOrientation: true,
        date: { 
            start: new Date('{{.FirstDate.Format "2006-01-02"}}'),
        },
        range: {{.Range}}
        });
        </script>
        `)
    if err != nil{
        fmt.Println(err)
        panic(err)
    }

    var anon struct{
        Years     map[int]Tyear
        FirstDate time.Time
        Range int
    }
    anon.Years = years
    anon.FirstDate = firstDate
    anon.Range = endDate.Year() - firstDate.Year() + 1

    var buf bytes.Buffer

    err = tmpl.Execute(&buf, anon)

    if err != nil{
        fmt.Println(err)
        panic(err)
    }
    return buf.String()
}

func (p *HtmlPrinter) RegisterReport(r report.Report) {
    p.reports = append(p.reports, r)
}

func (p HtmlPrinter) renderLineChart(c report.Report) string {
    tmpl, err := template.New("HeatMap").Parse(`
        <div style="width: 800px;"><canvas id="{{.ElementId}}"></canvas></div>
        <script>
        new Chart(
            document.getElementById("{{.ElementId}}"),
            {
              type: 'bar',
              data: {
                labels: [
                    {{range .Labels}}
                        {{.}},
                    {{end}}
                ],
                datasets: [
                  {
                    label: '{{.Title}}',
                    data: [
                        {{range .Data}}
                            {{.}},
                        {{end}}
                    ],
                  }
                ]
              },
        options: {
          scales: {
            x: {
              ticks: {
                display: true,
                autoSkip: false
              }
            }
          }
        }
            }
          );
        </script>
        `)
    if err != nil{
        fmt.Println(err)
        panic(err)
    }
    var anon struct{
        Title string
        Labels []string
        Data []int
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

    if err != nil{
        fmt.Println(err)
        panic(err)
    }
    return buf.String()
}

func (p HtmlPrinter) Print() {
    tmpl, err := template.New("header").Parse(`
<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="utf-8">
        <meta name="viewport" content="width=device-width, initial-scale=1">
        <meta name="description" content="">
        <meta name="author" content="Mark Otto, Jacob Thornton, and Bootstrap contributors">
        <meta name="generator" content="Hugo 0.84.0">
        <title>Git reports for {{.ProjectTitle}}</title>

        <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.0.2/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-EVSTQN3/azprG1Anm3QDgpJLIm9Nao0Yz1ztcQTwFspd3yD65VohhpuuCOmLASjC" crossorigin="anonymous">

        </head>
    <body>

        <header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
            <a class="navbar-brand col-md-12 me-0 px-3 text-center" href="#">{{.ProjectTitle}}</a>
        </header>

        <div class="container-fluid p-5">
            <div class="row">
                <main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
                    <script src="https://d3js.org/d3.v7.min.js"></script>
                    <script src="https://unpkg.com/cal-heatmap/dist/cal-heatmap.min.js"></script>
                    <link rel="stylesheet" href="https://unpkg.com/cal-heatmap/dist/cal-heatmap.css">
                    <script src="https://cdn.jsdelivr.net/npm/chart.js@4.4.8/dist/chart.umd.min.js"></script>
                    <div style="width: 800px;" id="cal-heatmap"></div>
                    {{.RenderedReports}}
                </main>
            </div>
        </div>
        <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.0.2/dist/js/bootstrap.bundle.min.js" integrity="sha384-MrcW6ZMFYlzcLA8Nl+NtUVF0sA7MsXsP1UyJoMp4YLEuNSfAP+JcXn/tWtIaxVXM" crossorigin="anonymous"></script>
    </body>
</html>
        `)

    if err != nil{
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

    var anon struct{
        ProjectTitle string
        RenderedReports template.HTML
    }
    anon.ProjectTitle = p.GetProjectTitle()
    anon.RenderedReports = template.HTML(renderedReports.String())
    err = tmpl.Execute(os.Stdout, anon)

    if err != nil{
        fmt.Println(err)
        panic(err)
    }
}

func (p *HtmlPrinter) SetProjectTitle(s string)  {
    p.projectTitle = s
}

func (p *HtmlPrinter) GetProjectTitle()  string{
    return p.projectTitle
}
