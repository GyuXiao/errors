package errors

import (
	"fmt"
	"net/http"
	"sync"
)

var (
	unknownCoder = defaultCoder{
		HTTP: http.StatusInternalServerError,
		Ext:  "An internal server error occurred",
		Ref:  "http://github.com/GyuXiao/errors/README.md",
		C:    1,
	}
)

type Coder interface {
	HttpStatus() int

	String() string

	Reference() string

	Code() int
}

type defaultCoder struct {
	HTTP int
	Ext  string
	Ref  string
	C    int
}

func (c defaultCoder) Code() int {
	return c.C
}

func (c defaultCoder) String() string {
	return c.Ext
}

func (c defaultCoder) HttpStatus() int {
	if c.HTTP == 0 {
		return 500
	}
	return c.HTTP
}

func (c defaultCoder) Reference() string {
	return c.Ref
}

var codes = map[int]Coder{}
var codeMux = &sync.Mutex{}

func Register(coder Coder) {
	if coder.Code() == 0 {
		panic(any("code `0` is reserved by `github.com/GyuXiao/errors` as unknownCode error code"))
	}

	codeMux.Lock()
	defer codeMux.Unlock()

	codes[coder.Code()] = coder
}

func MustRegister(coder Coder) {
	if coder.Code() == 0 {
		panic(any("code `0` is reserved by `github.com/GyuXiao/errors` as unknownCode error code"))
	}

	codeMux.Lock()
	defer codeMux.Unlock()

	if _, ok := codes[coder.Code()]; ok {
		panic(any(fmt.Sprintf("code: %d already exist", coder.Code())))
	}

	codes[coder.Code()] = coder
}

func ParseCode(err error) Coder {
	if err == nil {
		return nil
	}
	if v, ok := err.(*withCode); ok {
		if coder, ok := codes[v.code]; ok {
			return coder
		}
	}
	return unknownCoder
}

func IsCode(err error, code int) bool {
	if v, ok := err.(*withCode); ok {
		if v.code == code {
			return true
		}
		return false
	}
	return false
}

func init() {
	codes[unknownCoder.Code()] = unknownCoder
}
