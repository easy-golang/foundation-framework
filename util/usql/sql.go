package usql

import (
	"database/sql"
	"github.com/easy-golang/foundation-framework/err/comm"
	"github.com/easy-golang/foundation-framework/util/ref"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

func ToNullString(value any) sql.NullString {
	nullString := sql.NullString{}
	if ref.IsPointer(value) {
		value = ref.GetValueFromPointer(value)
	}
	err := nullString.Scan(value)
	if err != nil {
		logx.Error(comm.Wrap(err))
	}
	return nullString
}

func ToStringPointer(nullString sql.NullString) *string {
	if nullString.Valid {
		return &nullString.String
	}
	return nil
}

func ToNullTime(value any) sql.NullTime {
	nullTime := sql.NullTime{}
	if ref.IsPointer(value) {
		value = ref.GetValueFromPointer(value)
	}
	err := nullTime.Scan(value)
	if err != nil {
		logx.Error(comm.Wrap(err))
	}
	return nullTime
}

func ToTimePointer(nullTime sql.NullTime) *time.Time {
	if nullTime.Valid {
		return &nullTime.Time
	}
	return nil
}

func ToNullInt16(value any) sql.NullInt16 {
	nullInt := sql.NullInt16{}
	if ref.IsPointer(value) {
		value = ref.GetValueFromPointer(value)
	}
	err := nullInt.Scan(value)
	if err != nil {
		logx.Error(comm.Wrap(err))
	}
	return nullInt
}

func ToInt16Pointer(nullInt16 sql.NullInt16) *int16 {
	if nullInt16.Valid {
		return &nullInt16.Int16
	}
	return nil
}

func ToNullInt32(value any) sql.NullInt32 {
	nullInt := sql.NullInt32{}
	if ref.IsPointer(value) {
		value = ref.GetValueFromPointer(value)
	}
	err := nullInt.Scan(value)
	if err != nil {
		logx.Error(comm.Wrap(err))
	}
	return nullInt
}

func ToInt32Pointer(nullInt32 sql.NullInt32) *int32 {
	if nullInt32.Valid {
		return &nullInt32.Int32
	}
	return nil
}

func ToNullInt64(value any) sql.NullInt64 {
	nullInt := sql.NullInt64{}
	if ref.IsPointer(value) {
		value = ref.GetValueFromPointer(value)
	}
	err := nullInt.Scan(value)
	if err != nil {
		logx.Error(comm.Wrap(err))
	}
	return nullInt
}

func ToInt64Pointer(nullInt64 sql.NullInt64) *int64 {
	if nullInt64.Valid {
		return &nullInt64.Int64
	}
	return nil
}

func ToNullByte(value any) sql.NullByte {
	nullByte := sql.NullByte{}
	if ref.IsPointer(value) {
		value = ref.GetValueFromPointer(value)
	}
	err := nullByte.Scan(value)
	if err != nil {
		logx.Error(comm.Wrap(err))
	}
	return nullByte
}

func ToBytePointer(nullByte sql.NullByte) *byte {
	if nullByte.Valid {
		return &nullByte.Byte
	}
	return nil
}

func ToNullFloat64(value any) sql.NullFloat64 {
	nullFloat64 := sql.NullFloat64{}
	if ref.IsPointer(value) {
		value = ref.GetValueFromPointer(value)
	}
	err := nullFloat64.Scan(value)
	if err != nil {
		logx.Error(comm.Wrap(err))
	}
	return nullFloat64
}

func ToFloat64Pointer(nullFloat64 sql.NullFloat64) *float64 {
	if nullFloat64.Valid {
		return &nullFloat64.Float64
	}
	return nil
}

func ToNullBool(value any) sql.NullBool {
	nullBool := sql.NullBool{}
	if ref.IsPointer(value) {
		value = ref.GetValueFromPointer(value)
	}
	err := nullBool.Scan(value)
	if err != nil {
		logx.Error(comm.Wrap(err))
	}
	return nullBool
}

func ToBoolPointer(nullBool sql.NullBool) *bool {
	if nullBool.Valid {
		return &nullBool.Bool
	}
	return nil
}
