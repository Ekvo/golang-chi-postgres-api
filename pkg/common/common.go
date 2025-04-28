// common - universal utilities
package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mime"
	"net/http"
	"sort"
	"strings"
)

var ErrCommonInvalidMedia = errors.New("invalid media type")

// Message - body format for response
type Message map[string]any

// String - many keys in message
//
// sort keys - result was predictable
func (msg Message) String() string {
	lineMsg := make([]string, 0, len(msg))
	for k, v := range msg {
		lineMsg = append(lineMsg, fmt.Sprintf(`{%s:%v}`, k, v))
	}
	sort.Strings(lineMsg)
	return strings.Join(lineMsg, ",")
}

// MessageError - error body format for response
type MessageError struct {
	Msg Message `json:"errors"`
}

func NewMessageError(key string, err error) MessageError {
	errStore := Message{key: err.Error()}
	return MessageError{Msg: errStore}
}

// DecodeJSON - get object from 'Request'
// with checking an object for correct fields (look ~> ../../internal/servises/validator.go)
func DecodeJSON(r *http.Request, obj any) error {
	media := r.Header.Get("Content-Type")
	parse, _, err := mime.ParseMediaType(media)
	if err != nil || parse != "application/json" {
		return ErrCommonInvalidMedia
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("common: r.Body.Close error - %v", err)
		}
	}()
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(obj)
}

// EncodeJSON - we write the status and the object type of 'json' to 'ResponseWriter'
func EncodeJSON(w http.ResponseWriter, status int, obj any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(obj); err != nil {
		log.Printf("json.Encode error - %v", err)
	}
}
