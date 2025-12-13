package todo

import "testing"

func TestNewTitle(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		wantErr bool
		want    string
	}{
		{"ok trimmed", "  buy milk  ", false, "buy milk"},
		{"empty", "   ", true, ""},
		{"too long", string(make([]byte, 201)), true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTitle(tt.in)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err=%v wantErr=%v", err, tt.wantErr)
			}
			if !tt.wantErr && got.String() != tt.want {
				t.Fatalf("got=%q want=%q", got.String(), tt.want)
			}
		})
	}
}
