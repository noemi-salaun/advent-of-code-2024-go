package main

import "testing"

var testGrid = rows{
	{'M', 'M', 'M', 'S', 'X', 'X', 'M', 'A', 'S', 'M'},
	{'M', 'S', 'A', 'M', 'X', 'M', 'S', 'M', 'S', 'A'},
	{'A', 'M', 'X', 'S', 'X', 'M', 'A', 'A', 'M', 'M'},
	{'M', 'S', 'A', 'M', 'A', 'S', 'M', 'S', 'M', 'X'},
	{'X', 'M', 'A', 'S', 'A', 'M', 'X', 'A', 'M', 'M'},
	{'X', 'X', 'A', 'M', 'M', 'X', 'X', 'A', 'M', 'A'},
	{'S', 'M', 'S', 'M', 'S', 'A', 'S', 'X', 'S', 'S'},
	{'S', 'A', 'X', 'A', 'M', 'A', 'S', 'A', 'A', 'A'},
	{'M', 'A', 'M', 'M', 'M', 'X', 'M', 'M', 'M', 'M'},
	{'M', 'X', 'M', 'X', 'A', 'X', 'M', 'A', 'S', 'X'},
}

func Test_countXMAS(t *testing.T) {
	type args struct {
		rows rows
		x    int
		y    int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"noop1", args{testGrid, 0, 0}, 0},
		{"north & north west", args{testGrid, 9, 9}, 2},
		{"east", args{testGrid, 0, 4}, 1},
		{"north & south", args{testGrid, 9, 3}, 2},
		{"west", args{testGrid, 4, 1}, 1},
		{"north & west", args{testGrid, 6, 4}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := countXMAS(tt.args.rows, tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("countXMAS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkX_MAS(t *testing.T) {
	type args struct {
		rows rows
		x    int
		y    int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"noop1", args{testGrid, 5, 2}, false},
		{"2,1", args{testGrid, 2, 1}, true},
		{"6,2", args{testGrid, 6, 2}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkX_MAS(tt.args.rows, tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("checkX_MAS() = %v, want %v", got, tt.want)
			}
		})
	}
}
