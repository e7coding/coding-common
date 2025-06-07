package jcode

const (
	Nil          = 0
	OK           = 2000
	Unauthorized = 4001
)

const (
	InternalErr = 5000 + iota*1
	NotFoundErr
	NotSupportErr
	VerifyErr
	ParamErr
	ParamRequired
	LenErr
	OptErr
	ConfigErr
	ForbiddenErr
	DataExistErr
	UnknownErr
)
