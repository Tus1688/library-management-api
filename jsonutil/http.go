package jsonutil

import (
	"bytes"
	"fmt"
	"github.com/goccy/go-json"
	"io"
	"net/http"
	"reflect"
)

// Render writes the JSON-encoded response to the http.ResponseWriter.
// Parameters:
// - w: the http.ResponseWriter to write to.
// - statusCode: the HTTP status code to set in the response.
// - v: the value to encode as JSON.
// Returns an error if the encoding or writing fails.
func Render(w http.ResponseWriter, statusCode int, v interface{}) error {
	buf, err := MarshalJSON(v)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_, err = w.Write(buf)
	return err
}

// RenderJSONBytes writes the JSON-encoded byte slice to the http.ResponseWriter.
// Parameters:
// - w: the http.ResponseWriter to write to.
// - statusCode: the HTTP status code to set in the response.
// - v: the byte slice to write as the response body.
// Returns an error if the writing fails.
func RenderJSONBytes(w http.ResponseWriter, statusCode int, v []byte) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_, err := w.Write(v)
	return err
}

// MarshalJSON encodes the given value as a JSON byte slice.
// Parameters:
// - v: the value to encode as JSON.
// Returns the JSON-encoded byte slice and an error if the encoding fails.
func MarshalJSON(v interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnmarshalJSON decodes the JSON-encoded byte slice into the given value.
// Parameters:
// - data: the JSON-encoded byte slice.
// - v: the value to decode into.
// Returns an error if the decoding fails.
func UnmarshalJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// ShouldBind decodes the JSON-encoded request body into the given value and checks for required fields.
// Parameters:
// - r: the HTTP request containing the JSON-encoded body.
// - v: the value to decode into.
// Returns an error if the decoding or required field check fails.
func ShouldBind(r *http.Request, v interface{}) error {
	err := json.NewDecoder(r.Body).Decode(&v)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)
	if err != nil {
		return err
	}
	return checkRequiredFields(v)
}

// checkRequiredFields checks if the required fields in the given value are set.
// Parameters:
// - v: the value to check for required fields.
// Returns an error if any required field is not set.
func checkRequiredFields(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typeOfV := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		tag := typeOfV.Field(i).Tag.Get("binding")
		if tag == "required" && isZeroOfUnderlyingType(field.Interface()) {
			return fmt.Errorf("field %s is required", typeOfV.Field(i).Name)
		}
	}
	return nil
}

// isZeroOfUnderlyingType checks if the given value is the zero value of its underlying type.
// Parameters:
// - x: the value to check.
// Returns true if the value is the zero value, false otherwise.
func isZeroOfUnderlyingType(x any) bool {
	return x == nil || reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}
