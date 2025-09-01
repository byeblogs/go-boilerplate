package utils

type Util struct {
	Date	DateUtil
}

func NewUtil() *Util {
	return &Util{
		Date: DateUtilImpl{},
	}
}