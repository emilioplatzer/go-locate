package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"github.com/jessevdk/go-flags"
)

const ExtListDefault = "exe,com,bat,sh"

type Opts struct {
	Verbose bool `short:"v" long:"verbose" description:"Show verbose debug information. Default:false"`
	PathList string `short:"p" long:"path" description:"The list of path to search in. Default: PATH env var"` 
	ExtList string `short:"e" long:"extensions" description:"The list of extensions to search in. Default: exe,com,bat,sh"` 
}

var opts Opts

func locate (opts Opts, fileToFind string) bool {
	var found = false
	var paths = strings.Split(opts.PathList, ";")
	var exts = make(map[string]bool)
	for _, ext := range strings.Split(opts.ExtList, ",") { exts["." + ext] = true }
	for _, part := range paths {
		if name := strings.TrimSpace(part); name != "" {
			if dirList, err := os.ReadDir(name); err != nil {
				fmt.Println(err)
			} else {
				for _, dir := range dirList {
					if !dir.IsDir() && (dir.Name() == fileToFind || strings.HasPrefix(dir.Name(), fileToFind + ".")) {
						if opts.Verbose && !found {
							fmt.Println("\nFound:")
						}
						fmt.Println(filepath.Join(name, dir.Name()))
						found = true
					}
				}
			}
		}
	}
	return found
}

func main(){
	var args = os.Args[1:]
	args, err := flags.ParseArgs(&opts, args)
	if err != nil || len(args) == 0 {
		fmt.Println("")
		fmt.Println("Search executable in PATH")
		fmt.Println("")
		fmt.Println("Usage:")
		fmt.Println("    GO-LOCATE exe_to_find [--path path_list] [--extensions extension_list] ")
		fmt.Println("")
		fmt.Println("Examples:")
		fmt.Println("    GO-LOCATE bash")
		fmt.Println("    GO-LOCATE readme --path docs\\manual;docs\\FAQS --extensions md,txt,doc")
	} else {
		var fileToFind = args[0]
		if opts.PathList == "" {
			opts.PathList = os.Getenv("PATH")
		}
		if opts.Verbose {
			fmt.Println("Looking in ",opts.PathList)
		}
		if opts.ExtList == "" {
			opts.ExtList = ExtListDefault
			if opts.Verbose {
				fmt.Print("For executable ")
			}
		} else {
			if opts.Verbose {
				fmt.Print("For file")
			}
		}
		if opts.Verbose {
			fmt.Println(fileToFind, "with extensions", opts.ExtList)
		}
		var found = locate(opts, fileToFind)
		if !found {
			fmt.Println(fileToFind, "not found")
		}
	}
}