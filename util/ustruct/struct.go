package ustruct

import "reflect"

func GetStructName(obj any) string {
	t := reflect.TypeOf(obj)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}
