package queries

import (
	"context"
	"time"

	"github.com/rojanmagar2001/gotodo/internal/application/ports"
	"github.com/rojanmagar2001/gotodo/internal/application/result"
	"github.com/rojanmagar2001/gotodo/internal/domain/todo"
)

type Stats struct {
	Repo  ports.TodoRepository
	Clock ports.Clock
}

type StatsDTO struct {
	Total    int
	Active   int
	Done     int
	Archived int
	Deleted  int

	Overdue  int
	DueToday int
	DueSoon  int // next 7 days (optional but useful)
}

func (q Stats) Execute(ctx context.Context) result.Result[StatsDTO] {
	spec := ports.ListSpec{
		IncludeDeleted: true,
		SortBy:         ports.SortByCreated,
		SortOrder:      ports.OrderAsc,
	}
	tds, err := q.Repo.List(ctx, spec)
	if err != nil {
		return result.Fail[StatsDTO](err)
	}

	now := q.Clock.Now()
	today := dateOnlyUTC(now)

	var s StatsDTO
	s.Total = len(tds)

	for _, t := range tds {
		if t.DeletedAt != nil {
			s.Deleted++
			continue
		}
		switch t.Status {
		case todo.StatusActive:
			s.Active++
		case todo.StatusDone:
			s.Done++
		case todo.StatusArchived:
			s.Archived++
		}

		if t.Status == todo.StatusActive && t.DueDate != nil {
			due := t.DueDate.AsTimeUTC()
			if due.Before(today) {
				s.Overdue++
			} else if sameDayUTC(due, today) {
				s.DueToday++
			} else if due.Before(today.Add(7 * 24 * time.Hour)) {
				s.DueSoon++
			}
		}
	}

	return result.Ok(s)
}

func dateOnlyUTC(t time.Time) time.Time {
	y, m, d := t.UTC().Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func sameDayUTC(a, b time.Time) bool {
	ay, am, ad := a.UTC().Date()
	by, bm, bd := b.UTC().Date()
	return ay == by && am == bm && ad == bd
}
