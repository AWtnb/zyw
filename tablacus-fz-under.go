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
		depth   int
		exclude string
	)
	flag.StringVar(&cur, "cur", "", "current directory")
	flag.StringVar(&root, "root", "", "root directory")
	flag.IntVar(&depth, "depth", -1, "search depth")
	flag.StringVar(&filer, "filer", "explorer.exe", "filer")
	flag.StringVar(&exclude, "exclude", "", "path to skip searching (comma-separated)")
	flag.Parse()
	var fl Filer
	fl.setPath(filer)
	var cd CurrentDir
	cd.setInfo(cur, root, toValidFiler(filer), depth, exclude)
	os.Exit(run(fl, cd))
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
	_, err := os.Stat(path)
	if err != nil {
		exec.Command(fl.path).Start()
		return
	}
	if fi, err := os.Stat(path); err == nil && fi.IsDir() {
		exec.Command(fl.path, path).Start()
		return
	}
	exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", path).Start()
}

func toValidFiler(path string) string {
	if _, err := os.Stat(path); err == nil {
		return path
	}
	return "explorer.exe"
}

type CurrentDir struct {
	path       string
	root       string
	searchRoot string
	filer      string
	depth      int
	exclude    string
}

func (cur *CurrentDir) setInfo(curPath string, root string, filer string, depth int, exclude string) {
	cur.path = curPath
	cur.root = root
	cur.searchRoot, cur.depth = cur.configSearch()
	cur.filer = filer
	cur.exclude = exclude
}

func (cur CurrentDir) configSearch() (searchRoot string, depth int) {
	elems := strings.Split(cur.path, string(os.PathSeparator))
	for i := 0; i <= len(elems); i++ {
		ln := len(elems) - i
		p := strings.Join(elems[0:ln], string(os.PathSeparator))
		if filepath.Dir(p) == cur.root {
			searchRoot = p
			depth = -1
			return
		}
	}
	searchRoot = cur.path
	depth = 5
	return
}

func (cur CurrentDir) getChildItemsFromRoot() (found []string, err error) {
	de := walk.DirEntry{Root: cur.searchRoot, All: false, Depth: cur.depth, Exclude: cur.exclude}
	if strings.HasPrefix(cur.searchRoot, "C:") {
		return de.GetChildItem()
	}
	found, err = de.GetChildItemWithEverything()
	if err != nil || len(found) < 1 {
		found, err = de.GetChildItem()
	}
	return
}

func (cur CurrentDir) selectItem(childPaths []string) (string, error) {
	if len(childPaths) < 2 {
		return cur.searchRoot, nil
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

func run(fl Filer, cur CurrentDir) int {
	candidates, err := cur.getChildItemsFromRoot()
	if err != nil {
		fmt.Println(err)
		return 1
	}
	se, err := cur.selectItem(candidates)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	fl.open(se)
	return 0
}
