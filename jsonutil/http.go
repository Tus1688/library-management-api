package jsonutil

import (
	"bytes"
	"fmt"
	"github.com/goccy/go-json"
	"io"
	"net/http"
	"reflect"
)

func Render(w http.ResponseWriter, statusCode int, v interface{}) error {
	buf, err := MarshalJSON(v)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_, err = w.Write(buf)
	return err
}

func RenderJSONBytes(w http.ResponseWriter, statusCode int, v []byte) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_, err := w.Write(v)
	return err
}

func MarshalJSON(v interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func UnmarshalJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

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
func isZeroOfUnderlyingType(x any) bool {
	return x == nil || reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}
