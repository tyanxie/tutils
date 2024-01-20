package crypto

import (
	"reflect"
	"testing"
)

func TestPadding(t *testing.T) {
	type args struct {
		padding   string
		src       []byte
		blockSize int
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "pkcs7#1",
			args: args{
				padding:   PaddingPkcs7,
				src:       []byte{0, 1, 2, 3, 4, 5, 6},
				blockSize: 16,
			},
			want: []byte{0, 1, 2, 3, 4, 5, 6, 9, 9, 9, 9, 9, 9, 9, 9, 9},
		},
		{
			name: "pkcs7#2",
			args: args{
				padding:   PaddingPkcs7,
				src:       []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
				blockSize: 16,
			},
			want: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 16, 16, 16, 16, 16, 16,
				16, 16, 16, 16, 16, 16, 16, 16, 16},
		},
		{
			name: "zero#1",
			args: args{
				padding:   PaddingZero,
				src:       []byte{0, 1, 2, 3, 4},
				blockSize: 8,
			},
			want: []byte{0, 1, 2, 3, 4, 0, 0, 0},
		},
		{
			name: "zero#2",
			args: args{
				padding:   PaddingZero,
				src:       []byte{0, 1, 2, 3, 4, 5, 6, 7},
				blockSize: 8,
			},
			want: []byte{0, 1, 2, 3, 4, 5, 6, 7, 0, 0, 0, 0, 0, 0, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Padding(tt.args.padding, tt.args.src, tt.args.blockSize); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Padding() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnPadding(t *testing.T) {
	type args struct {
		padding string
		src     []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "pkcs7#1",
			args: args{
				padding: PaddingPkcs7,
				src:     []byte{0, 1, 2, 3, 4, 5, 6, 9, 9, 9, 9, 9, 9, 9, 9, 9},
			},
			want: []byte{0, 1, 2, 3, 4, 5, 6},
		},
		{
			name: "pkcs7#2",
			args: args{
				padding: PaddingPkcs7,
				src: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 16, 16, 16, 16, 16, 16,
					16, 16, 16, 16, 16, 16, 16, 16, 16},
			},
			want: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		},
		{
			name: "zero#1",
			args: args{
				padding: PaddingZero,
				src:     []byte{0, 1, 2, 3, 4, 0, 0, 0},
			},
			want: []byte{0, 1, 2, 3, 4},
		},
		{
			name: "zero#2",
			args: args{
				padding: PaddingZero,
				src:     []byte{0, 1, 2, 3, 4, 5, 6, 7, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			want: []byte{0, 1, 2, 3, 4, 5, 6, 7},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnPadding(tt.args.padding, tt.args.src); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnPadding() = %v, want %v", got, tt.want)
			}
		})
	}
}
