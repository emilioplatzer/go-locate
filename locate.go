package main

import (
	"os"
	"path/filepath"
	"strings"
)

type LocateOptions struct {
	fileToFind string
	paths []string
	extList string
	onFound func (string)
	onError func (error)
}

func Locate (opts LocateOptions) {
	var exts = make(map[string]bool)
	for _, ext := range strings.Split(opts.extList, ",") { exts["." + ext] = true }
	for _, part := range opts.paths {
		if name := strings.TrimSpace(part); name != "" {
			if dirList, err := os.ReadDir(name); err != nil {
				opts.onError(err)
			} else {
				for _, dir := range dirList {
					if !dir.IsDir() && (dir.Name() == opts.fileToFind || strings.HasPrefix(dir.Name(), opts.fileToFind + ".") && exts[filepath.Ext(dir.Name())]) {
						opts.onFound(filepath.Join(name, dir.Name()))
					}
				}
			}
		}
	}
}
