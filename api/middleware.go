package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ankeshnirala/go/go_services/constants"
	"github.com/ankeshnirala/go/go_services/types"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	res := &types.StandardResponse{
		Code: status,
		Data: v,
	}

	return json.NewEncoder(w).Encode(res)
}

func MakeHTTPHandleFunc(f types.ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// handle the error
			WriteJSON(w, http.StatusBadRequest, types.ApiError{Error: err.Error()})
		}
	}
}

func (s *Server) PathNotFound(w http.ResponseWriter, r *http.Request) error {

	resp := types.ApiError{
		Error: r.URL.Path + " " + constants.NOT_FOUND,
	}

	return WriteJSON(w, http.StatusNotFound, resp)
}

func (s *Server) MethodNotAllowed(w http.ResponseWriter, r *http.Request) error {
	return WriteJSON(w, http.StatusMethodNotAllowed, types.ApiError{Error: fmt.Errorf(constants.METHODNOTALLOWED, r.Method).Error()})
}
