package main

import (
	"fmt"
	"github.com/GyuXiao/errors"
	"github.com/novalagung/gubrak"
)

func main() {
	if err := bindUser(); err != nil {
		// %s: Returns the user-safe error string mapped to the error code or the error message if none is specified.
		fmt.Println("====================> %s <====================")
		fmt.Printf("%s\n\n", err)

		// %v: Alias for %s.
		fmt.Println("====================> %v <====================")
		fmt.Printf("%v\n\n", err)

		// %-v: Output caller details, useful for troubleshooting.
		fmt.Println("====================> %-v <====================")
		fmt.Printf("%-v\n\n", err)

		// %+v: Output full error stack details, useful for debugging.
		fmt.Println("====================> %+v <====================")
		fmt.Printf("%+v\n\n", err)

		// %#-v: Output caller details, useful for troubleshooting with JSON formatted output.
		fmt.Println("====================> %#-v <====================")
		fmt.Printf("%#-v\n\n", err)

		// %#+v: Output full error stack details, useful for debugging with JSON formatted output.
		fmt.Println("====================> %#+v <====================")
		fmt.Printf("%#+v\n\n", err)

		// do some business process based on the error type
		if errors.IsCode(err, ErrEncodingFailed) {
			fmt.Println("this is a ErrEncodingFailed error")
		}

		if errors.IsCode(err, ErrDatabase) {
			fmt.Println("this is a ErrDatabase error")
		}
	}
}

func bindUser() error {
	if err := getUser(); err != nil {
		return errors.WrapC(err, ErrEncodingFailed, "encoding User Gyu failed")
	}
	return nil
}

func getUser() error {
	if err := queryDatabase(); err != nil {
		return errors.Wrap(err, "get user failed")
	}
	return nil
}

func queryDatabase() error {
	return errors.WithCode(ErrDatabase, "user Gyu not found")
}

func init() {
	register(ErrUserNotFound, 404, "User not found")
	register(ErrUserAlreadyExist, 400, "User already exist")
	register(ErrReachMaxCount, 400, "Secret reach the max count")
	register(ErrSecretNotFound, 404, "Secret not found")
	register(ErrSuccess, 200, "OK")
	register(ErrUnknown, 500, "Internal server error")
	register(ErrBind, 400, "Error occurred while binding the request body to the struct")
	register(ErrValidation, 400, "Validation failed")
	register(ErrTokenInvalid, 401, "Token invalid")
	register(ErrDatabase, 500, "Database error")
	register(ErrEncrypt, 401, "Error occurred while encrypting the user password")
	register(ErrSignatureInvalid, 401, "Signature is invalid")
	register(ErrExpired, 401, "Token expired")
	register(ErrInvalidAuthHeader, 401, "Invalid authorization header")
	register(ErrMissingHeader, 401, "The `Authorization` header was empty")
	register(ErrorExpired, 401, "Token expired")
	register(ErrPasswordIncorrect, 401, "Password was incorrect")
	register(ErrPermissionDenied, 403, "Permission denied")
	register(ErrEncodingFailed, 500, "Encoding failed due to an error with the data")
	register(ErrDecodingFailed, 500, "Decoding failed due to an error with the data")
	register(ErrInvalidJSON, 500, "Data is not valid JSON")
	register(ErrEncodingJSON, 500, "JSON data could not be encoded")
	register(ErrDecodingJSON, 500, "JSON data could not be decoded")
	register(ErrInvalidYaml, 500, "Data is not valid Yaml")
	register(ErrEncodingYaml, 500, "Yaml data could not be encoded")
	register(ErrDecodingYaml, 500, "Yaml data could not be decoded")
}

func register(code int, httpStatus int, message string, refs ...string) {
	found, _ := gubrak.Includes([]int{200, 400, 401, 403, 404, 500}, httpStatus)
	if !found {
		panic(any("http code not in `200, 400, 401, 403, 404, 500`"))
	}

	var reference string
	if len(refs) > 0 {
		reference = refs[0]
	}

	coder := &ErrCode{
		C:    code,
		HTTP: httpStatus,
		Ext:  message,
		Ref:  reference,
	}

	errors.MustRegister(coder)
}

type ErrCode struct {
	C    int
	HTTP int
	Ext  string
	Ref  string
}

func (coder ErrCode) Code() int {
	return coder.C
}

func (coder ErrCode) String() string {
	return coder.Ext
}

func (coder ErrCode) Reference() string {
	return coder.Ref
}

func (coder ErrCode) HttpStatus() int {
	if coder.HTTP == 0 {
		return 500
	}
	return coder.HTTP
}
