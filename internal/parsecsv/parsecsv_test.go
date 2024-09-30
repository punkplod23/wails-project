package parsecsv_test

import (
	"testing"

	"github.com/punkplod23/wails-project/internal/parsecsv"
)

func TestFileName(t *testing.T) {

	name := "file"
	parser := parsecsv.NewCSVParser()
	parser.RunFile(name)
	//msg, err := Hello("Gladys")

}
