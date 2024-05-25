// Demo code for the TextView primitive.
package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	lib = flag.String("lib", "btea", "library to use; possible values: btea, tview")
	hp  = flag.Bool("hp", false, "whether to use high performance renderer")
)

func main() {
	flag.Parse()

	if *hp && *lib != "btea" {
		fmt.Printf("-hp can only be used with bubbletea\n")
		os.Exit(1)
	}

	switch *lib {
	case "tview":
		tviewTextView()
	case "btea":
		bteaVPort(*hp)
	}
}
