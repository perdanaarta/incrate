package log

import (
	"io"
	"os"
)

var DefaultLogWriter *LogWriter

func NewLogWriter() *LogWriter {
	return &LogWriter{
		MultiWriter: make([]io.Writer, 0),
	}
}

type LogWriter struct {
	Writer      io.Writer
	MultiWriter []io.Writer
}

func (w *LogWriter) AddConsoleWriter() {
	w.MultiWriter = append(w.MultiWriter, os.Stdout)
	w.setWriter()
}

func (w *LogWriter) AddFileWriter(filepath string) error {
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	w.MultiWriter = append(w.MultiWriter, file)
	w.setWriter()
	return nil
}

func (w *LogWriter) SetDefault() {
	DefaultLogWriter = w
}

func (w *LogWriter) setWriter() {
	w.Writer = io.MultiWriter(w.MultiWriter...)
}
