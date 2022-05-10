package utils

import "reflect"

func IsNotNil(Object interface{}) bool {
	return !IsNilObject(Object)
}

func IsNilObject(object interface{}) bool {
	if object == nil {
		return true
	}
	value := reflect.ValueOf(object)
	kind := value.Kind()
	if kind >= reflect.Chan && kind <= reflect.Slice && value.IsNil() {
		return true
	}

	return false
}
