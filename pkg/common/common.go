// common - universal utilities
package common

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"mime"
	"net/http"
)

// Message - body format for response
type Message map[string]any

// MessageError - error body format for response
type MessageError struct {
	Msg Message `json:"errors"`
}

func NewMessageError(key string, err error) MessageError {
	errStore := Message{key: err.Error()}
	return MessageError{Msg: errStore}
}

var ErrCommonInvalidMedia = errors.New("invalid media type")

// DecodeJSON - get object from 'Request'
// with checking an object for correct fields (look ~> ../../internal/servises/validator.go)
func DecodeJSON(r *http.Request, obj any) error {
	media := r.Header.Get("Content-Type")
	parse, _, err := mime.ParseMediaType(media)
	if err != nil || parse != "application/json" {
		return ErrCommonInvalidMedia
	}
	reqBody := r.Body
	defer reqBody.Close()
	dec := json.NewDecoder(reqBody)
	dec.DisallowUnknownFields()
	return dec.Decode(obj)
}

// EncodeJSON - we write the status and the object type of 'json' to 'ResponseWriter'
//
// if ctx.Err() == context.DeadlineExceeded - return us to 'func Timeout(timeout time.Duration)' ()
// (look ~> ../../internal/transport/transport.go)
func EncodeJSON(ctx context.Context, w http.ResponseWriter, status int, obj any) {
	if ctx.Err() == context.DeadlineExceeded {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(obj); err != nil {
		log.Printf("json.Encode error - %v", err)
	}
}
