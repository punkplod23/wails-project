package github.com/punkplod23/wails-project/csvparser

import (
	"fmt"
)

type CsvParser struct {
	Filename string
}

func NewCSVParser() *CsvParser {
	return &CsvParser{}
}

func (parser *CsvParser) RunFile(filename string) {
	parser.Filename = filename
}
