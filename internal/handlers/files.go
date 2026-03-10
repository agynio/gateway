package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/agynio/gateway/internal/filesclient"
)

const maxUploadBytes = 20 << 20

type FilesUploader interface {
	Upload(ctx context.Context, filename, contentType string, sizeBytes int64, body io.Reader) (filesclient.UploadResult, error)
}

type FilesHandler struct {
	client FilesUploader
}

func NewFilesHandler(client FilesUploader) *FilesHandler {
	if client == nil {
		panic("files client is required")
	}
	return &FilesHandler{client: client}
}

func (h *FilesHandler) Upload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadBytes)
	file, header, err := r.FormFile("file")
	if err != nil {
		writeUploadFormError(w, err)
		return
	}
	defer file.Close()

	contentType := strings.TrimSpace(header.Header.Get("Content-Type"))
	result, err := h.client.Upload(r.Context(), header.Filename, contentType, header.Size, file)
	if err != nil {
		writeUploadGRPCError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("failed to encode upload response: %v", err)
	}
}

func writeUploadFormError(w http.ResponseWriter, err error) {
	statusCode := http.StatusBadRequest
	message := strings.TrimSpace(err.Error())
	if message == "" {
		message = "invalid upload request"
	}

	var maxErr *http.MaxBytesError
	if errors.As(err, &maxErr) || errors.Is(err, multipart.ErrMessageTooLarge) {
		statusCode = http.StatusRequestEntityTooLarge
	}

	if errors.Is(err, http.ErrMissingFile) {
		message = "file is required"
	}

	problem := NewProblem(statusCode, http.StatusText(statusCode), message, nil)
	WriteProblem(w, problem)
}

func writeUploadGRPCError(w http.ResponseWriter, err error) {
	grpcStatus, ok := status.FromError(err)
	if !ok {
		log.Printf("files upload error: %v", err)
		problem := NewProblem(http.StatusBadGateway, http.StatusText(http.StatusBadGateway), err.Error(), nil)
		WriteProblem(w, problem)
		return
	}

	statusCode := http.StatusBadGateway
	// UploadFile is expected to return InvalidArgument, ResourceExhausted, or
	// Unavailable; other codes (NotFound, PermissionDenied, Unauthenticated,
	// DeadlineExceeded) are treated as upstream failures.
	switch grpcStatus.Code() {
	case codes.InvalidArgument:
		statusCode = http.StatusBadRequest
	case codes.ResourceExhausted:
		statusCode = http.StatusRequestEntityTooLarge
	case codes.Unavailable:
		statusCode = http.StatusServiceUnavailable
	}

	if statusCode >= http.StatusInternalServerError {
		log.Printf("files upload error: %v", err)
	}

	message := strings.TrimSpace(grpcStatus.Message())
	if message == "" {
		message = http.StatusText(statusCode)
	}
	problem := NewProblem(statusCode, http.StatusText(statusCode), message, nil)
	WriteProblem(w, problem)
}
