// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package field

import (
	"fmt"
	"strings"
	"testing"
)

func TestMakeFuncs(t *testing.T) {
	testCases := []struct {
		fn       func() *Error
		expected ErrorType
	}{
		{
			func() *Error { return Invalid(NewPath("f"), "v", "d") },
			ErrorTypeInvalid,
		},
		{
			func() *Error { return NotSupported(NewPath("f"), "v", nil) },
			ErrorTypeNotSupported,
		},
		{
			func() *Error { return Duplicate(NewPath("f"), "v") },
			ErrorTypeDuplicate,
		},
		{
			func() *Error { return NotFound(NewPath("f"), "v") },
			ErrorTypeNotFound,
		},
		{
			func() *Error { return Required(NewPath("f"), "d") },
			ErrorTypeRequired,
		},
		{
			func() *Error { return InternalError(NewPath("f"), fmt.Errorf("e")) },
			ErrorTypeInternal,
		},
	}

	for _, testCase := range testCases {
		err := testCase.fn()
		if err.Type != testCase.expected {
			t.Errorf("expected Type %q, got %q", testCase.expected, err.Type)
		}
	}
}

func TestErrorUsefulMessage(t *testing.T) {
	{
		s := Invalid(nil, nil, "").Error()
		t.Logf("message: %v", s)
		if !strings.Contains(s, "null") {
			t.Errorf("error message did not contain 'null': %s", s)
		}
	}

	s := Invalid(NewPath("foo"), "bar", "deet").Error()
	t.Logf("message: %v", s)
	for _, part := range []string{"foo", "bar", "deet", ErrorTypeInvalid.String()} {
		if !strings.Contains(s, part) {
			t.Errorf("error message did not contain expected part '%v'", part)
		}
	}

	type complicated struct {
		Baz   int
		Qux   string
		Inner interface{}
		KV    map[string]int
	}
	s = Invalid(
		NewPath("foo"),
		&complicated{
			Baz:   1,
			Qux:   "aoeu",
			Inner: &complicated{Qux: "asdf"},
			KV:    map[string]int{"Billy": 2},
		},
		"detail",
	).Error()
	t.Logf("message: %v", s)
	for _, part := range []string{
		"foo", ErrorTypeInvalid.String(),
		"Baz", "Qux", "Inner", "KV", "detail",
		"1", "aoeu", "Billy", "2",
		// "asdf", TODO: re-enable once we have a better nested printer
	} {
		if !strings.Contains(s, part) {
			t.Errorf("error message did not contain expected part '%v'", part)
		}
	}
}

func TestToAggregate(t *testing.T) {
	testCases := struct {
		ErrList         []ErrorList
		NumExpectedErrs []int
	}{
		[]ErrorList{
			nil,
			{},
			{Invalid(NewPath("f"), "v", "d")},
			{Invalid(NewPath("f"), "v", "d"), Invalid(NewPath("f"), "v", "d")},
			{Invalid(NewPath("f"), "v", "d"), InternalError(NewPath(""), fmt.Errorf("e"))},
		},
		[]int{
			0,
			0,
			1,
			1,
			2,
		},
	}

	if len(testCases.ErrList) != len(testCases.NumExpectedErrs) {
		t.Errorf("Mismatch: length of NumExpectedErrs does not match length of ErrList")
	}
	for i, tc := range testCases.ErrList {
		agg := tc.ToAggregate()
		numErrs := 0

		if agg != nil {
			numErrs = len(agg.Errors())
		}
		if numErrs != testCases.NumExpectedErrs[i] {
			t.Errorf("[%d] Expected %d, got %d", i, testCases.NumExpectedErrs[i], numErrs)
		}

		if len(tc) == 0 {
			if agg != nil {
				t.Errorf("[%d] Expected nil, got %#v", i, agg)
			}
		} else if agg == nil {
			t.Errorf("[%d] Expected non-nil", i)
		}
	}
}

func TestErrListFilter(t *testing.T) {
	list := ErrorList{
		Invalid(NewPath("test.field"), "", ""),
		Invalid(NewPath("field.test"), "", ""),
		Duplicate(NewPath("test"), "value"),
	}
	if len(list.Filter(NewErrorTypeMatcher(ErrorTypeDuplicate))) != 2 {
		t.Errorf("should not filter")
	}
	if len(list.Filter(NewErrorTypeMatcher(ErrorTypeInvalid))) != 1 {
		t.Errorf("should filter")
	}
}

func TestNotSupported(t *testing.T) {
	notSupported := NotSupported(NewPath("f"), "v", []string{"a", "b", "c"})
	expected := `Unsupported value: "v": supported values: "a", "b", "c"`
	if notSupported.ErrorBody() != expected {
		t.Errorf("Expected: %s\n, but got: %s\n", expected, notSupported.ErrorBody())
	}
}
