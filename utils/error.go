package utils

import (
	"log"
	"runtime"
	"strings"
	"time"

	"github.com/Stream-I-T-Consulting/stream-http-service-go/config"
	"github.com/getsentry/sentry-go"
)

const (
	// HTTP status codes were copied from https://en.wikipedia.org/wiki/List_of_HTTP_status_codes
	StatusInvalidToken = 498
)

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func HandleErrors(errs ...error) {
	if len(errs) > 0 {
		for _, err := range errs {
			_, fn, line, _ := runtime.Caller(1)
			fns := strings.Split(fn, "/")
			log.Println(
				"[",
				time.Now().Format("2006-01-02 15:04:05"),
				"]",
				"ServerError:",
				fns[len(fns)-1],
				line,
				"|",
				err.Error(),
			)

			if config.AppConfig.Env == "production" || config.AppConfig.Env == "staging" {
				sentry.CaptureException(err)
			}
		}
	}
}
