package tushare

import "fmt"

type ApiError struct {
	Code   int
	Status string
}

func (e ApiError) Error() string {
	return fmt.Sprintf("[%d]: %s", e.Code, e.Status)
}
