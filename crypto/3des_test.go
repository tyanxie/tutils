package crypto

import (
	"reflect"
	"testing"
)

func Test_TripleDesEcbEncrypt(t *testing.T) {
	type args struct {
		key       []byte
		plaintext []byte
		padding   string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "#1",
			args: args{
				key:       []byte("123456789012345678901234"),
				plaintext: []byte("Hello World"),
				padding:   PaddingPkcs7,
			},
			want:    []byte{40, 235, 249, 1, 31, 198, 94, 171, 56, 37, 33, 155, 158, 21, 12, 46},
			wantErr: false,
		},
		{
			name: "#2",
			args: args{
				key:       []byte("123456789012345678901234"),
				plaintext: []byte("Hello World"),
				padding:   PaddingZero,
			},
			want:    []byte{40, 235, 249, 1, 31, 198, 94, 171, 243, 99, 64, 122, 114, 64, 14, 181},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TripleDesEcbEncrypt(tt.args.key, tt.args.plaintext, tt.args.padding)
			if (err != nil) != tt.wantErr {
				t.Errorf("TripleDesEcbEncrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TripleDesEcbEncrypt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_TripleDesEcbDecrypt(t *testing.T) {
	type args struct {
		key        []byte
		ciphertext []byte
		padding    string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "#1",
			args: args{
				key:        []byte("123456789012345678901234"),
				ciphertext: []byte{40, 235, 249, 1, 31, 198, 94, 171, 56, 37, 33, 155, 158, 21, 12, 46},
				padding:    PaddingPkcs7,
			},
			want:    []byte("Hello World"),
			wantErr: false,
		},
		{
			name: "#2",
			args: args{
				key:        []byte("123456789012345678901234"),
				ciphertext: []byte{40, 235, 249, 1, 31, 198, 94, 171, 243, 99, 64, 122, 114, 64, 14, 181},
				padding:    PaddingZero,
			},
			want:    []byte("Hello World"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TripleDesEcbDecrypt(tt.args.key, tt.args.ciphertext, tt.args.padding)
			if (err != nil) != tt.wantErr {
				t.Errorf("TripleDesEcbDecrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TripleDesEcbDecrypt() got = %v, want %v", got, tt.want)
			}
		})
	}
}
