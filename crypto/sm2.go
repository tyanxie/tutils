// Package crypto sm2国密非对称加密工具包
package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"
)

/*
本仓库使用开源库：github.com/tjfoc/gmsm

该开源库中如果希望将sm2公私钥导出为字节数组、16进制字符串、base64字符串等，需要使用到该库的`x509`包的类似如下函数：
x509.ReadPrivateKeyFromHex
x509.WritePrivateKeyToHex
x509.ReadPublicKeyFromHex
x509.WritePublicKeyToHex

本仓库的CreateSm2PrivateKeyWithBase64和CreateSm2PublicKeyWithBase64以及sm2_test.go的init函数展示了相关用法
*/

// CreateSm2PrivateKeyWithBase64 通过base64编码字符串构造sm2私钥
func CreateSm2PrivateKeyWithBase64(privateKey string) (*sm2.PrivateKey, error) {
	// base64解析私钥
	k, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return nil, fmt.Errorf("decode with base64 failed: %w", err)
	}
	if len(k) != 32 {
		return nil, errors.New("private key bytes length must be 32")
	}
	// 返回结果
	return x509.ReadPrivateKeyFromHex(hex.EncodeToString(k))
}

// CreateSm2PublicKeyWithBase64 通过base64编码字符串构造sm2公钥
func CreateSm2PublicKeyWithBase64(publicKey string) (*sm2.PublicKey, error) {
	// base64解析公钥
	k, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, fmt.Errorf("decode with base64 failed: %w", err)
	}
	if len(k) != 64 {
		return nil, errors.New("public key bytes length must be 64")
	}
	// 返回结果
	return x509.ReadPublicKeyFromHex(hex.EncodeToString(k))
}

// Sm2Encrypt sm2加密
// @param publicKey 公钥
// @param plaintext 明文内容
// @param mode 密文顺序，参考github.com/tjfoc/gmsm/sm2包的枚举值，0为C1C3C2，1为C1C2C3
func Sm2Encrypt(publicKey *sm2.PublicKey, plaintext []byte, mode int) ([]byte, error) {
	ciphertext, err := sm2.Encrypt(publicKey, plaintext, rand.Reader, mode)
	if err != nil {
		return nil, fmt.Errorf("encrypt failed: %w", err)
	}
	ciphertext, err = sm2.CipherMarshal(ciphertext)
	if err != nil {
		return nil, fmt.Errorf("marshal cipher failed: %w", err)
	}
	return ciphertext, nil
}

// Sm2Decrypt sm2解密
// @param privateKey 私钥
// @param ciphertext 密文
// @param mode 密文顺序，参考github.com/tjfoc/gmsm/sm2包的枚举值，0为C1C3C2，1为C1C2C3
func Sm2Decrypt(privateKey *sm2.PrivateKey, ciphertext []byte, mode int) ([]byte, error) {
	ciphertext, err := sm2.CipherUnmarshal(ciphertext)
	if err != nil {
		return nil, fmt.Errorf("unmarshal cipher failed: %w", err)
	}
	plaintext, err := sm2.Decrypt(privateKey, ciphertext, mode)
	if err != nil {
		return nil, fmt.Errorf("decrypt failed: %w", err)
	}
	return plaintext, err
}
