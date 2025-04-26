package code

type RetCode interface {
	Code() int

	Msg() string

	Desc() interface{}
}
