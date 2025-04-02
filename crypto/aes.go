// Package crypto aes加密工具包
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

// AesCbcEncrypt cbc模式的aes加密
// @param key 密钥
// @param iv 初始偏移向量
// @param plaintext 明文
// @param padding 填充方式
func AesCbcEncrypt(key, iv, plaintext []byte, padding string) (ciphertext []byte, err error) {
	defer func() {
		if rc := recover(); rc != nil {
			err = fmt.Errorf("encrypt blocks failed: %+v", rc)
		}
	}()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create block failed: %w", err)
	}
	plaintext = Padding(padding, plaintext, block.BlockSize())
	encrypter := cipher.NewCBCEncrypter(block, iv)
	ciphertext = make([]byte, len(plaintext))
	encrypter.CryptBlocks(ciphertext, plaintext)
	return ciphertext, nil
}

// AesCbcDecrypt cbc模式的aes解密
// @param key 密钥
// @param iv 初始偏移向量
// @param ciphertext 密文
// @param padding 填充方式
func AesCbcDecrypt(key, iv, ciphertext []byte, padding string) (plaintext []byte, err error) {
	defer func() {
		if rc := recover(); rc != nil {
			err = fmt.Errorf("decrypt blocks failed: %+v", rc)
		}
	}()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create block failed: %w", err)
	}
	decrypter := cipher.NewCBCDecrypter(block, iv)
	plaintext = make([]byte, len(ciphertext))
	decrypter.CryptBlocks(plaintext, ciphertext)
	return UnPadding(padding, plaintext), nil
}
