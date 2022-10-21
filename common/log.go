package common

import "log"

type Logger interface {
	Error(args ...any)
	Errorf(template string, args ...any)
}

type DefaultLogger struct{}

func (l DefaultLogger) Error(args ...any) {
	log.Println(args...)
}

func (l DefaultLogger) Errorf(template string, args ...any) {
	log.Printf(template, args...)
}
