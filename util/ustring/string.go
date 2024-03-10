package ustring

import (
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"strings"
	"unicode"
	"unsafe"
)

func ToString(obj any) string {
	return fmt.Sprintf("%v", obj)
}

func HasEmptyString(strs []string) bool {
	for _, str := range strs {
		if IsEmptyString(str) {
			return true
		}
	}
	return false
}

func IsEmptyString(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}

func IsNotEmptyString(str string) bool {
	return len(strings.TrimSpace(str)) != 0
}

func IsJson(str string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(str), &js) == nil
}

func IsNotJson(str string) bool {
	return !IsJson(str)
}

func IsASCII(str string) bool {
	for i := 0; i < len(str); i++ {
		if str[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func ToBytes(str *string) []byte {
	if str == nil {
		return nil
	}
	return *(*[]byte)(unsafe.Pointer(str))
}

func GetUUID() string {
	return uuid.NewV4().String()
}
