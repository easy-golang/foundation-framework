package util

const (
	NORMAL = "2006-01-02 15:04:05"
	// ISO_NORMAL = "2006-01-02T03:04:05Z"  12小时制
	ISO_NORMAL = "2006-01-02T15:04:05Z"
	DATE       = "2006-01-02"
	TIME       = "15:04:05"
)

type dateUtil struct {
}

func NewDateUtil() *dateUtil {
	return &dateUtil{}
}

/*func ParseTimeToString(time time.Time, format string) (string, error) {
	time.Format(format)
}*/
