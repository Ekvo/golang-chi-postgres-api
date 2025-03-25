package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Ekvo/golang-postgres-chi-api/internal/source"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMessageErrorFromValidator(t *testing.T) {
	type Information struct {
		Body   string `form:"body" json:"body" binding:"required,min=1,max=20"`
		Access string `form:"access" json:"access" binding:"required,numeric,len=1"`
	}

	var validatorTestData = []struct {
		description    string
		bodyData       string
		expectedCode   int
		responseRegexp string
		msg            string
	}{
		{
			description:    "Valid information",
			bodyData:       `{"body":"so big secret tcc","access":"1"}`,
			expectedCode:   http.StatusOK,
			responseRegexp: `{"information":"valid"}`,
			msg:            "valid data and return status 200",
		},
		{
			description:    "Invalid validating (max)",
			bodyData:       `{"body":"so big secret tcc! but it's no work here'","access":"1"}`,
			expectedCode:   http.StatusUnprocessableEntity,
			responseRegexp: `{"errors":{"Body":"{max:20}"}}`,
			msg:            "invalid return status 422",
		},
		{
			description:    "Invalid validating (len)",
			bodyData:       `{"body":"so big secret tcc!","access":"12"}`,
			expectedCode:   http.StatusUnprocessableEntity,
			responseRegexp: `{"errors":{"Access":"{len:1}"}}`,
			msg:            "invalid return status 422",
		},
		{
			description:    "Invalid validating (numeric)",
			bodyData:       `{"body":"so big secret tcc!","access":"a"}`,
			expectedCode:   http.StatusUnprocessableEntity,
			responseRegexp: `{"errors":{"Access":"{numeric:Access}"}}`,
			msg:            "invalid return status 422",
		},
		{
			description:    "Wrong information",
			bodyData:       `{"body":"so big secret tcc!","access":"2"}`,
			expectedCode:   http.StatusForbidden,
			responseRegexp: `{"errors":{"information":"invalid"}}`,
			msg:            "invalid return status 403",
		},
	}

	r := chi.NewRouter()

	encode := func(w http.ResponseWriter, httpStatus int, obj any) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpStatus)
		_ = json.NewEncoder(w).Encode(obj)
	}

	r.Post("/secret", func(w http.ResponseWriter, r *http.Request) {
		info := Information{}
		if err := Bind(r, &info); err != nil {
			encode(w, http.StatusUnprocessableEntity, NewMessageErrorFromValidator(err))
			return
		}
		if info.Body != "so big secret tcc" || info.Access != "1" {
			encode(w, http.StatusForbidden, NewMessageError("information", errors.New("invalid")))
			return
		}
		encode(w, http.StatusOK, Message{"information": "valid"})
	})

	asserts := assert.New(t)
	requires := require.New(t)

	for i, test := range validatorTestData {
		log.Printf("\t %d test validator: %s", i+1, test.description)
		bodyData := strings.Replace(test.bodyData, "\n", "", -1)

		req, err := http.NewRequest(http.MethodPost, "/secret", bytes.NewBuffer([]byte(bodyData)))
		requires.NoError(err, fmt.Sprintf("http.NewRequest error - %v", err))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		asserts.Equal(test.expectedCode, w.Code, test.msg)
		requires.NotEmpty(w.Body, "body from Response is empty")
		asserts.Regexp(test.responseRegexp, w.Body.String(), test.msg)
	}
}

func TestNewMessageError(t *testing.T) {
	asserts := assert.New(t)

	msgError := NewMessageError("login", source.ErrSourceNotFound)
	asserts.IsType(MessageError{}, msgError, "should be type - MessageError")
	asserts.Equal(Message{"login": "not found"}, msgError.Msg, "should be")
}
