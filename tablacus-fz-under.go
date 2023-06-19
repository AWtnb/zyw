package main

import (
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/AWtnb/tablacus-fz-under/walk"
	"github.com/ktr0731/go-fuzzyfinder"
)

func main() {
	var (
		cur     string
		depth   int
		filer   string
		exclude string
	)
	flag.StringVar(&cur, "cur", "", "current directory")
	flag.IntVar(&depth, "depth", 1, "search depth")
	flag.StringVar(&filer, "filer", "explorer.exe", "filer")
	flag.StringVar(&exclude, "exclude", "", "path to skip searching (comma-separated)")
	flag.Parse()
	os.Exit(run(cur, depth, filer, exclude))
}

func run(cur string, depth int, filer string, exclude string) int {
	cs, err := walk.GetChildItems(cur, depth, false, toSlice(exclude, ","))
	if err != nil {
		return 1
	}
	cs = trimElem(cs, cur)
	if len(cs) < 1 {
		return 0
	}
	idx, err := fuzzyfinder.Find(cs, func(i int) string {
		rel, _ := filepath.Rel(cur, cs[i])
		return rel
	})
	if err != nil {
		return 1
	}
	src := cs[idx]
	if fi, err := os.Stat(src); err == nil && fi.IsDir() {
		exec.Command(filer, src).Start()
	} else {
		return 1
	}
	return 0
}

func toSlice(s string, sep string) []string {
	var ss []string
	for _, elem := range strings.Split(s, sep) {
		ss = append(ss, strings.TrimSpace(elem))
	}
	return ss
}

func trimElem(elems []string, search string) []string {
	var ss []string
	for _, s := range elems {
		if s != search {
			ss = append(ss, strings.TrimSpace(s))
		}
	}
	return ss
}
