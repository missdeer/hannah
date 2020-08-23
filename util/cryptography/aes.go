package cryptography

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

func pkcs5Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	paddingText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, paddingText...)
}

func pkcs5UnPadding(src []byte) []byte {
	n := len(src)
	unPadding := int(src[n-1])
	return src[:n-unPadding]
}

func AESCBCEncrypt(plainText, key, iv []byte) []byte {
	block, _ := aes.NewCipher(key)
	plainText = pkcs5Padding(plainText, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, iv)
	cipherText := make([]byte, len(plainText))
	blockMode.CryptBlocks(cipherText, plainText)
	return cipherText
}

func AESCBCDecrypt(cipherText, key, iv []byte) []byte {
	block, _ := aes.NewCipher(key)
	blockMode := cipher.NewCBCDecrypter(block, iv)
	plainText := make([]byte, len(cipherText))
	blockMode.CryptBlocks(plainText, cipherText)
	plainText = pkcs5UnPadding(plainText)
	return plainText
}

func AESECBEncrypt(plainText, key []byte) []byte {
	block, _ := aes.NewCipher(key)
	plainText = pkcs5Padding(plainText, block.BlockSize())
	blockMode := NewECBEncrypter(block)
	cipherText := make([]byte, len(plainText))
	blockMode.CryptBlocks(cipherText, plainText)
	return cipherText
}

func AESECBDecrypt(cipherText, key []byte) []byte {
	block, _ := aes.NewCipher(key)
	blockMode := NewECBDecrypter(block)
	plainText := make([]byte, len(cipherText))
	blockMode.CryptBlocks(plainText, cipherText)
	plainText = pkcs5UnPadding(plainText)
	return plainText
}

func AesEncryptCBCWithIv(origData []byte, key []byte, iv []byte) (encrypted []byte) {
	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, _ := aes.NewCipher(key)
	blockSize := block.BlockSize()                 // 获取秘钥块的长度
	origData = pkcs5Padding(origData, blockSize)   // 补全码
	blockMode := cipher.NewCBCEncrypter(block, iv) // 加密模式
	encrypted = make([]byte, len(origData))        // 创建数组
	blockMode.CryptBlocks(encrypted, origData)     // 加密
	return encrypted
}
