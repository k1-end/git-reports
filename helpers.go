package main

import (
	"fmt"
	"os"
)




func checkIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

func getColor(commitCount int) string {
    if commitCount == 0 {
        return "178;215;155" // #B2D79B (Light Green)
    } else if commitCount < 5 {
        return "139;195;74"  // #8BC34A (Medium Green)
    } else if commitCount < 10 {
        return "34;139;34"   // #228B22 (Forest Green)
    } else if commitCount < 15 {
        return "0;100;0"     // #006400 (Dark Green)
    } else if commitCount < 20 {
        return "0;128;128"   // #008080 (Emerald Green)
    } else {
        return "0;64;0"     // #004000 (Darker Green)
    }
}
