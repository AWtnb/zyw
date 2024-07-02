package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/AWtnb/go-walk"
	"github.com/ktr0731/go-fuzzyfinder"
)

func removeElem(sl []string, remove string) []string {
	ss := []string{}
	for _, s := range sl {
		if s == remove {
			continue
		}
		ss = append(ss, s)
	}
	return ss
}

type CurrentDir struct {
	path    string
	root    string
	exclude string
	all     bool
}

func (cur *CurrentDir) Init(path string, offset int, exclude string, all bool) {
	cur.path = path
	if -1 < offset {
		cur.setRootRel(offset)
	} else {
		cur.setRoot()
	}
	cur.exclude = exclude
	cur.all = all
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

func (cur CurrentDir) getChildItemsFromRoot() (assisted bool, found []string, err error) {
	var d walk.Dir
	d.Init(cur.root, cur.all, -1, cur.exclude)
	found, err = d.GetChildItemWithEverything()
	assisted = true
	if err != nil || len(found) < 1 {
		assisted = false
		found, err = d.GetChildItem()
	}
	found = removeElem(found, cur.path)
	return
}

func (cur CurrentDir) SelectItem() (string, error) {
	withEv, childPaths, err := cur.getChildItemsFromRoot()
	if err != nil || len(childPaths) < 1 {
		return "", nil
	}

	var prompt string
	if withEv {
		prompt = "#"
	} else {
		prompt = ">"
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
