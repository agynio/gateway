package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/agynio/gateway/internal/filesclient"
	"github.com/agynio/gateway/internal/gen"
)

type stubFilesUploader struct {
	t                *testing.T
	expectCall       bool
	expectedFilename string
	expectedType     string
	expectedSize     *int64
	expectBody       bool
	expectedBody     []byte
	result           filesclient.UploadResult
	err              error
	called           bool
}

func (s *stubFilesUploader) Upload(ctx context.Context, filename, contentType string, sizeBytes int64, body io.Reader) (filesclient.UploadResult, error) {
	s.called = true
	if !s.expectCall {
		s.t.Fatalf("unexpected upload call")
	}
	if ctx == nil {
		s.t.Fatalf("expected context")
	}
	if s.expectedFilename != "" && filename != s.expectedFilename {
		s.t.Fatalf("unexpected filename: got %q want %q", filename, s.expectedFilename)
	}
	if s.expectedType != "" && contentType != s.expectedType {
		s.t.Fatalf("unexpected content type: got %q want %q", contentType, s.expectedType)
	}
	if s.expectedSize != nil && sizeBytes != *s.expectedSize {
		s.t.Fatalf("unexpected size: got %d want %d", sizeBytes, *s.expectedSize)
	}
	if s.expectBody {
		data, err := io.ReadAll(body)
		if err != nil {
			s.t.Fatalf("read upload body: %v", err)
		}
		if !bytes.Equal(data, s.expectedBody) {
			s.t.Fatalf("unexpected body: got %q want %q", string(data), string(s.expectedBody))
		}
	} else {
		_, _ = io.Copy(io.Discard, body)
	}
	return s.result, s.err
}

func TestFilesHandlerUploadSuccess(t *testing.T) {
	data := []byte("hello")
	size := int64(len(data))
	createdAt := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)

	uploader := &stubFilesUploader{
		t:                t,
		expectCall:       true,
		expectedFilename: "hello.txt",
		expectedType:     "text/plain",
		expectedSize:     &size,
		expectBody:       true,
		expectedBody:     data,
		result: filesclient.UploadResult{
			ID:          "file-123",
			Filename:    "hello.txt",
			ContentType: "text/plain",
			SizeBytes:   size,
			CreatedAt:   createdAt,
		},
	}

	request := newMultipartRequest(t, "hello.txt", "text/plain", data)
	response := httptest.NewRecorder()

	NewFilesHandler(uploader).Upload(response, request)

	result := response.Result()
	if result.StatusCode != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, result.StatusCode)
	}
	if got := result.Header.Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected content type application/json, got %q", got)
	}
	if !uploader.called {
		t.Fatalf("expected uploader to be called")
	}

	var payload filesclient.UploadResult
	if err := json.NewDecoder(result.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if payload != uploader.result {
		t.Fatalf("unexpected payload: got %+v want %+v", payload, uploader.result)
	}
}

func TestFilesHandlerUploadMissingFile(t *testing.T) {
	uploader := &stubFilesUploader{t: t, expectCall: false}
	request := newMultipartRequest(t, "", "", nil)
	response := httptest.NewRecorder()

	NewFilesHandler(uploader).Upload(response, request)

	assertProblemResponse(t, response, http.StatusBadRequest, "file is required")
	if uploader.called {
		t.Fatalf("expected uploader not to be called")
	}
}

func TestFilesHandlerUploadTooLarge(t *testing.T) {
	data := bytes.Repeat([]byte("a"), int(maxUploadBytes)+1)
	uploader := &stubFilesUploader{t: t, expectCall: false}
	request := newMultipartRequest(t, "big.bin", "application/octet-stream", data)
	response := httptest.NewRecorder()

	NewFilesHandler(uploader).Upload(response, request)

	assertProblemResponse(t, response, http.StatusRequestEntityTooLarge, "")
	if uploader.called {
		t.Fatalf("expected uploader not to be called")
	}
}

func TestFilesHandlerUploadInvalidArgument(t *testing.T) {
	uploader := &stubFilesUploader{
		t:          t,
		expectCall: true,
		err:        status.Error(codes.InvalidArgument, "invalid upload"),
	}
	request := newMultipartRequest(t, "bad.txt", "text/plain", []byte("bad"))
	response := httptest.NewRecorder()

	NewFilesHandler(uploader).Upload(response, request)

	assertProblemResponse(t, response, http.StatusBadRequest, "invalid upload")
}

func TestFilesHandlerUploadResourceExhausted(t *testing.T) {
	uploader := &stubFilesUploader{
		t:          t,
		expectCall: true,
		err:        status.Error(codes.ResourceExhausted, "quota exceeded"),
	}
	request := newMultipartRequest(t, "big.txt", "text/plain", []byte("data"))
	response := httptest.NewRecorder()

	NewFilesHandler(uploader).Upload(response, request)

	assertProblemResponse(t, response, http.StatusRequestEntityTooLarge, "quota exceeded")
}

func TestFilesHandlerUploadUnavailable(t *testing.T) {
	uploader := &stubFilesUploader{
		t:          t,
		expectCall: true,
		err:        status.Error(codes.Unavailable, "service down"),
	}
	request := newMultipartRequest(t, "file.txt", "text/plain", []byte("data"))
	response := httptest.NewRecorder()

	NewFilesHandler(uploader).Upload(response, request)

	assertProblemResponse(t, response, http.StatusServiceUnavailable, "service down")
}

func TestFilesHandlerUploadNonGRPCError(t *testing.T) {
	uploader := &stubFilesUploader{
		t:          t,
		expectCall: true,
		err:        errors.New("kaboom"),
	}
	request := newMultipartRequest(t, "file.txt", "text/plain", []byte("data"))
	response := httptest.NewRecorder()

	NewFilesHandler(uploader).Upload(response, request)

	assertProblemResponse(t, response, http.StatusBadGateway, "kaboom")
}

func newMultipartRequest(t *testing.T, filename, contentType string, data []byte) *http.Request {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	if filename != "" {
		header := textproto.MIMEHeader{}
		header.Set("Content-Disposition", `form-data; name="file"; filename="`+filename+`"`)
		if contentType != "" {
			header.Set("Content-Type", contentType)
		}
		part, err := writer.CreatePart(header)
		if err != nil {
			t.Fatalf("create multipart part: %v", err)
		}
		if len(data) > 0 {
			if _, err := part.Write(data); err != nil {
				t.Fatalf("write multipart data: %v", err)
			}
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close multipart writer: %v", err)
	}

	request := httptest.NewRequest(http.MethodPost, "/files/v1/files", body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	return request
}

func assertProblemResponse(t *testing.T, response *httptest.ResponseRecorder, statusCode int, detail string) {
	t.Helper()
	result := response.Result()
	if result.StatusCode != statusCode {
		t.Fatalf("expected status %d, got %d", statusCode, result.StatusCode)
	}
	if got := result.Header.Get("Content-Type"); got != problemContentType {
		t.Fatalf("expected content type %q, got %q", problemContentType, got)
	}

	var problem gen.Problem
	if err := json.NewDecoder(result.Body).Decode(&problem); err != nil {
		t.Fatalf("decode problem: %v", err)
	}
	if problem.Status != statusCode {
		t.Fatalf("unexpected problem status: got %d want %d", problem.Status, statusCode)
	}
	if detail != "" {
		if problem.Detail == nil {
			t.Fatalf("expected detail %q, got nil", detail)
		}
		if *problem.Detail != detail {
			t.Fatalf("unexpected detail: got %q want %q", *problem.Detail, detail)
		}
	}
}
