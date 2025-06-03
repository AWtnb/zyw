package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/AWtnb/go-walk"
	fzf "github.com/junegunn/fzf/src"
)

type Root struct {
	path    string
	exclude string
	all     bool
}

func (r *Root) Init(path string, exclude string, all bool) {
	r.path = path
	r.exclude = exclude
	r.all = all
}

func (r Root) walk() (prompt string, found []string, err error) {
	var w walk.Walker
	w.Init(r.path, r.all, -1, r.exclude)
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
	var r Root
	r.Init(root, exclude, all)
	prompt, found, err := r.walk()
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

	outputChan := make(chan string)
	go func() {
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
