package helper

import "time"

func CurrentUTC7Date() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func UnixTimestampRangeForDay(date time.Time, add time.Duration) (int64, int64) {
	year, month, dt := date.Add(add).Date()
	startOfDay := time.Date(year, month, dt, 0, 0, 0, 0, date.Location())
	endOfDay := time.Date(year, month, dt, 23, 59, 59, 0, date.Location())
	return startOfDay.Unix(), endOfDay.Unix()
}

func UnixTimestampRangeForMonth(date time.Time, add time.Duration) (int64, int64) {
	year, month, _ := date.Add(add).Date()
	startOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, date.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Second)
	return startOfMonth.Unix(), endOfMonth.Unix()
}

func ThisMonthUTC7Range() (int64, int64) {
	return UnixTimestampRangeForMonth(time.Now(), time.Hour*7)
}

func TodayUTC7Range() (int64, int64) {
	return UnixTimestampRangeForDay(time.Now(), time.Hour*7)
}

func YesterdayUTC7Range() (int64, int64) {
	return UnixTimestampRangeForDay(time.Now().AddDate(0, 0, -1), time.Hour*7)
}

func PreviousMonthUTC7Range() (int64, int64) {
	return UnixTimestampRangeForMonth(time.Now().AddDate(0, -1, 0), time.Hour*7)
}
