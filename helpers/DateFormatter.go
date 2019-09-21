package helpers

import (
	"time"
)

func StrToTime(str string, layout string) time.Time {
	t, _ := time.Parse(layout, str)
	return t
}