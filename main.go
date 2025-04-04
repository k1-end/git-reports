package main

import (
	"github.com/k1-end/git-reports/cmd"
)

// nodemon --exec go run main.go . --signal SIGTERM

var version string

func main() {
    cmd.Version = version
	cmd.Execute()
}
