package errors

type Coder interface {
	HttpStatus() int

	String() string

	Reference() string

	Code() int
}
