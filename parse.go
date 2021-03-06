package rrule

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

// ParseRecurrence parses a whole recurrence from an iCalendar object. iCalendar
// properties recognized are DTSTART, RRULE, EXRULE, RDATE, EXDATE. Others are
// ignored.
//
// loc defines what "local" means to the parsed rules. Some patterns may
// specify a "floating" time, one without a timezone or offset, which matches
// a different actual time in different timezones. For example,
//
//	RRULE:FREQ=YEARLY;BYSECOND=-10
//
// might be useful to alert you to start counting down to the new year, no
// matter your timezone. A rule with a specific timezone, however,
//
//	DTSTART;TZID=America/New_York:19991231T000000
//	RRULE:FREQ=YEARLY;BYSECOND=-10
//
// would track a specific timezone and ignore loc. This example would alert you
// to count down to the ball dropping in New York's Times Square for each new year.
//
// If nil, time.UTC will be used.
func ParseRecurrence(src []byte, loc *time.Location) (*Recurrence, error) {
	scanner := bufio.NewScanner(bytes.NewBuffer(src))

	recurrence := &Recurrence{}

	for scanner.Scan() {
		text := scanner.Text()
		colonIdx := strings.IndexAny(text, ":;")

		if colonIdx < 0 || len(text)-1 == colonIdx {
			return nil, fmt.Errorf("misformatted line %q", text)
		}

		propName := text[:colonIdx]
		propVal := text[colonIdx+1:]

		switch propName {
		case "DTSTART":
			t, floating, err := parseTime(text, loc)
			if err != nil {
				return nil, err
			}
			recurrence.Dtstart = t
			recurrence.FloatingLocation = floating

		case "RRULE":
			rrule, err := ParseRRule(propVal)
			if err != nil {
				return nil, err
			}
			recurrence.RRules = append(recurrence.RRules, rrule)
		case "EXRULE":
			rrule, err := ParseRRule(propVal)
			if err != nil {
				return nil, err
			}
			recurrence.ExRules = append(recurrence.ExRules, rrule)
		case "RDATE":
			t, _, err := parseTime(propVal, loc)
			if err != nil {
				return nil, err
			}

			recurrence.RDates = append(recurrence.RDates, t)
		case "EXDATE":
			t, _, err := parseTime(propVal, loc)
			if err != nil {
				return nil, err
			}

			recurrence.ExDates = append(recurrence.ExDates, t)
		}
	}

	recurrence.setDtstart()

	return recurrence, nil
}

// ParseRRule parses a single RRule pattern.
func ParseRRule(str string) (RRule, error) {
	scanner := bufio.NewScanner(bytes.NewBufferString(str))
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		if idx := bytes.Index(data, []byte{';'}); idx >= 0 {
			return idx + 1, data[:idx], nil
		}

		if atEOF {
			return len(data), data, nil
		}

		return 0, nil, nil
	})

	rrule := RRule{}

	for scanner.Scan() {
		wholeComponent := scanner.Text()
		parts := strings.SplitN(wholeComponent, "=", 2)
		if len(parts) < 2 {
			return rrule, fmt.Errorf("rrule segment %q is invalid", scanner.Text())
		}

		directive, value := parts[0], parts[1]

		switch strings.ToUpper(directive) {
		case "FREQ":
			freq, err := strToFreq(value)
			if err != nil {
				return rrule, err
			}
			rrule.Frequency = freq
		case "UNTIL":
			t, floating, err := parseTime(wholeComponent, nil)
			if err != nil {
				return rrule, err
			}
			rrule.Until = t
			rrule.UntilFloating = floating

		case "COUNT":
			i, err := strconv.Atoi(value)
			if err != nil {
				return rrule, err
			}
			rrule.Count = uint64(i)
		case "INTERVAL":
			i, err := strconv.Atoi(value)
			if err != nil {
				return rrule, err
			}
			rrule.Interval = i
		case "BYSECOND":
			ints, err := parseInts(value)
			if err != nil {
				return rrule, err
			}
			rrule.BySeconds = ints
		case "BYMINUTE":
			ints, err := parseInts(value)
			if err != nil {
				return rrule, err
			}
			rrule.ByMinutes = ints
		case "BYHOUR":
			ints, err := parseInts(value)
			if err != nil {
				return rrule, err
			}
			rrule.ByHours = ints
		case "BYDAY":
			wds, err := parseQualifiedWeekdays(value)
			if err != nil {
				return rrule, err
			}
			rrule.ByWeekdays = wds
		case "BYMONTHDAY":
			ints, err := parseInts(value)
			if err != nil {
				return rrule, err
			}
			rrule.ByMonthDays = ints
		case "BYYEARDAY":
			ints, err := parseInts(value)
			if err != nil {
				return rrule, err
			}
			rrule.ByYearDays = ints
		case "BYWEEKNO":
			ints, err := parseInts(value)
			if err != nil {
				return rrule, err
			}
			rrule.ByWeekNumbers = ints
		case "BYMONTH":
			months, err := parseMonths(value)
			if err != nil {
				return rrule, err
			}
			rrule.ByMonths = months
		case "BYSETPOS":
			ints, err := parseInts(value)
			if err != nil {
				return rrule, err
			}
			rrule.BySetPos = ints
		case "WKST":
			wd, err := parseWeekday(value)
			if err != nil {
				return rrule, err
			}
			rrule.WeekStart = &wd
		default:
			return rrule, fmt.Errorf("%q is not a supported RRULE part", directive)
		}
	}

	err := rrule.Validate()
	return rrule, err
}

func parseInts(str string) ([]int, error) {
	if len(str) == 0 {
		return nil, nil
	}
	var err error
	parts := strings.Split(str, ",")
	ints := make([]int, len(parts))
	for i, p := range parts {
		ints[i], err = strconv.Atoi(p)
		if err != nil {
			return nil, err
		}
	}

	return ints, nil
}

func parseQualifiedWeekdays(str string) ([]QualifiedWeekday, error) {
	var err error
	parts := strings.Split(str, ",")
	wds := make([]QualifiedWeekday, len(parts))
	for i, p := range parts {
		if len(p) == 0 {
			return nil, errors.New("cannot have empty weekday segment in a comma-separated list")
		}

		idx := 0

		switch p[0] {
		case '-', '+':
			idx++
		}

		for _, r := range p[idx:] {
			if !unicode.IsDigit(r) {
				break
			}
			idx += utf8.RuneLen(r)
		}

		var digit int
		if idx > 0 {
			digit, err = strconv.Atoi(p[:idx])
			if err != nil {
				return nil, err
			}
		}

		wd, err := parseWeekday(p[idx:])
		if err != nil {
			return nil, err
		}

		wds[i] = QualifiedWeekday{N: digit, WD: wd}
	}

	return wds, nil
}

func parseWeekday(str string) (time.Weekday, error) {
	switch strings.ToLower(str) {
	case "mo":
		return time.Monday, nil
	case "tu":
		return time.Tuesday, nil
	case "we":
		return time.Wednesday, nil
	case "th":
		return time.Thursday, nil
	case "fr":
		return time.Friday, nil
	case "sa":
		return time.Saturday, nil
	case "su":
		return time.Sunday, nil
	default:
		return time.Sunday, fmt.Errorf("invalid day of week %q", str)
	}
}

func parseMonths(str string) ([]time.Month, error) {
	parts := strings.Split(str, ",")
	months := make([]time.Month, len(parts))
	for i, p := range parts {
		parsedInt, err := strconv.Atoi(p)
		if err != nil {
			return nil, err
		}
		months[i] = time.Month(parsedInt)
	}

	return months, nil
}

func strToFreq(str string) (Frequency, error) {
	switch strings.ToLower(str) {
	case "secondly":
		return Secondly, nil
	case "minutely":
		return Minutely, nil
	case "hourly":
		return Hourly, nil
	case "daily":
		return Daily, nil
	case "weekly":
		return Weekly, nil
	case "monthly":
		return Monthly, nil
	case "yearly":
		return Yearly, nil
	default:
		return Yearly, fmt.Errorf("frequency %q is not valid", str)
	}
}
