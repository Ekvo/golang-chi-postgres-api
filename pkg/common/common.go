// common - universal utilities
package common

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
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

// label for incorrect use - 'func (fv FiledValidator) FieldName() (string, error)'
// getFieldName(structNamespace string) (string, error)
var ErrCommonFiledValidatorIncorrect = errors.New("incorrect struct namespace")

// FiledValidator - creating custom error handling functions
// for object 'validator.FieldError'
type FiledValidator struct {
	validator.FieldError
}

// FieldName - get name field from validator.FieldError StructNamespace()
func (fv FiledValidator) FieldName() (string, error) {
	structNamespace := fv.StructNamespace()
	n := len(structNamespace)
	if n < 3 ||
		structNamespace[0] == '.' ||
		structNamespace[n-1] == '.' {
		return "", ErrCommonFiledValidatorIncorrect
	}
	for i := n - 1; i > -1; i-- {
		if structNamespace[i] == '.' {
			return structNamespace[i+1 : n], nil
		}
	}
	return "", ErrCommonFiledValidatorIncorrect
}

// NewMessageErrorFromValidator - handler error
// after ('func Bind(r *http.Request, obj any) error')
func NewMessageErrorFromValidator(err error) MessageError {
	dataErr := err.(validator.ValidationErrors)
	errStore := make(Message)
	for _, field := range dataErr {
		info := field.Param()
		if len(info) == 0 {
			fv := FiledValidator{field}
			info, _ = fv.FieldName()
		}
		errStore[field.Field()] = fmt.Sprintf("{%v:%v}", field.Tag(), info)
	}
	return MessageError{Msg: errStore}
}

// Bind - get object from 'Request'
// with checking an object for certain properties (look ~> ../../internal/servises/validator.go)
func Bind(r *http.Request, obj any) error {
	b := binding.Default(r.Method, r.Header.Get("Content-Type"))
	return b.Bind(r, obj)
}
