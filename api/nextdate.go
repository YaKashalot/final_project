package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	nowParam := r.FormValue("now")
	dateParam := r.FormValue("date")
	repeatParam := r.FormValue("repeat")

	var now time.Time
	var err error

	if nowParam == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(dateFormat, nowParam)
		if err != nil {
			http.Error(w, "invalid 'now' parameter", http.StatusBadRequest)
			return
		}
	}

	nextDate, err := NextDate(now, dateParam, repeatParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(nextDate))
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("repeat rule is empty")
	}

	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("invalid start date %q: %w", dstart, err)
	}

	parts := strings.Split(repeat, " ")

	switch parts[0] {
	case "d":
		return nextDateDays(now, date, parts)
	case "y":
		return nextDateYear(now, date), nil
	case "w":
		return nextDateWeek(now, date, parts)
	case "m":
		return nextDateMonth(now, date, parts)
	default:
		return "", fmt.Errorf("unsupported repeat format: %q", repeat)
	}
}

func afterNow(date, now time.Time) bool {
	y1, m1, d1 := date.Date()
	y2, m2, d2 := now.Date()
	d := time.Date(y1, m1, d1, 0, 0, 0, 0, time.UTC)
	n := time.Date(y2, m2, d2, 0, 0, 0, 0, time.UTC)
	return d.After(n)
}

func nextDateDays(now, date time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid 'd' rule: %q", strings.Join(parts, " "))
	}
	days, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", fmt.Errorf("invalid days interval: %q", parts[1])
	}
	if days < 1 || days > 400 {
		return "", fmt.Errorf("days interval out of range (1-400): %d", days)
	}

	for {
		date = date.AddDate(0, 0, days)
		if afterNow(date, now) {
			break
		}
	}
	return date.Format(dateFormat), nil
}

func nextDateYear(now, date time.Time) string {
	for {
		date = date.AddDate(1, 0, 0)
		if afterNow(date, now) {
			break
		}
	}
	return date.Format(dateFormat)
}

func nextDateWeek(now, date time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid 'w' rule: %q", strings.Join(parts, " "))
	}

	validDays := make(map[int]bool)
	for _, ds := range strings.Split(parts[1], ",") {
		d, err := strconv.Atoi(ds)
		if err != nil || d < 1 || d > 7 {
			return "", fmt.Errorf("invalid weekday value: %q", ds)
		}
		validDays[d] = true
	}

	for {
		date = date.AddDate(0, 0, 1)
		weekday := int(date.Weekday())
		if weekday == 0 { // time.Sunday == 0, а у нас воскресенье = 7
			weekday = 7
		}
		if validDays[weekday] && afterNow(date, now) {
			break
		}
	}
	return date.Format(dateFormat), nil
}

func nextDateMonth(now, date time.Time, parts []string) (string, error) {
	if len(parts) < 2 || len(parts) > 3 {
		return "", fmt.Errorf("invalid 'm' rule: %q", strings.Join(parts, " "))
	}

	validDays := make(map[int]bool)
	for _, ds := range strings.Split(parts[1], ",") {
		d, err := strconv.Atoi(ds)
		if err != nil || d == 0 || d < -2 || d > 31 {
			return "", fmt.Errorf("invalid day of month value: %q", ds)
		}
		validDays[d] = true
	}

	var validMonths map[int]bool
	if len(parts) == 3 {
		validMonths = make(map[int]bool)
		for _, ms := range strings.Split(parts[2], ",") {
			m, err := strconv.Atoi(ms)
			if err != nil || m < 1 || m > 12 {
				return "", fmt.Errorf("invalid month value: %q", ms)
			}
			validMonths[m] = true
		}
	}

	for {
		date = date.AddDate(0, 0, 1)

		if validMonths != nil && !validMonths[int(date.Month())] {
			continue
		}

		day := date.Day()
		last := lastDayOfMonth(date)

		matched := validDays[day] ||
			(day == last && validDays[-1]) ||
			(day == last-1 && validDays[-2])

		if matched && afterNow(date, now) {
			break
		}
	}
	return date.Format(dateFormat), nil
}

func lastDayOfMonth(date time.Time) int {
	firstOfNextMonth := time.Date(date.Year(), date.Month()+1, 1, 0, 0, 0, 0, date.Location())
	return firstOfNextMonth.AddDate(0, 0, -1).Day()
}
