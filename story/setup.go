package story

import "github.com/rs/zerolog"

type StoryHandlers struct {
	s Story
	l *zerolog.Logger
}

func NewStoryHandlers(s Story, l *zerolog.Logger) *StoryHandlers {
	return &StoryHandlers{s, l}
}
