package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/AWtnb/go-walk"
	"github.com/ktr0731/go-fuzzyfinder"
)

type CurrentDir struct {
	path string
	root string
}

func (cur *CurrentDir) Init(path string, offset int) {
	cur.path = path
	if -1 < offset {
		cur.setRootRel(offset)
	} else {
		cur.setRoot()
	}
}

func (cur CurrentDir) Path() string {
	return cur.path
}

func (cur CurrentDir) pathElems() []string {
	return strings.Split(cur.path, string(os.PathSeparator))
}

func (cur *CurrentDir) setRootRel(offset int) {
	elems := cur.pathElems()
	if 0 < offset && offset < len(elems) {
		ss := elems[0 : len(elems)-offset]
		cur.root = strings.Join(ss, string(os.PathSeparator))
		return
	}
	cur.root = cur.path
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
	cur.root = cur.path
}

func (cur CurrentDir) getChildItemsFromRoot(exclude string, all bool) (assisted bool, found []string, err error) {
	var d walk.Dir
	d.Init(cur.root, all, -1, exclude)
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
		rel, _ := filepath.Rel(cur.root, childPaths[i])
		return filepath.ToSlash(rel)
	}, fuzzyfinder.WithPromptString(prompt))
	if err != nil {
		return "", err
	}
	return childPaths[idx], nil
}
