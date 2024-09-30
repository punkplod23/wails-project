package github.com/punkplod23/wails-project/parsecsv/parsecsv

type CsvParser struct {
	Filename string
}

func NewCSVParser() *CsvParser {
	return &CsvParser{}
}

func (parser *CsvParser) RunFile(filename string) {
	parser.Filename = filename
}
