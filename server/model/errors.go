package model

import (
	"fmt"
	"net/http"
)

const (
	DDCOK                 = 0
	DDCErrInternal        = 1
	DDCErrInvalidArgument = 2
	DDCErrNotRenderable   = 3
)

// ErrMap code for standard argus
var ErrMap = map[int]string{
	DDCOK:                 "OK",
	DDCErrInvalidArgument: "INVALID_ARGUMENT",
	DDCErrInternal:        "INTERNAL",
	DDCErrNotRenderable:   "NOT_RENDERABLE",
}

// HTTPCodeMap ErrMap code for standard argus
var HTTPCodeMap = map[int]int{
	DDCOK:                 http.StatusOK,
	DDCErrInternal:        http.StatusInternalServerError,
	DDCErrInvalidArgument: http.StatusBadRequest,
	DDCErrNotRenderable:   http.StatusInternalServerError,
}

// ErrMsg 返回错误码对应的错误信息
func ErrMsg(code int) string {
	if msg, ok := ErrMap[code]; ok {
		return msg
	}
	return "unknown"
}

func ErrPrintf(code int, format string, args ...interface{}) string {
	if code == DDCOK {
		return ""
	}
	return fmt.Sprintf("%v: %v", ErrMsg(code), fmt.Sprintf(format, args...))
}

func ErrPrint(code int, args ...interface{}) string {
	if code == DDCOK {
		return ""
	}
	return fmt.Sprintf("%v: %v", ErrMsg(code), fmt.Sprint(args...))
}

// HTTPCode ...
func HTTPCode(code int) int {
	if msg, ok := HTTPCodeMap[code]; ok {
		return msg
	}
	return 500
}
