package crypto

import (
	"bytes"
	"strings"
)

/*
填充方式说明：https://www.cnblogs.com/midea0978/articles/1437257.html
经过比较，ANSIX923、ISO10126、PKCS5、PKCS7均可以按照最后一个字节来判断填充长度，
其余填充的字节分别为0（ANSIX923）、随机数（ISO10126）、与最后一个字节相同（PKCS5/PKCS7），
但这些并不影响判断填充长度，因此这四种方法都统一使用PKCS7的实现方式来做。
*/

const (
	PaddingPkcs5    = "pkcs5"    // 填充方式：pkcs5（块大小为8bytes的pkcs7填充）
	PaddingPkcs7    = "pkcs7"    // 填充方式：pkcs7
	PaddingZero     = "zero"     // 填充方式：0填充
	PaddingAnsix923 = "ansix923" // 填充方式：ansix923（除最后一个字节外全部填充0的pkcs7填充）
	PaddingIso10126 = "iso10126" // 填充方式：iso10126（除最后一个字节外全部填充随机数的pkcs7填充）
)

var (
	// paddingMap 填充方式映射
	paddingMap = map[string]func(src []byte, blockSize int) []byte{
		PaddingPkcs5:    pkcs7Padding,
		PaddingPkcs7:    pkcs7Padding,
		PaddingZero:     zeroPadding,
		PaddingAnsix923: pkcs7Padding,
		PaddingIso10126: pkcs7Padding,
	}

	// unPaddingMap 去填充方式映射
	unPaddingMap = map[string]func(src []byte) []byte{
		PaddingPkcs5:    pkcs7UnPadding,
		PaddingPkcs7:    pkcs7UnPadding,
		PaddingZero:     zeroUnPadding,
		PaddingAnsix923: pkcs7UnPadding,
		PaddingIso10126: pkcs7UnPadding,
	}
)

// Padding 填充
// @param padding 填充方式
// @param src 原文
// @param blockSize 块大小
func Padding(padding string, src []byte, blockSize int) []byte {
	if f, ok := paddingMap[strings.ToLower(padding)]; ok && f != nil {
		return f(src, blockSize)
	}
	return src
}

// UnPadding 去填充
// @param src 原文
// @param padding 填充方式
func UnPadding(padding string, src []byte) []byte {
	if f, ok := unPaddingMap[strings.ToLower(padding)]; ok && f != nil {
		return f(src)
	}
	return src
}

// pkcs7Padding pkcs7填充
func pkcs7Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

// pkcs7UnPadding pkcs7去填充
func pkcs7UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	end := length - unpadding
	if end < 0 {
		return src
	}
	return src[:end]
}

// zeroPadding 0填充
func zeroPadding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(src, padtext...)
}

// zeroPadding 0去填充
func zeroUnPadding(src []byte) []byte {
	return bytes.TrimRight(src, string([]byte{0}))
}
