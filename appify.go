package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("lack of input file")
		return
	}
	fixDylib(os.Args[1])
}
func fixDylib(fname string) {
	curr, _ := os.Getwd()
	fmt.Println("curr == ", curr)
	out, e := exec.Command("otool", "-L", fname).Output()
	if e != nil {
		fmt.Println(e)
		return
	}
	libs := strings.Split(string(out), "\n")
	fmt.Println(len(libs), "libs of ", fname)
	for _, str := range libs {
		lib := cutStr(str)
		if !strings.Contains(lib, "local") {
			continue
		}
		fmt.Println(lib)
		if !fileExists(getFileName(lib)) {
			e = copyFile(lib, curr+"/"+getFileName(lib))
			if e != nil {
				fmt.Println("while cp", lib, e)
				return
			}
		}
		e = exec.Command("install_name_tool", "-change", lib, `@executable_path/`+getFileName(lib), fname).Run()
		if e != nil {
			fmt.Println("install:", e)
			return
		}
	}
}
func cutStr(str string) string {
	f := strings.Split(str, " ")[0]
	for i := 0; i < len(f); i++ {
		if f[i:i+1] == "/" {
			return f[i:]
		}
	}
	return f
}
func getFileName(p string) string {
	for i := len(p) - 1; i > -1; i-- {
		if p[i:i+1] == "/" {
			return p[i+1:]
		}
	}
	return p
}
func fileExists(p string) bool {
	_, e := os.Stat(p)
	if os.IsNotExist(e) {
		return false
	}
	return true
}
func copyFile(from, to string) error {
	fi, e := os.Open(from)
	if e != nil {
		return e
	}
	defer fi.Close()
	fo, e := os.OpenFile(to, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if e != nil {
		return e
	}
	defer fo.Close()
	_, e = io.Copy(fo, fi)
	if e != nil {
		return e
	}
	return nil
}
