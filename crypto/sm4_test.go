package crypto

import (
	"reflect"
	"testing"
)

func Test_Sm4CbcEncrypt(t *testing.T) {
	type args struct {
		key       []byte
		iv        []byte
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
				key:       []byte("1234567890123456"),
				iv:        []byte("0987654321098765"),
				plaintext: []byte("Hello World"),
				padding:   PaddingPkcs7,
			},
			want: []byte{87, 171, 77, 110, 248, 173, 25, 57, 80, 88, 14, 120, 187, 122, 4, 24},
		},
		{
			name: "#2",
			args: args{
				key:       []byte("1234567890123456"),
				iv:        []byte("0987654321098765"),
				plaintext: []byte("Hello World"),
				padding:   PaddingZero,
			},
			want: []byte{154, 52, 83, 147, 136, 81, 255, 254, 236, 231, 251, 151, 109, 188, 206, 41},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Sm4CbcEncrypt(tt.args.key, tt.args.iv, tt.args.plaintext, tt.args.padding)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sm4CbcEncrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sm4CbcEncrypt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Sm4CbcDecrypt(t *testing.T) {
	type args struct {
		key        []byte
		iv         []byte
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
				key:        []byte("1234567890123456"),
				iv:         []byte("0987654321098765"),
				ciphertext: []byte{87, 171, 77, 110, 248, 173, 25, 57, 80, 88, 14, 120, 187, 122, 4, 24},
				padding:    PaddingPkcs7,
			},
			want: []byte("Hello World"),
		},
		{
			name: "#2",
			args: args{
				key:        []byte("1234567890123456"),
				iv:         []byte("0987654321098765"),
				ciphertext: []byte{154, 52, 83, 147, 136, 81, 255, 254, 236, 231, 251, 151, 109, 188, 206, 41},
				padding:    PaddingZero,
			},
			want: []byte("Hello World"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Sm4CbcDecrypt(tt.args.key, tt.args.iv, tt.args.ciphertext, tt.args.padding)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sm4CbcDecrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sm4CbcDecrypt() got = %v, want %v", got, tt.want)
			}
		})
	}
}
