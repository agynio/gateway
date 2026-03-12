package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func decodeJSONBody(r *http.Request, out any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(out); err != nil {
		if errors.Is(err, io.EOF) {
			return fmt.Errorf("request body is required")
		}
		return err
	}
	if decoder.Decode(&struct{}{}) != io.EOF {
		return fmt.Errorf("request body must contain a single JSON object")
	}
	return nil
}

func writeBadRequest(w http.ResponseWriter, err error) {
	message := strings.TrimSpace(err.Error())
	if message == "" {
		message = "invalid request"
	}
	writeValidationMessage(w, message)
}

func writeValidationMessage(w http.ResponseWriter, message string) {
	problem := NewProblem(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), message)
	WriteProblem(w, problem)
}

func writeJSONResponse(w http.ResponseWriter, status int, payload any, service string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("failed to encode %s response: %v", service, err)
	}
}
