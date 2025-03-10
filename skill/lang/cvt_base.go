package lang

import (
	"bytes"
	"errors"
	"strconv"
	"time"
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

var (
	errorNilValue      = errors.New("Value is nil.")
	errorConvertValue  = errors.New("Value convert error.")
	errorNumOutOfRange = errors.New("Number out of range.")
)

func ToBool(v any) (b bool, err error) {
	if v == nil {
		return false, errorNilValue
	}

	switch vt := v.(type) {
	case string:
		b, err = strconv.ParseBool(vt)
	case bool:
		b = vt
	case []byte:
		b, err = strconv.ParseBool(string(vt))
	default:
		err = errorConvertValue
	}
	return
}

func ToTime(layout string, v any) (tm *time.Time, err error) {
	if v == nil {
		return nil, errorNilValue
	}
	if layout == "" {
		layout = timeFormat
	}

	switch vt := v.(type) {
	case string:
		tm2, err2 := time.Parse(layout, vt)
		err = err2
		tm = &tm2
	case time.Time:
		tm = &vt
	case []byte:
		tm2, err2 := time.Parse(layout, string(vt))
		err = err2
		tm = &tm2
	default:
		err = errorConvertValue
	}
	return
}

func ToDuration(v any) (dr time.Duration, err error) {
	if v == nil {
		return 0, errorNilValue
	}

	switch vt := v.(type) {
	case string:
		dr, err = time.ParseDuration(vt)
	case time.Duration:
		dr = vt
	case []byte:
		dr, err = time.ParseDuration(string(vt))
	default:
		err = errorConvertValue
	}
	return
}

func Camel2Snake(s string) string {
	newS := bytes.Buffer{}
	for i := 0; i < len(s); i++ {
		if s[i] >= 65 && s[i] <= 90 {
			if i > 0 && s[i-1] >= 97 && s[i-1] <= 122 {
				newS.WriteByte('_')
			}
			newS.WriteByte(s[i] + 32)
		} else {
			newS.WriteByte(s[i])
		}
	}
	return newS.String()
}
