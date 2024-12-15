package mylog

import (
	"bytes"
	"fmt"
	"io"
)

// Colors per a la sortida del log
const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
)

// ColorWriter és un Writer personalitzat que afegeix colors a les sortides.
type ColorWriter struct {
	writer io.Writer
}

// NewColorWriter crea una nova instància de ColorWriter.
func NewColorWriter(w io.Writer) *ColorWriter {
	return &ColorWriter{writer: w}
}

// Write implementa l'interface io.Writer per ColorWriter.
func (cw *ColorWriter) Write(p []byte) (n int, err error) {
	// Aquí pots definir el color que vols per a cada nivell de log
	// Per exemple, podríem fer que els missatges d'error siguin vermells
	var color string
	if bytes.Contains(p, []byte(`"level":"error"`)) {
		color = Red
	} else if bytes.Contains(p, []byte(`"level":"info"`)) {
		color = Green
	} else if bytes.Contains(p, []byte(`"level":"debug"`)) {
		color = Yellow
	} else {
		color = Reset
	}
	// Escriure el missatge amb el color escollit
	return fmt.Fprintf(cw.writer, "%s%s%s", color, p, Reset)
}
