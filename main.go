// package main is the main package (silencing linter)
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

type config struct {
	listen   string
	min, max time.Duration
}

func configure() (c config, err error) {
	flag.DurationVar(&c.min, "min", time.Millisecond*250, "minimum sleep time")
	flag.DurationVar(&c.max, "max", time.Millisecond*2500, "maximum sleep time")
	flag.StringVar(&c.listen, "listen", ":3000", "[IP address and] and port to listen on")
	flag.Parse()
	return
}

// very barebones middleware
func requestLogger(logger *log.Logger, h http.Handler) http.Handler {
	logFn := func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.ServeHTTP(rw, r) // serve the original request
		logger.WithFields(log.Fields{
			"uri":      r.RequestURI,
			"method":   r.Method,
			"duration": time.Since(start).Seconds(),
		}).Info("processing request")
	}
	return http.HandlerFunc(logFn)
}

func sleepyHandler(min, max time.Duration) http.Handler {
	fn := func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		time.Sleep(min + time.Duration(rand.Int63n(int64(max-min))))
		fmt.Fprintln(rw, "*snore*")
	}
	return http.HandlerFunc(fn)
}

func health() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		fmt.Fprintln(rw, "OK")
	})
}

func main() {
	logger := &log.Logger{
		Out:   os.Stdout,
		Level: log.DebugLevel,
		Formatter: &log.JSONFormatter{
			TimestampFormat: time.RFC3339,
			FieldMap: log.FieldMap{
				log.FieldKeyTime:  "timestamp",
				log.FieldKeyLevel: "severity",
				log.FieldKeyMsg:   "message",
			},
		},
	}
	c, err := configure()
	if err != nil {
		log.Fatalf("configure: %v", err)
	}
	logger.Printf("listening on %s", c.listen)
	http.Handle("/health", requestLogger(logger, health()))
	http.Handle("/randomsleep", requestLogger(logger, sleepyHandler(c.min, c.max)))
	log.Fatal(http.ListenAndServe(c.listen, nil))
}
