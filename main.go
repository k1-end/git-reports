package main

import (
	"github.com/k1-end/git-reports/cmd"
)

var version string

func main() {
    cmd.Version = version
	cmd.Execute()
}
