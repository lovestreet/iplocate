package iplocate

import (
	"fmt"
	"testing"
)

func TestParseFile(t *testing.T) {

	var filePath = "data/ip.merge.csv"

	var mgr FileLoader
	if err := mgr.ParseFile(filePath, processLineFunction); err != nil {
		t.Errorf("%v", err)
	}
}

func processLineFunction(line string) bool {
	fmt.Println(string(line))
	return false
}

func TestLoadFile(t *testing.T) {

	var filePath = "data/ip.merge.csv"

	var mgr FileLoader
	lines, err := mgr.LoadFile(filePath)
	if err != nil {
		t.Errorf("%v", err)
	} else {
		for _, line := range lines {
			fmt.Println(string(line))
		}
	}
}
