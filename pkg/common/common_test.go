package common

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Ekvo/golang-chi-postgres-api/internal/source"
	"github.com/Ekvo/golang-chi-postgres-api/internal/variables"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestDecodeJSON(t *testing.T) {
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

	encode := func(w http.ResponseWriter, httpStatus int, obj any) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpStatus)
		_ = json.NewEncoder(w).Encode(obj)
	}

	r.Post("/secret", func(w http.ResponseWriter, r *http.Request) {
		info := information{}
		if err := DecodeJSON(r, &info); err != nil {
			encode(w, http.StatusUnprocessableEntity, NewMessageError(variables.Validator, err))
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

func TestEncodeJSON(t *testing.T) {
	var encodeTestData = []struct {
		description string

		// write in url number of seconds for check Deadline in EncodeJSON		/
		ctxTimeOut int

		expectedCode   int
		responseRegexp string
		msg            string
	}{
		{
			description:    "Valid encode",
			ctxTimeOut:     2,
			expectedCode:   http.StatusOK,
			responseRegexp: `"approve"`,
			msg:            `valid Response, Header "Content-Type" "application/json" - code 200`,
		},
		{
			description:  "invalid encode - timeout",
			ctxTimeOut:   0,
			expectedCode: http.StatusGatewayTimeout,
			msg:          "invalid Response empty body code 504",
		},
	}

	r := chi.NewRouter()
	r.Get("/task/{ctx}", func(w http.ResponseWriter, r *http.Request) {
		timeOut, _ := strconv.Atoi(chi.URLParam(r, "ctx"))
		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(timeOut)*time.Second)
		defer func() {
			if ctx.Err() == context.DeadlineExceeded {
				w.WriteHeader(http.StatusGatewayTimeout)
			}
			cancel()
		}()
		time.Sleep(1 * time.Second)
		EncodeJSON(ctx, w, http.StatusOK, "approve")
	})

	asserts := assert.New(t)
	requires := require.New(t)

	for i, test := range encodeTestData {
		log.Printf("\t %d test encode: %s", i+1, test.description)

		req, err := http.NewRequest(http.MethodGet, "/task/"+strconv.Itoa(test.ctxTimeOut), nil)
		requires.NoError(err, fmt.Sprintf("http.NewRequest error - %v", err))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		code := test.expectedCode
		asserts.Equal(code, w.Code, "Code "+test.msg)

		if code == http.StatusGatewayTimeout {
			asserts.Empty(w.Body, "Response Body should be empty:")
		} else {
			media := w.Header().Get("Content-Type")
			asserts.NotEmpty(media, "Media is empty")
			asserts.Regexp(media, "application/json", "")
			asserts.Regexp(test.responseRegexp, w.Body.String(), "Body response "+test.msg)
		}
	}
}

func TestNewMessageError(t *testing.T) {
	asserts := assert.New(t)

	msgError := NewMessageError("login", source.ErrSourceNotFound)
	asserts.IsType(MessageError{}, msgError, "should be type - MessageError")
	asserts.Equal(Message{"login": "not found"}, msgError.Msg, `shoud be - map[string]any{"login": "not found"}`)
}
