package errno

import (
	"fmt"
)

type Errno struct {
	Code    int
	Msg 	string
}

func (err Errno) Error() string {
	return err.Msg
}

type Err struct {
	Code    int
	Msg 	string
	Err     error
}

func (err *Err) Error() string {
	return fmt.Sprintf("Err - code: %d, message: %s, error: %s", err.Code, err.Msg, err.Err)
}

func New(errno *Errno, err error) *Err {
	return &Err{Code: errno.Code, Msg: errno.Msg, Err: err}
}

func (err *Err) Add(message string) error {
	err.Msg += " " + message
	return err
}

func (err *Err) Addf(format string, args ...interface{}) error {
	err.Msg += " " + fmt.Sprintf(format, args...)
	return err
}

func DecodeErr(err error) (int, string) {
	if err == nil {
		return OK.Code, OK.Msg
	}

	switch typed := err.(type) {
	case *Err:
		return typed.Code, typed.Msg
	case *Errno:
		return typed.Code, typed.Msg
	default:
	}

	return InternalServerError.Code, err.Error()
}
