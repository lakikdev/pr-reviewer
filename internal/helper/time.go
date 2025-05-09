package helper

import (
	"time"

	"pr-reviewer/internal/api/utils"
)

func GetEndOfCurrentWeekPST() *time.Time {
	nowPST := utils.TimePST()
	weekday := nowPST.Weekday()
	if weekday == 0 {
		weekday = 7
	}
	startWeekDatePST := nowPST.Add(-time.Duration(weekday-1) * 24 * time.Hour)
	nowYear, nowMonth, nowDay := startWeekDatePST.Date()
	endAt := time.Date(nowYear, nowMonth, nowDay+7, 0, 0, 0, 0, utils.TimeLocationPST()).UTC()

	return &endAt
}

func GetStartOfCurrentWeekPST() *time.Time {
	nowPST := utils.TimePST()
	weekday := nowPST.Weekday()
	if weekday == 0 {
		weekday = 7
	}
	startWeekDatePST := nowPST.Add(-time.Duration(weekday-1) * 24 * time.Hour)
	nowYear, nowMonth, nowDay := startWeekDatePST.Date()
	startAt := time.Date(nowYear, nowMonth, nowDay, 0, 0, 0, 0, utils.TimeLocationPST()).UTC()

	return &startAt
}
