package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"text/template"
	"time"

	"github.com/go-chi/chi"
	"github.com/joshuabl97/cyoa/story"
	"github.com/rs/zerolog"
)

func main() {
	portNum := flag.String("port_number", "8080", "The port number the server runs on")
	jsonPath := flag.String("jsonPath", "gopher.json", "the path to the choose your own adventure json")
	flag.Parse()

	// instantiate logger
	l := zerolog.New(os.Stderr).With().Timestamp().Logger()
	// make the logs look pretty
	l = l.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	// story is a map of [chapter_title string]Chapter
	book, err := story.ParseJSON(jsonPath)
	if err != nil {
		l.Fatal().Err(err).Msg("Failure to parse JSON - exiting program")
		os.Exit(0)
	}

	// generate a new 'StoryHandlers' struct for calling handler methods
	sh := story.NewStoryHandlers(book, &l)

	// registering the handlers on the serve mux (sm)
	sm := chi.NewRouter()
	sm.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("cyoa.html"))
		tmpl.Execute(w, book["intro"])
	})
	sm.Get("/{endpoint}", sh.UserChapterLoad)

	// create a custom logger that wraps the zerolog.Logger we instantiated/customized above
	errorLog := &zerologLogger{l}

	// create a new server
	s := http.Server{
		Addr:         ":" + *portNum,           // configure the bind address
		Handler:      sm,                       // set the default handler
		IdleTimeout:  120 * time.Second,        // max duration to wait for the next request when keep-alives are enabled
		ReadTimeout:  5 * time.Second,          // max duration for reading the request
		WriteTimeout: 10 * time.Second,         // max duration before returning the request
		ErrorLog:     log.New(errorLog, "", 0), // set the logger for the server
	}

	// this go function starts the server
	// when the function is done running, that means we need to shutdown the server
	// we can do this by killing the program, but if there are requests being processed
	// we want to give them time to complete
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal().Err(err)
		}
	}()

	// sending kill and interrupt signals to os.Signal channel
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	// does not invoke 'graceful shutdown' unless the signalChannel is closed
	<-sigChan

	l.Info().Msg("Received terminate, graceful shutdown")

	// this timeoutContext allows the server 30 seconds to complete all requests (if any) before shutting down
	timeoutCtx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err = s.Shutdown(timeoutCtx)
	if err != nil {
		l.Fatal().Err(err).Msg("Shutdown exceeded timeout")
		os.Exit(1)
	}
}

// custom logger type that wraps zerolog.Logger
type zerologLogger struct {
	logger zerolog.Logger
}

// implement the io.Writer interface for our custom logger.
func (l *zerologLogger) Write(p []byte) (n int, err error) {
	l.logger.Error().Msg(string(p))
	return len(p), nil
}
