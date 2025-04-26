package werr

import (
	"fmt"
	"github.com/coding-common/errs/code"
)

const (
	errMsgSplit = ", "
)

type EIs interface {
	Error() string
	Is(target error) bool
}

type EErr interface {
	Error() string
	Code() code.RetCode
}
type EErrCode interface {
	Error() string
	SetCode(rc code.RetCode)
}
type EUnwrap interface {
	Error() string
	Unwrap() error
}

type WErr struct {
	error
	code int
	msg  string
}

func (e *WErr) Error() string {
	msg := e.msg
	if e.error != nil {
		return fmt.Sprintf("%s: %s", msg, e.error.Error())
	}
	return msg
}
