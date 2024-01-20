package crypto

import (
	"encoding/base64"
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"
)

var (
	privateKeyBytes = []byte{106, 107, 32, 123, 12, 118, 197, 245, 11, 198, 177, 106, 208, 204, 28,
		98, 176, 154, 172, 22, 30, 130, 187, 219, 31, 192, 157, 145, 250, 36, 243, 132}
	publicKeyBytes = []byte{232, 189, 14, 147, 48, 62, 31, 120, 123, 112, 68, 197, 173, 115, 227,
		192, 36, 120, 157, 127, 144, 62, 188, 170, 78, 35, 172, 88, 8, 77, 249, 111, 73, 218, 90,
		161, 46, 179, 90, 24, 109, 120, 188, 12, 12, 194, 149, 212, 210, 210, 187, 99, 167, 108,
		45, 45, 70, 195, 235, 237, 44, 230, 86, 0}

	privateKey *sm2.PrivateKey
	publicKey  *sm2.PublicKey
)

func init() {
	var err error
	privateKey, err = x509.ReadPrivateKeyFromHex(hex.EncodeToString(privateKeyBytes))
	if err != nil {
		panic(err)
	}
	publicKey, err = x509.ReadPublicKeyFromHex(hex.EncodeToString(publicKeyBytes))
	if err != nil {
		panic(err)
	}
}

func Test_CreateSm2PrivateKeyWithBase64(t *testing.T) {
	tests := []struct {
		name       string
		privateKey string
		want       *sm2.PrivateKey
		wantErr    bool
	}{
		{
			name:       "#1",
			privateKey: base64.StdEncoding.EncodeToString(privateKeyBytes),
			want:       privateKey,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateSm2PrivateKeyWithBase64(tt.privateKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateSm2PrivateKeyWithBase64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateSm2PrivateKeyWithBase64() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_CreateSm2PublicKeyWithBase64(t *testing.T) {
	tests := []struct {
		name      string
		publicKey string
		want      *sm2.PublicKey
		wantErr   bool
	}{
		{
			name:      "#1",
			publicKey: base64.StdEncoding.EncodeToString(publicKeyBytes),
			want:      publicKey,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateSm2PublicKeyWithBase64(tt.publicKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateSm2PublicKeyWithBase64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateSm2PublicKeyWithBase64() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Sm2EncryptDecrypt(t *testing.T) {
	plaintext := []byte("Hello World")
	tests := []struct {
		name string
		mode int
	}{
		{name: "C1C3C2", mode: sm2.C1C3C2},
		{name: "C1C2C3", mode: sm2.C1C2C3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ciphertext, err := Sm2Encrypt(publicKey, plaintext, tt.mode)
			if err != nil {
				t.Errorf("Sm2Encrypt() error = %v", err)
				return
			}
			newPlaintext, err := Sm2Decrypt(privateKey, ciphertext, tt.mode)
			if err != nil {
				t.Errorf("Sm2Decrypt() error = %v", err)
				return
			}
			if !reflect.DeepEqual(plaintext, newPlaintext) {
				t.Errorf("Sm2Encrypt() and Sm2Decrypt() got = %v, want %v", newPlaintext, plaintext)
			}
		})
	}
}
