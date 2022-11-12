package common

import "errors"

var (
	FatalError        = errors.New("fatal error")
	ProjectTotalError = errors.New("get project total error")
	MissVersionError  = errors.New("miss version in swagger")
)

var ErrMaps = map[error]int{
	FatalError: 1000,
}

func GetCodeByErr(err error) int {
	return ErrMaps[err]
}
