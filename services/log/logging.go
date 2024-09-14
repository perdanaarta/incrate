package log

import (
	"log"
)

var Logger *log.Logger

func NewLogger(log_writer *LogWriter) *log.Logger {

	writer := log_writer.Writer

	Logger := log.New(writer, "[INCRATE]", log.Ldate|log.Ltime|log.Lshortfile)

	return Logger
}
