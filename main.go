package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ktr0731/go-fuzzyfinder"
)

func main() {
	var (
		cur     string
		offset  int
		filer   string
		all     bool
		exclude string
		stdout  bool
	)
	flag.StringVar(&cur, "cur", "", "current directory")
	flag.IntVar(&offset, "offset", -1, "Specify the directory to start file traversing, by the number of layers from the current directory.\n`0` for the current directory, `1` for the parent directory, `2` for its parent directory, and so on.\nIf this value is negative, the path is traversed back to the directory containing the file `.root`. If no `.root` file is found, the current directory is used as the root of the search.")
	flag.StringVar(&filer, "filer", "explorer.exe", "filer")
	flag.BoolVar(&all, "all", false, "switch to search including file")
	flag.StringVar(&exclude, "exclude", "", "path to skip searching (comma-separated)")
	flag.BoolVar(&stdout, "stdout", false, "switch to stdout")
	flag.Parse()

	var f Filer
	f.Init(filer)

	// cur = strings.TrimPrefix(cur, "\"")
	// cur = strings.TrimSuffix(cur, "\"")
	var d CurrentDir
	d.Init(cur, offset, exclude, all)

	os.Exit(run(f, d, stdout))
}

func run(fl Filer, cur CurrentDir, stdout bool) int {
	p, err := cur.SelectItem()
	if err != nil {
		if err != fuzzyfinder.ErrAbort {
			fmt.Println(err.Error())
		}
		return 1
	}
	if stdout {
		fmt.Print(p)
		return 0
	}
	if err := fl.OpenSmart(p, cur.Path()); err != nil {
		fmt.Println(err)
		return 1
	}
	return 0
}
