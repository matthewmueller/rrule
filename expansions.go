package rrule

import "time"

func expandBySeconds(tt []time.Time, seconds ...int) []time.Time {
	if len(seconds) == 0 {
		return tt
	}

	e := make([]time.Time, 0, len(tt)*len(seconds))
	for _, t := range tt {
		for _, s := range seconds {
			e = append(e, time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), s, t.Nanosecond(), t.Location()))
		}
	}

	return e
}

func expandByMinutes(tt []time.Time, minutes ...int) []time.Time {
	if len(minutes) == 0 {
		return tt
	}

	e := make([]time.Time, 0, len(tt)*len(minutes))
	for _, t := range tt {
		for _, m := range minutes {
			e = append(e, time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), m, t.Second(), t.Nanosecond(), t.Location()))
		}
	}

	return e
}

func expandByHours(tt []time.Time, hours ...int) []time.Time {
	if len(hours) == 0 {
		return tt
	}

	e := make([]time.Time, 0, len(tt)*len(hours))
	for _, t := range tt {
		for _, h := range hours {
			e = append(e, time.Date(t.Year(), t.Month(), t.Day(), h, t.Minute(), t.Second(), t.Nanosecond(), t.Location()))
		}
	}

	return e
}

func expandByWeekdays(tt []time.Time, weekStart time.Weekday, weekdays ...QualifiedWeekday) []time.Time {
	if len(weekdays) == 0 {
		return tt
	}

	e := make([]time.Time, 0, len(tt)*len(weekdays))
	for _, t := range tt {
		t = backToWeekday(t, weekStart)
		for _, wd := range weekdays {
			e = append(e, forwardToWeekday(t, wd.WD))
		}
	}

	return e
}
func expandByMonthDays(tt []time.Time, monthdays ...int) []time.Time {
	if len(monthdays) == 0 {
		return tt
	}

	e := make([]time.Time, 0, len(tt)*len(monthdays))
	for _, t := range tt {
		for _, md := range monthdays {
			e = append(e, time.Date(t.Year(), t.Month(), md, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location()))
		}
	}

	return e
}

func expandByYearDays(tt []time.Time, yeardays ...int) []time.Time {
	if len(yeardays) == 0 {
		return tt
	}

	e := make([]time.Time, 0, len(tt)*len(yeardays))
	for _, t := range tt {
		yearStart := time.Date(t.Year(), time.January, 1, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())

		for _, yd := range yeardays {
			e = append(e, yearStart.AddDate(0, 0, yd))
		}
	}

	return e
}

func expandByWeekNumbers(tt []time.Time, weekStarts time.Weekday, weekNumbers ...int) []time.Time {
	if len(weekNumbers) == 0 {
		return tt
	}

	e := make([]time.Time, 0, len(tt)*len(weekNumbers))
	for _, t := range tt {
		yearStart := time.Date(t.Year(), time.January, 1, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
		yearStart = forwardToWeekday(yearStart, t.Weekday())

		for _, w := range weekNumbers {
			e = append(e, yearStart.AddDate(0, 0, (w-1)*7))
		}
	}

	return e
}

func expandByMonths(tt []time.Time, ib InvalidBehavior, months ...time.Month) []time.Time {
	if len(months) == 0 {
		return tt
	}

	e := make([]time.Time, 0, len(tt)*len(months))
	for _, t := range tt {
		for _, m := range months {
			set := time.Date(t.Year(), m, t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
			if set.Month() != m {
				switch ib {
				case PrevInvalid:
					set = time.Date(t.Year(), t.Month()+1, -1, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
					e = append(e, set)
				case NextInvalid:
					set = time.Date(t.Year(), t.Month()+1, 1, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
					e = append(e, set)
				case OmitInvalid:
					// do nothing
				}
			} else {
				e = append(e, set)
			}
		}
	}

	return e
}

// expandMonthByWeekdays does a special expansion of the month by weekdays. If
// bySetPos is not nil, it is assumed tt is the full set of instances within the
// monthly iteration, and only the instances matching the posisions of bySetPos
// are returned. This is an optimization.
func expandMonthByWeekdays(tt []time.Time, ib InvalidBehavior, bySetPos []int, weekdays ...QualifiedWeekday) []time.Time {
	if len(weekdays) == 0 {
		return tt
	}

	e := make([]time.Time, 0, len(tt))
	for _, t := range tt {
		e = append(e, weekdaysInMonth(t, weekdays, bySetPos, ib)...)
	}

	return e
}

func expandYearByWeekdays(tt []time.Time, ib InvalidBehavior, weekdays ...QualifiedWeekday) []time.Time {
	if len(weekdays) == 0 {
		return tt
	}

	e := make([]time.Time, 0, len(tt))
	for _, t := range tt {
		for _, wd := range weekdays {
			e = append(e, weekdaysInYear(t, wd, ib)...)
		}
	}

	return e

}