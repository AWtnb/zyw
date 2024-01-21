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
	os.Exit(run(cur, root, filer, depth, exclude))
}

type CurrentDir struct {
	root       string
	path       string
	searchRoot string
	filer      string
	depth      int
	exclude    string
}

func (cur *CurrentDir) setInfo(curPath string, root string, filer string, depth int, exclude string) {
	cur.path = curPath
	cur.setRoot(root)
	cur.setSearchRoot()
	cur.setFiler(filer)
	cur.depth = depth
	cur.exclude = exclude
}

func (cur *CurrentDir) setRoot(path string) {
	if path == "..." {
		cur.root = filepath.Dir(filepath.Dir(cur.path))
		return
	}
	cur.root = path
}

func (cur *CurrentDir) setSearchRoot() {
	elems := strings.Split(cur.path, string(os.PathSeparator))
	for i := 0; i <= len(elems); i++ {
		ln := len(elems) - i
		p := strings.Join(elems[0:ln], string(os.PathSeparator))
		if filepath.Dir(p) == cur.root {
			cur.searchRoot = p
			return
		}
	}
	cur.searchRoot = cur.path
}

func (cur *CurrentDir) setFiler(path string) {
	if _, err := os.Stat(path); err == nil {
		cur.filer = path
		return
	}
	cur.filer = "explorer.exe"
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

func (cur CurrentDir) run(path string) {
	_, err := os.Stat(path)
	if err != nil {
		exec.Command(cur.filer).Start()
		return
	}
	if fi, err := os.Stat(path); err == nil && fi.IsDir() {
		exec.Command(cur.filer, path).Start()
		return
	}
	exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", path).Start()
}

func run(curPath string, root string, filer string, depth int, exclude string) int {

	var cur CurrentDir
	cur.setInfo(curPath, root, filer, depth, exclude)
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
	cur.run(se)
	return 0
}
