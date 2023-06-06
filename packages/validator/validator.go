package validator

type ArgsValidator interface {
	NumberInt64() (value int64, err error)
	NumberInt() (value int, err error)
	NumberUint64() (value uint64, err error)
	String() (value string, err error)
}
