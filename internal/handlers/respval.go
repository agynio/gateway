package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	legacy "github.com/getkin/kin-openapi/routers/legacy"
)

func NewResponseValidationMiddleware(swagger *openapi3.T) (func(http.Handler) http.Handler, error) {
	if swagger == nil {
		return nil, fmt.Errorf("swagger specification is required")
	}

	router, err := legacy.NewRouter(swagger)
	if err != nil {
		return nil, fmt.Errorf("build router: %w", err)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			recorder := httptest.NewRecorder()
			next.ServeHTTP(recorder, r)

			if err := validateResponse(r, recorder, router); err != nil {
				problem := NewProblem(http.StatusInternalServerError, "Response validation failed", err.Error(), nil)
				WriteProblem(w, problem)
				return
			}

			copyResponse(w, recorder)
		})
	}, nil
}

func validateResponse(r *http.Request, recorder *httptest.ResponseRecorder, router routers.Router) error {
	route, pathParams, err := router.FindRoute(r)
	if err != nil {
		var routeErr *routers.RouteError
		if errors.As(err, &routeErr) {
			return nil
		}
		return fmt.Errorf("match route: %w", err)
	}

	requestInput := &openapi3filter.RequestValidationInput{
		Request:    r,
		PathParams: pathParams,
		Route:      route,
	}

	responseInput := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: requestInput,
		Status:                 recorder.Code,
		Header:                 cloneHeader(recorder.Header()),
	}
	responseInput.SetBodyBytes(recorder.Body.Bytes())

	return openapi3filter.ValidateResponse(r.Context(), responseInput)
}

func copyResponse(dst http.ResponseWriter, src *httptest.ResponseRecorder) {
	for key, values := range src.Header() {
		copied := append([]string(nil), values...)
		dst.Header()[key] = copied
	}
	dst.WriteHeader(src.Code)
	if _, err := dst.Write(src.Body.Bytes()); err != nil {
		// nothing better can be done at this point
	}
}

func cloneHeader(header http.Header) http.Header {
	clone := make(http.Header, len(header))
	for key, values := range header {
		clone[key] = append([]string(nil), values...)
	}
	return clone
}
