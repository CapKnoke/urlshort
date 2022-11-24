package main

import (
	"flag"
	"fmt"

	// "github.com/CapKnoke/urlshort"
	"net/http"
	"os"

	"github.com/CapKnoke/urlshort"
)

var (
	yaml string
	json string
)

func init() {
	flag.StringVar(&yaml, "yaml", "redirects.yaml", "redirects file in yaml format")
	flag.StringVar(&json, "json", "redirects.json", "redirects file in json format")
	flag.Parse()
}

var pathsToUrls = map[string]string{
	"/json-godoc": "https://pkg.go.dev/encoding/json",
	"/yaml-godoc": "https://godoc.org/gopkg.in/yaml.v2",
}

func main() {
	mapHandler := createMapHandler()
	yamlHandler := createYAMLHandler(mapHandler)
	jsonHandler := createJSONHandler(yamlHandler)
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(writer http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(writer, "Hello, world!")
}

func createMapHandler() http.HandlerFunc {
	mux := defaultMux()
	return urlshort.MapHandler(pathsToUrls, mux)
}

func createYAMLHandler(fallback http.HandlerFunc) http.HandlerFunc {
	yamlFile, err := os.Open(yaml)
	if err != nil {
		panic(err)
	}
	yamlHandler, err := urlshort.YAMLHandler(yamlFile, fallback)
	if err != nil {
		panic(err)
	}
	return yamlHandler
}

func createJSONHandler(fallback http.HandlerFunc) http.HandlerFunc {
	jsonFile, err := os.Open(json)
	if err != nil {
		panic(err)
	}
	jsonHandler, err := urlshort.JSONHanlder(jsonFile, fallback)
	if err != nil {
		panic(err)
	}
	return jsonHandler
}
