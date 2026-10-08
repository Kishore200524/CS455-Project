package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/Kishore200524/CS455-Project/backend/internal/feedback"
)

const maxRequestBody = 16 * 1024

func (s *Server) createFeedback(w http.ResponseWriter, r *http.Request) {
	var input feedback.CreateInput
	if err := decodeRequestJSON(w, r, &input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "request must contain one valid JSON object"})
		return
	}

	input = feedback.Normalize(input)
	if err := feedback.Validate(input); err != nil {
		message := "invalid feedback"
		switch {
		case errors.Is(err, feedback.ErrInvalidCourseID):
			message = feedback.ErrInvalidCourseID.Error()
		case errors.Is(err, feedback.ErrInvalidContent):
			message = feedback.ErrInvalidContent.Error()
		}
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": message})
		return
	}

	created, err := s.feedbackStore.Create(r.Context(), input)
	if err != nil {
		log.Printf("create feedback: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not save feedback"})
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

func decodeRequestJSON(w http.ResponseWriter, r *http.Request, target any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBody)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("request must contain exactly one JSON value")
	}
	return nil
}
