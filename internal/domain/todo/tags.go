package todo

import (
	"sort"
	"strings"
)

type Tags []string

func NewTags(raw []string) Tags {
	m := make(map[string]struct{}, len(raw))
	for _, t := range raw {
		v := strings.ToLower(strings.TrimSpace(t))
		if v == "" {
			continue
		}
		m[v] = struct{}{}
	}
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out) // stable order for deterministic tests/UX
	return Tags(out)
}

func (t Tags) Contains(tag string) bool {
	needle := strings.ToLower(strings.TrimSpace(tag))
	for _, v := range t {
		if v == needle {
			return true
		}
	}
	return false
}
