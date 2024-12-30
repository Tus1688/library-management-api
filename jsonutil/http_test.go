package jsonutil

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestShouldBind(t *testing.T) {
	// Create a mock request object
	req := httptest.NewRequest("POST", "/path", strings.NewReader(`{"foo": "bar"}`))

	// Create a mock object to decode the request body into
	var v struct {
		Foo string `json:"foo" binding:"required"`
	}

	// Call the ShouldBind function with the mock request and object
	err := ShouldBind(req, &v)

	// Check that there was no error
	if err != nil {
		t.Errorf("ShouldBind returned an error: %v", err)
	}

	// Check that the object was correctly decoded
	if v.Foo != "bar" {
		t.Errorf("ShouldBind decoded the wrong value: got %q, want %q", v.Foo, "bar")
	}
}

func TestShouldBindFaulty(t *testing.T) {
	// Create a mock request object
	req := httptest.NewRequest("POST", "/path", strings.NewReader(`{"bla": "bar"}`))

	// Create a mock object to decode the request body into
	var v struct {
		Foo string `json:"foo" binding:"required"`
	}

	// Call the ShouldBind function with the mock request and object
	err := ShouldBind(req, &v)

	// Check that there was no error
	if err == nil {
		t.Errorf("ShouldBind cannot detect binding error")
	}
}
