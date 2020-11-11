package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gofrs/uuid"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello! you've requested %s\n", r.URL.Path)
	})

	http.HandleFunc("/cached", func(w http.ResponseWriter, r *http.Request) {
		maxAgeParams, ok := r.URL.Query()["max-age"]
		if ok && len(maxAgeParams) > 0 {
			maxAge, _ := strconv.Atoi(maxAgeParams[0])
			w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", maxAge))
		}
		requestID := uuid.Must(uuid.NewV4())
		fmt.Fprintf(w, requestID.String())
	})

	http.HandleFunc("/headers", func(w http.ResponseWriter, r *http.Request) {
		keys, ok := r.URL.Query()["key"]
		if ok && len(keys) > 0 {
			fmt.Fprintf(w, r.Header.Get(keys[0]))
			return
		}
		headers := []string{}
		for key, values := range r.Header {
			headers = append(headers, fmt.Sprintf("%s=%s", key, strings.Join(values, ",")))
		}
		fmt.Fprintf(w, strings.Join(headers, "\n"))
	})

	http.HandleFunc("/env", func(w http.ResponseWriter, r *http.Request) {
		keys, ok := r.URL.Query()["key"]
		if ok && len(keys) > 0 {
			fmt.Fprintf(w, os.Getenv(keys[0]))
			return
		}
		envs := []string{}
		for _, env := range os.Environ() {
			envs = append(envs, env)
		}
		fmt.Fprintf(w, strings.Join(envs, "\n"))
	})

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		codeParams, ok := r.URL.Query()["code"]
		if ok && len(codeParams) > 0 {
			statusCode, _ := strconv.Atoi(codeParams[0])
			if statusCode >= 200 && statusCode < 600 {
				w.WriteHeader(statusCode)
			}
		}
		requestID := uuid.Must(uuid.NewV4())
		fmt.Fprintf(w, requestID.String())
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	for _, encodedRoute := range strings.Split(os.Getenv("ROUTES"), ",") {
		if encodedRoute == "" {
			continue
		}
		pathAndBody := strings.SplitN(encodedRoute, "=", 2)
		path, body := pathAndBody[0], pathAndBody[1]
		http.HandleFunc("/"+path, func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, body)
		})
	}

	bindAddr := fmt.Sprintf(":%s", port)
	fmt.Println()
	fmt.Printf("==> Server listening at %s 🚀\n", bindAddr)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		panic(err)
	}
}
