package main

import (
	"net/url"
	"strings"
)

type formErrors map[string][]string

type Form struct {
	Data   url.Values
	Errors formErrors
}

func (e formErrors) Get(field string) string {
	efield := e[field]
	if len(efield) == 0 {
		return ""
	}
	return efield[0]
}

func (e *formErrors) Add(field string, value string) {
	concrete := *e
	concrete[field] = append(concrete[field], value)
}

func NewForm(data url.Values) *Form {
	return &Form{
		Data:   data,
		Errors: make(map[string][]string, 8),
	}
}

func (f *Form) Has(field string) bool {
	x := f.Data.Get(field)
	return x != ""
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Data.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

func (f *Form) Check(ok bool, key, message string) {
	if !ok {
		f.Errors.Add(key, message)
	}
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
