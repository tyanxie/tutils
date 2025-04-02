// Package crypto triple-des(3des)加密工具包
package crypto

import (
	"crypto/des"
	"fmt"
)

// TripleDesEcbEncrypt ecb模式的3des加密
// @param key 加密key
// @param plaintext 明文
// @param padding 填充模式
func TripleDesEcbEncrypt(key, plaintext []byte, padding string) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create block failed: %w", err)
	}
	plaintext = Padding(padding, plaintext, block.BlockSize())
	ciphertext, err := EcbEncrypt(block, plaintext)
	if err != nil {
		return nil, fmt.Errorf("ecb encrypt failed: %w", err)
	}
	return ciphertext, nil
}

// TripleDesEcbDecrypt ecb模式的3des解密
// @param key 加密key
// @param ciphertext 密文
// @param padding 填充模式
func TripleDesEcbDecrypt(key, ciphertext []byte, padding string) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create block failed: %w", err)
	}
	plaintext, err := EcbDecrypt(block, ciphertext)
	if err != nil {
		return nil, fmt.Errorf("ecb decrypt failed: %w", err)
	}
	return UnPadding(padding, plaintext), nil
}
