package main

import "testing"

func Test_check(t *testing.T) {
	type args struct {
		expectedValue int
		head          int
		tail          []int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"10: 2 x 5", args{10, 2, []int{5}}, true},
		{"10! 2 3", args{10, 2, []int{3}}, false},
		{"7290: 6 * 8 || 6 * 15", args{7290, 6, []int{8, 6, 15}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := check(tt.args.expectedValue, tt.args.head, tt.args.tail); got != tt.want {
				t.Errorf("check() = %v, want %v", got, tt.want)
			}
		})
	}
}
