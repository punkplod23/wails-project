package parsecsv

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/labstack/gommon/log"
)

type CsvParser struct {
	ctx           context.Context
	FilePath      string
	FileOutName   string
	FileWriter    *os.File
	header        []string
	row           [][]string
	record        []string
	CSVReader     *csv.Reader
	InvertedIndex *Index
}

type csvRow struct {
	key   string
	value string
}

func NewCSVParser(ctx context.Context) *CsvParser {
	CsvParser := &CsvParser{}
	CsvParser.ctx = ctx
	index := &Index{}
	index.ctx = ctx
	CsvParser.InvertedIndex = index
	return CsvParser
}

func (parser *CsvParser) Query(query string) string {
	result := parser.InvertedIndex.SearchForResults(query)
	fmt.Println(result)
	return result
}

func (parser *CsvParser) RunFile(filePath string) string {
	fmt.Println(filePath)
	parser.FilePath = filePath
	parser.reader()
	parser.processCSV()
	//ForResults("nzsioc")
	return parser.complete()
}

func (parser *CsvParser) reader() {
	parser.FileOutName = "test.json"
	fout, err := os.OpenFile(parser.FileOutName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		panic(err)
	}

	parser.FileWriter = fout
	if _, err = parser.FileWriter.Write([]byte("[")); err != nil {
		panic(err)
	}

	f, err := os.Open(parser.FilePath)
	if err != nil {
		log.Fatal(parser.ctx, "Unable to read input file "+parser.FilePath, err)
	}
	csvReader := csv.NewReader(f)
	parser.CSVReader = csvReader
}

func (parser *CsvParser) processCSV() {

	i := 0
	delimeter := ""
	json := ""

	for {
		record, err := parser.CSVReader.Read()
		if err != nil {
			// If we reached the end of the file, break out of the loop
			if err == io.EOF {
				break
			}
			log.Fatalf("err while reading CSV file: %s", err)
		}
		if i == 0 {
			parser.header = record

		} else {
			parser.record = record
		}
		if i > 0 {
			if i > 1 {
				delimeter = ","
			}
			json = delimeter + "{"
			var tokens = make(map[string]struct{})
			for itter, header := range parser.header {
				if itter > 0 {
					json += ","
				}

				json += string(header) + ":" + string(parser.record[itter])

			}
			json += "}"

			for key, token := range parser.InvertedIndex.createTokens(parser.record) {
				if _, ok := tokens[key]; ok {
					continue
				}
				if key == "" {
					continue
				}
				tokens[key] = token
			}

			if _, err := parser.FileWriter.Write([]byte(json)); err != nil {
				panic(err)
			}

			byteSize, err := parser.FileWriter.Seek(0, io.SeekCurrent)
			if err != nil {
				fmt.Println("Error getting file position:", err)
				return
			}
			parser.InvertedIndex.createIndexTokens(tokens, byteSize, len(json))
		}
		i++
	}
	parser.InvertedIndex.SaveIndexToFile()
}

func (parser *CsvParser) complete() string {
	if _, err := parser.FileWriter.Write([]byte("]")); err != nil {
		panic(err)
	}
	if err := parser.FileWriter.Close(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(parser.FileWriter)
	return "Parsing complete"

}
