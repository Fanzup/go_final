package date

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("repeat is empty")
	}

	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", err
	}

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", errors.New("repeat is empty")
	}
	switch parts[0] {
	case "d":
		return nextDaily(now, date, parts)
	case "w":
		return nextWeekly(now, date, parts)
	case "m":
		return nextMonthly(now, date, parts)
	case "y":
		return nextYearly(now, date, parts)
	default:
		return "", fmt.Errorf("invalid repeat rule: %s", repeat)
	}
}

func afterNow(date, now time.Time) bool {
	dy, dm, dd := date.Date()
	ny, nm, nd := now.Date()
	d := time.Date(dy, dm, dd, 0, 0, 0, 0, time.UTC)
	n := time.Date(ny, nm, nd, 0, 0, 0, 0, time.UTC)
	return d.After(n)
}

func lastDayOfMonth(t time.Time) int {
	firstNext := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location())
	return firstNext.AddDate(0, 0, -1).Day()
}

func nextDaily(now, date time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", errors.New("d: interval required")
	}

	num, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", err
	}
	if num < 1 || num > 400 {
		return "", fmt.Errorf("d - interval out of range [1:400]")
	}
	// сразу считаем число интервалов, без перебора дней
	today := calendarDay(now)
	steps := int64(1)
	if !date.After(today) {
		steps = (today.Unix()-date.Unix())/86400/int64(num) + 1
	}
	return formatDate(date.AddDate(0, 0, int(steps)*num))
}

func nextWeekly(now, date time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", errors.New("w - need more arguments!")
	}

	var day [8]bool
	for _, str := range strings.Split(parts[1], ",") {
		num, err := strconv.Atoi(str)
		if err != nil || num < 1 || num > 7 {
			return "", errors.New("w - incorrect weekday!")
		}
		day[num] = true
	}

	if afterNow(now, date) {
		date = calendarDay(now)
	}
	for {
		date = date.AddDate(0, 0, 1)
		wday := int(date.Weekday())
		if wday == 0 {
			wday = 7
		}
		if day[wday] && afterNow(date, now) {
			return formatDate(date)
		}
	}
}

func nextMonthly(now, date time.Time, parts []string) (string, error) {
	if len(parts) < 2 || len(parts) > 3 {
		return "", errors.New("m - days required")
	}

	var day [32]bool
	var lastDay, prevLastDay bool
	for _, s := range strings.Split(parts[1], ",") {
		n, err := strconv.Atoi(s)
		if err != nil {
			return "", fmt.Errorf("m: invalid day %q", s)
		}
		switch {
		case n >= 1 && n <= 31:
			day[n] = true
		case n == -1:
			lastDay = true
		case n == -2:
			prevLastDay = true
		default:
			return "", fmt.Errorf("m: day out of range: %d", n)
		}
	}

	var month [13]bool
	hasMonths := len(parts) == 3
	if hasMonths {
		for _, s := range strings.Split(parts[2], ",") {
			n, err := strconv.Atoi(s)
			if err != nil || n < 1 || n > 12 {
				return "", fmt.Errorf("m: invalid month %q", s)
			}
			month[n] = true
		}
	}

	// проверяем, есть ли такой день хотя бы в одном выбранном месяце
	possible := lastDay || prevLastDay
	for m := 1; m <= 12 && !possible; m++ {
		if hasMonths && !month[m] {
			continue
		}
		// високосный год учитывает допустимое 29 февраля
		last := lastDayOfMonth(time.Date(2000, time.Month(m), 1, 0, 0, 0, 0, time.UTC))
		for d := 1; d <= last; d++ {
			if day[d] {
				possible = true
				break
			}
		}
	}
	if !possible {
		return "", errors.New("m: no matching calendar date")
	}

	// календарь повторяется через 400 лет — дальше искать бессмысленно
	if afterNow(now, date) {
		date = calendarDay(now)
	}
	start := date.AddDate(0, 0, 1)
	cursor := time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 400*12; i++ {
		if !hasMonths || month[int(cursor.Month())] {
			last := lastDayOfMonth(cursor)
			for d := 1; d <= last; d++ {
				if !day[d] && !(lastDay && d == last) && !(prevLastDay && d == last-1) {
					continue
				}
				candidate := time.Date(cursor.Year(), cursor.Month(), d, 0, 0, 0, 0, time.UTC)
				if candidate.Before(start) {
					continue
				}
				return formatDate(candidate)
			}
		}
		cursor = cursor.AddDate(0, 1, 0)
	}
	return "", errors.New("m: no matching calendar date")
}

func nextYearly(now, date time.Time, parts []string) (string, error) {
	if len(parts) != 1 {
		return "", errors.New("y: unexpected arguments")
	}
	for {
		date = date.AddDate(1, 0, 0)
		if afterNow(date, now) {
			break
		}
	}

	return formatDate(date)
}

func calendarDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func formatDate(t time.Time) (string, error) {
	if t.Year() < 0 || t.Year() > 9999 {
		return "", errors.New("next date is outside the supported year range")
	}
	return t.Format(dateFormat), nil
}
