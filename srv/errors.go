package srv

import (
	"encoding/json"
	"net/http"
)

// APIError encapsulates error responses and corresponding status codes
type APIError struct {
	Err  string `json:"error"`
	Code int    `json:"-"`
}

func (e APIError) Error() string {
	return e.Err
}

// BadRequestError returns an APIError with the bad request status code
func BadRequestError(err string) APIError {
	return APIError{Err: err, Code: http.StatusBadRequest}
}

// NotFoundError returns an APIError with the bad request status code
func NotFoundError(err string) APIError {
	return APIError{Err: err, Code: http.StatusNotFound}
}

// InternalError returns an APIError with the bad request status code
func InternalError(err string) APIError {
	return APIError{Err: err, Code: http.StatusInternalServerError}
}

// encapsulate handling of statuscode and return errors as JSON
func errorHandler(f myHandlerFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			if apiError, ok := err.(APIError); ok {
				w.WriteHeader(apiError.Code)
				errResp, _ := json.Marshal(apiError) // just a string, should be marshallable
				w.Header().Set("Content-Type", "application/json")
				w.Write(errResp)
			} else {
				w.WriteHeader(500)
				w.Write([]byte(err.Error()))
			}
		}
	}
}
