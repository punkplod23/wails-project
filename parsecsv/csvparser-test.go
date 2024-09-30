package github.com/punkplod23/wails-project/parsecsv/parsecsvtest

import (
	"testing"

	"github.com/punkplod23/wails-project/parsecsv/parsecsv"
)

func TestFileName(t *testing.T) {

	name := "file"
	parser := parsecsv.NewCSVParser()
	parser.RunFile(name)
	//msg, err := Hello("Gladys")

}
