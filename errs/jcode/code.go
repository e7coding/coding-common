package jcode

type RetCode interface {
	Code() int

	Msg() string

	Desc() interface{}
}
