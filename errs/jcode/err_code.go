package jcode

import "fmt"

type errCode struct {
	code int
	msg  string
	desc interface{}
}

func (ec *errCode) Code() int {
	return ec.code
}
func (ec *errCode) Msg() string {
	return ec.msg
}

func (ec *errCode) Desc() interface{} {
	return ec.desc
}
func (ec *errCode) String() string {
	if ec.desc != nil {
		return fmt.Sprintf(`{Code: %d, Msg: %s} %v`, ec.code, ec.msg, ec.desc)
	}
	if ec.msg != "" {
		return fmt.Sprintf(`{Code: %d, Msg: %s}`, ec.code, ec.msg)
	}
	return fmt.Sprintf(`Code: %d`, ec.code)
}

func NewErrCode(code int) *errCode {
	return &errCode{
		code: code,
	}
}
func NewWithMsg(msg string) *errCode {
	return &errCode{
		code: Nil,
		msg:  msg,
	}
}
func NewWithCodeMsg(code int, msg string) *errCode {
	return &errCode{
		code: code,
		msg:  msg,
	}
}

func NewWithCode(code RetCode) *errCode {
	return &errCode{
		code: code.Code(),
		msg:  code.Msg(),
		desc: code.Desc(),
	}
}

func NewWithDesc(code RetCode, desc interface{}) *errCode {
	return &errCode{
		code: code.Code(),
		msg:  code.Msg(),
		desc: desc,
	}
}

func (ec *errCode) WithCode(code int) *errCode {
	ec.code = code
	return ec
}

func (ec *errCode) WithMsg(msg string) *errCode {
	ec.msg = msg
	return ec
}
func (ec *errCode) WithDesc(desc string) *errCode {
	ec.desc = desc
	return ec
}

func ToMsg(code int) string {
	switch code {
	case Nil:
		return "Unknown code"
	case OK:
		return "Response ok"
	case Unauthorized:
		return "Not authorized"
	case InternalErr:
		return "Internal error"
	case NotFoundErr:
		return "Data not found"
	case NotSupportErr:
		return "Not support"
	case VerifyErr:
		return "Verify error"
	case ParamErr:
		return "Param error"
	case ParamRequired:
		return "Param required"
	case LenErr:
		return "Length error"
	case OptErr:
		return "Operate error"
	case ConfigErr:
		return "Config error"
	case ForbiddenErr:
		return "Forbidden error"
	case DataExistErr:
		return "Data already exist"
	case UnknownErr:
		return "Unknown error"
	default:
		return ""
	}
}
