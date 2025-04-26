package werr

import (
	"errors"
	"fmt"
	"github.com/coding-common/errs/code"
	"strings"
)

func WithCode(rc code.RetCode, msg ...string) *WErr {
	return &WErr{
		msg:  fmt.Sprintf("[%s] %s", rc.Msg(), strings.Join(msg, errMsgSplit)),
		code: rc.Code(),
	}
}

func WithCodeF(rc code.RetCode, format string, args ...interface{}) *WErr {
	return &WErr{
		msg:  fmt.Sprintf("[%s] %s", rc.Msg(), fmt.Sprintf(format, args...)),
		code: rc.Code(),
	}
}

func WithCodeErr(rc code.RetCode, err error, msg ...string) *WErr {
	if err == nil {
		return nil
	}
	return &WErr{
		error: err,
		msg:   fmt.Sprintf("[%s] %s", rc.Msg(), strings.Join(msg, errMsgSplit)),
		code:  rc.Code(),
	}
}

func WithCodeErrF(rc code.RetCode, err error, format string, args ...interface{}) *WErr {
	if err == nil {
		return nil
	}
	return &WErr{
		error: err,
		msg:   fmt.Sprintf("[%s] %s", rc.Msg(), fmt.Sprintf(format, args...)),
		code:  rc.Code(),
	}
}

func ToRetCode(err error) code.RetCode {
	if err == nil {
		return code.NewErrCode(code.Nil)
	}
	var e1 EErr
	if errors.As(err, &e1) {
		return e1.Code()
	}
	var e2 EUnwrap
	if errors.As(err, &e2) {
		return ToRetCode(e2.Unwrap())
	}
	return code.NewErrCode(code.Nil)
}
func HasCode(err error, rc code.RetCode) bool {
	if err == nil {
		return false
	}
	var e1 EErr
	if errors.As(err, &e1) && rc == e1.Code() {
		return true
	}
	var e2 EUnwrap
	if errors.As(err, &e2) {
		return HasCode(e2.Unwrap(), rc)
	}
	return false
}

func (e *WErr) Code() code.RetCode {
	if e == nil {
		return code.NewErrCode(code.Nil)
	}
	return code.NewErrCode(e.code)
}

func (e *WErr) SetCode(rc code.RetCode) {
	if e == nil {
		return
	}
	e.code = rc.Code()
}
