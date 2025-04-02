// Package crypto sm4国密对称加密工具包
package crypto

import (
	"crypto/cipher"
	"fmt"

	"github.com/tjfoc/gmsm/sm4"
)

// Sm4CbcEncrypt cbc模式的sm4加密
// @param key 密钥
// @param iv 初始向量
// @param plaintext 明文内容
// @param padding 填充方式
func Sm4CbcEncrypt(key, iv, plaintext []byte, padding string) (ciphertext []byte, err error) {
	defer func() {
		if rc := recover(); rc != nil {
			err = fmt.Errorf("encrypt blocks failed: %+v", rc)
		}
	}()
	block, err := sm4.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create sm4 ciphter failed, err: %w", err)
	}
	plaintext = Padding(padding, plaintext, block.BlockSize())
	encrypter := cipher.NewCBCEncrypter(block, iv)
	ciphertext = make([]byte, len(plaintext))
	encrypter.CryptBlocks(ciphertext, plaintext)
	return ciphertext, nil
}

// Sm4CbcDecrypt cbc模式的sm4解密
// @param key 密钥
// @param iv 初始向量
// @param ciphertext 密文
// @param padding 填充方式
func Sm4CbcDecrypt(key, iv, ciphertext []byte, padding string) (plaintext []byte, err error) {
	defer func() {
		if rc := recover(); rc != nil {
			err = fmt.Errorf("decrypt blocks failed: %+v", rc)
		}
	}()
	block, err := sm4.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create sm4 ciphter failed: %w", err)
	}
	decrypter := cipher.NewCBCDecrypter(block, iv)
	plaintext = make([]byte, len(ciphertext))
	decrypter.CryptBlocks(plaintext, ciphertext)
	return UnPadding(padding, plaintext), nil
}
