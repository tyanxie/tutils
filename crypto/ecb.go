// Package crypto ecb加密模式工具包
package crypto

import (
	"crypto/cipher"
	"errors"
)

// EcbEncrypt ecb模式加密
func EcbEncrypt(block cipher.Block, plaintext []byte) ([]byte, error) {
	if len(plaintext)%block.BlockSize() != 0 {
		return nil, errors.New("plaintext not full blocks")
	}
	blockSize := block.BlockSize()
	ciphertext := make([]byte, len(plaintext))
	for start := 0; start < len(plaintext); start += blockSize {
		end := start + blockSize
		block.Encrypt(ciphertext[start:], plaintext[start:end])
	}
	return ciphertext, nil
}

// EcbDecrypt ecb模式解密
func EcbDecrypt(block cipher.Block, ciphertext []byte) ([]byte, error) {
	if len(ciphertext)%block.BlockSize() != 0 {
		return nil, errors.New("ciphertext not full blocks")
	}
	blockSize := block.BlockSize()
	plaintext := make([]byte, len(ciphertext))
	for start := 0; start < len(plaintext); start += blockSize {
		end := start + blockSize
		block.Decrypt(plaintext[start:], ciphertext[start:end])
	}
	return plaintext, nil
}
