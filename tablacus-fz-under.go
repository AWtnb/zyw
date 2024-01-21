package main

import (
	"flag"
	"fmt"
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
		root    string
		filer   string
		exclude string
	)
	flag.StringVar(&cur, "cur", "", "current directory")
	flag.StringVar(&root, "root", "", "root directory")
	flag.StringVar(&filer, "filer", "explorer.exe", "filer")
	flag.StringVar(&exclude, "exclude", "", "path to skip searching (comma-separated)")
	flag.Parse()
	var fl Filer
	fl.setPath(filer)
	var cd CurrentDir
	cd.setInfo(cur, root)
	os.Exit(run(fl, cd, exclude))
}

type Filer struct {
	path string
}

func (fl *Filer) setPath(path string) {
	if _, err := os.Stat(path); err == nil {
		fl.path = path
		return
	}
	fl.path = "explorer.exe"
}

func (fl Filer) open(path string) {
	if _, err := os.Stat(path); err != nil {
		exec.Command(fl.path).Start()
		return
	}
	if fi, err := os.Stat(path); err == nil && fi.IsDir() {
		exec.Command(fl.path, path).Start()
		return
	}
	exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", path).Start()
}

type CurrentDir struct {
	path       string
	searchRoot string
	depth      int
}

func (cur *CurrentDir) setInfo(curPath string, root string) {
	cur.path = curPath
	cur.searchRoot, cur.depth = cur.configSearch(root)
}

func (cur CurrentDir) configSearch(root string) (searchRoot string, depth int) {
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
	depth = 2
	return
}

func (cur CurrentDir) getChildItemsFromRoot(exclude string) (found []string, err error) {
	de := walk.DirEntry{Root: cur.searchRoot, All: false, Depth: cur.depth, Exclude: exclude}
	if strings.HasPrefix(cur.searchRoot, "C:") {
		return de.GetChildItem()
	}
	found, err = de.GetChildItemWithEverything()
	if err != nil || len(found) < 1 {
		found, err = de.GetChildItem()
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

func (cur CurrentDir) selectItem(childPaths []string) (string, error) {
	if len(childPaths) < 1 {
		return "", nil
	}
	idx, err := fuzzyfinder.Find(childPaths, func(i int) string {
		rel, _ := filepath.Rel(cur.searchRoot, childPaths[i])
		return rel
	})
	if err != nil {
		return "", err
	}
	return childPaths[idx], nil
}

func run(fl Filer, cur CurrentDir, exclude string) int {
	candidates, err := cur.getChildItemsFromRoot(exclude)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	se, err := cur.selectItem(cur.dropCurrent(candidates))
	if err != nil {
		fmt.Println(err)
		return 1
	}
	if 0 < len(se) {
		fl.open(se)
	}
	return 0
}
