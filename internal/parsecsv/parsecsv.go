package parsecsv

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/labstack/gommon/log"
)

type FilePosition struct {
	bytePosition int64
	offset       int
}

type InvertedIndex struct {
	Term          string
	FilePositions []FilePosition
}

type CsvParser struct {
	ctx                context.Context
	FilePath           string
	FileOutName        string
	FileWriter         *os.File
	header             []string
	row                [][]string
	record             []string
	CSVReader          *csv.Reader
	WordListAjdectives map[string]struct{}
	Tokens             map[string]InvertedIndex
}

type csvRow struct {
	key   string
	value string
}

func NewCSVParser(ctx context.Context) *CsvParser {
	CsvParser := &CsvParser{}
	CsvParser.ctx = ctx
	return CsvParser
}

func (parser *CsvParser) createWordListAjdectives() map[string]struct{} {
	file, err := os.Open("assets/word-list.txt")

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	wordListAdjectives := make(map[string]struct{})

	for scanner.Scan() {
		line := scanner.Text()
		// Process the line here
		wordListAdjectives[strings.ToLower(line)] = struct{}{}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	return wordListAdjectives
}

func (parser *CsvParser) RunFile(filePath string) string {
	fmt.Println(filePath)
	parser.FilePath = filePath
	parser.WordListAjdectives = parser.createWordListAjdectives()
	//fmt.Println(parser.WordListAjdectives)
	parser.reader()
	parser.processCSV()
	//fmt.Println(readStringFromPositionWithOffset("C:\\github\\wails-project\\test.json", 2609164, 362))
	parser.searchForResults("BDCQ.SF8RS1CA 2024.03")
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
	parser.Tokens = make(map[string]InvertedIndex)
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

				for key, token := range parser.createTokens(string(parser.record[itter])) {
					tokens[key] = token
				}

				json += string(header) + ":" + string(parser.record[itter])

			}
			json += "}"

			if _, err := parser.FileWriter.Write([]byte(json)); err != nil {
				panic(err)
			}

			byteSize, err := parser.FileWriter.Seek(0, io.SeekCurrent)
			if err != nil {
				fmt.Println("Error getting file position:", err)
				return
			}

			for key, _ := range tokens {
				if _, ok := parser.Tokens[key]; ok {
					indexCopy := parser.Tokens[key]
					indexCopy.FilePositions = append(indexCopy.FilePositions, FilePosition{bytePosition: byteSize, offset: len(json)})
					parser.Tokens[key] = indexCopy
				} else {
					newIndex := InvertedIndex{
						Term: key,
						FilePositions: []FilePosition{
							{bytePosition: byteSize, offset: len(json)},
						},
					}
					parser.Tokens[key] = newIndex
				}

			}

		}

		i++
	}
	//fmt.Println(parser.Tokens)
}

func (parser *CsvParser) searchForResults(search string) {

	result := parser.createTokens(search)
	termResults := make(map[FilePosition]int)

	var kvPairs []struct {
		key   FilePosition
		value int
	}

	for term, _ := range result {
		if _, ok := parser.Tokens[term]; ok {
			for _, FP := range parser.Tokens[term].FilePositions {
				if _, ok := termResults[FP]; ok {
					termResults[FP] += 1
				} else {
					termResults[FP] = 1
				}
			}
		}
	}

	for k, v := range termResults {
		kvPairs = append(kvPairs, struct {
			key   FilePosition
			value int
		}{k, v})
	}

	sort.Slice(kvPairs, func(i, j int) bool {
		return kvPairs[i].value > kvPairs[j].value
	})

	limitedPairs := kvPairs[:10]
	for _, pair := range limitedPairs {
		go fmt.Println(readStringFromPositionWithOffset("C:\\github\\wails-project\\test.json", pair.key.bytePosition, int64(pair.key.offset)))
	}
}

func (parser *CsvParser) createTokens(record string) map[string]struct{} {

	replaceStr := strings.Replace(string(record), ",", " ", -1)
	replaceStr = strings.Replace(string(record), ".", " ", -1)
	replaceStr = strings.Replace(string(record), ":", " ", -1)
	result := strings.Split(string(replaceStr), " ")
	tokens := make(map[string]struct{})
	for _, line := range result {
		if line == "" {
			continue
		}
		if _, ok := parser.WordListAjdectives[line]; ok {
			continue
		}

		tokens[strings.ToLower(line)] = struct{}{}

	}

	return tokens

}

func readStringFromPosition(filename string, bytePosition int64) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = file.Seek(bytePosition, io.SeekStart)
	if err != nil {
		return "", err
	}

	buf := make([]byte, 100) // Adjust buffer size as needed
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return "", err
	}

	return string(buf[:n]), nil
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

func readStringFromPositionWithOffset(filename string, bytePosition int64, offset int64) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	_, err = file.Seek(bytePosition-offset, io.SeekStart)
	if err != nil {
		return "", fmt.Errorf("failed to seek to byte position: %w", err)
	}

	buf := make([]byte, offset) // Adjust buffer size as needed
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to read from file: %w", err)
	}

	return string(buf[:n]), nil
}
