package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var (
		cur     string
		offset  int
		filer   string
		exclude string
		all     bool
	)
	flag.StringVar(&cur, "cur", "", "current directory")
	flag.IntVar(&offset, "offset", -1, "Specify the directory to start file traversing, by the number of layers from the current directory.\n`0` for the current directory, `1` for the parent directory, `2` for its parent directory, and so on.\nIf this value is negative, the path is traversed back to the directory containing the file `.root`. If no `.root` file is found, the current directory is used as the root of the search.")
	flag.StringVar(&filer, "filer", "explorer.exe", "filer")
	flag.StringVar(&exclude, "exclude", "", "path to skip searching (comma-separated)")
	flag.BoolVar(&all, "all", false, "switch to search including file")
	flag.Parse()

	var f Filer
	f.Init(filer)

	var d CurrentDir
	d.Init(cur, offset)

	os.Exit(run(f, d, exclude, all))
}

func run(fl Filer, cur CurrentDir, exclude string, all bool) int {
	withEv, candidates, err := cur.getChildItemsFromRoot(exclude, all)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	var prompt string
	if withEv {
		prompt = "#"
	} else {
		prompt = ">"
	}
	se, err := cur.selectItem(cur.dropCurrent(candidates), prompt)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	if len(se) < 1 {
		return 0
	}
	if err := fl.OpenSmart(se, cur.Path()); err != nil {
		fmt.Println(err)
		return 1
	}
	return 0
}
