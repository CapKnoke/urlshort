package urlshort

import (
	"io"
	"net/http"

	"encoding/json"

	yaml "gopkg.in/yaml.v2"
)

type pathMap struct {
	Path string `yaml:"path" json:"path"`
	URL  string `yaml:"url" json:"url"`
}

type decoder interface {
	Decode(v interface{}) error
}

func decode(decoder decoder) ([]pathMap, error) {
	var paths []pathMap
	for {
		err := decoder.Decode(&paths)
		if err == io.EOF {
			return paths, nil
		} else if err != nil {
			return nil, err
		}
	}
}

func buildMap(pathURLs []pathMap) map[string]string {
	paths := make(map[string]string)
	for _, y := range pathURLs {
		paths[y.Path] = y.URL
	}
	return paths
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		url := pathsToUrls[req.URL.Path]
		if url != "" {
			http.Redirect(writer, req, url, http.StatusPermanentRedirect)
		} else {
			fallback.ServeHTTP(writer, req)
		}
	})
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
func YAMLHandler(yml io.Reader, fallback http.Handler) (http.HandlerFunc, error) {
	decoder := yaml.NewDecoder(yml)
	parsedYaml, err := decode(decoder)
	if err != nil {
		return nil, err
	}
	paths := buildMap(parsedYaml)
	return MapHandler(paths, fallback), nil
}

func JSONHanlder(j io.Reader, fallback http.Handler) (http.HandlerFunc, error) {
	decoder := json.NewDecoder(j)
	parsedJSON, err := decode(decoder)
	if err != nil {
		return nil, err
	}
	paths := buildMap(parsedJSON)
	return MapHandler(paths, fallback), nil
}
