package util

import (

)

func CounterToInt(c interface{}) int {
	if val, ok := c.(int); ok {
		return val
	} else if val, ok := c.(float64); ok {
		return int(val)
	}
	return -1
}