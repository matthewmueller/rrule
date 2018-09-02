package rrule

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	rrule "github.com/teambition/rrule-go"
)

var now = time.Date(2018, 8, 25, 9, 8, 7, 6, time.UTC)

var cases = []struct {
	Name     string
	String   string
	RRule    RRule
	Dates    []string
	Terminal bool

	NoBenchmark            bool
	NoTest                 bool
	NoTeambitionComparison bool
}{
	{
		Name: "simple secondly",
		RRule: RRule{
			Frequency: Secondly,
			Count:     3,
			Dtstart:   now,
		},
		Dates:    []string{"2018-08-25T09:08:07Z", "2018-08-25T09:08:08Z", "2018-08-25T09:08:09Z"},
		Terminal: true,
	},
	{
		Name: "simple minutely",
		RRule: RRule{
			Frequency: Minutely,
			Count:     3,
			Dtstart:   now,
		},
		Dates:    []string{"2018-08-25T09:08:07Z", "2018-08-25T09:09:07Z", "2018-08-25T09:10:07Z"},
		Terminal: true,
	},

	{
		Name: "simple hourly",
		RRule: RRule{
			Frequency: Hourly,
			Count:     3,
			Dtstart:   now,
		},
		Dates:    []string{"2018-08-25T09:08:07Z", "2018-08-25T10:08:07Z", "2018-08-25T11:08:07Z"},
		Terminal: true,
	},

	{
		Name: "simple daily",
		RRule: RRule{
			Frequency: Daily,
			Count:     3,
			Dtstart:   now,
		},
		Dates:    []string{"2018-08-25T09:08:07Z", "2018-08-26T09:08:07Z", "2018-08-27T09:08:07Z"},
		Terminal: true,
	},

	{
		Name: "secondly setpos",
		RRule: RRule{
			Frequency: Secondly,
			Count:     4,
			Dtstart:   now,
			BySeconds: []int{1, 2, 3},
			ByMonths:  []time.Month{time.August, time.September},
			BySetPos:  []int{1, 3, -1},
		},
		Dates:    []string{"2018-08-25T09:09:01Z", "2018-08-25T09:09:02Z", "2018-08-25T09:09:03Z", "2018-08-25T09:10:01Z"},
		Terminal: true,
	},
	{
		Name: "minutely setpos",
		RRule: RRule{
			Frequency: Minutely,
			Count:     4,
			Dtstart:   now,
			BySeconds: []int{1, 2, 3},
			ByMonths:  []time.Month{time.August, time.September},
			BySetPos:  []int{1, 3, -1},
		},
		Dates:    []string{"2018-08-25T09:09:01Z", "2018-08-25T09:09:03Z", "2018-08-25T09:10:01Z", "2018-08-25T09:10:03Z"},
		Terminal: true,
	},

	{
		Name: "hourly setpos",
		RRule: RRule{
			Frequency: Hourly,
			Count:     4,
			Dtstart:   now,
			ByMinutes: []int{1, 2, 3},
			ByMonths:  []time.Month{time.August, time.September},
			BySetPos:  []int{1, 3, -1},
		},
		Dates:    []string{"2018-08-25T10:01:07Z", "2018-08-25T10:03:07Z", "2018-08-25T11:01:07Z", "2018-08-25T11:03:07Z"},
		Terminal: true,
	},

	{
		Name: "daily setpos",
		RRule: RRule{
			Frequency: Daily,
			Count:     4,
			Dtstart:   now,
			ByHours:   []int{1, 2, 3},
			ByMonths:  []time.Month{time.August, time.September},
			BySetPos:  []int{1, 3, -1},
		},
		Dates:    []string{"2018-08-26T01:08:07Z", "2018-08-26T03:08:07Z", "2018-08-27T01:08:07Z", "2018-08-27T03:08:07Z"},
		Terminal: true,
	},
	{
		Name: "weekly setpos",
		RRule: RRule{
			Frequency: Weekly,
			Count:     4,
			Dtstart:   now,
			ByHours:   []int{1, 2, 3},
			ByMonths:  []time.Month{time.August, time.September},
			BySetPos:  []int{1, 3, -1},
		},
		Dates:    []string{"2018-09-01T01:08:07Z", "2018-09-01T03:08:07Z", "2018-09-08T01:08:07Z", "2018-09-08T03:08:07Z"},
		Terminal: true,
	},

	{
		Name: "monthly setpos",
		RRule: RRule{
			Frequency:  Monthly,
			ByWeekdays: []QualifiedWeekday{{N: 0, WD: time.Monday}, {N: 0, WD: time.Tuesday}, {N: 0, WD: time.Wednesday}, {N: 0, WD: time.Thursday}, {N: 0, WD: time.Friday}, {N: 0, WD: time.Saturday}, {N: 0, WD: time.Sunday}},
			Count:      4,
			Dtstart:    now,
			ByMonths:   []time.Month{time.August, time.September},
			BySetPos:   []int{1, 3, -1},
		},
		Dates:    []string{"2018-08-31T09:08:07Z", "2018-09-01T09:08:07Z", "2018-09-03T09:08:07Z", "2018-09-30T09:08:07Z"},
		Terminal: true,
	},

	{
		Name: "yearly setpos",
		RRule: RRule{
			Frequency:  Yearly,
			ByWeekdays: []QualifiedWeekday{{N: 0, WD: time.Monday}, {N: 0, WD: time.Tuesday}, {N: 0, WD: time.Wednesday}, {N: 0, WD: time.Thursday}, {N: 0, WD: time.Friday}, {N: 0, WD: time.Saturday}, {N: 0, WD: time.Sunday}},
			Count:      4,
			Dtstart:    now,
			ByMonths:   []time.Month{time.August, time.September},
			BySetPos:   []int{1, 3, -1},
		},
		Dates:    []string{"2018-09-30T09:08:07Z", "2019-08-01T09:08:07Z", "2019-08-03T09:08:07Z", "2019-09-30T09:08:07Z"},
		Terminal: true,
	},

	{
		Name: "daily until",
		RRule: RRule{
			Frequency: Daily,
			Until:     time.Date(2018, 8, 30, 0, 0, 0, 0, time.UTC),
			Dtstart:   now,
		},
		Dates:    []string{"2018-08-25T09:08:07Z", "2018-08-26T09:08:07Z", "2018-08-27T09:08:07Z", "2018-08-28T09:08:07Z", "2018-08-29T09:08:07Z"},
		Terminal: true,
	},

	{
		Name: "simple monthly",
		RRule: RRule{
			Frequency: Monthly,
			Count:     3,
			Dtstart:   now,
		},
		Dates:    []string{"2018-08-25T09:08:07Z", "2018-09-25T09:08:07Z", "2018-10-25T09:08:07Z"},
		Terminal: true,
	},

	{
		Name: "long monthly",
		RRule: RRule{
			Frequency: Monthly,
			Count:     300,
			Dtstart:   now,
		},
		Terminal: true,
		NoTest:   true,
	},

	{
		Name: "monthly by weekday",
		RRule: RRule{
			Frequency:  Monthly,
			Count:      3,
			Dtstart:    now,
			ByWeekdays: []QualifiedWeekday{{N: 1, WD: time.Tuesday}},
		},
		Dates:    []string{"2018-09-04T09:08:07Z", "2018-10-02T09:08:07Z", "2018-11-06T09:08:07Z"},
		Terminal: true,
	},

	{
		Name: "simple weekly",
		RRule: RRule{
			Frequency: Weekly,
			Count:     3,
			Dtstart:   now,
		},
		Dates:    []string{"2018-08-25T09:08:07Z", "2018-09-01T09:08:07Z", "2018-09-08T09:08:07Z"},
		Terminal: true,
	},

	{
		Name: "weekly by weekday",
		RRule: RRule{
			Frequency:  Weekly,
			Count:      3,
			Dtstart:    now,
			ByWeekdays: []QualifiedWeekday{{WD: time.Tuesday}},
		},
		Dates:    []string{"2018-08-28T09:08:07Z", "2018-09-04T09:08:07Z", "2018-09-11T09:08:07Z"},
		Terminal: true,
	},

	{
		Name:   "yearly by weekday",
		String: "FREQ=YEARLY;COUNT=3;BYDAY=TU,+35WE,-17MO",
		RRule: RRule{
			Frequency:  Yearly,
			Count:      4,
			Dtstart:    now,
			ByWeekdays: []QualifiedWeekday{{WD: time.Tuesday}, {N: 35, WD: time.Wednesday}, {N: -17, WD: time.Monday}},
		},
		Dates:    []string{"2018-08-28T09:08:07Z", "2018-08-29T09:08:07Z", "2018-09-04T09:08:07Z", "2018-09-10T09:08:07Z"},
		Terminal: true,

		// I'm not sure if I'm reading the spec wrong or if they have a bug, but they return no results.
		// lib-recur agrees with my implementation.
		NoTeambitionComparison: true,
	},
}

func TestRRule(t *testing.T) {
	for _, tc := range cases {
		if tc.NoTest {
			continue
		}

		t.Run(tc.Name, func(t *testing.T) {
			dates := tc.RRule.All(0)
			assert.Equal(t, tc.Dates, rfcAll(dates))
		})
	}
}

// TestAgainstTeambition checks that our test case expectations match against
// an existing RRULE library.
func TestAgainstTeambition(t *testing.T) {
	for _, tc := range cases {
		if tc.NoTest {
			continue
		}
		if tc.NoTeambitionComparison {
			continue
		}

		t.Run(tc.Name, func(t *testing.T) {
			ro := rruleToROption(tc.RRule)
			teambitionRRule, err := rrule.NewRRule(ro)
			require.NoError(t, err)

			dates := teambitionRRule.All()
			assert.Equal(t, tc.Dates, rfcAll(dates))
			t.Log(ro.String())
		})
	}
}

func BenchmarkRRule(b *testing.B) {
	for _, tc := range cases {
		if tc.NoBenchmark {
			continue
		}

		b.Run(tc.Name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				tc.RRule.All(0)
			}
		})
	}
}

func rruleToROption(r RRule) rrule.ROption {
	converted := rrule.ROption{
		Dtstart: r.Dtstart,

		Until:    r.Until,
		Count:    int(r.Count),
		Interval: r.Interval,

		Bysecond:   r.BySeconds,
		Byminute:   r.ByMinutes,
		Byhour:     r.ByHours,
		Bymonthday: r.ByMonthDays,
		Byweekno:   r.ByWeekNumbers,
		Byyearday:  r.ByYearDays,
		Bysetpos:   r.BySetPos,

		Bymonth:   make([]int, 0, len(r.ByMonths)),
		Byweekday: make([]rrule.Weekday, 0, len(r.ByWeekdays)),
	}

	switch r.Frequency {
	case Secondly:
		converted.Freq = rrule.SECONDLY
	case Minutely:
		converted.Freq = rrule.MINUTELY
	case Hourly:
		converted.Freq = rrule.HOURLY
	case Daily:
		converted.Freq = rrule.DAILY
	case Weekly:
		converted.Freq = rrule.WEEKLY
	case Monthly:
		converted.Freq = rrule.MONTHLY
	case Yearly:
		converted.Freq = rrule.YEARLY
	}

	for _, m := range r.ByMonths {
		converted.Bymonth = append(converted.Bymonth, int(m))
	}
	for _, wd := range r.ByWeekdays {
		switch wd.WD {
		case time.Sunday:
			converted.Byweekday = append(converted.Byweekday, rrule.SU.Nth(wd.N))
		case time.Monday:
			converted.Byweekday = append(converted.Byweekday, rrule.MO.Nth(wd.N))
		case time.Tuesday:
			converted.Byweekday = append(converted.Byweekday, rrule.TU.Nth(wd.N))
		case time.Wednesday:
			converted.Byweekday = append(converted.Byweekday, rrule.WE.Nth(wd.N))
		case time.Thursday:
			converted.Byweekday = append(converted.Byweekday, rrule.TH.Nth(wd.N))
		case time.Friday:
			converted.Byweekday = append(converted.Byweekday, rrule.FR.Nth(wd.N))
		case time.Saturday:
			converted.Byweekday = append(converted.Byweekday, rrule.SA.Nth(wd.N))
		}
	}

	return converted
}

func BenchmarkTeambition(b *testing.B) {
	for _, tc := range cases {

		ro := rruleToROption(tc.RRule)
		if tc.NoBenchmark {
			continue
		}
		if tc.NoTeambitionComparison {
			continue
		}

		b.Run(tc.Name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				teambitionRRule, _ := rrule.NewRRule(ro)
				teambitionRRule.All()
			}
		})
	}

}

func rfcAll(times []time.Time) []string {
	strs := make([]string, len(times))
	for i, t := range times {
		strs[i] = t.Format(time.RFC3339)
	}
	return strs
}
