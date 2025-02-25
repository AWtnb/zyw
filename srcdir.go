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

type SrcDir struct {
	path    string
	root    string
	exclude string
	all     bool
}

func (sd *SrcDir) Init(path string, offset int, exclude string, all bool) {
	sd.path = path
	if -1 < offset {
		sd.setRootRel(offset)
	} else {
		sd.setRoot()
	}
	sd.exclude = exclude
	sd.all = all
}

func (sd SrcDir) pathElems() []string {
	return strings.Split(sd.path, string(os.PathSeparator))
}

func (sd *SrcDir) setRootRel(offset int) {
	elems := sd.pathElems()
	if 0 < offset && offset < len(elems) {
		ss := elems[0 : len(elems)-offset]
		sd.root = strings.Join(ss, string(os.PathSeparator))
		return
	}
	sd.root = sd.path
}

func (sd *SrcDir) setRoot() {
	elems := sd.pathElems()
	for i := 0; i < len(elems); i++ {
		ln := len(elems) - i
		p := strings.Join(elems[0:ln], string(os.PathSeparator))
		if _, err := os.Stat(filepath.Join(p, ".root")); err == nil {
			sd.root = p
			return
		}
	}
	sd.root = sd.path
}

func (sd SrcDir) getChildItemsFromRoot() (assisted bool, found []string, err error) {
	var d walk.Dir
	d.Init(sd.root, sd.all, -1, sd.exclude)
	found, err = d.GetChildItemWithEverything()
	assisted = true
	if err != nil || len(found) < 2 {
		assisted = false
		found, err = d.GetChildItem()
	}
	found = removeElem(found, sd.path)
	return
}

func (sd SrcDir) SelectItem() (string, error) {
	withEv, childPaths, err := sd.getChildItemsFromRoot()
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
		rel, _ := filepath.Rel(sd.root, childPaths[i])
		return filepath.ToSlash(rel)
	}, fuzzyfinder.WithPromptString(prompt))
	if err != nil {
		return "", err
	}
	return childPaths[idx], nil
}
