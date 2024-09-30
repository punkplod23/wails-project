package parsecsv

import (
	"testing"
)

func TestFileName(t *testing.T) {

	name := "file"
	parser := parsecsv.NewCSVParser()
	parser.RunFile(name)
	//msg, err := Hello("Gladys")

}
