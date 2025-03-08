package report

type Data struct {
    IntValue    int
    StringValue string
    IsInt       bool
}

type Report struct {
    data []Data
    labels []string
    title string
    reportType string
}

func (r *Report) SetTitle(t string) {
    r.title = t
}

func (r *Report) SetData(d []Data) {
    r.data = d
}

func (r *Report) SetLabels(l []string) {
    r.labels = l
}

func (r *Report) SetReportType(l string) {
    r.reportType = l
}

func (r Report) GetTitle() string {
    return r.title
}

func (r Report) GetData() []Data {
    return r.data
}

func (r Report) GetLabels() []string {
    return r.labels
}

func (r Report) GetReportType() string {
    return r.reportType
}
