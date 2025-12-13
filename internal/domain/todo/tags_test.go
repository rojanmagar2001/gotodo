package todo

import "testing"

func TestNewTags_NormalizesAndDedupes(t *testing.T) {
	tags := NewTags([]string{" Work ", "work", "HOME", "", "  "})
	if len(tags) != 2 {
		t.Fatalf("len=%d want=2 (%v)", len(tags), tags)
	}
	if tags[0] != "home" || tags[1] != "work" {
		t.Fatalf("got=%v want=[home work]", tags)
	}
	if !tags.Contains("WORK") {
		t.Fatalf("expected Contains(WORK)=true")
	}
}
