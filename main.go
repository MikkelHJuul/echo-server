package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Echo server listening on port %s.\n", port)

	err := http.ListenAndServe(
		":"+port,
		http.HandlerFunc(handler),
	)
	if err != nil {
		panic(err)
	}
}

func handler(wr http.ResponseWriter, req *http.Request) {
	fmt.Println(req.RemoteAddr + " | " + req.Method + " " + req.URL.String())
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	if len(bodyBytes) != 0 {
		fmtFunct(req.URL.Path)(bodyBytes)
		// Replace original body with buffered version so it's still sent to the
		// browser.
		req.Body.Close()
		req.Body = ioutil.NopCloser(
			bytes.NewReader(bodyBytes),
		)
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
