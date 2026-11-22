package apperror

import (
	"fmt"
	"runtime"
)

type Attr struct {
	Key   string
	Value any
}

type AppError struct {
	msg    string
	err    error
	attrs  []Attr
	caller string // Tutaj zapiszemy plik i linię (np. service.go:42)
}

func (e *AppError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %v", e.msg, e.err)
	}
	return e.msg
}

func (e *AppError) Unwrap() error {
	return e.err
}

func (e *AppError) Attrs() []Attr {
	return e.attrs
}

// Caller zwraca informację o miejscu wywołania błędu
func (e *AppError) Caller() string {
	return e.caller
}

func New(msg string, args ...any) error {
	return &AppError{
		msg:    msg,
		attrs:  parseAttrs(args...),
		caller: getCallerLocation(1), // Pobieramy miejsce wywołania AppError.New
	}
}

func Wrap(err error, msg string, args ...any) error {
	if err == nil {
		return nil
	}
	return &AppError{
		msg:    msg,
		err:    err,
		attrs:  parseAttrs(args...),
		caller: getCallerLocation(1), // Pobieramy miejsce wywołania AppError.Wrap
	}
}

// Pomocnik do wyciągania lokalizacji (pomijamy wewnętrzne ramki)
func getCallerLocation(skip int) string {
	// skip + 1 bo wchodzimy do funkcji pomocniczej
	_, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return "unknown"
	}
	// Możesz zostawić pełną ścieżkę lub filepath.Base(file) jeśli wolisz samą nazwę pliku
	return fmt.Sprintf("%s:%d", file, line)
}

func parseAttrs(args ...any) []Attr {
	var attrs []Attr
	for i := 0; i < len(args); i += 2 {
		key, ok := args[i].(string)
		if !ok {
			continue
		}
		var val any
		if i+1 < len(args) {
			val = args[i+1]
		}
		attrs = append(attrs, Attr{Key: key, Value: val})
	}
	return attrs
}
