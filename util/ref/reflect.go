package ref

import "reflect"

func IsPointer(value any) bool {
	if value == nil {
		return false
	}
	// 使用reflect包获取类型信息
	t := reflect.TypeOf(value)
	// 判断类型是否是指针类型
	return t.Kind() == reflect.Ptr
}

func GetValueFromPointer(ptr any) any {
	// 使用reflect包获取指针的值
	v := reflect.ValueOf(ptr)
	if v.IsNil() {
		return nil
	}
	// 取得指针指向的值
	return v.Elem().Interface()
}
