package usql

import (
	"fmt"
	"testing"
)

func TestToNullInt64(t *testing.T) {
	var value *int64
	nullInt64 := ToNullInt64(value)
	fmt.Println(nullInt64)
}
