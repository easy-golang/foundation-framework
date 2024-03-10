package collection

import "github.com/emirpasic/gods/sets/hashset"

func IsNotEmptySet(value *hashset.Set) bool {
	if value != nil && value.Size() != 0 {
		return true
	}
	return false
}

func IsEmptySet(value *hashset.Set) bool {
	return !IsNotEmptySet(value)
}

func IsNotEmptySlice[T any](value []T) bool {
	if value != nil && len(value) != 0 {
		return true
	}
	return false
}

func IsEmptySlice[T any](value []T) bool {
	return !IsNotEmptySlice(value)
}

func ChanTransformToSlice[T any](chanValue chan T) []T {
	result := make([]T, 0)
	for value := range chanValue {
		result = append(result, value)
	}
	return result
}

func SplitArray[T any](originalArray []T, chunkSize int) [][]T {
	if originalArray == nil {
		return nil
	}
	numChunks := (len(originalArray) + chunkSize - 1) / chunkSize
	splitArray := make([][]T, numChunks)
	// 循环遍历原始数组，将元素分配到子数组中
	for i := 0; i < len(originalArray); i += chunkSize {
		end := i + chunkSize
		if end > len(originalArray) {
			end = len(originalArray)
		}
		splitArray[i/chunkSize] = originalArray[i:end]
	}
	return splitArray
}
