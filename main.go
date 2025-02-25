package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ktr0731/go-fuzzyfinder"
)

func main() {
	var (
		src     string
		offset  int
		all     bool
		exclude string
	)
	flag.StringVar(&src, "src", "", "source directory")
	flag.IntVar(&offset, "offset", -1, "Specify the directory to start file traversal, by the number of layers from the current directory.\n`0` for the current directory, `1` for the parent directory, `2` for its parent directory, and so on.\nIf this value is negative, the path is traversed back to the directory containing the file `.root`. If no `.root` file is found, the current directory is used as the root of the search.")
	flag.BoolVar(&all, "all", false, "switch to search including file")
	flag.StringVar(&exclude, "exclude", "", "path to skip searching (comma-separated)")
	flag.Parse()

	os.Exit(run(src, offset, exclude, all))
}

func run(src string, offset int, exclude string, all bool) int {
	var sd SrcDir
	sd.Init(src, offset, exclude, all)
	p, err := sd.SelectItem()
	if err != nil {
		if err != fuzzyfinder.ErrAbort {
			fmt.Println(err.Error())
		}
		return 1
	}
	if 0 < len(p) {
		fmt.Print(p)
	}
	return 0
}
