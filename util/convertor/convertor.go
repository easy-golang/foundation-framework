package convertor

import "strconv"

func StringToInt64(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}

func Int64ToString(value int64) string {
	return strconv.FormatInt(value, 10)
}

func IntToString(value int) string {
	return strconv.Itoa(value)
}

func StringToInt32(value string) (int32, error) {
	parseInt, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(parseInt), nil
}

func Int32ToString(value int32) string {
	return strconv.FormatInt(int64(value), 10)
}

func StringToFloat64(value string) (float64, error) {
	return strconv.ParseFloat(value, 64)
}

func Float64ToString(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func StringToFloat32(value string) (float32, error) {
	float, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return 0, err
	}
	return float32(float), nil
}

func Float32ToString(value float32) string {
	return strconv.FormatFloat(float64(value), 'f', -1, 32)
}
