package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Echo server listening on port %s.\n", port)

	err := http.ListenAndServe(
		":"+port,
		h2c.NewHandler(
			http.HandlerFunc(handler),
			&http2.Server{},
		),
	)
	if err != nil {
		panic(err)
	}
}

func handler(wr http.ResponseWriter, req *http.Request) {
	if req.Body != nil {
		bodyBytes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			panic(err)
		}

		fmtFunct(req.URL.Path)(bodyBytes)
		// Replace original body with buffered version so it's still sent to the
		// browser.
		req.Body.Close()
		req.Body = ioutil.NopCloser(
			bytes.NewReader(bodyBytes),
		)
	} else {
		fmt.Printf("%s | %s %s\n", req.RemoteAddr, req.Method, req.URL)
	}
	serveHTTP(wr, req)
}

func fmtFunct(path string) func(b interface{}) {
	if path == "/err" {
		return func(b interface{}) { fmt.Errorf("%s\n", b) }
	}
	return func(b interface{}) { fmt.Printf("%s\n", b) }
}

func serveHTTP(wr http.ResponseWriter, req *http.Request) {
	wr.Header().Add("Content-Type", "text/plain")
	wr.WriteHeader(200)

	host, err := os.Hostname()
	if err == nil {
		fmt.Fprintf(wr, "Request served by %s\n\n", host)
	} else {
		fmt.Fprintf(wr, "Server hostname unknown: %s\n\n", err.Error())
	}

	fmt.Fprintf(wr, "%s %s %s\n", req.Proto, req.Method, req.URL)
	fmt.Fprintln(wr, "")
	fmt.Fprintf(wr, "Host: %s\n", req.Host)
	for key, values := range req.Header {
		for _, value := range values {
			fmt.Fprintf(wr, "%s: %s\n", key, value)
		}
	}

	fmt.Fprintln(wr, "")
	io.Copy(wr, req.Body)
}
