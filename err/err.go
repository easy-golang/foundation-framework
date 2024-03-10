package err

import (
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

const (
	ServerError = "抱歉，系统繁忙，请稍后重试！"
)

type IError interface {
	error
	ICause
	GetMessage() string
	GetStackTrace() string
}

type ICause interface {
	GetCause() error
}

func String(err IError) string {
	sb := strings.Builder{}
	sb.WriteString(GetFullStructPath(err))
	sb.WriteString(":")
	sb.WriteString(err.GetMessage())
	sb.WriteString(err.GetStackTrace())
	causeError, ok := err.GetCause().(IError)
	if ok {
		sb.WriteString(String(causeError))
	} else {
		if causeError != nil {
			sb.WriteString(causeError.Error())
		}
	}
	return sb.String()
}

func GetFullStructPath(o any) string {
	t := reflect.TypeOf(o)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	pkgPath := t.PkgPath()
	structName := t.Name()
	if pkgPath != "" {
		return pkgPath + "." + structName
	}
	return structName
}

func GetStackTrace(skip int) string {
	pc := callers(skip)
	frames := runtime.CallersFrames(pc)
	sb := strings.Builder{}
	sb.WriteString("\n")
	for {
		//    at D:/Program Files/Go/IdeaProject/wlj/abc/main.go (main.go:10)
		frame, more := frames.Next()
		sb.WriteString("    at ")
		sb.WriteString(frame.File)
		sb.WriteString(" ( ")
		_, file := filepath.Split(frame.File)
		sb.WriteString(file)
		sb.WriteString(":")
		sb.WriteString(strconv.Itoa(frame.Line))
		sb.WriteString(" )\n")
		if !more {
			break
		}
	}
	return sb.String()
}

func callers(skip int) []uintptr {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(4+skip, pcs[:])
	return pcs[0:n]
}
