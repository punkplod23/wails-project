package parsecsv

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/labstack/gommon/log"
)

// FilePosition struct
type FilePosition struct {
	BytePosition int64 `json:"bytePosition"`
	Offset       int   `json:"offset"`
}

// Index struct
type Index struct {
	ctx                context.Context           `json:"-"` // Exclude from JSON output
	WordListAjdectives map[string]struct{}       `json:"wordListAdjectives"`
	Tokens             map[string][]FilePosition `json:"tokens"`
}

func NewInvertedIndex(ctx context.Context) *Index {
	WordList := &Index{}
	WordList.ctx = ctx

	WordList.WordListAjdectives = createWordListAjdectives()
	return WordList
}

func createWordListAjdectives() map[string]struct{} {
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

func sanitizeString(record string) string {
	re := regexp.MustCompile("[^a-zA-Z0-9. ]+")
	record = re.ReplaceAllString(record, " ")
	return record
}

func (index *Index) createTokens(records []string) map[string]struct{} {

	tokens := make(map[string]struct{})
	for _, record := range records {
		record = sanitizeString(record)
		result := strings.Split(string(record), " ")

		for _, line := range result {
			if line == "" {
				continue
			}
			if _, ok := index.WordListAjdectives[line]; ok {
				continue
			}

			tokens[strings.ToLower(line)] = struct{}{}

		}
	}

	return tokens
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

func (index *Index) SaveIndexToFile() {

	jsonData, err := json.Marshal(index.Tokens)
	if err != nil {
		// Handle error
		log.Fatal(err)
	}
	// Write the JSON data to a file
	err = os.WriteFile("index.json", jsonData, 0644)
	if err != nil {
		// Handle error
		log.Fatal(err)
	}
}

func (index *Index) createIndexTokens(tokens map[string]struct{}, BytePosition int64, Offset int) {
	if len(index.Tokens) < 1 {
		index.Tokens = make(map[string][]FilePosition)
	}

	for key, _ := range tokens {
		index.Tokens[key] = append(index.Tokens[key], FilePosition{BytePosition: BytePosition, Offset: Offset})
	}

}

func (index *Index) reloadIndex() {

	// Read the JSON data from a file
	jsonData, err := os.ReadFile("index.json")
	if err != nil {
		// Handle error
		log.Fatal(err)
	}

	err = json.Unmarshal(jsonData, &index.Tokens)
	if err != nil {
		// Handle error
		log.Fatal(err)
	}
}

func (index *Index) SearchForResults(search string) string {

	index.reloadIndex()
	result := index.createTokens(strings.Split(string(search), "%"))
	termResults := make(map[FilePosition]int)

	var kvPairs []struct {
		key   FilePosition
		value int
	}

	for term, _ := range result {
		if _, ok := index.Tokens[term]; ok {
			for _, FP := range index.Tokens[term] {
				if _, ok := termResults[FP]; ok {
					termResults[FP] += 1
				} else {
					termResults[FP] = 1
				}
			}
		}
	}

	if len(termResults) < 1 {

		return "No results found"

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
	var builder strings.Builder
	for _, pair := range limitedPairs {
		result, err := readStringFromPositionWithOffset("C:\\github\\wails-project\\test.json", pair.key.BytePosition, int64(pair.key.Offset))
		if err != nil {
			// Handle error
			log.Fatal(err)
		}
		builder.WriteString(result)

	}

	return builder.String()
}
