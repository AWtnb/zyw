package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/AWtnb/go-walk"
	"github.com/ktr0731/go-fuzzyfinder"
)

type WalkRoot struct {
	path    string
	exclude string
	all     bool
}

func (wr *WalkRoot) Init(path string, exclude string, all bool) {
	wr.path = path
	wr.exclude = exclude
	wr.all = all
}

func (wr WalkRoot) walk() (prompt string, found []string, err error) {
	var d walk.Dir
	d.Init(wr.path, wr.all, -1, wr.exclude)
	found, err = d.GetChildItemWithEverything()
	if err != nil || len(found) < 1 {
		prompt = ">"
		found, err = d.GetChildItem()
	} else {
		prompt = "#"
	}
	return
}

func (wr WalkRoot) SelectItem() (string, error) {
	prompt, childPaths, err := wr.walk()
	if err != nil || len(childPaths) < 1 {
		return "", nil
	}

	idx, err := fuzzyfinder.Find(childPaths, func(i int) string {
		rel, _ := filepath.Rel(wr.path, childPaths[i])
		return filepath.ToSlash(rel)
	}, fuzzyfinder.WithPromptString(prompt))
	if err != nil {
		return "", err
	}
	return childPaths[idx], nil
}

func run(root string, exclude string, all bool) int {
	var r WalkRoot
	r.Init(root, exclude, all)
	p, err := r.SelectItem()
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

func main() {
	var (
		root    string
		all     bool
		exclude string
	)
	flag.StringVar(&root, "root", "", "root of traversal")
	flag.BoolVar(&all, "all", false, "switch to search including file")
	flag.StringVar(&exclude, "exclude", "", "directory name to skip (comma-separated)")
	flag.Parse()

	os.Exit(run(root, exclude, all))
}
