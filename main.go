package main

import (
	"fmt"
	"os"
	"strings"
	"github.com/jessevdk/go-flags"
)

var Opts struct {
	Verbose bool `short:"v" long:"verbose" description:"Show verbose debug information. Default:false"`
	PathList string `short:"p" long:"path" description:"The list of path to search in. Default: PATH env var"` 
	ExtList string `short:"e" long:"extensions" description:"The list of extensions to search in. Default: exe,com,bat,sh"` 
	PathSeparator string `long:"separator" description:"The path separator. Default: depends on S.O."` 
}

const ExtListDefault = "exe,com,bat,sh"

func main(){
	var args = os.Args[1:]
	args, err := flags.ParseArgs(&Opts, args)
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
		if Opts.PathSeparator == "" {
			Opts.PathSeparator = string(os.PathListSeparator)
		}
		if Opts.PathList == "" {
			Opts.PathList = os.Getenv("PATH")
		}
		if Opts.Verbose {
			fmt.Println("Looking in ",Opts.PathList)
		}
		if Opts.ExtList == "" {
			Opts.ExtList = ExtListDefault
			if Opts.Verbose {
				fmt.Print("For executable ")
			}
		} else {
			if Opts.Verbose {
				fmt.Print("For file")
			}
		}
		if Opts.Verbose {
			fmt.Println(fileToFind, "with extensions", Opts.ExtList)
		}
		var found = false
		var paths = strings.Split(Opts.PathList, Opts.PathSeparator)
		Locate(
			fileToFind,
			paths,
			Opts.ExtList,
			func(path string){
				if Opts.Verbose && !found {
					fmt.Println("\nFound:")
				}
				fmt.Println(path)
				found = true
			},
			func(err error){
				fmt.Println(err)
			},
		)
		if !found {
			fmt.Println(fileToFind, "not found")
		}
	}
}