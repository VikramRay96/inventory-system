package dto

import (
	"strings"
)

type StatusCode string

func GetStatusDetails(code StatusCode) ErrorResponseDto {
	errStrings := strings.Split(string(code), ":")
	return ErrorResponseDto{
		StatusCode: StatusCode(errStrings[0]),
		Message:    errStrings[1],
	}
}
