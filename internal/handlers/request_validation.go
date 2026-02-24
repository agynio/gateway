package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	legacy "github.com/getkin/kin-openapi/routers/legacy"
)

var allowedMethodCandidates = []string{
	http.MethodGet,
	http.MethodHead,
	http.MethodPost,
	http.MethodPut,
	http.MethodDelete,
	http.MethodPatch,
	http.MethodOptions,
	http.MethodTrace,
}

func NewRequestValidationMiddleware(swagger *openapi3.T) (func(http.Handler) http.Handler, error) {
	if swagger == nil {
		return nil, fmt.Errorf("swagger specification is required")
	}

	router, err := legacy.NewRouter(swagger)
	if err != nil {
		return nil, fmt.Errorf("build router: %w", err)
	}

	opts := &openapi3filter.Options{AuthenticationFunc: openapi3filter.NoopAuthenticationFunc}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			input, statusCode, err := buildRequestValidationInput(r, router, opts)
			if err != nil {
				if statusCode == http.StatusMethodNotAllowed {
					if allow := allowedMethodsForRequest(r, router); len(allow) > 0 {
						w.Header().Set("Allow", strings.Join(allow, ", "))
					}
				}
				RequestValidationError(w, err.Error(), statusCode)
				return
			}

			if err := openapi3filter.ValidateRequest(r.Context(), input); err != nil {
				if status, handledErr := handleRequestValidationError(err); handledErr != nil {
					RequestValidationError(w, handledErr.Error(), status)
					return
				}

				RequestValidationError(w, fmt.Sprintf("error validating request: %s", err.Error()), http.StatusInternalServerError)
				return
			}

			next.ServeHTTP(w, r)
		})
	}, nil
}

func buildRequestValidationInput(r *http.Request, router routers.Router, opts *openapi3filter.Options) (*openapi3filter.RequestValidationInput, int, error) {
	route, pathParams, err := router.FindRoute(r)
	if err != nil {
		status := http.StatusNotFound
		var routeErr *routers.RouteError
		if errors.As(err, &routeErr) && strings.EqualFold(routeErr.Error(), routers.ErrMethodNotAllowed.Error()) {
			status = http.StatusMethodNotAllowed
		}
		return nil, status, err
	}

	input := &openapi3filter.RequestValidationInput{
		Request:    r,
		PathParams: pathParams,
		Route:      route,
	}

	if opts != nil {
		input.Options = opts
	}

	return input, 0, nil
}

func allowedMethodsForRequest(r *http.Request, router routers.Router) []string {
	allowed := make([]string, 0, len(allowedMethodCandidates))

	for _, method := range allowedMethodCandidates {
		if method == r.Method {
			continue
		}

		clone := r.Clone(r.Context())
		clone.Method = method
		clone.Body = http.NoBody

		if _, _, err := router.FindRoute(clone); err == nil {
			allowed = append(allowed, method)
		}
	}

	return allowed
}

func handleRequestValidationError(err error) (int, error) {
	var multi openapi3.MultiError
	if errors.As(err, &multi) {
		status, multiErr := RequestValidationMultiError(multi)
		return status, multiErr
	}

	var reqErr *openapi3filter.RequestError
	if errors.As(err, &reqErr) {
		message := reqErr.Error()
		if idx := strings.IndexRune(message, '\n'); idx >= 0 {
			message = message[:idx]
		}
		return http.StatusBadRequest, fmt.Errorf(message)
	}

	var secErr *openapi3filter.SecurityRequirementsError
	if errors.As(err, &secErr) {
		return http.StatusUnauthorized, secErr
	}

	return 0, nil
}
