package reportprinter

import (
	"fmt"
	"html/template"
	"os"
	"sort"
	"time"

	"github.com/k1-end/git-visualizer/src/report"
)


type HtmlPrinter struct {
    reports []report.Report
}

func (p HtmlPrinter) PrintDateHeatMapChart(c report.Report) {
    keys := c.GetLabels()
    data := c.GetData()
    if len(data) == 0 {
        fmt.Println("No commits where found!")
        return
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
        const data = [
        {{range .Years}}
            {{range .Tmonths}}
                {{range .Tdays}}
                    { date: '{{.Date.Format "2006-01-02"}}', value: {{.CommitCount}} },
                {{end}}
            {{end}}
        {{end}}
        ];

        const cal = new CalHeatmap();
        cal.paint({
        data: { source: data, x: 'date', y: 'value' },
        domain: { type: 'year'},
        subDomain: { type: 'day'},
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
    err = tmpl.Execute(os.Stdout, anon)

    if err != nil{
        fmt.Println(err)
        panic(err)
    }
}

func (p *HtmlPrinter) RegisterReport(r report.Report) {
    p.reports = append(p.reports, r)
}

func (p HtmlPrinter) Print() {
    fmt.Println(`
        <!DOCTYPE html>
        <html lang="en">
            <head>
                <meta charset="UTF-8">
                <meta name="viewport" content="width=device-width, initial-scale=1">
                <title></title>
                <link href="css/style.css" rel="stylesheet">
                <script src="https://d3js.org/d3.v7.min.js"></script>
                <script src="https://unpkg.com/cal-heatmap/dist/cal-heatmap.min.js"></script>
                <link rel="stylesheet" href="https://unpkg.com/cal-heatmap/dist/cal-heatmap.css">
            </head>
            <body>
            <div id="cal-heatmap"></div>
        `)
    for k := range p.reports {
        switch p.reports[k].GetReportType() {
        case "date_heatmap":
            p.PrintDateHeatMapChart(p.reports[k])
        }
    }

    fmt.Println(`
            </body>
        </html>
        `)
}
