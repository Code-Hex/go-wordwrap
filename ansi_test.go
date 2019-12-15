package wordwrap

import (
	"testing"
)

func Test_trimAnsiColor(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{
			name: "valid 1",
			arg:  "\x1b[33m\x1b[0m",
			want: "",
		},
		{
			name: "valid 2",
			arg:  "\x1b[33mHello, World\x1b[0m",
			want: "Hello, World",
		},
		{
			name: "valid 3",
			arg:  "\x1b[38;5;82m\x1b[0m",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := trimAnsiColor(tt.arg); got != tt.want {
				t.Errorf("trimAnsiColor() = %q, want %q", got, tt.want)
			}
		})
	}
}

func Test_rangeAnsiColor(t *testing.T) {
	tests := []struct {
		name  string
		arg   string
		want  int
		want1 int
	}{
		{
			name:  "valid 1",
			arg:   "hello\x1b[38;5;82m",
			want:  5,
			want1: 15,
		},
		{
			name:  "valid 2",
			arg:   "\x1b[38mhello",
			want:  0,
			want1: 5,
		},
		{
			name:  "invalid ansi format 1",
			arg:   "hello\x1b[10",
			want:  -1,
			want1: -1,
		},
		{
			name:  "invalid ansi format 2",
			arg:   "hello\x1b[m",
			want:  -1,
			want1: -1,
		},
		{
			name:  "does not contains ansi",
			arg:   "hello",
			want:  -1,
			want1: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := rangeAnsiColor(tt.arg)
			if got != tt.want {
				t.Errorf("rangeAnsiColor() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("rangeAnsiColor() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
