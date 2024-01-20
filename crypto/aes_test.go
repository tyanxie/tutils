package crypto

import (
	"reflect"
	"testing"
)

func Test_AesCbcEncrypt(t *testing.T) {
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
			want:    []byte{3, 172, 4, 45, 97, 108, 114, 100, 116, 62, 29, 14, 15, 236, 210, 220},
			wantErr: false,
		},
		{
			name: "#2",
			args: args{
				key:       []byte("1234567890123456"),
				iv:        []byte("0987654321098765"),
				plaintext: []byte("Hello World"),
				padding:   PaddingZero,
			},
			want:    []byte{25, 249, 95, 102, 144, 132, 154, 131, 26, 55, 42, 159, 250, 239, 53, 245},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AesCbcEncrypt(tt.args.key, tt.args.iv, tt.args.plaintext, tt.args.padding)
			if (err != nil) != tt.wantErr {
				t.Errorf("AesCbcEncrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AesCbcEncrypt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_AesCbcDecrypt(t *testing.T) {
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
				ciphertext: []byte{3, 172, 4, 45, 97, 108, 114, 100, 116, 62, 29, 14, 15, 236, 210, 220},
				padding:    PaddingPkcs7,
			},
			want:    []byte("Hello World"),
			wantErr: false,
		},
		{
			name: "#2",
			args: args{
				key:        []byte("1234567890123456"),
				iv:         []byte("0987654321098765"),
				ciphertext: []byte{25, 249, 95, 102, 144, 132, 154, 131, 26, 55, 42, 159, 250, 239, 53, 245},
				padding:    PaddingZero,
			},
			want:    []byte("Hello World"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AesCbcDecrypt(tt.args.key, tt.args.iv, tt.args.ciphertext, tt.args.padding)
			if (err != nil) != tt.wantErr {
				t.Errorf("AesCbcDecrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AesCbcDecrypt() got = %v, want %v", got, tt.want)
			}
		})
	}
}
