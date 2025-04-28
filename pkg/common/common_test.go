package common

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Ekvo/golang-chi-postgres-api/internal/variables"
)

func Test_Message_String(t *testing.T) {
	msg := Message{
		"id":      "1111",
		"title":   "some woed",
		"error":   "nil",
		"created": "empty",
	}
	msgLine := msg.String()
	assert.Equal(t, `{created:empty},{error:nil},{id:1111},{title:some woed}`, msgLine)
	assert.NotEqual(t, `{error:nil},{created:empty},{id:1111},{title:some woed}`, msgLine)
}

func Test_DecodeJSON_EncodeJSON(t *testing.T) {
	type information struct {
		Body   string `json:"body"`
		Access string `json:"access"`
	}

	var decodeTestData = []struct {
		description       string
		bodyData          string
		headerContentType []string // [0]-key,[1]-value
		expectedCode      int
		responseRegexp    string
		msg               string
	}{
		{
			description:       "Valid information",
			bodyData:          `{"body":"so big secret tcc","access":"1"}`,
			headerContentType: []string{"Content-Type", "application/json"},
			expectedCode:      http.StatusOK,
			responseRegexp:    `{"information":"valid"}`,
			msg:               "valid data and return status 200",
		},
		{
			description:       "Invalid media type - value",
			bodyData:          `{"body":"so big secret tcc! but it's no work here'","access":"1"}`,
			headerContentType: []string{"Content-Type", "qwe"},
			expectedCode:      http.StatusUnprocessableEntity,
			responseRegexp:    `{"errors":{"validator":"invalid media type"}}`,
			msg:               "invalid return status 422",
		},
		{
			description:       "Invalid media type - key",
			bodyData:          `{"body":"so big secret tcc!","access":"12"}`,
			headerContentType: []string{"some-Type", "application/json"},
			expectedCode:      http.StatusUnprocessableEntity,
			responseRegexp:    `{"errors":{"validator":"invalid media type"}}`,
			msg:               "invalid return status 422",
		},
		{
			description:       "Wrong information",
			bodyData:          `{"body":"so big secret tcc!","access":"2"}`,
			headerContentType: []string{"Content-Type", "application/json"},
			expectedCode:      http.StatusForbidden,
			responseRegexp:    `{"errors":{"information":"invalid"}}`,
			msg:               "invalid return status 403",
		},
		{
			description:       "data with 'alien' field",
			bodyData:          `{"body":"so big secret tcc!","access":"1","alien":"UFO"}`,
			headerContentType: []string{"Content-Type", "application/json"},
			expectedCode:      http.StatusUnprocessableEntity,
			responseRegexp:    `{"errors":{"validator":"json: unknown field \\"alien\\""}}`,
			msg:               "invalid return status 403",
		},
	}

	r := chi.NewRouter()

	r.Post("/secret", func(w http.ResponseWriter, r *http.Request) {
		info := information{}
		if err := DecodeJSON(r, &info); err != nil {
			EncodeJSON(w, http.StatusUnprocessableEntity, NewMessageError(variables.Validator, err))
			return
		}
		if info.Body != "so big secret tcc" || info.Access != "1" {
			EncodeJSON(w, http.StatusForbidden, NewMessageError("information", errors.New("invalid")))
			return
		}
		EncodeJSON(w, http.StatusOK, Message{"information": "valid"})
	})

	asserts := assert.New(t)
	requires := require.New(t)

	for i, test := range decodeTestData {
		log.Printf("\t %d test decode: %s", i+1, test.description)
		bodyData := strings.Replace(test.bodyData, "\n", "", -1)

		req, err := http.NewRequest(http.MethodPost, "/secret", bytes.NewBuffer([]byte(bodyData)))
		requires.NoError(err, fmt.Sprintf("http.NewRequest error - %v", err))
		req.Header.Set(test.headerContentType[0], test.headerContentType[1])

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		asserts.Equal(test.expectedCode, w.Code, test.msg)
		requires.NotEmpty(w.Body, "body from Response is empty")
		asserts.Regexp(test.responseRegexp, w.Body.String(), test.msg)
	}
}

func TestNewMessageError(t *testing.T) {
	asserts := assert.New(t)

	msgError := NewMessageError("param", ErrCommonInvalidMedia)
	asserts.IsType(MessageError{}, msgError, "should be type - MessageError")
	asserts.Equal(Message{"param": "invalid media type"}, msgError.Msg, `shoud be - map[string]any{"param": "invalid media type"}`)
}
