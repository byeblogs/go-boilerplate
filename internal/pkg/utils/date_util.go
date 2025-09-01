package utils

import "time"

// DateUtil is the interface that defines the methods for dates utility
type DateUtil interface {
	GetTimeNowJakarta() (time.Time, error)
}

// DateUtilImpl is the implementation of interface DateUtil
type DateUtilImpl struct {}

func (du DateUtilImpl) GetTimeNowJakarta() (time.Time, error) {
	jakarta, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return time.Time{}, err
	}
	return time.Now().In(jakarta), nil
}