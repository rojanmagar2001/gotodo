package todo

import "time"

type DueDate struct {
	year  int
	month time.Month
	day   int
}

func ParseDueDate(iso string) (DueDate, error) {
	// Expect YYYY-MM-DD strictly
	t, err := time.Parse("2006-01-02", iso)
	if err != nil {
		return DueDate{}, ErrInvalidDueDate
	}
	y, m, d := t.Date()
	return DueDate{year: y, month: m, day: d}, nil
}

func (d DueDate) String() string {
	return time.Date(d.year, d.month, d.day, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
}

func (d DueDate) AsTimeUTC() time.Time {
	return time.Date(d.year, d.month, d.day, 0, 0, 0, 0, time.UTC)
}

func (d DueDate) IsBefore(other DueDate) bool {
	return d.AsTimeUTC().Before(other.AsTimeUTC())
}
