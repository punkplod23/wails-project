package parsecsv_test

import (
	"context"
	"testing"

	"github.com/punkplod23/wails-project/internal/parsecsv"
)

func TestFileName(t *testing.T) {
	//Will sort at a later date
	name := "file"
	parser := parsecsv.NewCSVParser(context.Background())
	parser.RunFile(name)
	//msg, err := Hello("Gladys")

}
