package todo

import "strings"

type Title string

const (
	minTitleLen = 1
	maxTitleLen = 200
)

func NewTitle(raw string) (Title, error) {
	v := strings.TrimSpace(raw)
	if len(v) < minTitleLen || len(v) > maxTitleLen {
		return "", ErrInvalidTitle
	}
	return Title(v), nil
}

func (t Title) String() string { return string(t) }
