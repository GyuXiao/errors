package errors

import (
	"errors"
	"fmt"
	"io"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		err  string
		want error
	}{
		{"", fmt.Errorf("")},
		{"foo", fmt.Errorf("foo")},
		{"foo", New("foo")},
		{"string with format specifiers: %v", errors.New("string with format specifiers: %v")},
	}

	for _, tt := range tests {
		got := New(tt.err)
		if got.Error() != tt.want.Error() {
			t.Errorf("New.Error(): got: %q, want %q", got, tt.want)
		}
	}
}

func TestErrorf(t *testing.T) {
	tests := []struct {
		err  error
		want string
	}{
		{Errorf("read error without format specifiers"), "read error without format specifiers"},
		{Errorf("read error with %d format specifier", 1), "read error with 1 format specifier"},
	}

	for _, tt := range tests {
		got := tt.err.Error()
		if got != tt.want {
			t.Errorf("Errorf(%v): got: %q, want %q", tt.err, got, tt.want)
		}
	}
}

func TestWithStackNil(t *testing.T) {
	got := WithStack(nil)
	if got != nil {
		t.Errorf("WithStack(nil): got %#v, expected nil", got)
	}
}

func TestWithStack(t *testing.T) {
	tests := []struct {
		err  error
		want string
	}{
		{io.EOF, "EOF"},
		{WithStack(io.EOF), "EOF"},
	}

	for _, tt := range tests {
		got := WithStack(tt.err).Error()
		if got != tt.want {
			t.Errorf("WithStack(%v): got: %v, want %v", tt.err, got, tt.want)
		}
	}
}

func TestWithMessageNil(t *testing.T) {
	got := WithMessage(nil, "no error")
	if got != nil {
		t.Errorf("WithMessage(nil, \"no error\"): got %#v, expected nil", got)
	}
}

func TestWithMessage(t *testing.T) {
	tests := []struct {
		err     error
		message string
		want    string
	}{
		{io.EOF, "read error", "read error"},
		{WithMessage(io.EOF, "read error"), "client error", "client error"},
	}

	for _, tt := range tests {
		got := WithMessage(tt.err, tt.message).Error()
		if got != tt.want {
			t.Errorf("WithMessage(%v, %q): got: %q, want %q", tt.err, tt.message, got, tt.want)
		}
	}
}

func TestWithMessagefNil(t *testing.T) {
	got := WithMessagef(nil, "no error")
	if got != nil {
		t.Errorf("WithMessage(nil, \"no error\"): got %#v, expected nil", got)
	}
}

func TestWithMessagef(t *testing.T) {
	tests := []struct {
		err     error
		message string
		want    string
	}{
		{io.EOF, "read error", "read error"},
		{WithMessagef(io.EOF, "read error without format specifier"), "client error", "client error"},
		{WithMessagef(io.EOF, "read error with %d format specifier", 1), "client error", "client error"},
	}

	for _, tt := range tests {
		got := WithMessagef(tt.err, tt.message).Error()
		if got != tt.want {
			t.Errorf("WithMessage(%v, %q): got: %q, want %q", tt.err, tt.message, got, tt.want)
		}
	}
}

func TestWithCode(t *testing.T) {
	tests := []struct {
		code     int
		message  string
		wantType string
		wantCode int
	}{
		{ConfigurationNotValid, "ConfigurationNotValid error", "*withCode", ConfigurationNotValid},
	}

	for _, tt := range tests {
		got := WithCode(tt.code, tt.message)
		err, ok := got.(*withCode)
		if !ok {
			t.Errorf("WithCode(%v, %q): error type got: %T, want %s", tt.code, tt.message, got, tt.wantType)
		}

		if err.code != tt.wantCode {
			t.Errorf("WithCode(%v, %q): got: %v, want %v", tt.code, tt.message, err.code, tt.wantCode)
		}
	}
}

func TestWithCodef(t *testing.T) {
	tests := []struct {
		code       int
		format     string
		args       string
		wantType   string
		wantCode   int
		wantString string
	}{
		{ConfigurationNotValid, "Configuration %s", "failed", "*withCode", ConfigurationNotValid, `Configuration failed`},
	}

	for _, tt := range tests {
		got := WithCode(tt.code, tt.format, tt.args)
		err, ok := got.(*withCode)
		if !ok {
			t.Errorf("WithCode(%v, %q %q): error type got: %T, want %s", tt.code, tt.format, tt.args, got, tt.wantType)
		}

		if err.code != tt.wantCode {
			t.Errorf("WithCode(%v, %q %q): got: %v, want %v", tt.code, tt.format, tt.args, err.code, tt.wantCode)
		}

		if got.Error() != tt.wantString {
			t.Errorf("WithCode(%v, %q %q): got: %v, want %v", tt.code, tt.format, tt.args, got.Error(), tt.wantString)
		}
	}
}

func TestErrorEquality(t *testing.T) {
	vals := []error{
		nil,
		io.EOF,
		errors.New("EOF"),
		New("EOF"),
		Errorf("EOF"),
		WithMessage(nil, "whoops"),
		WithMessage(io.EOF, "whoops"),
		WithStack(io.EOF),
		WithStack(nil),
	}

	for i := range vals {
		for j := range vals {
			_ = vals[i] == vals[j] // mustn't panic
		}
	}
}

func TestParseCoder(t *testing.T) {
	tests := []struct {
		err           error
		wantHTTPCode  int
		wantString    string
		wantCode      int
		wantReference string
	}{
		{fmt.Errorf("yes error"), 500, "An internal server error occurred", 1, "http://github.com/GyuXiao/errors/README.md"},
		{WithCode(unknownCoder.Code(), "internal error message"), 500, "An internal server error occurred", 1, "http://github.com/GyuXiao/errors/README.md"},
	}

	for i, tt := range tests {
		coder := ParseCode(tt.err)
		if coder.HttpStatus() != tt.wantHTTPCode {
			t.Errorf("TestCodeParse(%d): got %q, want: %q", i, coder.HttpStatus(), tt.wantHTTPCode)
		}

		if coder.String() != tt.wantString {
			t.Errorf("TestCodeParse(%d): got %q, want: %q", i, coder.String(), tt.wantString)
		}

		if coder.Code() != tt.wantCode {
			t.Errorf("TestCodeParse(%d): got %q, want: %q", i, coder.Code(), tt.wantCode)
		}

		if coder.Reference() != tt.wantReference {
			t.Errorf("TestCodeParse(%d): got %q, want: %q", i, coder.Reference(), tt.wantReference)
		}
	}
}