package main

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

func Test_parseOperations(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    []operationType
		wantErr bool
	}{
		{"one", args{strings.NewReader("xmul(2,4)%&mul[3,7]!@^do_not_mul(5,5)+mul(32,64]then(mul(11,8)mul(8,5))")}, []operationType{{2, 4}, {5, 5}, {11, 8}, {8, 5}}, false},
		{"two", args{strings.NewReader("xmul(2,4)&mul[3,7]!^don't()_mul(5,5)+mul(32,64](mul(11,8)undo()?mul(8,5))")}, []operationType{{2, 4}, {8, 5}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseOperations(tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseOperations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseOperations() got = %v, want %v", got, tt.want)
			}
		})
	}
}
