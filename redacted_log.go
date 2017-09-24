package main

import (
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/kovetskiy/lorg"
)

type redactedLog struct {
	*lorg.Log

	writer *redactedWriter
}

func NewRedactedLog() *redactedLog {
	log := &redactedLog{
		Log:    lorg.NewLog(),
		writer: &redactedWriter{},
	}

	log.SetOutput(log.writer)

	return log
}

func (log *redactedLog) Hide(value string) {
	log.writer.values = append(log.writer.values, value)
}

func (log *redactedLog) HideFromConfig(config Config) {
	log.Hide(config.Secret)
	log.Hide(config.UserID)
	log.Hide(config.AccountID)
}

func (log *redactedLog) GetWriter() io.Writer {
	return log.writer
}

type redactedWriter struct {
	values []string
}

func (writer redactedWriter) Write(buffer []byte) (int, error) {
	output := string(buffer)

	for _, value := range writer.values {
		output = regexp.MustCompile(
			fmt.Sprintf(
				"(%s)%s",
				regexp.QuoteMeta(value[:3]),
				regexp.QuoteMeta(value[3:]),
			),
		).ReplaceAllString(output, `$1***`)
	}

	return os.Stderr.Write([]byte(output))
}
