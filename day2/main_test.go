package main

import "testing"

func Test_report_isSafe(t *testing.T) {
	tests := []struct {
		name string
		r    report
		want bool
	}{
		{"7 6 4 2 1", []int{7, 6, 4, 2, 1}, true},
		{"1 2 7 8 9", []int{1, 2, 7, 8, 9}, false},
		{"9 7 6 2 1", []int{9, 7, 6, 2, 1}, false},
		{"1 3 2 4 5", []int{1, 3, 2, 4, 5}, false},
		{"8 6 4 4 1", []int{8, 6, 4, 4, 1}, false},
		{"1 3 6 7 9", []int{1, 3, 6, 7, 9}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.isSafe(); got != tt.want {
				t.Errorf("isSafe() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_report_isAlmostSafe(t *testing.T) {
	tests := []struct {
		name string
		r    report
		want bool
	}{
		{"7 6 4 2 1", []int{7, 6, 4, 2, 1}, true},
		{"1 2 7 8 9", []int{1, 2, 7, 8, 9}, false},
		{"9 7 6 2 1", []int{9, 7, 6, 2, 1}, false},
		{"1 3 2 4 5", []int{1, 3, 2, 4, 5}, true},
		{"8 6 4 4 1", []int{8, 6, 4, 4, 1}, true},
		{"1 3 6 7 9", []int{1, 3, 6, 7, 9}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.isAlmostSafe(); got != tt.want {
				t.Errorf("isAlmostSafe() = %v, want %v", got, tt.want)
			}
		})
	}
}
