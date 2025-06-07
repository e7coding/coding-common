package jerr

import (
	"errors"
	"fmt"
	"github.com/e7coding/coding-common/errs/jcode"
	"strings"
)

func WithCode(rc jcode.RetCode, msg ...string) *WErr {
	return &WErr{
		msg:  fmt.Sprintf("[%s] %s", rc.Msg(), strings.Join(msg, errMsgSplit)),
		code: rc.Code(),
	}
}

func WithCodeF(rc jcode.RetCode, format string, args ...interface{}) *WErr {
	return &WErr{
		msg:  fmt.Sprintf("[%s] %s", rc.Msg(), fmt.Sprintf(format, args...)),
		code: rc.Code(),
	}
}

func WithCodeErr(rc jcode.RetCode, err error, msg ...string) *WErr {
	if err == nil {
		return nil
	}
	return &WErr{
		error: err,
		msg:   fmt.Sprintf("[%s] %s", rc.Msg(), strings.Join(msg, errMsgSplit)),
		code:  rc.Code(),
	}
}

func WithCodeErrF(rc jcode.RetCode, err error, format string, args ...interface{}) *WErr {
	if err == nil {
		return nil
	}
	return &WErr{
		error: err,
		msg:   fmt.Sprintf("[%s] %s", rc.Msg(), fmt.Sprintf(format, args...)),
		code:  rc.Code(),
	}
}

func ToRetCode(err error) jcode.RetCode {
	if err == nil {
		return jcode.NewErrCode(jcode.Nil)
	}
	var e1 EErr
	if errors.As(err, &e1) {
		return e1.Code()
	}
	var e2 EUnwrap
	if errors.As(err, &e2) {
		return ToRetCode(e2.Unwrap())
	}
	return jcode.NewErrCode(jcode.Nil)
}
func HasCode(err error, rc jcode.RetCode) bool {
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

func (e *WErr) Code() jcode.RetCode {
	if e == nil {
		return jcode.NewErrCode(jcode.Nil)
	}
	return jcode.NewErrCode(e.code)
}

func (e *WErr) SetCode(rc jcode.RetCode) {
	if e == nil {
		return
	}
	e.code = rc.Code()
}
