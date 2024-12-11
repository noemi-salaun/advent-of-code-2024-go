package main

import (
	"reflect"
	"testing"
)

func Test_getAntinodes(t *testing.T) {
	type args struct {
		t1         vector2
		t2         vector2
		boundaries vector2
	}
	tests := []struct {
		name string
		args args
		want []vector2
	}{
		{
			"",
			args{vector2{1, 2}, vector2{2, 4}, vector2{9, 9}},
			[]vector2{{1, 2}, {2, 4}, {3, 6}, {4, 8}, {0, 0}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getAntinodes(tt.args.t1, tt.args.t2, tt.args.boundaries); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getAntinodes() = %v, want %v", got, tt.want)
			}
		})
	}
}
