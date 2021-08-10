package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func getOutputFileName(file os.FileInfo) string {
	if file.IsDir() {
		return file.Name()
	}

	size := file.Size()
	if size == 0 {
		return fmt.Sprintf("%s (empty)", file.Name())
	} else {
		return fmt.Sprintf("%s (%db)", file.Name(), size)
	}
}

func dirWalk(output io.Writer, path string, printFiles bool, prefixFromParent string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	var listFiles []os.FileInfo
	for _, file := range files {
		if file.IsDir() {
			listFiles = append(listFiles, file)
		} else if printFiles {
			listFiles = append(listFiles, file)
		}
	}

	for idx, file := range listFiles {
		var prefix, prefixForChild string
		if idx == len(listFiles)-1 {
			prefixForChild = "\t"
			prefix = "└───"
		} else {
			prefixForChild = "│\t"
			prefix = "├───"
		}

		fmt.Fprintln(output, fmt.Sprintf("%s%s%s", prefixFromParent, prefix, getOutputFileName(file)))
		dirWalk(output, fmt.Sprintf("%s%c%s", path, os.PathSeparator, file.Name()),
			printFiles, prefixFromParent+prefixForChild)
	}

	return nil
}

func dirTree(output io.Writer, path string, printFiles bool) error {
	return dirWalk(output, path, printFiles, "")
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
