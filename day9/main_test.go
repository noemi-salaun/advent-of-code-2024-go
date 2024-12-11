package main

import (
	"reflect"
	"testing"
)

func Test_expandDiskMap(t *testing.T) {
	type args struct {
		dm []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			"12345",
			args{[]int{1, 2, 3, 4, 5}},
			[]int{0, -1, -1, 1, 1, 1, -1, -1, -1, -1, 2, 2, 2, 2, 2},
		},
		{
			"2333133121414131402",
			args{[]int{2, 3, 3, 3, 1, 3, 3, 1, 2, 1, 4, 1, 4, 1, 3, 1, 4, 0, 2}},
			[]int{0, 0, -1, -1, -1, 1, 1, 1, -1, -1, -1, 2, -1, -1, -1, 3, 3, 3, -1, 4, 4, -1, 5, 5, 5, 5, -1, 6, 6, 6, 6, -1, 7, 7, 7, -1, 8, 8, 8, 8, 9, 9},
		},
		{
			"90909",
			args{[]int{9, 0, 9, 0, 9}},
			[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 2},
		},
		{
			"1111111111111111111111111",
			args{[]int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}},
			[]int{0, -1, 1, -1, 2, -1, 3, -1, 4, -1, 5, -1, 6, -1, 7, -1, 8, -1, 9, -1, 10, -1, 11, -1, 12},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := expandDiskMap(tt.args.dm); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("expandDiskMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defrag(t *testing.T) {
	type args struct {
		input []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{"20", args{[]int{0, 0}}, []int{0, 0}},
		{"22", args{[]int{0, 0, -1, -1}}, []int{0, 0, -1, -1}},
		{
			"12345",
			args{[]int{0, -1, -1, 1, 1, 1, -1, -1, -1, -1, 2, 2, 2, 2, 2}},
			[]int{0, 2, 2, 1, 1, 1, 2, 2, 2, -1, -1, -1, -1, -1, -1},
		},
		{
			"2333133121414131402",
			args{[]int{0, 0, -1, -1, -1, 1, 1, 1, -1, -1, -1, 2, -1, -1, -1, 3, 3, 3, -1, 4, 4, -1, 5, 5, 5, 5, -1, 6, 6, 6, 6, -1, 7, 7, 7, -1, 8, 8, 8, 8, 9, 9}},
			[]int{0, 0, 9, 9, 8, 1, 1, 1, 8, 8, 8, 2, 7, 7, 7, 3, 3, 3, 6, 4, 4, 6, 5, 5, 5, 5, 6, 6, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := defrag(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("defrag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_compact(t *testing.T) {
	type args struct {
		input []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			"2333133121414131402",
			args{[]int{0, 0, -1, -1, -1, 1, 1, 1, -1, -1, -1, 2, -1, -1, -1, 3, 3, 3, -1, 4, 4, -1, 5, 5, 5, 5, -1, 6, 6, 6, 6, -1, 7, 7, 7, -1, 8, 8, 8, 8, 9, 9}},
			[]int{0, 0, 9, 9, 2, 1, 1, 1, 7, 7, 7, -1, 4, 4, -1, 3, 3, 3, -1, -1, -1, -1, 5, 5, 5, 5, -1, 6, 6, 6, 6, -1, -1, -1, -1, -1, 8, 8, 8, 8, -1, -1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compact(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("compact() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checksum(t *testing.T) {
	type args struct {
		input []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"2333133121414131402",
			args{[]int{0, 0, 9, 9, 8, 1, 1, 1, 8, 8, 8, 2, 7, 7, 7, 3, 3, 3, 6, 4, 4, 6, 5, 5, 5, 5, 6, 6, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}},
			1928,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checksum(tt.args.input); got != tt.want {
				t.Errorf("checksum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findGroup(t *testing.T) {
	type args struct {
		input   []int
		id      int
		startAt int
	}
	tests := []struct {
		name string
		args args
		want groupPos
	}{
		{
			"0 in 12345",
			args{[]int{0, -1, -1, 1, 1, 1, -1, -1, -1, -1, 2, 2, 2, 2, 2}, 0, 0},
			groupPos{0, 1},
		},
		{
			"1 in 12345",
			args{[]int{0, -1, -1, 1, 1, 1, -1, -1, -1, -1, 2, 2, 2, 2, 2}, 1, 0},
			groupPos{3, 3},
		},
		{
			"2 in 12345",
			args{[]int{0, -1, -1, 1, 1, 1, -1, -1, -1, -1, 2, 2, 2, 2, 2}, 2, 0},
			groupPos{10, 5},
		},
		{
			"first -1 in 12345",
			args{[]int{0, -1, -1, 1, 1, 1, -1, -1, -1, -1, 2, 2, 2, 2, 2}, -1, 0},
			groupPos{1, 2},
		},
		{
			"second -1 in 12345",
			args{[]int{0, -1, -1, 1, 1, 1, -1, -1, -1, -1, 2, 2, 2, 2, 2}, -1, 3},
			groupPos{6, 4},
		},
		{
			"not in 12345",
			args{[]int{0, -1, -1, 1, 1, 1, -1, -1, -1, -1, 2, 2, 2, 2, 2}, 3, 0},
			groupPos{-1, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findGroup(tt.args.input, tt.args.id, tt.args.startAt); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findEmptyGroup(t *testing.T) {
	type args struct {
		input     []int
		minLength int
	}
	tests := []struct {
		name string
		args args
		want groupPos
	}{
		{
			"length of 1 in 12345",
			args{[]int{0, -1, -1, 1, 1, 1, -1, -1, -1, -1, 2, 2, 2, 2, 2}, 1},
			groupPos{1, 2},
		},
		{
			"length of 2 in 12345",
			args{[]int{0, -1, -1, 1, 1, 1, -1, -1, -1, -1, 2, 2, 2, 2, 2}, 2},
			groupPos{1, 2},
		},
		{
			"length of 3 in 12345",
			args{[]int{0, -1, -1, 1, 1, 1, -1, -1, -1, -1, 2, 2, 2, 2, 2}, 3},
			groupPos{6, 4},
		},
		{
			"length of 4 in 12345",
			args{[]int{0, -1, -1, 1, 1, 1, -1, -1, -1, -1, 2, 2, 2, 2, 2}, 4},
			groupPos{6, 4},
		},
		{
			"not found length of 5 in 12345",
			args{[]int{0, -1, -1, 1, 1, 1, -1, -1, -1, -1, 2, 2, 2, 2, 2}, 5},
			groupPos{-1, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findEmptyGroup(tt.args.input, tt.args.minLength); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findEmptyGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}
