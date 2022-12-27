package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/yargevad/filepathx" // improved Glob function
)

var (
	wg         sync.WaitGroup
	addonsPath string      = "./Interface/AddOns"
	jobsPerms  chan string = make(chan string, 256)
	jobsTocs   chan string = make(chan string, 256)
)

func main() {

	cpuc := runtime.NumCPU()
	runtime.GOMAXPROCS(cpuc)

	// Check startup from the game folder
	if _, err := os.Stat(addonsPath); os.IsNotExist(err) {
		fmt.Println("[!] > This utility must be run from the game folder!")
		fmt.Println("[!] > Make sure this file is next to Wow.exe")
		fmt.Scanln()
		os.Exit(3)
	}
	fmt.Println("[?] > Checking AddOns folder..")
	fixReadOnly(addonsPath)

	dirs, _ := filepathx.Glob(addonsPath + "/*")
	for _, dir := range dirs {
		jobsPerms <- dir
		if !checkForTOCFile(dir) {
			jobsTocs <- dir
		}
	}

	close(jobsPerms)
	close(jobsTocs)

	for i := 0; i <= min(cpuc-1, len(jobsPerms)); i++ {
		wg.Add(1)
		go fixReadOnlyGlob(jobsPerms, &wg)
	}

	wg.Wait()

	for i := 0; i <= min(cpuc-1, len(jobsTocs)); i++ {
		wg.Add(1)
		go processAddOnFolder(jobsTocs, &wg)
	}

	wg.Wait()

	fmt.Println("[+] > All AddOns seems to be in good state!")
	fmt.Scanln()
}

func min(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func fixReadOnly(dir string) {
	err := os.Chmod(dir, 0666)
	if err != nil {
		fmt.Println("[!] > Can't fix permissions for \"" + dir + "\"")
	}
}

func fixReadOnlyGlob(jobs <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for dir := range jobs {
		fixReadOnly(dir)
		files, _ := filepathx.Glob(dir + "/**/*")
		for _, obj := range files {
			err := os.Chmod(obj, 0666)
			if err != nil {
				fmt.Println("[!] > Can't fix permissions for \"" + obj + "\"")
			}
		}
	}
}

func checkForTOCFile(dir string) bool {
	files, _ := filepathx.Glob(dir + "/*.toc")
	for _, file := range files {
		return strings.TrimSuffix(filepath.Base(file), filepath.Ext(file)) == filepath.Base(dir)
	}
	return false
}

func processAddOnFolder(jobs <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for dir := range jobs {
		err := os.Rename(dir, dir+"_forfix") // prepare directory name
		if err != nil {
			fmt.Println("[x] > Can't fix directory \"" + dir + "\", skipped..")
			continue
		}
		dir = dir + "_forfix" // new prepared directory name
		files, _ := filepathx.Glob(dir + "/**/*.toc")
		tocs := make(map[string]bool)
		errn := 0
	OUTER:
		for _, file := range files {
			for _, ancestor := range strings.Split(file, "\\") {
				if tocs[ancestor] {
					continue OUTER // it's a library, need to skip
				}
			}
			cpdir := filepath.Dir(file)
			cndir := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
			tocs[filepath.Base(cpdir)] = true
			fmt.Println("[!] > AddOn \"" + strings.Replace(cpdir, "_forfix", "", -1) + "\" is broken, trying to fix it..")
			err := os.Rename(cpdir, "./Interface/AddOns/"+cndir)
			if err != nil {
				fmt.Println("[x] > Can't fix AddOn \"" + strings.Replace(cpdir, "_forfix", "", -1) + "\", skipped..")
				errn += 1
				continue
			}
			fmt.Println("[+] > AddOn \"" + cndir + "\" has been successfully fixed!")
		}
		if errn == 0 {
			os.RemoveAll(dir) // only if no errors in entire broken directory
		}
	}
}
