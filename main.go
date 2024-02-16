package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AWtnb/tablacus-fz-under/walk"
	"github.com/ktr0731/go-fuzzyfinder"
)

func main() {
	var (
		cur     string
		root    string
		filer   string
		exclude string
		all     bool
	)
	flag.StringVar(&cur, "cur", "", "current directory")
	flag.StringVar(&root, "root", "", "root directory")
	flag.StringVar(&filer, "filer", "explorer.exe", "filer")
	flag.StringVar(&exclude, "exclude", "", "path to skip searching (comma-separated)")
	flag.BoolVar(&all, "all", false, "switch to search including file")
	flag.Parse()
	var f Filer
	f.SetPath(filer)
	var cd CurrentDir
	cd.setInfo(cur, root)
	os.Exit(run(f, cd, exclude, all))
}

type CurrentDir struct {
	path       string
	searchRoot string
	depth      int
}

func (cur *CurrentDir) setInfo(curPath string, root string) {
	cur.path = curPath
	cur.searchRoot, cur.depth = cur.configure(root)
}

func (cur CurrentDir) configure(root string) (searchRoot string, depth int) {
	if root == "." {
		searchRoot = cur.path
		depth = -1
		return
	}
	if root == ".." {
		searchRoot = filepath.Dir(cur.path)
		depth = -1
		return
	}
	elems := strings.Split(cur.path, string(os.PathSeparator))
	for i := 0; i <= len(elems); i++ {
		ln := len(elems) - i
		p := strings.Join(elems[0:ln], string(os.PathSeparator))
		if filepath.Dir(p) == root {
			searchRoot = p
			depth = -1
			return
		}
	}
	searchRoot = cur.path
	depth = -1
	return
}

func (cur CurrentDir) getChildItemsFromRoot(exclude string, all bool) (withEverything bool, found []string, err error) {
	d := walk.Dir{All: all, Root: cur.searchRoot}
	d.SetWalkDepth(cur.depth)
	d.SetWalkException(exclude)
	if strings.HasPrefix(cur.searchRoot, "C:") && (2 < walk.GetDepth(cur.path)) {
		return d.GetChildItem()
	}
	withEverything, found, err = d.GetChildItemWithEverything()
	if err != nil || len(found) < 1 {
		withEverything, found, err = d.GetChildItem()
	}
	return
}

func (cur CurrentDir) dropCurrent(childPaths []string) (paths []string) {
	for _, p := range childPaths {
		if p == cur.path {
			continue
		}
		paths = append(paths, p)
	}
	return
}

func (cur CurrentDir) selectItem(childPaths []string, prompt string) (string, error) {
	if len(childPaths) < 1 {
		return "", nil
	}
	idx, err := fuzzyfinder.Find(childPaths, func(i int) string {
		rel, _ := filepath.Rel(cur.searchRoot, childPaths[i])
		return rel
	}, fuzzyfinder.WithPromptString(prompt))
	if err != nil {
		return "", err
	}
	return childPaths[idx], nil
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
	fl.OpenSmart(se)
	return 0
}
