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
		cur        string
		depth      int
		filer      string
		exclude    string
		fromParent bool
	)
	flag.StringVar(&cur, "cur", "", "current directory")
	flag.IntVar(&depth, "depth", 1, "search depth")
	flag.StringVar(&filer, "filer", "explorer.exe", "filer")
	flag.StringVar(&exclude, "exclude", "", "path to skip searching (comma-separated)")
	flag.BoolVar(&fromParent, "from-parent", false, "search from parent of current directory")
	flag.Parse()
	os.Exit(run(cur, depth, filer, exclude, fromParent))
}

func run(cur string, depth int, filer string, exclude string, fromParent bool) int {
	if fromParent {
		exclude = exclude + "," + filepath.Base(cur)
		cur = filepath.Dir(cur)
		if !isValidPath(cur) {
			return 0
		}
	}
	cs, err := walk.GetChildItems(cur, depth, false, toSlice(exclude, ","))
	if err != nil {
		return 1
	}
	cs = removeElem(cs, cur)
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
	return openDir(filer, cs[idx])
}

func isValidPath(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func openDir(filer string, path string) int {
	if fi, err := os.Stat(path); err == nil && fi.IsDir() {
		exec.Command(filer, path).Start()
		return 0
	}
	return 1
}

func toSlice(s string, sep string) []string {
	var ss []string
	for _, elem := range strings.Split(s, sep) {
		ss = append(ss, strings.TrimSpace(elem))
	}
	return ss
}

func removeElem(elems []string, target string) []string {
	var ss []string
	for _, s := range elems {
		if s != target {
			ss = append(ss, strings.TrimSpace(s))
		}
	}
	return ss
}
