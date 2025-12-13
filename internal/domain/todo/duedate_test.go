package todo

import "testing"

func TestParseDueDate(t *testing.T) {
	d, err := ParseDueDate("2025-12-13")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if d.String() != "2025-12-13" {
		t.Fatalf("got=%s", d.String())
	}

	_, err = ParseDueDate("13-12-2025")
	if err == nil {
		t.Fatalf("expected error")
	}
}
