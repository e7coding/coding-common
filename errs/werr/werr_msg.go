package werr

import (
	"fmt"
	"github.com/coding-common/errs/code"
)

func WithMsg(msg string) *WErr {
	return &WErr{
		code: code.Nil,
		msg:  msg,
	}
}

func WithMsgF(format string, args ...interface{}) *WErr {
	return &WErr{
		code: code.Nil,
		msg:  fmt.Sprintf(format, args...),
	}
}

func WithMsgErr(err error, msg string) error {
	if err == nil {
		return nil
	}
	return &WErr{
		error: err,
		msg:   msg,
		code:  ToRetCode(err).Code(),
	}
}

func WithMsgErrF(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &WErr{
		error: err,
		msg:   fmt.Sprintf(format, args...),
		code:  ToRetCode(err).Code(),
	}
}
