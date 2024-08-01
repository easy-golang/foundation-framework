package umap

import (
	"github.com/easy-golang/foundation-framework/constant/symbol"
	"github.com/easy-golang/foundation-framework/util/ustring"
	"reflect"
	"strings"
)

func IsMap(value any) bool {
	valueType := reflect.TypeOf(value)
	for valueType.Kind() == reflect.Pointer {
		valueType = valueType.Elem()
	}
	return valueType.Kind() == reflect.Map
}

func GetStringValue(value map[string]any, key string) string {
	if value == nil || value[key] == nil {
		return symbol.EmptyString
	}
	return ustring.ToString(value[key])
}

func GetValueByKeyPath(value map[string]any, keyPath string) any {
	keys := strings.Split(keyPath, ".")
	for index, key := range keys {
		if value == nil {
			return nil
		}
		if index == len(keys)-1 {
			return value[key]
		} else {
			childValue, flag := value[key].(map[string]any)
			if flag {
				value = childValue
				continue
			}
			break
		}
	}
	return nil
}

func MergeMap[T any](mapArray ...map[string]T) map[string]T {
	mergedMap := make(map[string]T)
	for _, m := range mapArray {
		for key, value := range m {
			mergedMap[key] = value
		}
	}
	return mergedMap
}
