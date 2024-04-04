package main

import (
	"net/url"
	"testing"
)

func TestForm_Has(t *testing.T) {
	form := NewForm(nil)
	has := form.Has("whatever")
	if has {
		t.Error("Form shows that it has a key that definitively does not exist")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	form = NewForm(postedData)
	has = form.Has("a")
	if !has {
		t.Error("Could not find the added key")
	}
}

func TestForm_required(t *testing.T) {
	tests := []struct {
		name           string
		requiredFields []string
		data           map[string][]string
		passes         bool
	}{
		{"nothing required, nothing possessed", make([]string, 0), make(map[string][]string), true},
		{"all requirements met", []string{"a", "b"}, map[string][]string{"a": {"a"}, "b": {"b"}}, true},
		{"lacking data", []string{"a", "q"}, map[string][]string{"a": {"a"}}, false},
	}

	for _, test := range tests {
		var data url.Values = test.data
		form := NewForm(data)
		form.Required(test.requiredFields...)
		if form.Valid() != test.passes {
			t.Errorf("test %s: Expected to pass {%v} but got {%v} instead", test.name, test.passes, form.Valid())
		} else if !form.Valid() {
			var errorMsg string
			for _, field := range test.requiredFields {
				if len(errorMsg) == 0 {
					errorMsg = form.Errors.Get(field)
				}
			}
			expected := "This field cannot be blank"
			if errorMsg != expected {
				t.Errorf("test %s: Expected error message {%s} but got {%s}", test.name, expected, errorMsg)
			}
		}
	}
}

func TestForm_Check(t *testing.T) {
	form := NewForm(nil)
	field := "password"
	expectedMsg := "password is required"
	form.Check(false, field, expectedMsg)
	if form.Valid() {
		t.Error("Valid() returns true though it should return false as the check is false")
	}
	actual := form.Errors.Get(field)
	if actual != expectedMsg {
		t.Errorf("Expected to find message %s on error key %s, but got %s",
			expectedMsg, field, actual)
	}
}
