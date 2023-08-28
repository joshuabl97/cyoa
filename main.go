package main

import (
	"flag"
	"net/http"
	"os"
	"text/template"

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
		os.Exit(0)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("cyoa.html"))
		tmpl.Execute(w, story["intro"])
	})

	http.ListenAndServe(":8080", nil)
}
