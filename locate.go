package main

import (
	"os"
	"path/filepath"
	"strings"
)

func Locate (fileToFind string, paths []string, extList string, onFound func (string), onError func (error)) {
	var exts = make(map[string]bool)
	for _, ext := range strings.Split(extList, ",") { exts["." + ext] = true }
	for _, part := range paths {
		if name := strings.TrimSpace(part); name != "" {
			if dirList, err := os.ReadDir(name); err != nil {
				onError(err)
			} else {
				for _, dir := range dirList {
					if !dir.IsDir() && (dir.Name() == fileToFind || strings.HasPrefix(dir.Name(), fileToFind + ".") && exts[filepath.Ext(dir.Name())]) {
						onFound(filepath.Join(name, dir.Name()))
					}
				}
			}
		}
	}
}
