package main

import (
	"flag"
	"fmt"

	"github.com/joshuabl97/cyoa/story"
	"github.com/rs/zerolog"
)

func main() {
	jsonPath := flag.String("jsonPath", "gopher.json", "the path to the choose your own adventure json")
	flag.Parse()

	var l *zerolog.Logger

	story, err := story.ParseJSON(jsonPath)
	if err != nil {
		l.Fatal().Err(err).Msg("Failure to parse JSON - exiting program")
	}

	fmt.Printf("Story: %+v", story)
}
