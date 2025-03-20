package object

import "reflect"

// Default : 传入多个值, 返回第一次出现的有效值
func Default[T comparable](vals ...T) T {
	for _, val := range vals {
		if reflect.DeepEqual(val, reflect.Zero(reflect.TypeOf(val)).Interface()) {
			continue
		}
		return val
	}
	return vals[len(vals)-1]
}

func In[T comparable](target T, arr []T) bool {
	for _, v := range arr {
		if target == v {
			return true
		}
	}
	return false
}
