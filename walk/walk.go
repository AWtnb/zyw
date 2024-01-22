package walk

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/AWtnb/tablacus-fz-under/everything"
)

type WalkException struct {
	names []string
}

func (wex *WalkException) setNames(s string, sep string) {
	if len(s) < 1 {
		return
	}
	for _, elem := range strings.Split(s, sep) {
		wex.names = append(wex.names, strings.TrimSpace(elem))
	}
}

func (wex WalkException) contains(name string) bool {
	for _, n := range wex.names {
		if n == name {
			return true
		}
	}
	return false
}

func (wex WalkException) isSkippablePath(path string) bool {
	sep := string(os.PathSeparator)
	if strings.Contains(path, sep+".") {
		return true
	}
	for _, n := range wex.names {
		if strings.Contains(path, sep+n+sep) || strings.HasSuffix(path, n) {
			return true
		}
	}
	return false
}

func (wex WalkException) filter(paths []string) []string {
	if len(wex.names) < 1 {
		return paths
	}
	sl := []string{}
	for i := 0; i < len(paths); i++ {
		p := paths[i]
		if wex.isSkippablePath(p) {
			continue
		}
		sl = append(sl, p)
	}
	return sl
}

type ChildItems struct {
	rootDepth int
	maxDepth  int
	sep       string
}

func (ci *ChildItems) setRoot(path string) {
	ci.rootDepth = ci.getDepth(path)
}

func (ci ChildItems) getDepth(path string) int {
	return strings.Count(strings.TrimSuffix(path, ci.sep), ci.sep)
}

func (ci ChildItems) isSkippableDepth(path string) bool {
	return 0 < ci.maxDepth && ci.maxDepth < ci.getDepth(path)-ci.rootDepth
}

func (ci ChildItems) filterByDepth(paths []string) (filteredPaths []string) {
	if ci.maxDepth < 0 {
		filteredPaths = paths
		return
	}
	for i := 0; i < len(paths); i++ {
		p := paths[i]
		if ci.isSkippableDepth(p) {
			continue
		}
		filteredPaths = append(filteredPaths, p)
	}
	return
}

type DirWalker struct {
	All        bool
	Root       string
	childItems ChildItems
	exeception WalkException
}

func (dw *DirWalker) ChildItemsHandler(depth int) {
	ci := ChildItems{maxDepth: depth, sep: string(os.PathSeparator)}
	ci.setRoot(dw.Root)
	dw.childItems = ci
}

func (dw *DirWalker) ExceptionHandler(exclude string) {
	var wex WalkException
	wex.setNames(exclude, ",")
	dw.exeception = wex
}

func (dw DirWalker) GetChildItemWithEverything() (found []string, err error) {
	if dw.childItems.maxDepth == 0 {
		return
	}
	found, err = everything.Scan(dw.Root, !dw.All)
	if err != nil {
		return
	}
	if 0 < len(found) {
		found = dw.childItems.filterByDepth(dw.exeception.filter(found))
	}
	return
}

func (dw DirWalker) GetChildItem() (found []string, err error) {
	if dw.childItems.maxDepth == 0 {
		return
	}
	err = filepath.WalkDir(dw.Root, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if dw.childItems.isSkippableDepth(path) {
			return filepath.SkipDir
		}
		if dw.exeception.contains(info.Name()) {
			return filepath.SkipDir
		}
		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			found = append(found, path)
		} else {
			if dw.All {
				found = append(found, path)
			}
		}
		return nil
	})
	return
}
