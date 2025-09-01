package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
)

type Response struct {
	Status string `json:"status"`
	Data   any    `json:"data"`
	Error  string `json:"error"`
}

const (
	StatusOk    = "OK"
	StatusError = "ERROR"
)

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func GeneralError(err error) Response {
	return Response{
		Status: StatusError,
		Data:   nil,
		Error:  err.Error(),
	}
}

func ValidationError(err validator.ValidationErrors) Response {
	var errMsgs []string
	for _, e := range err {
		switch e.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("%s is required", e.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("%s is invalid", e.Field()))
		}
	}
	return Response{
		Status: StatusError,
		Data:   nil,
		Error:  strings.Join(errMsgs, ","),
	}
}
