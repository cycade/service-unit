package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/cycade/service-unit/images"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type interperter struct {
	http.ResponseWriter
	request       *http.Request
	apiVersion    string
	status        int
	headerWritten bool
}

func (it *interperter) log() {
	info := map[string]string{}
	info["status-code"] = fmt.Sprintf("%d", it.status)
	info["last-hop"] = it.request.RemoteAddr

	xff := it.request.Header.Get("X-Forwarded-For")
	if xff != "" {
		info["origin"] = strings.Split(xff, ", ")[0]
	}

	content := fmt.Sprintf("[CUSTOM-INFO] %s", it.request.RequestURI)
	for k, v := range info {
		content += fmt.Sprintf(" | %s: %s", k, v)
	}

	log.Println(content)
}

func (it *interperter) Header() http.Header {
	return it.ResponseWriter.Header()
}

func (it *interperter) Write(raw []byte) (int, error) {
	if !it.headerWritten {
		it.WriteHeader(http.StatusOK)
	}

	return it.ResponseWriter.Write(raw)
}

func (it *interperter) WriteHeader(code int) {
	if it.headerWritten {
		return
	}

	it.ResponseWriter.Header().Add("VERSION", it.apiVersion)
	for key, value := range it.request.Header {
		it.ResponseWriter.Header().Add(key, strings.Join(value, ", "))
	}

	it.ResponseWriter.WriteHeader(code)
	it.status = code
	it.headerWritten = true

	it.log()
}

type UnitMux struct {
	version string
}

func (u *UnitMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := interperter{ResponseWriter: w, request: r, apiVersion: u.version}

	if handler.request.RequestURI == "/healthz" {
		handler.Header().Add("content-type", "application/json")
		handler.Write([]byte("{\"status\": \"UP\"}"))
	} else {
		handler.WriteHeader(http.StatusOK)
		message := fmt.Sprintf("receive request from %s", handler.request.RequestURI)
		handler.Write([]byte(message))
	}
}

func main() {
	RegisterMetricsFunc(*images.FunctionLatency)
	ver := os.Getenv("VERSION")

	mux := http.NewServeMux()
	mux.Handle("/", &UnitMux{version: ver})
	mux.HandleFunc("/images", images.Handler)
	mux.Handle("/metrics", promhttp.Handler())

	log.Println("beeper start ...")
	http.ListenAndServe(":8080", mux)
}
