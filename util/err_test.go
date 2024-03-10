package util

import "testing"

func TestException(t *testing.T) {
	defer Recover(nil)
	Abc()
	for {

	}
}

func Abc() {
	panic("1221")
}
