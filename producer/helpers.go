package main

import (
	"encoding/json"
	"net/http"
	"os"
)

func lookupEnvOrDefault(key, defaultValue string) string {
	val, isPresent := os.LookupEnv(key)
	if !isPresent {
		return defaultValue
	}
	return val
}

func writeJSON(w http.ResponseWriter, status int, v any, headers http.Header) error {
	js, err := json.Marshal(v)
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for k, v := range headers {
		w.Header()[k] = v
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	w.Write(js)

	return nil
}
