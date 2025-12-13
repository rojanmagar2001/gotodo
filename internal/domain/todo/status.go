package todo

type Status string

const (
	StatusActive   Status = "active"
	StatusDone     Status = "done"
	StatusArchived Status = "archived"
)

func (s Status) Valid() bool {
	switch s {
	case StatusActive, StatusDone, StatusArchived:
		return true
	default:
		return false
	}
}
