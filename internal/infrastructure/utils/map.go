package utils

import (
	"sync"
)

func IsEmpty(val *sync.Map) bool {
	isEmpty := true
	val.Range(func(_, _ interface{}) bool {
		isEmpty = false

		return false
	})

	return isEmpty
}
