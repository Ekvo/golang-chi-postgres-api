// common - универсальные инструменты
package common

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Message - формат тела для Response
type Message map[string]any

// MessageError - формат тела ошибки для Response
type MessageError struct {
	Msg Message `json:"errors"`
}

func NewMessageError(key string, err error) MessageError {
	errStore := Message{key: err.Error()}
	return MessageError{Msg: errStore}
}

// ошибка для защиты от некоректного использования:
// getFieldName(structNamespace string) (string, error)
var ErrCommonFiledValidatorIncorrect = errors.New("incorrect struct namespace")

// FiledValidator - для создания кастомных функций обработки ошибок
// объекта validator.FieldError
type FiledValidator struct {
	validator.FieldError
}

// FieldName - получении имени поля из validator.FieldError StructNamespace()
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

// NewMessageErrorFromValidator - обработка ошибки
// после ('func Bind(r *http.Request, obj any) error')
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

// Bind - получение объекта из 'Request'
// с проверкой объекта на определенный свойсва (см. servises)
func Bind(r *http.Request, obj any) error {
	b := binding.Default(r.Method, r.Header.Get("Content-Type"))
	return b.Bind(r, obj)
}
