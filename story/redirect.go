package story

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi"
)

func (sh *StoryHandlers) UserChapterLoad(w http.ResponseWriter, r *http.Request) {
	endpoint := chi.URLParam(r, "endpoint")
	sh.l.Info().Str("endpoint", endpoint).Msg("user chapter being processed")

	chapter, exists := sh.s[endpoint]
	if !exists {
		sh.l.Error().Msg("unable top find selected chapter")
		http.Error(w, "unable top find selected chapter", http.StatusNotFound)
		return
	}

	tmpl := template.Must(template.ParseFiles("cyoa.html"))
	tmpl.Execute(w, chapter)
}
