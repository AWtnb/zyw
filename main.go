package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AWtnb/go-walk"
	"github.com/ktr0731/go-fuzzyfinder"
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
	f.SetPath(filer)
	cd := CurrentDir{Path: cur}
	if -1 < offset {
		cd.setRootRel(offset)
	} else {
		cd.setRoot()
	}
	os.Exit(run(f, cd, exclude, all))
}

type CurrentDir struct {
	Path string
	root string
}

func (cur CurrentDir) pathElems() []string {
	return strings.Split(cur.Path, string(os.PathSeparator))
}

func (cur *CurrentDir) setRootRel(offset int) {
	elems := cur.pathElems()
	if 0 < offset && offset < len(elems) {
		ss := elems[0 : len(elems)-offset]
		cur.root = strings.Join(ss, string(os.PathSeparator))
		return
	}
	cur.root = cur.Path
}

func (cur *CurrentDir) setRoot() {
	elems := cur.pathElems()
	for i := 0; i < len(elems); i++ {
		ln := len(elems) - i
		p := strings.Join(elems[0:ln], string(os.PathSeparator))
		if _, err := os.Stat(filepath.Join(p, ".root")); err == nil {
			cur.root = p
			return
		}

	}
	cur.root = cur.Path
}

func (cur CurrentDir) getChildItemsFromRoot(exclude string, all bool) (assisted bool, found []string, err error) {
	d := walk.Dir{All: all, Root: cur.root}
	d.SetWalkDepth(-1)
	d.SetWalkException(exclude)
	found, err = d.GetChildItemWithEverything()
	assisted = true
	if err != nil || len(found) < 1 {
		assisted = false
		found, err = d.GetChildItem()
	}
	return
}

func (cur CurrentDir) dropCurrent(childPaths []string) (paths []string) {
	for _, p := range childPaths {
		if p == cur.Path {
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
		rel, _ := filepath.Rel(cur.root, childPaths[i])
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
	fl.OpenSmart(se, cur.Path)
	return 0
}
