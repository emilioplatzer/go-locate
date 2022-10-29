package main

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func makeRecorder[T any](toString func(T) string) (result *[]string, recorder func(T)){
	var firstArray = make([]string,0,1)
	result = &firstArray
	recorder = func(received T){
		*result = append(*result, toString(received))
	}
	return result, recorder
}

func makeRejecter[T any](t *testing.T) (rejecter func(T)){
	rejecter = func(received T){
		t.Fatal(received)
	}
	return rejecter
}

func commonLocateOpts(t *testing.T, fileToFind string, ok bool) (opts LocateOptions, obtained *[]string){
	opts.paths = []string {"./fixtures/one", "./fixtures/two"}
	opts.extList = "md,docx"
	opts.fileToFind = fileToFind
	if ok {
		var rejecter = makeRejecter[error](t)
		var recorder func(string)
		obtained, recorder = makeRecorder(func(t string) string { return t })
		opts.onFound = recorder
		opts.onError = rejecter
	} else {
		var rejecter = makeRejecter[string](t)
		var recorder func(error)
		obtained, recorder = makeRecorder(func(err error) string { return err.Error() })
		opts.onFound = rejecter
		opts.onError = recorder
	}
	return opts, obtained
}

func compareOk(t *testing.T, obtained *[]string, expected []string) {
	if !reflect.DeepEqual(*obtained, expected) {
		t.Fatalf("obtained %v expected %v", *obtained, expected)
	}
}

func TestFindOne(t *testing.T) {
	var opts, result = commonLocateOpts(t, "one", true)
	Locate(opts)
	compareOk(t, result, []string{ filepath.Join("./fixtures/one", "one.md")})
}

func TestFindTwo(t *testing.T) {
	var opts, result = commonLocateOpts(t, "both", true)
	Locate(opts)
	compareOk(t, result, []string{ filepath.Join("./fixtures/one", "both.docx"),  filepath.Join("./fixtures/two", "both.md")})
}

func TestErrorDir(t *testing.T) {
	var opts, result = commonLocateOpts(t, "one", false)
	opts.paths =  []string {"./fixtures/cuac!", "./fixtures/two"}
	Locate(opts)
	if len(*result) != 1 {
		t.Fatalf("One error expected, obtained %v ", *result)
	}
	if !strings.HasPrefix((*result)[0], "open ./fixtures/cuac!") {
		t.Fatalf("No prefix %v ", (*result)[0])
	}
}

/*
go run . cmd
go test -v
go test . -v -cover -coverprofile=coverage.out
go tool cover -html=coverage 
*/