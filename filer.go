package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/AWtnb/tablacus-fz-under/sys"
)

type Filer struct {
	path string
}

func (flr *Filer) SetPath(path string) {
	if _, err := os.Stat(path); err == nil {
		flr.path = path
		return
	}
	flr.path = "explorer.exe"
}

func (flr Filer) Open(path string) {
	exec.Command(flr.path, path).Start()
}

func (flr Filer) OpenSmart(path string) {
	if fi, err := os.Stat(path); err == nil && fi.IsDir() {
		flr.Open(path)
		return
	}
	fmt.Printf("'%s' is a file.\nopen itself? (y/N): ", path)
	sc := bufio.NewScanner(os.Stdin)
	sc.Scan()
	s := sc.Text()
	if strings.EqualFold(s, "y") {
		fmt.Println("[ACCEPTED] opening in default app.")
		sys.Open(path)
		return
	}
	fmt.Println("[DENIED] opening its directory on filer.")
	d := filepath.Dir(path)
	flr.Open(d)
}
