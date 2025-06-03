package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/AWtnb/go-walk"
	fzf "github.com/junegunn/fzf/src"
)

func find(path string, exclude string, all bool) (prompt string, found []string, err error) {
	var w walk.Walker
	w.Init(path, all, -1, exclude)
	found, err = w.EverythingTraverse()
	if err != nil || len(found) < 2 {
		prompt = ">"
		found, err = w.Traverse()
	} else {
		prompt = "#"
	}
	return
}

// https://gist.github.com/junegunn/193990b65be48a38aac6ac49d5669170
func run(root string, exclude string, all bool) int {
	prompt, found, err := find(root, exclude, all)
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}

	inputChan := make(chan string)
	go func() {
		for _, p := range found {
			rel, _ := filepath.Rel(root, p)
			inputChan <- rel
		}
		close(inputChan)
	}()

	var wg sync.WaitGroup
	outputChan := make(chan string)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for s := range outputChan {
			fmt.Print(filepath.Join(root, s))
		}
	}()

	options, err := fzf.ParseOptions(
		true,
		[]string{fmt.Sprintf("--prompt=%s", prompt)},
	)
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}

	options.Input = inputChan
	options.Output = outputChan

	code, err := fzf.Run(options)
	close(outputChan)
	wg.Wait()
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}

	return code
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
