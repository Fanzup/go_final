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
		return "", errors.New("Repeat is empty!")
	}

	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", err
	}

	parts := strings.Fields(repeat)
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
		return "", fmt.Errorf("Incorrect format data: %s", repeat)
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
	for {
		date = date.AddDate(0, 0, num)
		if afterNow(date, now) {
			break
		}
	}

	return date.Format(dateFormat), nil
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

	for {
		date = date.AddDate(0, 0, 1)
		wday := int(date.Weekday())
		if wday == 0 {
			wday = 7
		}
		if day[wday] && afterNow(date, now) {
			return date.Format(dateFormat), nil
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

	for {
		date = date.AddDate(0, 0, 1)
		if !afterNow(date, now) {
			continue
		}
		if hasMonths && !month[int(date.Month())] {
			continue
		}

		d := date.Day()
		last := lastDayOfMonth(date)

		ok := day[d] ||
			(lastDay && d == last) ||
			(prevLastDay && d == last-1)
		if ok {
			return date.Format(dateFormat), nil
		}
	}
}

func nextYearly(now, date time.Time, parts []string) (string, error) {
	if len(parts) != 1 {
		return "", errors.New("y - need more arguments!")
	}
	for {
		date = date.AddDate(1, 0, 0)
		if afterNow(date, now) {
			break
		}
	}

	return date.Format(dateFormat), nil
}
