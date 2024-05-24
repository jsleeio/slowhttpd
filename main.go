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
	listen                   string
	listenhttps              string
	min, max                 time.Duration
	certificatefile, keyfile string
	http01tokenpath          string
	servehttps               bool
	servehttp01token         bool
}

func configure() (c config, err error) {
	flag.StringVar(&c.http01tokenpath, "http01", "", "path to directory containing ACME HTTP-01 challenge tokens")
	flag.StringVar(&c.certificatefile, "certificate", "", "path to PEM-encoded certificate")
	flag.StringVar(&c.keyfile, "key", "", "path to PEM-encoded private key")
	flag.DurationVar(&c.min, "min", time.Millisecond*250, "minimum sleep time")
	flag.DurationVar(&c.max, "max", time.Millisecond*2500, "maximum sleep time")
	flag.StringVar(&c.listen, "listen", ":3000", "[IP address and] and port to listen on for HTTP requests")
	flag.StringVar(&c.listenhttps, "listen-https", ":443", "[IP address and] and port to listen on for HTTPS requests")
	flag.Parse()
	c.servehttps = c.certificatefile != "" && c.keyfile != ""
	c.servehttp01token = c.http01tokenpath != ""
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

func logger() *log.Logger {
	return &log.Logger{
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
}

func main() {
	c, err := configure()
	if err != nil {
		log.Fatalf("configure: %v", err)
	}
	logger := logger()
	mux := http.NewServeMux()
	mux.Handle("/health", requestLogger(logger, health()))
	mux.Handle("/randomsleep", requestLogger(logger, sleepyHandler(c.min, c.max)))
	if c.servehttp01token {
		fs := http.FileServer(http.Dir(c.http01tokenpath))
		mux.Handle("/.well-known/acme-challenge/", requestLogger(logger, http.StripPrefix("/.well-known/acme-challenge", fs)))
	}
	if c.servehttps {
		go func() {
			logger.Printf("listening on %s", c.listenhttps)
			log.Fatal(http.ListenAndServeTLS(c.listenhttps, c.certificatefile, c.keyfile, mux))
		}()
	}
	logger.Printf("listening on %s", c.listen)
	log.Fatal(http.ListenAndServe(c.listen, mux))
}
