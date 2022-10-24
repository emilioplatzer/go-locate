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

func FixtureFind(t *testing.T, fileToFind string, expected []string) {
	var pathList =  []string {"./fixtures/one", "./fixtures/two"}
	var extList = "md,docx";
	var rejecter = makeRejecter[error](t)
	var obtained, recorder = makeRecorder(func(t string) string { return t })
	Locate(fileToFind, pathList, extList, recorder, rejecter)
	if !reflect.DeepEqual(*obtained, expected) {
		t.Fatalf("obtained %v expected %v", *obtained, expected)
	}
}

func TestFindOne(t *testing.T) {
	FixtureFind(t, "one", []string{ filepath.Join("./fixtures/one", "one.md")})
}

func TestFindTwo(t *testing.T) {
	FixtureFind(t, "both", []string{ filepath.Join("./fixtures/one", "both.docx"),  filepath.Join("./fixtures/two", "both.md")})
}

func FixtureError(t *testing.T, fileToFind string) []string {
	var pathList =  []string {"./fixtures/cuac!", "./fixtures/two"}
	var extList = "md,docx";
	var rejecter = makeRejecter[string](t)
	var obtained, recorder = makeRecorder(func(err error) string { return err.Error() })
	Locate(fileToFind, pathList, extList, rejecter, recorder)
	return *obtained;
}

func TestErrorDir(t *testing.T) {
	var obtained = FixtureError(t, "one")
	if len(obtained) != 1 {
		t.Fatalf("One erro expected, obtained %v ", obtained)
	}
	if !strings.HasPrefix(obtained[0], "open ./fixtures/cuac!") {
		t.Fatalf("No prefix %v ", obtained[0])
	}
}

/*
go test . -v -cover -coverprofile=coverage.out
go tool cover -html=coverage 
*/