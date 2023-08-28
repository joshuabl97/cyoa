package story

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Story map[string]Chapter

type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

func ParseJSON(jsonPath *string) (Story, error) {
	var story Story
	// open the file
	file, err := os.Open(*jsonPath)
	if err != nil {
		return story, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// read JSON data from the file
	var jsonData []byte
	jsonData, err = io.ReadAll(file)
	if err != nil {
		return story, fmt.Errorf("error reading file: %v", err)
	}

	// unmarshal the json into a type Story map[ChapterTitle](type Chapter)
	err = json.Unmarshal(jsonData, &story)
	if err != nil {
		return story, fmt.Errorf("error unmarshalling json: %v", err)
	}

	return story, nil
}
